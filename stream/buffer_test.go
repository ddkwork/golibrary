package stream_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/ddkwork/golibrary/gen"
	"github.com/ddkwork/golibrary/safeType"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	gen.New().FileAction()
	return
	s := stream.New(mylog.Body())
	fmt.Println(s.String())
	fmt.Println(stream.New(mylog.Body()).String())
}

func TestNewHexString(t *testing.T) {
	b := stream.New("1122")
	mylog.HexDump("", b.Bytes())
}

func TestNewReadFile(t *testing.T) {
}

func TestNewReadFrom(t *testing.T) {
}

func Test_hasAddress(t *testing.T) {
}

func Test_object_StreamBuffer(t *testing.T) {
}

func TestCopyDir(t *testing.T) {
}

func TestCreatDirectory(t *testing.T) {
}

func TestCutUint16(t *testing.T) {
}

func TestGetPackageName(t *testing.T) {
}

func TestReadFileAndWriteTruncate(t *testing.T) {
}

func TestWriteAppend(t *testing.T) {
}

func TestWriteBinary(t *testing.T) {
}

func TestWriteGoCode(t *testing.T) {
}

func TestWriteTruncate(t *testing.T) {
}

func Test_object_AppendByteSlice(t *testing.T) {
}

func Test_object_BigNumXorWithAlign(t *testing.T) {
}

func Test_object_CutUint16(t *testing.T) {
}

func Test_object_CutWithIndex(t *testing.T) {
}

func Test_object_DesBlockSizeSizeCheck(t *testing.T) {
}

func Test_object_HexString(t *testing.T) {
}

func Test_object_HexStringUpper(t *testing.T) {
}

func Test_object_Indent(t *testing.T) {
}

func Test_object_InsertBytes(t *testing.T) {
}

func Test_object_InsertString(t *testing.T) {
}

func Test_object_LinesToString(t *testing.T) {
}

func Test_object_MergeByte(t *testing.T) {
}

func Test_object_NewLine(t *testing.T) {
}

func Test_object_ObjectBegin(t *testing.T) {
}

func Test_object_ObjectEnd(t *testing.T) {
}

func Test_object_Quote(t *testing.T) {
}

func Test_object_QuoteWith(t *testing.T) {
}

func Test_object_ReadAny(t *testing.T) {
}

func Test_object_ReadToLines(t *testing.T) {
}

func Test_object_SerialNumber(t *testing.T) {
}

func Test_object_SliceBegin(t *testing.T) {
}

func Test_object_SliceEnd(t *testing.T) {
}

func Test_object_SplitBytes(t *testing.T) {
}

func Test_object_SplitString(t *testing.T) {
}

func Test_object_ToLines(t *testing.T) {
}

func Test_object_WriteAny(t *testing.T) {
}

func Test_object_WriteBytesLn(t *testing.T) {
}

func Test_object_WritePackageName(t *testing.T) {
}

func Test_object_WriteStringLn(t *testing.T) {
}

func TestStream_AppendByteSlice(t *testing.T) {
	s := stream.New("")
	s.AppendByteSlice([]byte{0x11}, []byte{0x22})
	mylog.HexDump("", s.Bytes())
}

func TestStream_BigNumXorWithAlign(t *testing.T) { // todo test
	// s := stream.newObject("")
	// s.BigNumXorWithAlign()
}

func TestStream_CutWithIndex(t *testing.T) {
	mylog.HexDump("", stream.New([]byte{0x11, 0x22, 0x33, 0x44, 0x55}).CutWithIndex(2, 4))
}

func TestStream_DesBlockSizeSizeCheck(t *testing.T) {
	stream.New([]byte{0x11, 0x22, 0x33, 0x44, 0x55}).CheckDesBlockSize()
}

func TestStream_HexString(t *testing.T) {
	println(stream.New([]byte{0x11, 0x22, 0x33, 0x44, 0x55}).HexString())
}

