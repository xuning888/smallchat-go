package cmd

import (
	"net"
	"os"
)

type Exit struct {
}

func (e *Exit) Execute(conn net.Conn) {
	os.Exit(0)
}

func (e *Exit) Description() string {
	return "关闭客户端"
}
