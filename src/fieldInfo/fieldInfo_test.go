package fieldInfo_test

import (
	"github.com/ddkwork/golibrary/src/fieldInfo"
	"runtime/debug"
	"testing"
)

func TestFiled(t *testing.T) {
	f := fieldInfo.New()
	f.Append(fieldInfo.Row{
		Number:   "1",
		Deep:     "",
		Name:     "info",
		Kind:     "int",
		Size:     "8",
		ValueFmt: "0x1122334455667788|1234605616436508552",
		Value:    nil,
		Stream:   "",
		Json:     "",
	})
	f.Append(fieldInfo.Row{
		Number:   "2",
		Deep:     "",
		Name:     "xx",
		Kind:     "int32",
		Size:     "4",
		ValueFmt: "0x1122334455667788|1234605616436508552",
		Value:    nil,
		Stream:   "",
		Json:     "",
	})
	f.Append(fieldInfo.Row{
		Number:   "3",
		Deep:     "",
		Name:     "",
		Kind:     "",
		Size:     "",
		ValueFmt: "0x1122334455667788|1234605616436508552",
		Value:    nil,
		Stream:   "",
		Json:     "",
	})
	f.Append(fieldInfo.Row{
		Number:   "4",
		Deep:     "",
		Name:     "",
		Kind:     "",
		Size:     "",
		ValueFmt: "0x1122334455667788|1234605616436508552",
		Value:    nil,
		Stream: `
00000000  01 02 03 04 05 06 07 07  08 01 02 03 04 05 06 07  |................|
00000010  07 08 01 02 03 04 05 06  07 07 08 01 02 03 04 05  |................|
00000020  06 07 07 08 01 02 03 04  05 06 07 07 08 01 02 03  |................|
00000030  04 05 06 07 07 08 01 02  03 04 05 06 07 07 08 01  |................|
00000040  02 03 04 05 06 07 07 08  01 02 03 04 05 06 07 07  |................|
00000050  08 01 02 03 04 05 06 07  07 08 01 02 03 04 05 06  |................|
00000060  07 07 08                                          |...|

`,
		Json: "",
	})
	f.Append(fieldInfo.Row{
		Number:   "5",
		Deep:     "",
		Name:     "",
		Kind:     "",
		Size:     "",
		ValueFmt: "0x1122334455667788|1234605616436508552",
		Value:    nil,
		Stream:   "",
		Json:     "",
	})
	f.Append(fieldInfo.Row{
		Number:   "6",
		Deep:     "",
		Name:     "",
		Kind:     "",
		Size:     "",
		ValueFmt: "0x1122334455667788|1234605616436508552",
		Value:    nil,
		Stream:   string(debug.Stack()),
		Json:     "",
	})
	println(f.Gen())
}
