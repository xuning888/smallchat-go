package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	cmd2 "smallchat/tcp/client/cmd"
	"syscall"
)

type SmallChatClient struct {
	addr      string
	port      int
	conn      net.Conn
	commands  map[string]cmd2.Command
	closeChan chan struct{}
}

func CreateChatClient(addr string, port int) *SmallChatClient {
	closeChan := createCloseChan()
	conn := createTCP(addr, port)
	return &SmallChatClient{
		addr:      addr,
		port:      port,
		conn:      conn,
		commands:  make(map[string]cmd2.Command),
		closeChan: closeChan,
	}
}

func createTCP(address string, port int) net.Conn {
	serveAddr := fmt.Sprintf("%s:%d", address, port)
	listen, err := net.Dial("tcp", serveAddr)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to server at %s:%d, error: %v", address, port, err)
		panic(errMsg)
	}
	log.Printf("Connected to server at %s:%d", address, port)
	return listen
}

func createCloseChan() chan struct{} {
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

func (s *SmallChatClient) handleUserInput() {
	scanner := bufio.NewScanner(os.Stdin)
	log.Printf("please input:")
	for scanner.Scan() {
		cmd := scanner.Text()
		command, ok := s.GetCmd(cmd)
		if ok {
			command.Execute(s.conn)
		} else {
			log.Printf("can't found cmd:%s", cmd)
		}
	}
}

func (s *SmallChatClient) RegisterCmd(cmd string, command cmd2.Command) {
	s.commands[cmd] = command
}

func (h *SmallChatClient) GetCmd(cmd string) (cmd2.Command, bool) {
	command, ok := h.commands[cmd]
	if ok {
		return command, ok
	}
	return nil, false
}

func (s *SmallChatClient) Spin() {
	go func() {
		s.handleUserInput()
	}()
	select {
	case <-s.closeChan:
		log.Printf("get exit singnal\n")
		log.Printf("close...")
		_ = s.conn.Close()
		return
	}
}
