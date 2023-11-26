package v1

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ChatServer struct {
	Port      int
	maxClient int
	Address   string
	server    net.Listener
	Handlers  []*ChatHandler
	closeChan <-chan struct{}
}

func CreateServer(address string, port int, maxClient int) *ChatServer {
	closeChan := createCloseChan()
	tcpListener := createTCP(address, port)
	server := &ChatServer{
		Port:      port,
		maxClient: maxClient,
		Address:   address,
		server:    tcpListener,
		Handlers:  make([]*ChatHandler, 0),
		closeChan: closeChan,
	}
	return server
}

func createTCP(address string, port int) net.Listener {
	serveAddr := fmt.Sprintf("%s:%d", address, port)
	listen, err := net.Listen("tcp", serveAddr)
	if err != nil {
		panic(fmt.Sprintf("create listener has err:%s", err.Error()))
	}
	log.Printf("bind: %d, start listening...", port)
	return listen
}

func createCloseChan() <-chan struct{} {
	closeChan := make(chan struct{})
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()
	return closeChan
}

func (s *ChatServer) Spin() {
	errCh := make(chan error, 1)
	defer func() {
		close(errCh)
	}()
	go func() {
		select {
		case <-s.closeChan:
			log.Printf("get exit singnal\n")
		case err := <-errCh:
			log.Printf("accept error:%s\n", err.Error())
		}
		log.Printf("shutting down...")
		_ = s.server.Close()
		for _, c := range s.Handlers {
			_ = c.Close()
		}
	}()
	var waitDone sync.WaitGroup
	for {
		conn, err := s.server.Accept()
		if err != nil {
			errCh <- err
			break
		}
		log.Printf("accept link %s\n", conn.RemoteAddr().String())
		numClient := len(s.Handlers)
		if numClient == s.maxClient {
			_, err = conn.Write([]byte(fmt.Sprintf("can't accept\n")))
			if err != nil {
				log.Printf("can't accept, wirte msg err:%s", err)
			}
			_ = conn.Close()
		}
		ctx := context.Background()
		log.Printf("connected %v\n", conn.RemoteAddr().String())
		handler := CreateHandler(ctx, conn, conn.RemoteAddr().String(), s)
		n, err := handler.WriteMsg(fmt.Sprintf("Welcome to Simple Chat! Use /nick <nick> to set your nick!\n"))
		if err != nil {
			_ = handler.Close()
			log.Printf("write %d byte to %s", n, conn.RemoteAddr().String())
		}
		s.Handlers = append(s.Handlers, handler)
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle()
		}()
	}
	waitDone.Wait()
}

func (s *ChatServer) RemoveC(c *ChatHandler) {
	clients := s.Handlers
	if clients == nil {
		return
	}
	rIndex := -1
	for idx, client := range s.Handlers {
		if client == c {
			rIndex = idx
			break
		}
	}
	if rIndex != -1 {
		s.Handlers[rIndex] = nil
		left := clients[:rIndex]
		right := clients[rIndex+1:]
		temp := make([]*ChatHandler, 0)
		temp = append(temp, left...)
		temp = append(temp, right...)
		s.Handlers = temp
	}
}
