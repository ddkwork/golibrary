package stream

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ddkwork/golibrary/mylog"

	"math/big"
	"strings"
)

func (s *Stream) NewLine()            { s.WriteString("\n") }
func (s *Stream) QuoteWith(ss string) { s.WriteString(ss) }
func (s *Stream) WriteBytesLn(p []byte) {
	s.Write(p)
	s.NewLine()
}
func (s *Stream) WriteStringLn(ss string) {
	s.WriteString(ss)
	s.NewLine()
}
func (s *Stream) Quote()                 { s.WriteByte('"') }
func (s *Stream) ObjectBegin()           { s.WriteByte('{') }
func (s *Stream) ObjectEnd()             { s.WriteByte('}') }
func (s *Stream) SliceBegin()            { s.WriteByte('[') }
func (s *Stream) SliceEnd()              { s.WriteByte(']') }
func (s *Stream) Indent(deep int) string { return strings.Repeat(" ", deep) }
func (s *Stream) HexString() string      { return hex.EncodeToString(s.Bytes()) }
func (s *Stream) HexStringUpper() string { return fmt.Sprintf("%#X", s.Bytes())[2:] }
func (s *Stream) SizeCheck() bool {
	switch s.Len() {
	case 0:
		return mylog.Error("buffer len == 0")
	default:
		if s.Len()%8 != 0 {
			return mylog.Error(" len%8 != 0")
		}
	}
	return true
}
func (s *Stream) ErrorInfo() string { return s.String() }
func (s *Stream) Append(buffer ...*Stream) {
	s.Reset()
	for _, b2 := range buffer {
		b2.WriteBytesLn(b2.Bytes())
	}
}
func (s *Stream) BigNumXorWithAlign(arg1, arg2 []byte, align int) (xorStream []byte) {
	xor := new(big.Int).Xor(new(big.Int).SetBytes(arg1), new(big.Int).SetBytes(arg2))
	alignBuf := make([]byte, align-len(xor.Bytes()))
	switch len(xor.Bytes()) {
	case 0:
		xorStream = alignBuf
	case align:
		xorStream = xor.Bytes()
	default:
		xorStream = s.Merge(alignBuf, xor.Bytes()).Bytes()
	}
	return
}

func (s *Stream) InsertString(size int, separate string) (ss string) {
	b := new(strings.Builder)
	for i, v := range s.String() {
		b.WriteRune(v)
		if (i+1)%size == 0 {
			b.WriteString(separate)
		}
	}
	ss = b.String()
	ss = ss[:b.Len()-1]
	return
}
func (s *Stream) SplitBytes(size int) (blocks [][]byte) {
	blocks = make([][]byte, 0)
	quantity := s.Len() / size
	remainder := s.Len() % size
	i := 0
	for i = 0; i < quantity; i++ {
		blocks = append(blocks, s.Bytes()[i*size:(i+1)*size])
	}
	if remainder != 0 {
		blocks = append(blocks, s.Bytes()[i*size:i*size+remainder])
	}
	return
}

func (s *Stream) Merge(bytesSlice ...[]byte) *Stream {
	b := bytes.NewBuffer(nil)
	b.Write(s.Bytes())
	for i := 0; i < len(bytesSlice); i++ {
		if !mylog.Error2(b.Write(bytesSlice[i])) {
			return nil
		}
	}
	return NewBuffer(b)
}