func TestStream_HexStringUpper(t *testing.T) {
	println(stream.New([]byte{0xaa, 0xbb, 0x33, 0x44, 0x55}).HexStringUpper())
}

func TestStream_Indent(t *testing.T) {
	s := stream.New("1111")
	s.Indent(3)
	s.WriteString("3344")
	println(s.String())
}

func TestStream_LinesToString(t *testing.T) {
	// newObject().LinesToString(tt.args.lines)
}

func TestStream_Merge(t *testing.T) {
	mylog.HexDump("", stream.New([]byte{0xaa, 0xbb, 0x33, 0x44, 0x55}).MergeByte([]byte{0xaa, 0xbb, 0x33, 0x44, 0x55}))
	concat := stream.Concat([]byte{0xaa, 0xbb, 0x33, 0x44, 0x55}, []byte{0xaa, 0xbb, 0x33, 0x44, 0x55})
	mylog.HexDump("", concat)
}

func TestStream_NewLine(t *testing.T) {
	s := stream.New("111")
	s.NewLine()
	println(s.String())
}

func TestStream_ObjectBegin(t *testing.T) {
	s := stream.New("111")
	s.NewLine()
	s.ObjectBegin()
	println(s.String())
}

func TestStream_ObjectEnd(t *testing.T) {
	s := stream.New("111")
	s.NewLine()
	s.ObjectEnd()
	println(s.String())
}

func TestStream_Quote(t *testing.T) {
	s := stream.New("111")
	s.Quote()
	println(s.String())
}

func TestStream_QuoteWith(t *testing.T) {
	s := stream.New("111")
	s.NewLine()
	s.QuoteWith("//")
	println(s.String())
}

func TestStream_SliceBegin(t *testing.T) {
	s := stream.New("111")
	s.NewLine()
	s.SliceBegin()
	println(s.String())
}

func TestStream_SliceEnd(t *testing.T) {
	s := stream.New("111")
	s.NewLine()
	s.SliceEnd()
	println(s.String())
}

func TestStream_SplitBytes(t *testing.T) {
}

func TestStream_SplitString(t *testing.T) {
}

func TestStream_WriteBytesLn(t *testing.T) {
}

func TestStream_WriteStringLn(t *testing.T) {
}

func TestStream_InsertString(t *testing.T) {
	s := stream.New("B3EBDDAF2EA789C4")
	s.InsertStringWithSplit(4, "-")
	assert.Equal(t, "B3EB-DDAF-2EA7-89C4", s.String())

	s = stream.New("aaaaa")
	s.InsertString(1, "22")
	println(s.String())
}

func ToLines[T string | []byte](data T) (lines []string, ok bool) {
	var r io.Reader
	switch data := any(data).(type) {
	case string:
		r = strings.NewReader(data)
	case []byte:
		r = bytes.NewReader(data)
	}

	lines = make([]string, 0)
	newReader := bufio.NewReader(r)

	for {
		line, _, err := newReader.ReadLine()
		switch err {
		case io.EOF:
			return lines, true
		default:
			if !mylog.Error(err) {
				return
			}
		}
		lines = append(lines, string(line))
	}
}

func ToLinesByBuffer(b *bytes.Buffer) (lines []string, ok bool) {
	lines = make([]string, 0)
	newReader := bufio.NewReader(b)
	for {
		line, _, err := newReader.ReadLine()
		switch err {
		case io.EOF:
			return lines, true
		default:
			if !mylog.Error(err) {
				return
			}
		}
		lines = append(lines, string(line))
	}
}

func main() {
	// 打开文件
	file, err := os.Open("D:\\clone\\HyperDbg\\hyperdbg\\demo\\xxx.js")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// 创建Scanner来扫描文件
	scanner := bufio.NewScanner(file)

	// 创建一个切片来存储文件的每一行
	var lines []string

	// 循环读取文件的每一行并存储到切片中
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// 检查扫描过程中是否出现错误
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 输出切片中的内容
	//for _, line := range lines {
	//	fmt.Println(line)
	//}
}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main()
	}
}

