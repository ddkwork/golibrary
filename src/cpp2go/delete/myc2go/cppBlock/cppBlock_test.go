package cppBlock

//
//import (
//	"fmt"
//	"github.com/ddkwork/golibrary/mylog"
//	"github.com/ddkwork/golibrary/src/stream/tool"
//	"strings"
//	"testing"
//)
//
//func Test2(t *testing.T) {
//	println(strings.Index("} ZydisDecoderState;", "}"))
//	println(strings.Contains("} ZydisDecoderState;", "}"))
//	println(strings.Contains(`#elif defined(ZYAN_WINDOWS)`, `#define`))
//	println(strings.Contains(`{`, `{`))
//}
//
//func TestFindAll(t *testing.T) {
//	//p := "Decoder.back"
//	//p := "ntpebteb.h.back"
//	p := "ScriptEngineCommonDefinitions.h.back"
//	//p := "Definition.h" //[DebuggerOutputSourceMaximumRemoteSourceForSingleEvent]; // tags of
//	lines, ok := tool.File().ReadToLines(p)
//	if !ok {
//		panic(111)
//	}
//	l := FindStruct(lines)
//	for _, info := range l {
//		mylog.Info(fmt.Sprint(info.Col), info.Line) //51 - 137
//	}
//	//l := FindEnum(lines)
//	//for _, info := range l {
//	//	mylog.Info(fmt.Sprint(info.Col), info.Line) //147 - 222
//	//}
//}
//
//func TestFindExtern(t *testing.T) {
//	p := "dt-struct.cpp.back"
//	lines, ok := tool.File().ReadToLines(p)
//	if !ok {
//		panic(111)
//	}
//	l := FindExtern(lines)
//	for _, info := range l {
//		mylog.Info(fmt.Sprint(info.Col), info.Line) //147 - 222
//	}
//}
//
//func TestFindDefine(t *testing.T) {
//	//p := "Thread.h.back"
//	p := "ScriptEngineCommonDefinitions.h.back"
//	lines, ok := tool.File().ReadToLines(p)
//	if !ok {
//		panic(111)
//	}
//	l := FindDefine(lines)
//	for _, info := range l {
//		mylog.Info(fmt.Sprint(info.Col), info.Line) //147 - 222
//	}
//}
//
//func TestFindMethod(t *testing.T) {
//	//p := "common.cpp.back"
//	p := "ntrtl.h.back"
//	lines, ok := tool.File().ReadToLines(p)
//	if !ok {
//		panic(111)
//	}
//	l := FindMethod(lines)
//	for _, info := range l {
//		mylog.Info(fmt.Sprint(info.Col), info.Line) //147 - 222
//	}
//}
//func TestFindMethod2(t *testing.T) {
//	lines, ok := tool.File().ToLines(api)
//	if !ok {
//		panic(111)
//	}
//	l := FindMethod(lines)
//	for _, info := range l {
//		mylog.Info(fmt.Sprint(info.Col), info.Line) //147 - 222
//	}
//}
//
//var api = `
//
//
//static ZyanStatus ZydisInputPeek(ZydisDecoderContext* context,
//    ZydisDecodedInstruction* instruction, ZyanU8* value)
//{
//    ZYAN_ASSERT(context);
//    ZYAN_ASSERT(instruction);
//    ZYAN_ASSERT(value);
//
//    if (instruction->length >= ZYDIS_MAX_INSTRUCTION_LENGTH)
//    {
//        return ZYDIS_STATUS_INSTRUCTION_TOO_LONG;
//    }
//
//    if (context->buffer_len > 0)
//    {
//        *value = context->buffer[0];
//        return ZYAN_STATUS_SUCCESS;
//    }
//
//    return ZYDIS_STATUS_NO_MORE_DATA;
//}
//
//
//
//
//`
