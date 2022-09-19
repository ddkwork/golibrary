package swap

import (
	"bytes"
	"encoding/hex"
	"github.com/ddkwork/golibrary/mylog"
)

type (
	//1 每两个字符串交换顺序
	//2 交换端序Endian
	Interface interface {
		SerialNumber(s string) string   //每两个字符串交换顺序用于硬盘序列号
		HexString(in string) []byte     //解码16进制字符串为buf并交换端序
		String(src []byte) (dst string) //交换字符串端序
		Bytes(src []byte) (dst []byte)  //交换字节切片端序
		CutUint16(v uint16) uint8       //从int16截取int8并交换端序，用于数据恢复软件算法
	}
	object struct{}
)

func (o *object) SerialNumber(s string) string {
	src := bytes.NewBufferString(s)
	to := bytes.Buffer{}
	for k, v := range src.Bytes() {
		if k%2 == 1 {
			to.WriteByte(v)
			to.WriteByte(src.Bytes()[k-1])
		}
	}
	return to.String()
}

func New() Interface { return &object{} }

func (o *object) HexString(hexStr string) []byte {
	decodeString, err := hex.DecodeString(hexStr)
	if !mylog.Error(err) {
		return nil
	}
	return o.Bytes(decodeString)
}

func (o *object) String(src []byte) (dst string) { return string(o.Bytes(src)) }
func (o *object) Bytes(src []byte) (dst []byte) {
	//mybinary.BigEndian.PutUint64()三次才行，10字节的话，如果是更多字节那就不通用
	to := bytes.Buffer{}
	for i := range src {
		to.WriteByte(src[len(src)-i-1])
	}
	return to.Bytes()
}

func (o *object) CutUint16(v uint16) uint8 { //6613-->16
	tmp := uint16(int32(v) << uint64(int32(4)) >> uint64(int32(8)))
	a := uint8(int32(tmp) << uint64(int32(12)) >> uint64(int32(12)))
	b := uint8(int32(a) >> uint64(int32(4)))
	return uint8(int32(a)<<uint64(int32(4)) | int32(b))
}
