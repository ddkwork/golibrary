package stream

import (
	"bufio"
	"bytes"
	"crypto/des"
	"encoding/hex"
	"fmt"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/src/mybinary"
	"io"
	"math/big"
	"os"
	"strings"
)

type (
	Type interface {
		~string | ~[]byte | ~*bytes.Buffer //todo and test type rename
	}
	Stream struct{ *bytes.Buffer }
)

var Default = New()

func New() *Stream                        { return &Stream{Buffer: &bytes.Buffer{}} }
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
	return NewString(fmt.Sprintf("%t\t", data))
}

func NewReadFile(path string) *Stream {
	b, err := os.ReadFile(path)
	if !mylog.Error(err) {
		return nil
	}
	return &Stream{bytes.NewBuffer(b)}
}

func (s *Stream) LinesToString(lines []string) string {
	for _, line := range lines {
		s.WriteStringLn(line)
	}
	return s.String()
}

func NewHexDump(hexdump string) (buf []byte) {
	const (
		address    = "00000000  "
		sep        = "|"
		newLine    = "\n"
		addressLen = len(address)
	)
	defer func() {
		s := New()
		//s.WriteStringLn("buf:=" + fmt.Sprintf("%#v", buf))
		cut := `[]byte`
		cxx := fmt.Sprintf("%#v", buf)
		cxx = cxx[len(cut):]
		s.WriteStringLn("char buf[] = " + cxx + ";")
		mylog.Json("gen c++ code", s.String())
		mylog.HexDump("recovery go buffer", buf)
	}()
	hexdump = strings.TrimSuffix(hexdump, newLine)
	if !strings.Contains(hexdump, address) && !strings.Contains(hexdump, sep) {
		hexdump = strings.ReplaceAll(hexdump, " ", "")
		decodeString, err := hex.DecodeString(hexdump)
		if !mylog.Error(err) {
			return
		}
		buf = decodeString
		return
	}
	split := strings.Split(hexdump, newLine)
	noAddres := make([]string, len(split))

	hexString := new(bytes.Buffer)
	for i, s := range split {
		if s == "" {
			continue
		}
		noAddres[i] = s[addressLen:strings.Index(s, sep)]
		noAddres[i] = strings.ReplaceAll(noAddres[i], " ", "")
		hexString.WriteString(noAddres[i])
	}
	decodeString, err := hex.DecodeString(hexString.String())
	if !mylog.Error(err) {
		return
	}
	buf = decodeString
	return
}

func (s *Stream) buffer(data any) *bytes.Buffer { //todo replaced as stream pkg
	switch data.(type) {
	case string:
		return bytes.NewBufferString(data.(string))
	case []byte:
		return bytes.NewBuffer(data.([]byte))
	}
	return bytes.NewBufferString("error file data type " + fmt.Sprintf("%t", data))
}

func (s *Stream) ToLines(data any) (lines []string, ok bool) {
	newReader := bufio.NewReader(s.buffer(data))
	for {
		line, _, err := newReader.ReadLine()
		switch err {
		case io.EOF:
			return lines, true
		default:
			if !mylog.Error(err) {
				return
			}
		}
		lines = append(lines, string(line))
	}
}

func (s *Stream) ReadToLines(path string) (lines []string, ok bool) {
	file, err := os.ReadFile(path)
	if !mylog.Error(err) {
		return
	}
	return s.ToLines(file)
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

func (s *Stream) CutWithIndex(x, y int) *Stream { return NewBytes(s.Bytes()[x:y]) }
func (s *Stream) NewLine() *Stream {
	s.WriteString("\n")
	return s
}
func (s *Stream) QuoteWith(ss string) *Stream {
	s.WriteString(ss)
	return s
}
func (s *Stream) WriteAny(order mybinary.ByteOrder, data any) *Stream {
	if !mylog.Error(mybinary.Write(s, order, data)) {
		return nil
	}
	return s
}
func (s *Stream) ReadAny(order mybinary.ByteOrder) (data any) {
	if !mylog.Error(mybinary.Read(s, order, data)) {
		return nil
	}
	return s
}
func (s *Stream) WriteBytesLn(p []byte) *Stream {
	s.Write(p)
	s.NewLine()
	return s
}
func (s *Stream) WriteStringLn(ss string) *Stream {
	s.WriteString(ss)
	s.NewLine()
	return s
}
func (s *Stream) Quote() *Stream {
	s.WriteByte('"')
	return s
}
func (s *Stream) ObjectBegin() *Stream {
	s.WriteByte('{')
	return s
}
func (s *Stream) ObjectEnd() *Stream {
	s.WriteByte('}')
	return s
}
func (s *Stream) SliceBegin() *Stream {
	s.WriteByte('[')
	return s
}
func (s *Stream) SliceEnd() *Stream {
	s.WriteByte(']')
	return s
}
func (s *Stream) Indent(deep int) *Stream {
	s.WriteString(strings.Repeat(" ", deep))
	return s
}
func (s *Stream) HexString() string      { return hex.EncodeToString(s.Bytes()) }
func (s *Stream) HexStringUpper() string { return fmt.Sprintf("%#X", s.Bytes())[2:] }
func (s *Stream) DesBlockSizeSizeCheck() bool {
	switch s.Len() {
	case 0:
		return mylog.Error("buffer len == 0")
	default:
		if s.Len()%des.BlockSize != 0 {
			return mylog.Error("Len()%des.BlockSize != 0")
		}
	}
	return true
}
func (s *Stream) Merge(streams ...*Stream) *Stream {
	for _, b2 := range streams {
		s.Write(b2.Bytes())
	}
	return s
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
		xorStream = s.AppendByteSlice(alignBuf, xor.Bytes()).Bytes()
	}
	return
}

func (s *Stream) InsertString(index int, insert string) *Stream {
	start := s.String()[:index]
	end := s.String()[index:]
	s.Reset()
	s.WriteString(start)
	s.WriteString(insert)
	s.WriteString(end)
	return s
}
func (s *Stream) InsertBytes(index int, insert []byte) *Stream {
	start := s.Bytes()[:index]
	end := s.Bytes()[index:]
	s.Reset()
	s.Write(start)
	s.Write(insert)
	s.Write(end)
	return s
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

func (s *Stream) AppendByteSlice(bytesSlice ...[]byte) *Stream {
	for _, slice := range bytesSlice {
		s.Write(slice)
	}
	return s
}
