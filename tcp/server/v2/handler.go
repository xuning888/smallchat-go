package v2

import (
	"bufio"
	"context"
	"log"
	"net"
	"smallchat/tcp/internal/protocol/v2"
)

type ChatHandler struct {
	Ctx    context.Context
	Conn   net.Conn
	Addr   string
	server *ChatServer
}

func CreateHandler(ctx context.Context,
	conn net.Conn, server *ChatServer) *ChatHandler {
	return &ChatHandler{
		Ctx:    ctx,
		Conn:   conn,
		Addr:   conn.RemoteAddr().String(),
		server: server,
	}
}

func (c *ChatHandler) WriteStringMsg(msg string) (int, error) {
	return c.Conn.Write([]byte(msg))
}

func (c *ChatHandler) WriteBytes(msg []byte) (int, error) {
	return c.Conn.Write(msg)
}

func (c *ChatHandler) Close() error {
	return c.Conn.Close()
}

func (c *ChatHandler) Handle() {
	reader := bufio.NewReader(c.Conn)
	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("err:%s", err)
		}
		message, err := v2.CreateMsgFromBytes(msg)
		if err != nil {
			log.Printf("err:%s", err)

		} else {
			log.Printf("message: %v\n", message)
		}
	}
}
