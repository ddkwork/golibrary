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

func TestHexDumpToGoBytes(t *testing.T) {
	ss := `00 00 00 1A 00 00 00 09 00 01 00 00 0B 00 00 00 8E 6A 64 01 15 4F 53 44 4B 5F 41 42 55 53 45 5F 52 45 50 4F 52 54 49 4E 47 00`
	NewHexDump(ss)
	//NewHexDump(dump)
	NewHexDump(`00 00 00 1A 00 00 00 09 00 01 00 00 0B 00 00 00`)             //16 byte header
	NewHexDump(`8E 6A 64 01`)                                                 //tag
	NewHexDump(`01`)                                                          //
	NewHexDump(`15`)                                                          //strinf type id
	NewHexDump(`4F 53 44 4B 5F 41 42 55 53 45 5F 52 45 50 4F 52 54 49 4E 47`) //OSDK_ABUSE_REPORTING
	NewHexDump(`00`)                                                          //string end
}

var dump = `00000000  7e 15 00 80 0b 00 00 00  09 25 ce f7 3d 01 00 10  |~........%..=...|
00000010  01 10 00 08 ac 80 04 1a  b2 01 32 00 00 00 04 00  |..........2.....|
00000020  00 00 25 ce f7 3d 01 00  10 01 07 00 00 00 83 7c  |..%..=.........||
00000030  39 6a 97 2b a8 c0 00 00  00 00 00 a9 25 63 80 58  |9j.+........%c.X|
00000040  41 63 01 00 00 00 00 00  00 00 00 00 ad 93 9b 27  |Ac.............'|
00000050  eb b6 3d dc 1d 57 fe b2  d1 86 79 de a1 41 61 eb  |..=..W....y..Aa.|
00000060  04 70 81 ce 35 f5 28 6a  05 52 d9 7b 7d 6c f9 2e  |.p..5.(j.R.{}l..|
00000070  5c b9 5e 8a b6 a5 87 dc  da 25 03 0b 00 48 76 7b  |\.^......%...Hv{|
00000080  66 ba f9 0b 48 78 62 09  bf 88 be 49 de 09 36 52  |f...Hxb....I..6R|
00000090  57 42 8d 69 34 8b 80 ac  e9 0b 8f ef e1 dd a2 0b  |WB.i4...........|
000000a0  25 0c cf 26 f9 0f dc 30  df 21 46 8f b6 8d c2 56  |%..&...0.!F....V|
000000b0  78 88 ef 2a 97 8c 50 c7  e2 9b 42 6f 53 09 82 42  |x..*..P...BoS..B|
000000c0  cc d4 3e 57 b5 ef b4 23  2c 54 13 97 20 d1 cf f0  |..>W...#,T.. ...|
000000d0  a7 b2 98 85 d3 54                                 |.....T|`
