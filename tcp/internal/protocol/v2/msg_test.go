package v2

import (
	"github.com/stretchr/testify/assert"
	"smallchat/tcp/internal/util"
	"testing"
)

func Test_Decode(t *testing.T) {
	testCases := []struct {
		name          string
		headerLength  []byte
		CmdId         []byte
		Seq           []byte
		ClientVersion []byte
		body          []byte
		bodyLen       []byte
		msg           string
		WantError     error
		WantMsgHeader *MsgHeader
	}{
		{
			name:          "Decode消息成功",
			msg:           "你好啊啊啊啊啊啊啊啊啊啊啊",
			headerLength:  util.IntToByteArray(20),
			ClientVersion: util.IntToByteArray(1),
			CmdId:         util.IntToByteArray(1001),
			Seq:           util.IntToByteArray(0),
			body:          []byte("你好啊啊啊啊啊啊啊啊啊啊啊"),
			bodyLen:       util.IntToByteArray(len([]byte("你好啊啊啊啊啊啊啊啊啊啊啊"))),
			WantError:     nil,
			WantMsgHeader: &MsgHeader{
				HeaderLength:  20,
				ClientVersion: 1,
				CmdId:         1001,
				Seq:           0,
				BodyLength:    len([]byte("你好啊啊啊啊啊啊啊啊啊啊啊")),
			},
		},
		{
			name:          "空消息",
			msg:           "",
			headerLength:  make([]byte, 0),
			ClientVersion: make([]byte, 0),
			CmdId:         make([]byte, 0),
			Seq:           make([]byte, 0),
			body:          make([]byte, 0),
			bodyLen:       make([]byte, 0),
			WantMsgHeader: nil,
			WantError:     EmptyMsgError,
		},
		{
			name:          "消息长度小于20个字节",
			msg:           "",
			headerLength:  util.IntToByteArray(20),
			ClientVersion: make([]byte, 0),
			CmdId:         make([]byte, 0),
			Seq:           make([]byte, 0),
			body:          make([]byte, 0),
			bodyLen:       make([]byte, 0),
			WantMsgHeader: nil,
			WantError:     ErrMsgHeader,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msgSequence := util.ConcatBytes(tc.headerLength, tc.ClientVersion, tc.CmdId, tc.Seq,
				tc.bodyLen, tc.body)
			msgHeader, err := CreateMsgHeaderFromBytes(msgSequence)
			if err != nil {
				assert.Nil(t, msgHeader)
				assert.Equal(t, err, tc.WantError)
			} else {
				assert.Equal(t, tc.WantMsgHeader, msgHeader)
			}
		})
	}
}

func TestCreateMsgFromBytes(t *testing.T) {
	testCases := []struct {
		name          string
		headerLength  []byte
		CmdId         []byte
		Seq           []byte
		ClientVersion []byte
		body          []byte
		bodyLen       []byte
		msg           string
		WantError     error
		WantMsg       *Msg
	}{
		{
			name:          "Decode消息成功",
			msg:           "你好啊啊啊啊啊啊啊啊啊啊啊",
			headerLength:  util.IntToByteArray(20),
			ClientVersion: util.IntToByteArray(1),
			CmdId:         util.IntToByteArray(1001),
			Seq:           util.IntToByteArray(0),
			body:          []byte("你好啊啊啊啊啊啊啊啊啊啊啊"),
			bodyLen:       util.IntToByteArray(len([]byte("你好啊啊啊啊啊啊啊啊啊啊啊"))),
			WantError:     nil,
			WantMsg: &Msg{
				Header: &MsgHeader{
					HeaderLength:  20,
					ClientVersion: 1,
					CmdId:         1001,
					Seq:           0,
					BodyLength:    len([]byte("你好啊啊啊啊啊啊啊啊啊啊啊")),
				},
				Body: []byte("你好啊啊啊啊啊啊啊啊啊啊啊"),
			},
		},
		{
			name:          "空消息",
			msg:           "",
			headerLength:  make([]byte, 0),
			ClientVersion: make([]byte, 0),
			CmdId:         make([]byte, 0),
			Seq:           make([]byte, 0),
			body:          make([]byte, 0),
			bodyLen:       make([]byte, 0),
			WantMsg:       nil,
			WantError:     EmptyMsgError,
		},
		{
			name:          "消息长度小于20个字节",
			msg:           "",
			headerLength:  util.IntToByteArray(20),
			ClientVersion: make([]byte, 0),
			CmdId:         make([]byte, 0),
			Seq:           make([]byte, 0),
			body:          make([]byte, 0),
			bodyLen:       make([]byte, 0),
			WantMsg:       nil,
			WantError:     ErrMsgHeader,
		},
		{
			name:          "消息头长度 + 消息体长度不等于消息长度",
			msg:           "你好1",
			headerLength:  util.IntToByteArray(20),
			ClientVersion: util.IntToByteArray(1),
			CmdId:         util.IntToByteArray(1001),
			Seq:           util.IntToByteArray(0),
			body:          []byte("你好"),
			bodyLen:       util.IntToByteArray(len([]byte("你好1"))),
			WantError:     ErrorMsgLen,
			WantMsg:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msgSequence := util.ConcatBytes(
				tc.headerLength,
				tc.ClientVersion,
				tc.CmdId,
				tc.Seq,
				tc.bodyLen,
				tc.body,
			)
			msg, err := CreateMsgFromBytes(msgSequence)
			if err != nil {
				assert.Nil(t, msg)
				assert.Equal(t, err, tc.WantError)
			} else {
				assert.Equal(t, tc.WantMsg, msg)
			}
		})
	}
}
