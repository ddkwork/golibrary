package mylog

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ddkwork/golibrary/src/stream/indent"
	"github.com/fatih/color"
	"net/http"
	"net/http/httputil"
	"os"
	"reflect"
	"strings"
)

func Reason(err any) string {
	switch err.(type) {
	case error:
		return err.(error).Error()
	case string:
		return err.(string)
	}
	return ""
}
func (o *object) Error(err any) bool {
	if err == nil {
		return true
	}
	return o.error("", err)
}

func (o *object) Error2(_ any, err error) bool {
	if err == nil {
		return true
	}
	return o.error("", err)
}
func (o *object) error(title string, err any) bool {
	*o = object{
		kind:  ErrorKind,
		title: title,
		msg:   Reason(err),
		body:  "",
		debug: o.debug,
	}
	o.printAndWrite()
	return false
}
func (o *object) HexDump(title string, msg any) {
	b := msg.([]byte)
	if len(b) > 257 {
		o.Warning("big data", len(b))
		b = b[:257]
	}
	*o = object{
		kind:  HexDumpKind,
		title: title,
		msg:   hex.Dump(b),
		body:  "",
		debug: o.debug,
	}
	o.printAndWrite()
}

func (o *object) Hex(title string, msg any) {
	*o = object{
		kind:  HexKind,
		title: title,
		msg:   fmt.Sprintf("%#x", msg),
		debug: o.debug,
	}
	o.printAndWrite()
}

func (o *object) Info(title string, msg ...any) {
	*o = object{
		kind:  InfoKind,
		title: title,
		msg:   Sprint(msg...),
		debug: o.debug,
	}
	o.printAndWrite()
}

func (o *object) Trace(title string, msg ...any) {
	*o = object{
		kind:  TraceKind,
		title: title,
		msg:   Sprint(msg...),
		debug: o.debug,
	}
	o.printAndWrite()
}

func (o *object) Warning(title string, msg ...any) {
	*o = object{
		kind:  WarningKind,
		title: title,
		msg:   Sprint(msg...),
		debug: o.debug,
	}
	o.printAndWrite()
}

func (o *object) Json(title string, msg ...any) {
	*o = object{
		kind:  JsonKind,
		title: title,
		msg:   Sprint(msg...), //[]
		debug: o.debug,
	}
	o.printAndWrite()
}

func (o *object) Success(title string, msg ...any) {
	*o = object{
		kind:  SuccessKind,
		title: title,
		msg:   Sprint(msg...),
		debug: o.debug,
	}
	o.printAndWrite()
}

const (
	inHttp  = ">>>--------------->>>--------------->>>--------------->>>\n"
	outHttp = "<<<---------------<<<---------------<<<---------------<<<\n"
	endHttp = "----------------------------------------------------------------------------------------------------------------------------------------\n"
)

func (o *object) DumpRequest(Request *http.Request, body bool) string {
	dumpRequest, err := httputil.DumpRequest(Request, body)
	if !o.Error(err) {
		return ""
	}
	s := inHttp + Request.URL.String() + "\n"
	s += strings.TrimSuffix(string(dumpRequest), "\n")
	s += endHttp
	*o = object{
		kind:   JsonKind,
		title:  "",
		msg:    s,
		body:   "",
		debug:  o.debug,
		isHttp: true,
	}
	o.printAndWrite()
	return s
}

func (o *object) DumpResponse(Response *http.Response, body bool) string {
	if Response == nil {
		return ""
	}
	dumpResponse, err := httputil.DumpResponse(Response, body)
	if !o.Error(err) {
		return ""
	}
	s := outHttp + strings.TrimSuffix(string(dumpResponse), "\n")
	s += endHttp
	*o = object{
		kind:   JsonKind,
		title:  "",
		msg:    s,
		body:   "",
		debug:  o.debug,
		isHttp: true,
	}
	o.printAndWrite()
	return s
}

