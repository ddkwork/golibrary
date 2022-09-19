package stream

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ddkwork/golibrary/mylog"

	"os"
)

type (
	Type interface {
		~string | ~[]byte | ~*bytes.Buffer //todo and test type rename
	}
	_Interface interface { //todo  合并tool包
		NewLine()
		Quote() //手动Quote字符串避免造成换行失效
		QuoteWith(ss string)
		ObjectBegin()
		ObjectEnd()
		SliceBegin()
		SliceEnd()
		Indent(deep int) string
		WriteBytesLn(p []byte)
		WriteStringLn(ss string)
		HexString() string
		HexStringUpper() string
		Append(buffer ...*Stream)
		WriteXMakeBody(key string, values ...string)
		SizeCheck() bool
		ErrorInfo() string
		CutWithIndex(x, y int)
		BigNumXorWithAlign(arg1, arg2 []byte, align int) (xorStream []byte)
		Merge(Bytes ...[]byte) *Stream
		InsertString(size int, separate string) (ss string)
		SplitBytes(size int) (blocks [][]byte)
		SplitString(size int) (blocks []string)
		RemoveHexDumpNewLine(dump string) (newDump string)
		LinesToString(lines []string) string
	}
	Stream struct{ *bytes.Buffer }
)

func (s *Stream) LinesToString(lines []string) string {
	for _, line := range lines {
		s.WriteStringLn(line)
	}
	return s.String()
}

var Default = New()

func newInterface() _Interface {
	return &Stream{
		Buffer: &bytes.Buffer{},
	}
}

func New() *Stream {
	return &Stream{
		Buffer: &bytes.Buffer{},
	}
}
func (s *Stream) RemoveHexDumpNewLine(dump string) (newDump string) {
	//strings.TrimSuffix()
	panic("implement me")
}

func (s *Stream) SplitString(size int) (blocks []string) {
	blocks = make([]string, 0)
	splitBytes := s.SplitBytes(size)
	for _, splitByte := range splitBytes {
		blocks = append(blocks, string(splitByte))
	}
	return
}

func (s *Stream) CutWithIndex(x, y int) {
	//TODO implement me
	panic("implement me")
}

func NewReadFile(path string) *Stream {
	b, err := os.ReadFile(path)
	if !mylog.Error(err) {
		return nil
	}
	return &Stream{bytes.NewBuffer(b)}
}

func NewBytes(b []byte) *Stream           { return &Stream{bytes.NewBuffer(b)} }
func NewBuffer(buf *bytes.Buffer) *Stream { return &Stream{buf} }
func NewString(ss string) *Stream         { return &Stream{Buffer: bytes.NewBufferString(ss)} }
func NewHexString(ss string) (b *Stream) {
	b = New()
	decodeString, err := hex.DecodeString(ss)
	if !mylog.Error(err) {
		b.WriteString(err.Error())
		return
	}
	b.Write(decodeString)
	return
}
func NewHexStringOrBytes(data any) (b *Stream) {
	switch data.(type) {
	case string:
		return NewHexString(data.(string))
	case []byte:
		return NewBytes(data.([]byte))
	}
	return NewErrorInfo(fmt.Sprintf("%t\t", data))
}
func NewNil() *Stream                 { return New() }
func NewErrorInfo(err string) *Stream { return NewString(err) }
