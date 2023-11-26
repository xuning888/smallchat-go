package v1

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"smallchat/tcp/internal/box"
	"strings"
)

type ChatHandler struct {
	Ctx    context.Context
	Conn   net.Conn
	Nick   string
	server *ChatServer
	box    *box.MessageBox
}

func CreateHandler(ctx context.Context,
	conn net.Conn, nick string, server *ChatServer) *ChatHandler {
	return &ChatHandler{
		Ctx:    ctx,
		Conn:   conn,
		Nick:   nick,
		server: server,
		box:    box.CreateMsgBox(20),
	}
}

func (c *ChatHandler) Close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}
	return nil
}

func (c *ChatHandler) WriteMsg(msg string) (int, error) {
	return c.Conn.Write([]byte(msg))
}

func (c *ChatHandler) Handle() {
	reader := bufio.NewReader(c.Conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			c.server.RemoveC(c)
			log.Printf("%s connection close, clients:%d\n", c.Conn.RemoteAddr().String(), len(c.server.Handlers))
			_ = c.Close()
			return
		}
		if strings.HasPrefix(msg, "/nick") {
			segs := strings.Split(msg, " ")
			if len(segs) < 2 {
				return
			}
			newNick := strings.TrimSpace(segs[1])
			oldNick := c.Nick
			c.Nick = newNick
			_, _ = c.WriteMsg(fmt.Sprintf("set nick success nickName: %s\n", newNick))
			c.box.PushMsg(fmt.Sprintf("%s reset nick from %s to %s\n", c.Nick, oldNick, newNick))
		} else {
			c.box.PushMsg(msg)
		}
		c.SendMsgToAll()
	}
}

func (c *ChatHandler) SendMsgToAll() {
	length := c.box.Len()
	for i := 0; i < length; i++ {
		msg := c.box.PopMsg()
		if msg == nil {
			return
		}
		sendMsg := msg.Msg(c.Nick)
		clients := c.server.Handlers
		sendSuccess := false
		if clients != nil {
			for _, client := range clients {
				if client != nil && client != c {
					_, _ = client.WriteMsg(sendMsg)
					sendSuccess = true
				}
			}
		}
		if !sendSuccess {
			c.box.PushMsgWithTime(msg)
		}
	}
}
