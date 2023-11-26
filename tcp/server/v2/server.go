package v2

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
	Handlers  map[string]*ChatHandler
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
		Handlers:  make(map[string]*ChatHandler),
		closeChan: closeChan,
	}
	return server
}

func (s *ChatServer) RemoveHandler(addr string) {
	_, ok := s.Handlers[addr]
	if ok {
		delete(s.Handlers, addr)
	}
}

func (s *ChatServer) GetChatHandler(addr string) (handler *ChatHandler, exist bool) {
	handler, ok := s.Handlers[addr]
	if ok {
		return handler, true
	}
	return nil, false
}

func (s *ChatServer) addChatHandler(addr string, handler *ChatHandler) {
	s.Handlers[addr] = handler
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
		handler := CreateHandler(ctx, conn, s)
		n, err := handler.WriteStringMsg(fmt.Sprintf("Welcome to Simple Chat! Use /nick <nick> to set your nick!\n"))
		if err != nil {
			_ = handler.Close()
			log.Printf("write %d byte to %s", n, conn.RemoteAddr().String())
		}
		s.addChatHandler(conn.RemoteAddr().String(), handler)
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
