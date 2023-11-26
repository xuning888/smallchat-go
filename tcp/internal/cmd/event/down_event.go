package event

import (
	"net"
	"smallchat/tcp/internal/cmd"
	"smallchat/tcp/internal/protocol/v2"
	"time"
)

type DownEvent struct {
	conn      net.Conn
	msg       *v2.Msg
	startTime int64
}

func CreateDownEvent(conn net.Conn, msg *v2.Msg) cmd.Event {
	return &DownEvent{
		conn:      conn,
		msg:       msg,
		startTime: time.Now().Unix(),
	}
}

func (d *DownEvent) GetMsg() *v2.Msg {
	return d.msg
}

func (d *DownEvent) GetStartTime() int64 {
	return d.startTime
}

func (d *DownEvent) GetConn() net.Conn {
	return d.conn
}

func (d *DownEvent) Up() bool {
	return false
}
