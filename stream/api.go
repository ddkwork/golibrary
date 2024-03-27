package stream

import (
	"os"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safeType"
)

type Stream struct{ *safeType.Data }

func New[T safeType.Type](data T) *Stream         { return &Stream{Data: safeType.New(data)} }
func NewHexString(s safeType.HexString) *Stream   { return &Stream{Data: safeType.NewHexString(s)} }
func NewHexDump(h safeType.HexDumpString) *Stream { return &Stream{Data: safeType.NewHexDump(h)} }
func NewBinaryType[T safeType.BinaryType](b T) *Stream {
	return &Stream{Data: safeType.NewBinaryType(b)}
}

func NewReadFile(path string) *Stream {
	b, err := os.ReadFile(path)
	if !mylog.Error(err) {
		return &Stream{safeType.New(err.Error())}
	}
	return &Stream{safeType.New(b)}
}
func ReadFileToLines(path string) (lines []string, ok bool) { return NewReadFile(path).ToLines() }

func (o *Stream) LinesToString(lines []string) string {
	for _, line := range lines {
		o.WriteStringLn(line)
	}
	return o.String()
}
