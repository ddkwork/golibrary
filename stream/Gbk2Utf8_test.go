package stream_test

import (
	"testing"

	"github.com/ddkwork/golibrary/stream"
)

func Test(t *testing.T) {
	return
	utf8 := stream.Gbk2Utf8("hello,世界")
	println(utf8)
	gbk := stream.Utf82Gbk(utf8)
	println(gbk)
}
