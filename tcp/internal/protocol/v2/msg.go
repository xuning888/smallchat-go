package v2

import "errors"

var (
	ErrorMsgLen = errors.New("headerLen + bodyLen != msgLen")
)

type Msg struct {
	Header *MsgHeader
	Body   []byte
}

func CreateMsgFromHeaderAndBoy(header *MsgHeader, body []byte) *Msg {
	return &Msg{
		Header: header,
		Body:   body,
	}
}

func CreateMsgFromBytes(bytes []byte) (*Msg, error) {
	if bytes == nil || len(bytes) == 0 {
		return nil, EmptyMsgError
	}
	msgLen := len(bytes)
	msgHeader, err := CreateMsgHeaderFromBytes(bytes)
	if err != nil {
		return nil, err
	}
	if msgHeader.HeaderLength+msgHeader.BodyLength != msgLen {
		return nil, ErrorMsgLen
	}
	body := make([]byte, msgHeader.BodyLength)
	copy(body, bytes[20:20+msgHeader.BodyLength])
	return &Msg{
		Header: msgHeader,
		Body:   body,
	}, nil
}

func (m *Msg) ToBytes() []byte {
	bytes := make([]byte, 0, 20+len(m.Body))
	headerBytes := m.Header.ToBytes()
	bytes = append(bytes, headerBytes...)
	bytes = append(bytes, m.Body...)
	return bytes
}
