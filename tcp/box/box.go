package box

import (
	"math"
	"smallchat/tcp/protocol"
	"sort"
)

type MessageBox struct {
	Messages []*protocol.Message
	limit    int
}

func CreateMsgBox(limit int) *MessageBox {
	limit = int(math.Max(float64(limit), 10))
	return &MessageBox{
		Messages: make([]*protocol.Message, 0),
		limit:    limit,
	}
}

func (mbx *MessageBox) PushMsg(msg string) {
	curLen := len(mbx.Messages)
	if curLen == mbx.limit {
		return
	}
	curMsg := protocol.CreateMsg(msg)
	mbx.Messages = append(mbx.Messages, curMsg)
}

func (mbx *MessageBox) PullMsg() *protocol.Message {
	messages := mbx.Messages
	msgLen := len(messages)
	if msgLen == 0 {
		return nil
	}
	// 按照时间排序
	sort.Slice(mbx.Messages, func(i, j int) bool {
		return mbx.Messages[i].SendTime.Before(mbx.Messages[j].SendTime)
	})
	curMg := mbx.Messages[0]
	mbx.Messages = mbx.Messages[1:]
	return curMg
}
