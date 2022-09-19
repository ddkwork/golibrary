package packet

import (
	"bytes"
	"io"
)

type (
	Interface interface {
		HeadIsNil() bool
		Bytes() []byte
		String() string
		//empty() bool
		Len() int
		//Cap() int
		//Truncate(n int)
		Reset()
		// tryGrowByReslice(n int) (int, bool)
		// grow(n int) int
		//Grow(n int)
		Write(p []byte) (n int, err error)
		WriteString(s string) (n int, err error)
		//ReadFrom(r io.Reader) (n int64, err error)
		// WriteTo(w io.Writer) (n int64, err error)
		// WriteByte(c byte) error
		// WriteRune(r rune) (n int, err error)
		// Read(p []byte) (n int, err error)
		//Next(n int) []byte
		//ReadByte() (byte, error)
		//ReadRune() (r rune, size int, err error)
		//UnreadRune() error
		//UnreadByte() error
		//ReadBytes(delim byte) (line []byte, err error)
		//readSlice(delim byte) (line []byte, err error)
		//ReadString(delim byte) (line string, err error)
	}
	object struct {
		b *bytes.Buffer
	}
)

var Default = New()

func (o *object) HeadIsNil() bool {
	readByte, err := o.b.ReadByte()
	if err != nil {
		return false
	}
	//if !mylog.Error(err) {
	//	return false
	//}
	return readByte == 0
}

func (o *object) Bytes() []byte {
	return o.b.Bytes()
}

func (o *object) String() string {
	return o.b.String()
}

func (o *object) empty() bool {
	panic("implement me")
}

func (o *object) Len() int {
	return o.b.Len()
}

func (o *object) Cap() int {
	panic("implement me")
}

func (o *object) Truncate(n int) {
	panic("implement me")
}

func (o *object) Reset() {
	o.b.Reset()
}

func (o *object) tryGrowByReslice(n int) (int, bool) {
	panic("implement me")
}

func (o *object) grow(n int) int {
	panic("implement me")
}

func (o *object) Grow(n int) {
	panic("implement me")
}

func (o *object) Write(p []byte) (n int, err error) {
	return o.b.Write(p)
}

func (o *object) WriteString(s string) (n int, err error) {
	return o.b.WriteString(s)
}

func (o *object) ReadFrom(r io.Reader) (n int64, err error) {
	panic("implement me")
}

func (o *object) WriteTo(w io.Writer) (n int64, err error) {
	panic("implement me")
}

func (o *object) WriteByte(c byte) error {
	panic("implement me")
}

func (o *object) WriteRune(r rune) (n int, err error) {
	panic("implement me")
}

func (o *object) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (o *object) Next(n int) []byte {
	panic("implement me")
}

func (o *object) ReadByte() (byte, error) {
	panic("implement me")
}

func (o *object) ReadRune() (r rune, size int, err error) {
	panic("implement me")
}

func (o *object) UnreadRune() error {
	panic("implement me")
}

func (o *object) UnreadByte() error {
	panic("implement me")
}

func (o *object) ReadBytes(delim byte) (line []byte, err error) {
	panic("implement me")
}

func (o *object) readSlice(delim byte) (line []byte, err error) {
	panic("implement me")
}

func (o *object) ReadString(delim byte) (line string, err error) {
	panic("implement me")
}

func New() Interface {
	return &object{
		b: bytes.NewBuffer(make([]byte, 512)),
	}
}
