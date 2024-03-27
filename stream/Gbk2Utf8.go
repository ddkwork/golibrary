package stream

import (
	"github.com/axgle/mahonia"
)

func Utf82Gbk(utf8 string) (gbk string) {
	enc := mahonia.NewEncoder("gbk")
	return enc.ConvertString(utf8)
}

func Gbk2Utf8(gbk string) (utf8 string) {
	decoder := mahonia.NewDecoder("gbk")
	return decoder.ConvertString(gbk)
}
