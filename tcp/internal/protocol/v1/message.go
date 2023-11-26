package v1

import (
	"fmt"
	"time"
)

type Message struct {
	msg      string
	SendTime time.Time
}

func CreateMsg(msg string) *Message {
	return &Message{
		msg:      msg,
		SendTime: time.Now(),
	}
}

func (m *Message) Msg(nick string) string {
	sendTimeStr := m.SendTime.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s#%s -> %s", sendTimeStr, nick, m.msg)
}
