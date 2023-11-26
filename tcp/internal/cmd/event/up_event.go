package event

import (
	"net"
	"smallchat/tcp/internal/cmd"
	"smallchat/tcp/internal/protocol/v2"
	"time"
)

type UpEvent struct {
	conn      net.Conn
	msg       *v2.Msg
	startTime int64
}

func CreateUpEvent(conn net.Conn, msg *v2.Msg) cmd.Event {
	return &UpEvent{
		conn:      conn,
		msg:       msg,
		startTime: time.Now().Unix(),
	}
}

func (u *UpEvent) GetMsg() *v2.Msg {
	return u.msg
}

func (u *UpEvent) GetStartTime() int64 {
	return u.startTime
}

func (u *UpEvent) GetConn() net.Conn {
	return u.conn
}

func (u *UpEvent) Up() bool {
	return true
}
