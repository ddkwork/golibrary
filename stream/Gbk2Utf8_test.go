package stream

import "testing"

func Test(t *testing.T) {
	//Gbk2Utf8Vs2022("D:\\clone\\EWDK_quickstart-master\\New folder\\VT1-master\\VT")
	return
	utf8 := Gbk2Utf8("hello,世界")
	println(utf8)
	gbk := Utf82Gbk(utf8)
	println(gbk)
}
