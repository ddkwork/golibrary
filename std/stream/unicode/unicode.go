package unicode

import (
	"encoding/binary"
	"fmt"
	"unicode/utf16"

	"github.com/ddkwork/golibrary/std/mylog"
)

type (
	Interface interface {
		FromString(s string)
		ToString(unicode []byte) string
		ToUtf16(unicode []byte) []uint16
		fmt.Stringer
		Bytes() []byte
		Utf16() []uint16
	}
	object struct {
		s       string
		utf16   []uint16
		unicode []byte
	}
)

func (o *object) SplitArrayByStep(buf []byte, splitSize int64) (blocks [][]byte) {
	size := int64(len(buf))
	blocks = make([][]byte, 0)
	quantity := size / splitSize
	remainder := size % splitSize
	for i := range quantity {
		blocks = append(blocks, buf[i*splitSize:(i+1)*splitSize])
		if remainder != 0 {
			blocks = append(blocks, buf[i*splitSize:i*splitSize+remainder])
		}
	}
	return
}

func (o *object) ToUtf16(unicode []byte) []uint16 {
	o.utf16 = make([]uint16, 0)
	blocks := o.SplitArrayByStep(unicode, 2)
	for _, block := range blocks {
		o.utf16 = append(o.utf16, binary.LittleEndian.Uint16(block))
	}
	return o.utf16
}

func (o *object) ToString(unicode []byte) string {
	o.s = ""
	o.utf16 = o.ToUtf16(unicode)
	o.s = string(utf16.Decode(o.utf16))
	return o.String()
}

func (o *object) String() string { return o.s }
func (o *object) Utf16() []uint16 {
	return o.utf16
}

func (o *object) FromString(s string) {
	for i := range len(s) {
		mylog.Check(s[i] != 0)
	}
	o.utf16 = utf16.Encode([]rune(s + "\x00"))
	o.utf16 = o.utf16[:len(o.utf16)-1]
	for _, u := range o.utf16 {

		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, u)
		o.unicode = append(o.unicode, b...)
	}
}

func (o *object) Bytes() []byte { return o.unicode }

func New() Interface {
	return &object{
		s:       "",
		utf16:   nil,
		unicode: make([]byte, 0),
	}
}
