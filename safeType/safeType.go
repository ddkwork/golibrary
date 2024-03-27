package safeType

import (
	"bytes"
)

type (
	BinaryType interface {
		HexString | HexDumpString | ~[]byte | ~*bytes.Buffer
	}
	Type interface { // todo 加入大数 for rsa，加入全部普通类型 for 二进制协议编解码约束，pb之类的
		string | BinaryType
	}
)

type Data struct{ *bytes.Buffer }

func New[T Type](s T) *Data {
	switch s := any(s).(type) {
	case []byte:
		return &Data{bytes.NewBuffer(s)}
	case *bytes.Buffer:
		return &Data{s}
	case string:
		return &Data{Buffer: bytes.NewBufferString(s)}
	case HexString:
		return NewHexString(s)
	case HexDumpString:
		return NewHexDump(s)
	default:
		return &Data{Buffer: &bytes.Buffer{}}
	}
}

// NewBinaryType for asm
func NewBinaryType[T BinaryType](s T) *Data {
	switch s := any(s).(type) {
	case []byte:
		return &Data{bytes.NewBuffer(s)}
	case *bytes.Buffer:
		return &Data{s}
	case HexString:
		return NewHexString(s)
	case HexDumpString:
		return NewHexDump(s)
	default:
		return &Data{Buffer: &bytes.Buffer{}}
	}
}
