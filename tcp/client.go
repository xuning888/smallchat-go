package tcp

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"sort"
	"strings"
	"time"
)

type message struct {
	msg      string
	sendTime time.Time
}

func createMsg(msg string) *message {
	return &message{
		msg:      msg,
		sendTime: time.Now(),
	}
}

func (m *message) Msg(nick string) string {
	sendTimeStr := m.sendTime.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s#%s -> %s", sendTimeStr, nick, m.msg)
}

type msgBox struct {
	messages []*message
	size     int
}

func createMsgBox(size int) *msgBox {
	size = int(math.Max(float64(size), 10))
	return &msgBox{
		messages: make([]*message, 0),
		size:     size,
	}
}

func (mbx *msgBox) pushMsg(msg string) {
	curLen := len(mbx.messages)
	if curLen == mbx.size {
		return
	}
	curMsg := createMsg(msg)
	mbx.messages = append(mbx.messages, curMsg)
}

func (mbx *msgBox) pullMsg() *message {
	messages := mbx.messages
	msgLen := len(messages)
	if msgLen == 0 {
		return nil
	}
	// 按照时间排序
	sort.Slice(mbx.messages, func(i, j int) bool {
		return mbx.messages[i].sendTime.Before(mbx.messages[j].sendTime)
	})
	curMg := mbx.messages[0]
	mbx.messages = mbx.messages[1:]
	return curMg
}

type ChatClient struct {
	Ctx    context.Context
	Conn   net.Conn
	Nick   string
	server *ChatServer
	box    *msgBox
}

func CreateClient(ctx context.Context,
	conn net.Conn, nick string, server *ChatServer) *ChatClient {
	return &ChatClient{
		Ctx:    ctx,
		Conn:   conn,
		Nick:   nick,
		server: server,
		box:    createMsgBox(20),
	}
}

func (c *ChatClient) Close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}
	return nil
}

func (c *ChatClient) WriteMsg(msg string) (int, error) {
	return c.Conn.Write([]byte(msg))
}

func (c *ChatClient) Handle() {
	reader := bufio.NewReader(c.Conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			c.server.RemoveC(c)
			log.Printf("%s connection close, clients:%d\n", c.Conn.RemoteAddr().String(), len(c.server.Clients))
			_ = c.Close()
			return
		}
		if strings.HasPrefix(msg, "/nick") {
			segs := strings.Split(msg, " ")
			newNick := strings.TrimSpace(segs[1])
			oldNick := c.Nick
			c.Nick = newNick
			_, _ = c.WriteMsg(fmt.Sprintf("set nick success nickName: %s\n", newNick))
			c.box.pushMsg(fmt.Sprintf("%s reset nick from %s to %s\n", c.Nick, oldNick, newNick))
		} else {
			c.box.pushMsg(msg)
		}
		c.SendMsgToAll()
	}
}

func (c *ChatClient) SendMsgToAll() {
	msg := c.box.pullMsg()
	if msg == nil {
		return
	}
	sendMsg := msg.Msg(c.Nick)
	clients := c.server.Clients
	if clients != nil {
		for _, client := range clients {
			if client != nil && client != c {
				_, _ = client.WriteMsg(sendMsg)
			}
		}
	}
}
