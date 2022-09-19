package unicode

import (
	"encoding/binary"
	"fmt"
	"github.com/ddkwork/golibrary/mylog"

	"unicode/utf16"
)

type (
	Interface interface {
		FromString(s string) (ok bool)
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

func (o *object) SplitArrayByStep(buf []byte, splitSize int64) (blocks [][]byte) { //循环引用，操
	size := int64(len(buf))
	blocks = make([][]byte, 0)
	quantity := size / splitSize
	remainder := size % splitSize
	i := int64(0)
	for i = int64(0); i < quantity; i++ {
		blocks = append(blocks, buf[i*splitSize:(i+1)*splitSize])
	}
	if remainder != 0 {
		blocks = append(blocks, buf[i*splitSize:i*splitSize+remainder])
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

func (o *object) FromString(s string) (ok bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == 0 {
			return mylog.Error("unicode字符串不能包含0值字节，0只能是固定的拓展位，就是说返回的unicode不应该有0x0000这种宽字符集")
		}
	}
	o.utf16 = utf16.Encode([]rune(s + "\x00")) //bug,这里会多了一组尾巴的0X0000，其实还可以了直接用类型拓展用，1字节转2字节即可
	o.utf16 = o.utf16[:len(o.utf16)-1]         //cut 0X0000
	for _, u := range o.utf16 {
		//fmt.Printf("%#04X\n", u)
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, u)
		o.unicode = append(o.unicode, b...)
	}
	return true
}

func (o *object) Bytes() []byte { return o.unicode }

var Default = New()

func New() Interface {
	return &object{
		s:       "",
		utf16:   nil,
		unicode: make([]byte, 0),
	}
}
