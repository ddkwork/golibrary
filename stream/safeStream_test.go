package stream_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

func TestNewHexDump(t *testing.T) {
	stream.NewHexDump(stream.HexDumpString(dump))
	stream.NewHexDump(stream.HexDumpString(bugBuf))
}

func TestHexDumpToGoBytes(t *testing.T) {
	ss := `00 00 00 1A 00 00 00 09 00 01 00 00 0B 00 00 00 8E 6A 64 01 15 4F 53 44 4B 5F 41 42 55 53 45 5F 52 45 50 4F 52 54 49 4E 47 00`
	stream.NewBuffer(stream.HexDumpString(ss))
	stream.NewHexDump(stream.HexDumpString(dump))
	stream.NewHexDump(`00 00 00 1A 00 00 00 09 00 01 00 00 0B 00 00 00`)
	stream.NewHexDump(`8E 6A 64 01`)
	stream.NewHexDump(`01`)
	stream.NewHexDump(`15`)
	stream.NewHexDump(`4F 53 44 4B 5F 41 42 55 53 45 5F 52 45 50 4F 52 54 49 4E 47`)
	stream.NewHexDump(`00`)
}

var bugBuf = `
08A73200 57 61 72 68 61 6D 6D 65 72 20 34 30 2C 30 30 30  Warhammer 40,000  
08A73210 20 44 61 77 6E 20 6F 66 20 57 61 72 20 49 49 20   Dawn of War II   
08A73220 52 65 74 72 69 62 75 74 69 6F 6E 20 2D 20 49 6D  Retribution - Im  
08A73230 70 65 72 69 61 6C 20 46 69 73 74 73 20 43 68 61  perial Fists Cha  
08A73240 70 74 65 72 20 50 61 63 6B 00 00 00 00 00 00 00  pter Pack.......  
08A73250 57 61 72 68 61 6D 6D 65 72 20 34 30 2C 30 30 30  Warhammer 40,000  
08A73260 3A 20 44 61 77 6E 20 6F 66 20 57 61 72 20 49 49  : Dawn of War II  
08A73270 20 2D 20 52 65 74 72 69 62 75 74 69 6F 6E 20 2D   - Retribution -  
08A73280 20 49 6D 70 65 72 69 61 6C 20 46 69 73 74 73 20   Imperial Fists   
08A73290 43 68 61 70 74 65 72 20 50 61 63 6B 00 00 00 00  Chapter Pack....  
08A732A0 50 65 6E 6E 79 20 41 72 63 61 64 65 20 41 64 76  Penny Arcade Adv  
08A732B0 65 6E 74 75 72 65 73 20 4F 6E 20 74 68 65 20 52  entures On the R  
08A732C0 61 69 6E 2D 53 6C 69 63 6B 20 50 72 65 63 69 70  ain-Slick Precip  

`

var dump = `
00000000  7e 15 00 80 0b 00 00 00  09 25 ce f7 3d 01 00 10  |~........%..=...|
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
000000d0  a7 b2 98 85 d3 54                                 |.....T|
`

func TestIsFilePath(t *testing.T) {
	assert.False(t, stream.IsFilePath("wss://alive.github.com/_sockets/u/19886504/ws?ses"))
}

func Test_trimSlash(t *testing.T) {
	stream.TrimSlash("asserts\\1.jpg")
}

func TestAlignString(t *testing.T) {
	fmt.Println(strconv.Quote(stream.AlignString("中文SetHan═╬═dles(ha电═╬═锅锅ndles []Handle)", 55)))
	fmt.Println(strconv.Quote(stream.AlignString("Handlesjk═╬═js 看见你地方df() []Handf的 dle", 55)))
	fmt.Println(strconv.Quote(stream.AlignString("en═╬═flish", 55)))
}

func TestIsDirDeep1(t *testing.T) {
	println(stream.IsDirRoot("pkg\\cpp2go\\cpp"))
	println(stream.IsDirRoot(".git"))
}

func TestSubDays(t *testing.T) {
	println(stream.GetDiffDays("2024-05-26"))
}

func Test_getUserConfigDirs(t *testing.T) {
	userConfigDirs := stream.GetUserConfigDirs()
	for username, ConfigDir := range userConfigDirs {
		fmt.Println(username + ": " + ConfigDir)
	}
}

func TestNew(t *testing.T) {
	return
	s := stream.NewBuffer(mylog.Body())
	fmt.Println(s.String())
	fmt.Println(stream.NewBuffer(mylog.Body()).String())
}

func TestNewHexString(t *testing.T) {
	b := stream.NewBuffer("1122")
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
	s := stream.NewBuffer("")
	s.AppendByteSlice([]byte{0x11}, []byte{0x22})
	mylog.HexDump("", s.Bytes())
}

func TestStream_BigNumXorWithAlign(t *testing.T) {
}

func TestStream_CutWithIndex(t *testing.T) {
	mylog.HexDump("", stream.NewBuffer([]byte{0x11, 0x22, 0x33, 0x44, 0x55}).CutWithIndex(2, 4))
}

func TestStream_DesBlockSizeSizeCheck(t *testing.T) {
	mylog.Call(func() {
		stream.NewBuffer([]byte{0x11, 0x22, 0x33, 0x44, 0x55}).CheckDesBlockSize()
	})
}

