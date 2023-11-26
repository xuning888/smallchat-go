package cmd

import (
	"errors"
	"log"
	"net"
	"smallchat/tcp/internal/cmd/handler"
	"smallchat/tcp/internal/protocol/v2"
	"sync"
)

var (
	ErrNotFoundHandler                = errors.New("not found handler")
	Holder             *HolderHandler = nil
)

type Event interface {
	GetMsg() *v2.Msg
	GetStartTime() int64
	GetConn() net.Conn
	Up() bool
}

type EventHandler interface {
	HandleMessage(event Event) error
}

func init() {
	Holder = &HolderHandler{
		handlers: sync.Map{},
	}
	Holder.register(Echo, &handler.EchoHandler{})
}

type HolderHandler struct {
	handlers sync.Map
}

func (h *HolderHandler) HandleMessage(event Event) error {
	cmdId := event.GetMsg().Header.CmdId
	value, exist := h.handlers.Load(cmdId)
	if exist {
		eventHandler, ok := value.(EventHandler)
		if ok {
			return eventHandler.HandleMessage(event)
		}
		return ErrNotFoundHandler
	}
	return ErrNotFoundHandler
}

func (h *HolderHandler) register(cmdId int, handler EventHandler) {
	log.Printf("register handler. cmdId:%d, handler:%v", cmdId, handler)
	h.handlers.Store(cmdId, handler)
}
