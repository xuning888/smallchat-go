package cmd

import (
	"log"
	"net"
)

type Help struct {
	commonds map[string]Command
}

func CreateHelp(cmds map[string]Command) *Help {
	return &Help{
		commonds: cmds,
	}
}

func (h *Help) register(cmd string, command Command) {
	log.Printf("register cmd %s", cmd)
}

func (h *Help) Execute(conn net.Conn) {
	for cmd, cmdExe := range h.commonds {
		log.Printf("%s: %s\n", cmd, cmdExe.Description())
	}
}

func (h *Help) Description() string {
	return "获取此客户端的使用方法"
}
