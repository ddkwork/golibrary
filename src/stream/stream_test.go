package stream

import (
	"github.com/ddkwork/golibrary/mylog"
	"testing"
)

func TestNewHexString(t *testing.T) {
	b := NewHexString("1122")
	mylog.HexDump("", b.Bytes())
}

func TestStream_AppendByteSlice(t *testing.T) {
	mylog.HexDump("", New().AppendByteSlice([]byte{0x11}, []byte{0x22}).Bytes())
}

func TestStream_BigNumXorWithAlign(t *testing.T) {
	//New().BigNumXorWithAlign()
}

func TestStream_CutWithIndex(t *testing.T) {
	mylog.HexDump("", NewBytes([]byte{0x11, 0x22, 0x33, 0x44, 0x55}).CutWithIndex(2, 4).Bytes())
}

func TestStream_DesBlockSizeSizeCheck(t *testing.T) {
	NewBytes([]byte{0x11, 0x22, 0x33, 0x44, 0x55}).DesBlockSizeSizeCheck()
}

func TestStream_HexString(t *testing.T) {
	println(NewBytes([]byte{0x11, 0x22, 0x33, 0x44, 0x55}).HexString())
}

func TestStream_HexStringUpper(t *testing.T) {
	println(NewBytes([]byte{0xaa, 0xbb, 0x33, 0x44, 0x55}).HexStringUpper())

}

func TestStream_Indent(t *testing.T) {
	indent := NewString("1111").Indent(3)
	indent.WriteString("3344")
	println(indent.String())
}

func TestStream_InsertString(t *testing.T) {
	println(NewString("aaaaa").InsertString(1, "22").String())
}

func TestStream_LinesToString(t *testing.T) {
	//New().LinesToString(tt.args.lines)
}

func TestStream_Merge(t *testing.T) {
	mylog.HexDump("", NewBytes([]byte{0xaa, 0xbb, 0x33, 0x44, 0x55}).Merge(NewBytes([]byte{0xaa, 0xbb, 0x33, 0x44, 0x55})).Bytes())
}

func TestStream_NewLine(t *testing.T) {
	println(NewString("111").NewLine().String())
}

func TestStream_ObjectBegin(t *testing.T) {
	println(NewString("111").NewLine().ObjectBegin().String())

}

func TestStream_ObjectEnd(t *testing.T) {
	println(NewString("111").NewLine().ObjectEnd().String())
}

func TestStream_Quote(t *testing.T) {
	println(NewString("111").NewLine().Quote().String())
}

func TestStream_QuoteWith(t *testing.T) {
	println(NewString("111").NewLine().QuoteWith("//").String())
}

func TestStream_RemoveHexDumpNewLine(t *testing.T) {

}

func TestStream_SliceBegin(t *testing.T) {
	println(NewString("111").NewLine().SliceBegin().String())
}

func TestStream_SliceEnd(t *testing.T) {
	println(NewString("111").NewLine().SliceEnd().String())
}

func TestStream_SplitBytes(t *testing.T) {
	//NewBytes().SplitBytes(tt.args.size)
}

func TestStream_SplitString(t *testing.T) {
	//NewBytes().SplitString(tt.args.size)
}

func TestStream_WriteBytesLn(t *testing.T) {

}

func TestStream_WriteStringLn(t *testing.T) {

}

func TestStream_InsertString1(t *testing.T) {
	s := NewString("B3EBDDAF2EA789C4")
	code := s.InsertString(4, "-")
	println(code.String())
}
