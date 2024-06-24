package stream

import (
	"github.com/axgle/mahonia"
)

func Utf82Gbk(utf8 string) (gbk string) {
	return mahonia.NewEncoder("gbk").ConvertString(utf8)
}

func Gbk2Utf8(gbk string) (utf8 string) {
	return mahonia.NewDecoder("gbk").ConvertString(gbk)
}