func BenchmarkByBufioReaderReadLine(b *testing.B) {
	path := "D:\\clone\\HyperDbg\\hyperdbg\\demo\\xxx.js"
	path = "C:\\Users\\Admin\\go\\pkg\\mod\\fyne.io\\fyne\\v2@v2.4.4\\theme\\bundled-emoji.go"
	s := stream.NewReadFile(path)
	s.String()

	//file, err := os.ReadFile("D:\\clone\\HyperDbg\\hyperdbg\\demo\\xxx.js")
	//if err != nil {
	//	return
	//}

	for i := 0; i < b.N; i++ {
		s.ToLines() // good
		// BenchmarkReaderReadLine-8         513292              2254 ns/op
		// BenchmarkReaderReadLine-8         348691              3156 ns/op
		// BenchmarkReaderReadLine-8         390274              3068 ns/op
		// BenchmarkReaderReadLine-8         471577              2296 ns/op
		// BenchmarkReaderReadLine-8         325484              3568 ns/op
		// BenchmarkReaderReadLine-8         370736              2860 ns/op
		// BenchmarkReaderReadLine-8         390140              3054 ns/op

		// ToLines(s.Bytes())//bad, why
		// ToLines(file) //bad, why
		// ToLinesByBuffer(bytes.NewBuffer(file)) //bad, why
		// BenchmarkReaderReadLine-8            181           6340871 ns/op
		// BenchmarkReaderReadLine-8             98          15149509 ns/op
		// BenchmarkReaderReadLine-8            138           7323254 ns/op
		// BenchmarkReaderReadLine-8            100          10153358 ns/op
	}
}

func BenchmarkByStringsSplitAfter(b *testing.B) {
	s := stream.NewReadFile("D:\\clone\\HyperDbg\\hyperdbg\\demo\\xxx.js")
	for i := 0; i < b.N; i++ {
		strings.SplitAfter(s.String(), "\n") // about
		// BenchmarkByStringsSplitAfter-8               355           3617550 ns/op
		// BenchmarkByStringsSplitAfter-8               332           4132401 ns/op
		// BenchmarkByStringsSplitAfter-8               248           4461010 ns/op
		// BenchmarkByStringsSplitAfter-8               370           3778242 ns/op
		// BenchmarkByStringsSplitAfter-8               247           4556421 ns/op
		// BenchmarkByStringsSplitAfter-8               303           4366858 ns/op
	}
}

func BenchmarkByBytessSplitAfter(b *testing.B) {
	s := stream.NewReadFile("D:\\clone\\HyperDbg\\hyperdbg\\demo\\xxx.js")
	for i := 0; i < b.N; i++ {
		bytes.SplitAfter(s.Bytes(), []byte("\n")) // about
		// BenchmarkByBytessSplitAfter-8                613           1812570 ns/op
		// BenchmarkByBytessSplitAfter-8                450           2775832 ns/op
		// BenchmarkByBytessSplitAfter-8                468           2581962 ns/op
		// BenchmarkByBytessSplitAfter-8                460           3021967 ns/op
		// BenchmarkByBytessSplitAfter-8                430           3010534 ns/op
	}
}

func BenchmarkByRegexp(b *testing.B) { // todo
	s := stream.NewReadFile("D:\\clone\\HyperDbg\\hyperdbg\\demo\\xxx.js")
	re := regexp.MustCompile(`\n`)
	for i := 0; i < b.N; i++ {
		re.Split(s.String(), -1) // about
		// BenchmarkByRegexp-8           48          22047369 ns/op
		// BenchmarkByRegexp-8           36          32086147 ns/op
		// BenchmarkByRegexp-8           38          32875274 ns/op
		// BenchmarkByRegexp-8           62          18963963 ns/op
		// BenchmarkByRegexp-8           34          32556691 ns/op
	}
}

