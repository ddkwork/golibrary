package mylog

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/constraints"
	"net/http"
	"os"
)

type (
	object struct {
		kind    kind
		title   string
		message string
		body    string
		debug   bool
		isHttp  bool // not use time and line number
		w       *os.File
	}
)

const logFileName = "log.log" // todo support android and termux

func newObject() *object {
	f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0o666)
	if err != nil {
		panic(err)
	}
	return &object{
		-1,
		"",
		"",
		"",
		true,
		false,
		f,
	}
}

func init() {
	if IsAndroid() {
		SetDebug(false)
	}
	TruncateLogFile()
	Trace("--------- title ---------", "------------------ info ------------------")
}
func TruncateLogFile() { _ = defaultObject.w.Truncate(0) }

var defaultObject = newObject()

// True 又要安全的接口断言，又不能影响可读性，只能检测ok强制不返回了，但是可以捕捉到错误
func True(b bool) {
	if !b {
		Error("maybe interface assertion was failed ")
	}
}

// func Assert(t *testing.T) *assert.Assertions { return assert.New(t) }
func Error(err any) bool             { return defaultObject.Error(err) }
func Error2(_ any, err error) bool   { return defaultObject.Error2(nil, err) }
func Reason() (reason string)        { return defaultObject.Reason() }
func HexDump(title string, b []byte) { defaultObject.hexDump(title, b) }
func HexInteger[T constraints.Integer](msg T) string {
	return Hex("", msg) + "|" + fmt.Sprintf("%d", msg)
}

func Hex[T constraints.Integer | []byte | *bytes.Buffer](title string, msg T) string {
	return defaultObject.Hex(title, msg)
}
func Info(title string, msg ...any)                   { defaultObject.Info(title, msg...) }
func Trace(title string, msg ...any)                  { defaultObject.Trace(title, msg...) }
func Warning(title string, msg ...any)                { defaultObject.Warning(title, msg...) }
func MarshalJson(title string, msg any)               { defaultObject.MarshalJson(title, msg) }
func Json(title string, msg ...any)                   { defaultObject.Json(title, msg...) }
func Success(title string, msg ...any)                { defaultObject.Success(title, msg...) }
func Struct(msg any)                                  { defaultObject.Struct(msg) }
func SetDebug(debug bool)                             { defaultObject.debug = debug }
func Request(Request *http.Request, body bool)        { defaultObject.Request(Request, body) }
func Response(Response *http.Response, body bool)     { defaultObject.Response(Response, body) }
func DumpRequest(req *http.Request, body bool) string { return defaultObject.DumpRequest(req, body) }
func DumpResponse(resp *http.Response, body bool) string {
	return defaultObject.DumpResponse(resp, body)
}

// Body log包不能依赖任何包，目前的引用逻辑是这样的:
// stream -->safeType-->log safeType包单独提出来以便后续的更多封装调用它，比如文本编辑器等所有输入类型的控件
// widget -->safeType-->log  文本编辑器设置数据借safeType支持了更多类似的设置支持
// safeType -->log
func Body() string    { return defaultObject.body }    // 至于这里不能被任何包的对象实例化，否则就会访问冲突 c0005
func Message() string { return defaultObject.message } // 用于抓包的http dump显示
