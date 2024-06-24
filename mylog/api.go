package mylog

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/exp/constraints"
)

type (
	object struct {
		kind    kind
		title   string
		message string
		body    string
		debug   bool
		isHttp  bool
		w       *os.File
	}
)

const logFileName = "log.log"

func newObject() *object {
	f, e := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, os.ModePerm)
	if e != nil {
		panic(e)
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
	FormatAllFiles()
}

func TruncateLogFile() { Check(os.Truncate(logFileName, io.SeekStart)) }

var defaultObject = newObject()

func Reason() (reason string)        { return defaultObject.Reason() }
func HexDump(title string, b []byte) { defaultObject.hexDump(title, b) }
func HexInteger[T constraints.Integer](msg T) string {
	return Hex("", msg) + "|" + fmt.Sprintf("%d", msg)
}

func Todo(body string) {
	Warning("TODO", body)
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

func Body() string    { return defaultObject.body }
func Message() string { return defaultObject.message }
