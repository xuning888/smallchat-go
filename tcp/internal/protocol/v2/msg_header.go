package v2

import (
	"errors"
	"smallchat/tcp/internal/util"
)

const headerLength = 20

var (
	EmptyMsgError = errors.New("msg byte array is null")
	ErrMsgHeader  = errors.New("msg length is less than 20")
)

type MsgHeader struct {
	HeaderLength  int
	ClientVersion int
	CmdId         int
	Seq           int
	BodyLength    int
}

func CreateMsgHeader() *MsgHeader {
	return &MsgHeader{}
}

func CreateMsgHeaderFromBytes(bytes []byte) (*MsgHeader, error) {
	header := CreateMsgHeader()
	err := header.Decode(bytes)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (h *MsgHeader) Decode(msg []byte) error {
	if msg == nil || len(msg) <= 0 {
		return EmptyMsgError
	}

	if len(msg) < headerLength {
		return ErrMsgHeader
	}

	offset := 0
	hLength, _ := util.ByteArrayToInt(msg, offset)
	h.HeaderLength = hLength

	offset += 4
	hCv, _ := util.ByteArrayToInt(msg, offset)
	h.ClientVersion = hCv

	offset += 4
	cmdId, _ := util.ByteArrayToInt(msg, offset)
	h.CmdId = cmdId

	offset += 4
	seq, _ := util.ByteArrayToInt(msg, offset)
	h.Seq = seq

	offset += 4
	bodyLen, _ := util.ByteArrayToInt(msg, offset)
	h.BodyLength = bodyLen
	return nil
}

func (h *MsgHeader) ToBytes() []byte {
	return util.ConcatBytes(
		util.IntToByteArray(h.HeaderLength),
		util.IntToByteArray(h.ClientVersion),
		util.IntToByteArray(h.CmdId), util.IntToByteArray(h.Seq),
		util.IntToByteArray(h.BodyLength))
}