func (o *object) Struct(msg any) {
	msg = reflect.Indirect(reflect.ValueOf(msg)).Interface()
	marshalIndent, err := json.MarshalIndent(msg, "", " ")
	if !o.Error(err) {
		return
	}
	body := string(marshalIndent)
	if body == "{}" || reflect.TypeOf(msg).Kind() == reflect.Slice {
		body = fmt.Sprintf("%#v", msg) //not export
	}
	*o = object{
		kind:  StructKind,
		title: "",
		msg:   body,
		debug: o.debug,
	}
	o.printAndWrite()
}
func isTermux() bool {
	dir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	return strings.Contains(dir, "termux")
}

func Sprint(msg_ ...any) string {
	msg := fmt.Sprint(msg_...) //遇到[]byte会加上[],
	switch {
	case msg[0] == '[' && msg[len(msg)-1] == ']':
		msg = msg[1 : len(msg)-1]
	}
	return msg
}

func (o *object) printAndWrite() {

	///data/data/com.termux/files/home/
	if IsAndroid() {
		//go run .                --> Android
		//run or debug main b     --> linux
		//run or debug unit test  --> linux
		if !isTermux() {
			return
		}
	}
	indentTitle := o.kind.String() + ` [` + o.GetTimeNowString() + `] ` + indent.New().Left(o.title)
	if o.msg == "" {
		o.msg = `""`
	}
	caller := " //" + Caller()
	switch o.kind {
	case JsonKind, HexDumpKind, StructKind:
		indentTitle += caller + "\n"
		o.body = indentTitle + o.msg
	default:
		o.body = indentTitle + o.msg + caller
	}
	if o.isHttp {
		o.body = o.msg
	}
	o.printColorBody()
	o.WriteAppend("log.log", o.body) //todo set apk path as log path
}

const (
	colorFormat = "\x1b[1m\x1b[%dm%s\x1b[0m"
)

func (o *object) printColorBody() {
	ColorBody := ""
	switch o.kind {
	case HexKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiBlue, o.body)
	case HexDumpKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiBlue, o.body)
	case JsonKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiBlue, o.body)
	case StructKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiBlue, o.body)
	case InfoKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiCyan, o.body)
	case TraceKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiMagenta, o.body)
	case ErrorKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiRed, o.body)
	case WarningKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiYellow, o.body)
	case SuccessKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiGreen, o.body)
	}
	if o.debug {
		fmt.Println(ColorBody)
	}
}

type kind int

var nameKind kind

const (
	HexKind kind = iota
	HexDumpKind
	JsonKind
	StructKind
	InfoKind
	TraceKind
	SuccessKind
	WarningKind
	ErrorKind
)

func (k kind) String() string {
	const (
		Hex     = "INFO hex    "
		HexDump = "INFO hexDump"
		Json    = "INFO json   "
		Struct  = "INFO struct "
		Info    = "INFO        "
		Trace   = "INFO trace  "
		Success = "INFO Success"
		Warning = "WARN Warning"
		Error   = "ERROR       "
	)
	switch k {
	case HexKind:
		return Hex
	case HexDumpKind:
		return HexDump
	case JsonKind:
		return Json
	case StructKind:
		return Struct
	case InfoKind:
		return Info
	case TraceKind:
		return Trace
	case SuccessKind:
		return Success
	case WarningKind:
		return Warning
	case ErrorKind:
		return Error
	}
	return ""
}

//https://github.com/JetBrains/ideolog/wiki/Custom-Log-Formats
//goland 默认安装的日志高亮插件，迎合它的level
//C:\Users\Admin\Downloads\ideolog-master\src\main\kotlin\com\intellij\ideolog\highlighting\LogEvent.kt
//    level = when (rawLevel.toUpperCase()) {
//      "E" -> "ERROR"
//      "W" -> "WARN"
//      "I" -> "INFO"
//      "V" -> "VERBOSE"
//      "D" -> "DEBUG"
//      "T" -> "TRACE"
//      else -> rawLevel.toUpperCase()
//    }