func TestStream_HexString(t *testing.T) {
	println(stream.NewBuffer([]byte{0x11, 0x22, 0x33, 0x44, 0x55}).HexString())
}

func TestStream_HexStringUpper(t *testing.T) {
	println(stream.NewBuffer([]byte{0xaa, 0xbb, 0x33, 0x44, 0x55}).HexStringUpper())
}

func TestStream_Indent(t *testing.T) {
	s := stream.NewBuffer("1111")
	s.Indent(3)
	s.WriteString("3344")
	println(s.String())
}

func TestStream_LinesToString(t *testing.T) {
}

func TestStream_Merge(t *testing.T) {
	mylog.HexDump("", stream.NewBuffer([]byte{0xaa, 0xbb, 0x33, 0x44, 0x55}).MergeByte([]byte{0xaa, 0xbb, 0x33, 0x44, 0x55}))
	concat := stream.Concat([]byte{0xaa, 0xbb, 0x33, 0x44, 0x55}, []byte{0xaa, 0xbb, 0x33, 0x44, 0x55})
	mylog.HexDump("", concat)
}

func TestStream_NewLine(t *testing.T) {
	s := stream.NewBuffer("111")
	s.NewLine()
	println(s.String())
}

func TestStream_ObjectBegin(t *testing.T) {
	s := stream.NewBuffer("111")
	s.NewLine()
	s.ObjectBegin()
	println(s.String())
}

func TestStream_ObjectEnd(t *testing.T) {
	s := stream.NewBuffer("111")
	s.NewLine()
	s.ObjectEnd()
	println(s.String())
}

func TestStream_Quote(t *testing.T) {
	s := stream.NewBuffer("111")
	s.Quote()
	println(s.String())
}

func TestStream_QuoteWith(t *testing.T) {
	s := stream.NewBuffer("111")
	s.NewLine()
	s.QuoteWith("//")
	println(s.String())
}

func TestStream_SliceBegin(t *testing.T) {
	s := stream.NewBuffer("111")
	s.NewLine()
	s.SliceBegin()
	println(s.String())
}

func TestStream_SliceEnd(t *testing.T) {
	s := stream.NewBuffer("111")
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
	s := stream.NewBuffer("B3EBDDAF2EA789C4")
	s.InsertStringWithSplit(4, "-")
	assert.Equal(t, "B3EB-DDAF-2EA7-89C4", s.String())

	s = stream.NewBuffer("aaaaa")
	s.InsertString(1, "22")
	println(s.String())
}

func main() {
	file := mylog.Check2(os.Open("D:\\clone\\HyperDbg\\hyperdbg\\demo\\xxx.js"))
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	mylog.Check(scanner.Err())
}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main()
	}
}

func BenchmarkByBufioReaderReadLine(b *testing.B) {
	path := "D:\\clone\\HyperDbg\\hyperdbg\\demo\\xxx.js"
	path = "C:\\Users\\Admin\\go\\pkg\\mod\\fyne.io\\fyne\\v2@v2.4.4\\theme\\bundled-emoji.go"
	s := stream.NewBuffer(path)
	for i := 0; i < b.N; i++ {
		s.ToLines()
	}
}

func BenchmarkByStringsSplitAfter(b *testing.B) {
	s := stream.NewBuffer("D:\\clone\\HyperDbg\\hyperdbg\\demo\\xxx.js")
	for i := 0; i < b.N; i++ {
		strings.SplitAfter(s.String(), "\n")
	}
}

func BenchmarkByBytessSplitAfter(b *testing.B) {
	s := stream.NewBuffer("D:\\clone\\HyperDbg\\hyperdbg\\demo\\xxx.js")
	for i := 0; i < b.N; i++ {
		bytes.SplitAfter(s.Bytes(), []byte("\n"))
	}
}

func BenchmarkByRegexp(b *testing.B) {
	s := stream.NewBuffer("D:\\clone\\HyperDbg\\hyperdbg\\demo\\xxx.js")
	re := regexp.MustCompile(`\n`)
	for i := 0; i < b.N; i++ {
		re.Split(s.String(), -1)
	}
}

func SplitLines(s string) []string {
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
	stream.WriteTruncate("clion64.vmoptions", body)

	stream.WriteAppend("1.txt", "111")
	stream.WriteAppend("1.txt", "222")
	mylog.Check(os.Remove("1.txt"))
	mylog.Check(os.Remove("clion64.vmoptions"))
}

func TestSwap(t *testing.T) {
	assert.Equal(t, byte(0x16), stream.CutUint16(0x6613))
	assert.Equal(t, byte(0x16), stream.CutUint16(0x6613))

	assert.Equal(t, "ET5AA5Q3N2KTR8      ", stream.SerialNumber("TEA55A3Q2NTK8R      "))
	assert.Equal(t, "TA1591503892      ", stream.SerialNumber("AT5119058329      "))

	assert.Equal(t, stream.HexString("8877665544332211"), stream.SwapBytes2HexString(stream.NewHexString("1122334455667788").Bytes()))

	println(stream.SwapBytes2HexString2("1122334455667788"))
}

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

func TestToCamelUpper(t *testing.T) {
	println(stream.ToCamelUpper("PAGE_SIZE", false))
}
