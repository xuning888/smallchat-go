package handler

import (
	"log"
	"smallchat/tcp/internal/cmd"
)

type EchoHandler struct {
}

func (e *EchoHandler) HandleMessage(event cmd.Event) error {
	log.Printf("echo dispatch\n")
	bytes := event.GetMsg().ToBytes()
	log.Printf("echo msg size: %d\n", len(bytes))
	_, err := event.GetConn().Write(bytes)
	return err
}