// Split a string on "\n" while preserving them. The output can be used
// as input for UnifiedDiff and ContextDiff structures.
func SplitLines(s string) []string {
	// bytes.Split()
	// bytes.SplitAfter()
	lines := strings.SplitAfter(s, "\n")
	lines[len(lines)-1] += "\n"
	return lines
}

func TestWrite(t *testing.T) {
	body := `
-Xms256m
-Xmx2000m
-XX:ReservedCodeCacheSize=512m
-Xss2m
-XX:NewSize=128m
-XX:MaxNewSize=128m
-XX:+IgnoreUnrecognizedVMOptions
-XX:+UseG1GC
-XX:SoftRefLRUPolicyMSPerMB=50
-XX:CICompilerCount=2
-XX:+HeapDumpOnOutOfMemoryError
-XX:-OmitStackTraceInFastThrow
-ea
-Dsun.io.useCanonCaches=false
-Djdk.http.auth.tunneling.disabledSchemes=""
-Djdk.attach.allowAttachSelf=true
-Djdk.module.illegalAccess.silent=true
-Dkotlinx.coroutines.debug=off
-Dsun.tools.attach.tmp.only=true
`
	ok := stream.WriteTruncate("clion64.vmoptions", body)
	assert.True(t, ok)

	assert.True(t, stream.WriteAppend("1.txt", "111"))
	assert.True(t, stream.WriteAppend("1.txt", "222"))
	os.Remove("1.txt")
	os.Remove("clion64.vmoptions")

	// stream.CopyDir("include/Common/EABase", "new/include")
}

func TestSwap(t *testing.T) {
	s := stream.New("")
	assert.Equal(t, byte(0x16), s.CutUint16(0x6613))
	assert.Equal(t, byte(0x16), s.CutUint16(0x6613))

	assert.Equal(t, "ET5AA5Q3N2KTR8      ", stream.New("TEA55A3Q2NTK8R      ").SerialNumber())
	assert.Equal(t, "TA1591503892      ", stream.New("AT5119058329      ").SerialNumber())

	assert.Equal(t, safeType.HexString("8877665544332211"), stream.SwapBytes2HexString(stream.NewHexString("1122334455667788").Bytes()))

	println(stream.SwapBytes2HexString2("1122334455667788"))
}

//	func TestInference(t *testing.T) {
//		s1 := []int{1, 2, 3}
//		apply(s1, slices.Reverse)
//		if want := []int{3, 2, 1}; !slices.Equal(s1, want) {
//			t.Errorf("Reverse(%v) = %v, want %v", []int{1, 2, 3}, s1, want)
//		}
//
//		type S []int
//		s2 := S{4, 5, 6}
//		apply(s2, slices.Reverse)
//		if want := (S{6, 5, 4}); !slices.Equal(s2, want) {
//			t.Errorf("Reverse(%v) = %v, want %v", S{4, 5, 6}, s2, want)
//		}
//	}
func apply[T any](v T, f func(T)) {
	f(v)
}

func TestInsert(t *testing.T) {
	mylog.Struct(slices.Insert([]int{1, 2, 3}, 0, 1))
	mylog.Struct(slices.Insert([]byte{1, 2, 3}, 2, 1))
}

func TestCaseconv(t *testing.T) {
	for _, s := range name {
		println(stream.ToCamelUpper(s, true))
	}
}

