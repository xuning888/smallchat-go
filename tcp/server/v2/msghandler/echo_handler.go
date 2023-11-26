package msghandler

import (
	"smallchat/tcp/internal/cmd"
	"smallchat/tcp/internal/protocol/v2"
)

var _ MsgHandler = &EchoMsgHandler{}

type EchoMsgHandler struct {
}

func (e *EchoMsgHandler) CmdId() int {
	return cmd.Echo
}

func (e *EchoMsgHandler) Process(msg *v2.Msg) error {
	return nil
}

func init() {
	handler := &EchoMsgHandler{}
	HandlerHolder.RegisterHandler(cmd.Echo, handler)
}
