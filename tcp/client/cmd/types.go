package cmd

import "net"

type Command interface {
	Execute(conn net.Conn)
	Description() string
}
