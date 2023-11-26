package util

import "errors"

var (
	ErrBytesIsNull = errors.New("bytes array is null")
)

func ByteArrayToInt(bytes []byte, offset int) (int, error) {
	if bytes == nil || len(bytes) <= 0 {
		return 0, ErrBytesIsNull
	}
	return fromBytes(bytes[offset], bytes[offset+1], bytes[offset+2], bytes[offset+3]), nil
}

func fromBytes(b1, b2, b3, b4 byte) int {
	b1Convert := int(b1) << 24
	b2Convert := (int(b2) & 0xFF) << 16
	b3Convert := (int(b3) & 0xFF) << 8
	b4Convert := int(b4) & 0xFF
	return b1Convert | b2Convert | b3Convert | b4Convert
}

func IntToByteArray(value int) []byte {
	return []byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	}
}

func ConcatBytes(bytess ...[]byte) []byte {
	result := make([]byte, 0)
	for _, bytes := range bytess {
		result = append(result, bytes...)
	}
	return result
}