var name = []string{
	"	HIDDEN_HOOK_READ_AND_WRITE                                                 ",
	"	HIDDEN_HOOK_READ                                                           ",
	"	HIDDEN_HOOK_WRITE                                                          ",
	"	HIDDEN_HOOK_EXEC_DETOURS                                                   ",
	"	HIDDEN_HOOK_EXEC_CC                                                        ",
	"	SYSCALL_HOOK_EFER_SYSCALL                                                  ",
	"	SYSCALL_HOOK_EFER_SYSRET                                                   ",
	"	CPUID_INSTRUCTION_EXECUTION                                                ",
	"	RDMSR_INSTRUCTION_EXECUTION                                                ",
	"	WRMSR_INSTRUCTION_EXECUTION                                                ",
	"	IN_INSTRUCTION_EXECUTION                                                   ",
	"	OUT_INSTRUCTION_EXECUTION                                                  ",
	"	EXCEPTION_OCCURRED                                                         ",
	"	EXTERNAL_INTERRUPT_OCCURRED                                                ",
	"	DEBUG_REGISTERS_ACCESSED                                                   ",
	"	TSC_INSTRUCTION_EXECUTION                                                  ",
	"	PMC_INSTRUCTION_EXECUTION                                                  ",
	"	VMCALL_INSTRUCTION_EXECUTION                                               ",
	"	CONTROL_REGISTER_MODIFIED                                                  ",
	"	DEBUGGER_EVENT_TYPE_ENUM                                                   ",
	"	BREAK_TO_DEBUGGER                                                          ",
	"	RUN_SCRIPT                                                                 ",
	"	RUN_CUSTOM_CODE                                                            ",
	"	DEBUGGER_EVENT_SYSCALL_SYSRET_SAFE_ACCESS_MEMORY                           ",
	"	DEBUGGER_EVENT_SYSCALL_SYSRET_HANDLE_ALL_UD                                ",
	"	DEBUGGER_MODIFY_EVENTS_QUERY_STATE                                         ",
	"	DEBUGGER_MODIFY_EVENTS_ENABLE                                              ",
	"	DEBUGGER_MODIFY_EVENTS_DISABLE                                             ",
	"	DEBUGGER_MODIFY_EVENTS_CLEAR                                               ",
	"	DEBUGGER_MODIFY_EVENTS_TYPE                                                ",
	"	struct__DEBUGGER_MODIFY_EVENTS                                             ",
	"	VirtualAddress                                                             ",
	"	ProcessId                                                                  ",
	"	Pml4eVirtualAddress                                                        ",
	"	Pml4eValue                                                                 ",
	"	PdpteVirtualAddress                                                        ",
	"	PdpteValue                                                                 ",
	"	PdeVirtualAddress                                                          ",
	"	PdeValue                                                                   ",
	"	PteVirtualAddress                                                          ",
	"	PteValue                                                                   ",
	"	KernelStatus                                                               ",
	"	DEBUGGER_READ_PAGE_TABLE_ENTRIES_DETAILS                                   ",
	"	PDEBUGGER_READ_PAGE_TABLE_ENTRIES_DETAILS                                  ",
	"	struct__DEBUGGER_VA2PA_AND_PA2VA_COMMANDS                                  ",
	"	DEBUGGER_VA2PA_AND_PA2VA_COMMANDS                                          ",
	"	PDEBUGGER_VA2PA_AND_PA2VA_COMMANDS                                         ",
	"	struct__DEBUGGER_DT_COMMAND_OPTIONS                                        ",
	"	ypeName                                                                    ",
	"	DEBUGGER_SHOW_COMMAND_DT                                                   ",
	"	DEBUGGER_SHOW_COMMAND_DISASSEMBLE64                                        ",
	"	DEBUGGER_SHOW_COMMAND_DISASSEMBLE32                                        ",
	"	DEBUGGER_SHOW_COMMAND_DB                                                   ",
	"	DEBUGGER_SHOW_COMMAND_DC                                                   ",
	"	DEBUGGER_SHOW_COMMAND_DQ                                                   ",
	"	DEBUGGER_SHOW_COMMAND_DD                                                   ",
}

func TestCutUint161(t *testing.T) {
	type args struct {
		u uint16
	}
	var tests []struct {
		name string
		args args
		want uint8
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, stream.CutUint16(tt.args.u), "CutUint16(%v)", tt.args.u)
		})
	}
}
