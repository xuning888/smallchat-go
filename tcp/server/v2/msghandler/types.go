package msghandler

import (
	"log"
	"smallchat/tcp/internal/protocol/v2"
	"sync"
)

var HandlerHolder *Holder

func init() {
	HandlerHolder = &Holder{
		msgHandlers: sync.Map{},
	}
}

type MsgHandler interface {
	CmdId() int
	Process(msg *v2.Msg) error
}

type Holder struct {
	msgHandlers sync.Map
}

func (m *Holder) RegisterHandler(cmdId int, handler MsgHandler) {
	log.Printf("register msgHandler, cmdId:%d, handler:%v", cmdId, handler)
	m.msgHandlers.Store(cmdId, handler)
}

func (m *Holder) GetMsgHandler(cmdId int) (MsgHandler, bool) {
	value, ok := m.msgHandlers.Load(cmdId)
	if ok {
		return value.(MsgHandler), ok
	}
	return nil, false
}
