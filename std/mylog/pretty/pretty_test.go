package pretty

import (
	"bytes"
	"testing"
)

func TestPrettyFormat2(test *testing.T) {
	key := "Path"
	value := `C:\Program Files\Eclipse Adoptium\jdk-23.0.2.7-hotspot\bin;D:\todo\ewdk\dist\sdk\bin\Hostx64\x64;E:\Program Files\Microsoft Visual Studio\2022\BuildTools\VC\Tools\MSVC\14.44.35207\bin\HostX64\x64;E:\Program Files\Windows Kits\10\bin\10.0.28000.0\x64;C:\Windows\system32;C:\Windows;C:\Windows\System32\Wbem;C:\Windows\System32\WindowsPowerShell\v1.0\;C:\Windows\System32\OpenSSH\;C:\Program Files\Go\bin;C:\Program Files\Git\cmd;C:\Program Files\CMake\bin;C:\TDM-GCC-64\bin;C:\Program Files\LLVM\bin;C:\Program Files\CodeArts Agent\bin`
	m := map[string]string{key: value}
	type Struct struct {
		Path  string
		Value string
	}
	obj := Struct{
		Path:  key,
		Value: value,
	}
	test.Log(Format(m))
	test.Log(Format(obj))
}

type Bag map[string]any

type Struct struct {
	N int
	S string
	B bool
	A []int
	Z []int
}

var (
	ch chan string

	s = struct {
		n int
		s string
	}{
		42,
		"hello world",
	}

	x = struct{}{}

	array = []Bag{bag, bag, bag}

	strutty = Struct{N: 42, S: "Hello", B: true, A: []int{1, 2, 3}}

	bag = Bag{
		"a": 1,
		"b": false,
		"c": "some stuff",
		"d": []float64{0.0, 0.1, 1.2, 1.23, 1.23456, 999999999999},
		"e": Bag{
			"e1": "here",
			"e2": []int{1, 2, 3, 4},
			"e3": nil,
			"e4": s,
		},
		"s":   s,
		"x":   x,
		"z":   []int{},
		"bad": ch,
	}
)

func TestPrettyPrint(test *testing.T) {
	Print(array)
}

func TestPrettyFormat(test *testing.T) {
	test.Log(Format(bag))
}

func TestStruct(test *testing.T) {
	test.Log(Format(strutty))
}

func TestPretty(test *testing.T) {
	var out bytes.Buffer
	p := Pretty{Indent: "", Out: &out, NilString: "nil"}
	p.Print(strutty)
	test.Log(out.String())
}

func TestPrettyCompact(test *testing.T) {
	var out bytes.Buffer
	p := Pretty{Indent: "", Out: &out, NilString: "nil", Compact: true}
	p.Print(strutty)
	test.Log(out.String())
}

func TestPrettyLevel(test *testing.T) {
	var out bytes.Buffer
	p := Pretty{Indent: "", Out: &out, NilString: "nil", Compact: true, MaxLevel: 2}
	p.Print(bag)
	test.Log(out.String())
}

func Example_tabPrint() {
	tp := NewTabPrinter(8)

	for i := range 33 {
		tp.Print(i)
	}

	tp.Println()

	for _, v := range []string{"one", "two", "three", "four", "five", "six"} {
		tp.Print(v)
	}

	tp.Println()

	// Output:
	// 0	1	2	3	4	5	6	7
	// 8	9	10	11	12	13	14	15
	// 16	17	18	19	20	21	22	23
	// 24	25	26	27	28	29	30	31
	// 32
	// one	two	three	four	five	six
}

func Example_tabPrintTwoFullLines() {
	tp := NewTabPrinter(4)

	for _, v := range []string{"one", "two", "three", "four", "five", "six", "seven", "eight"} {
		tp.Print(v)
	}

	tp.Println()

	// Output:
	// one	two	three	four
	// five	six	seven	eight
	//
}

func Example_tabPrintWider() {
	tp := NewTabPrinter(2)
	tp.TabWidth(10)

	for _, v := range []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "larger", "largest", "even more", "enough"} {
		tp.Print(v)
	}

	tp.Println()

	// Output:
	// one       two
	// three     four
	// five      six
	// seven     eight
	// larger    largest
	// even more enough
}
