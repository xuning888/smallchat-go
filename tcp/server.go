package tcp

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
	Clients   []*ChatClient
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
		Clients:   make([]*ChatClient, 0),
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
		for _, c := range s.Clients {
			_ = c.Close()
		}
	}()
	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		conn, err := s.server.Accept()
		if err != nil {
			errCh <- err
			break
		}
		log.Printf("accept link %s\n", conn.RemoteAddr().String())
		numClient := len(s.Clients)
		if numClient == s.maxClient {
			_, err = conn.Write([]byte(fmt.Sprintf("can't accept")))
			if err != nil {
				log.Printf("can't accept, wirte msg err:%s", err)
			}
			_ = conn.Close()
		}
		log.Printf("connected %v\n", conn.RemoteAddr().String())
		client := CreateClient(ctx, conn, "", s)
		n, err := client.WriteMsg(fmt.Sprintf("Welcome to Simple Chat! Use /nick <nick> to set your nick!\n"))
		if err != nil {
			_ = client.Close()
			log.Printf("write %d byte to %s", n, conn.RemoteAddr().String())
		}
		s.Clients = append(s.Clients, client)
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			client.Handle()
		}()
	}
	waitDone.Wait()
}

func (s *ChatServer) RemoveC(c *ChatClient) {
	clients := s.Clients
	if clients == nil {
		return
	}
	rIndex := -1
	for idx, client := range s.Clients {
		if client == c {
			rIndex = idx
			break
		}
	}
	if rIndex != -1 {
		s.Clients[rIndex] = nil
		left := clients[:rIndex]
		right := clients[rIndex+1:]
		temp := make([]*ChatClient, 0)
		temp = append(temp, left...)
		temp = append(temp, right...)
		s.Clients = temp
	}
}
