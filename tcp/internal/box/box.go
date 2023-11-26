package box

import (
	"container/heap"
	"math"
	"smallchat/tcp/internal/protocol/v1"
)

type MessageHeap []*v1.Message

func (h MessageHeap) Len() int {
	return len(h)
}

func (h MessageHeap) Less(i, j int) bool {
	return h[i].SendTime.Before(h[j].SendTime)
}

func (h MessageHeap) Swap(i, j int) {
	if i == j {
		return
	}
	h[i], h[j] = h[j], h[i]
}

func (h *MessageHeap) Push(x any) {
	message, ok := x.(*v1.Message)
	if ok {
		*h = append(*h, message)
	}
}

func (h *MessageHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = nil
	if n-1 == 0 {
		*h = make(MessageHeap, 0)
	} else {
		*h = old[0 : n-1]
	}
	return x
}

type MessageBox struct {
	Messages *MessageHeap
	limit    int
}

func CreateMsgBox(limit int) *MessageBox {
	limit = int(math.Max(float64(limit), 10))
	messages := make(MessageHeap, 0)
	heap.Init(&messages)
	return &MessageBox{
		Messages: &messages,
		limit:    limit,
	}
}

func (mbx *MessageBox) PushMsg(msg string) {
	curLen := mbx.Messages.Len()
	if curLen == mbx.limit {
		return
	}
	heap.Push(mbx.Messages, v1.CreateMsg(msg))
}

func (mbx *MessageBox) PushMsgWithTime(message *v1.Message) {
	messages := mbx.Messages
	heap.Push(messages, message)
}

func (mbx *MessageBox) PopMsg() *v1.Message {
	messages := mbx.Messages
	msgLen := messages.Len()
	if msgLen == 0 {
		return nil
	}
	pop := heap.Pop(messages)
	curMsg, ok := pop.(*v1.Message)
	if ok {
		return curMsg
	}
	return nil
}

func (mbx *MessageBox) Len() int {
	return mbx.Messages.Len()
}
