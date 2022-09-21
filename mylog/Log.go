package mylog

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ddkwork/golibrary/src/stream/indent"
	"github.com/fatih/color"
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
		tag:   tagError,
		title: title,
		msg:   Reason(err),
		body:  "",
		debug: o.debug,
	}
	o.do()
	return false
}
func (o *object) HexDump(title string, msg any) {
	*o = object{
		tag:   tagHexDump,
		title: title,
		msg:   hex.Dump(msg.([]byte)),
		body:  "",
		debug: o.debug,
	}
	o.do()
}

func (o *object) Hex(title string, msg any) {
	*o = object{
		tag:   tagHex,
		title: title,
		msg:   fmt.Sprintf("%#x", msg),
		debug: o.debug,
	}
	o.do()
}

func (o *object) Info(title string, msg ...any) {
	*o = object{
		tag:   tagInfo,
		title: title,
		msg:   fmt.Sprint(msg...),
		debug: o.debug,
	}
	o.do()
}

func (o *object) Trace(title string, msg ...any) {
	*o = object{
		tag:   tagTrace,
		title: title,
		msg:   fmt.Sprint(msg...),
		debug: o.debug,
	}
	o.do()
}

func (o *object) Warning(title string, msg ...any) {
	*o = object{
		tag:   tagWarning,
		title: title,
		msg:   fmt.Sprint(msg...),
		debug: o.debug,
	}
	o.do()
}

func (o *object) Json(title string, msg ...any) {
	*o = object{
		tag:   tagJson,
		title: title,
		msg:   fmt.Sprint(msg...),
		debug: o.debug,
	}
	o.do()
}

func (o *object) Success(title string, msg ...any) {
	*o = object{
		tag:   tagSuccess,
		title: title,
		msg:   fmt.Sprint(msg...),
		debug: o.debug,
	}
	o.do()
}

func (o *object) Struct(msg any) {
	msg = reflect.Indirect(reflect.ValueOf(msg)).Interface()
	marshalIndent, err := json.MarshalIndent(msg, "", " ")
	if !o.Error(err) {
		return
	}
	*o = object{
		tag:   tagStruct,
		title: "",
		//msg:   fmt.Sprintf("%#v", msg),
		msg:   string(marshalIndent),
		debug: o.debug,
	}
	o.do()
}
func isTermux() bool {
	dir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	return strings.Contains(dir, "termux")
}
func (o *object) do() (ok bool) {
	///data/data/com.termux/files/home/
	if IsAndroid() {
		//go run .                --> Android
		//run or debug main b     --> linux
		//run or debug unit test  --> linux
		if !isTermux() {
			return
		}
	}
	//2021-05-08 08:42:51 [STRC]                             | struct { a int; b string; c []uint8 }{a:89, b:"jhjsbdd", c:[]uint8{0x11, 0x22, 0x33, 0x44}}
	indentTitle := o.level() + `[` + o.GetTimeNowString() + `]\t` + indent.New().Left(o.title)
	o.cleanMessageStyle()
	if o.msg == "" {
		o.msg = `""`
	}
	caller := " //" + Caller()
	switch o.tag {
	case tagJson, tagHexDump, tagStruct:
		indentTitle += caller + "\n"
		o.body = indentTitle + o.msg
	default:
		o.body = indentTitle + o.msg + caller
	}
	o.printColorBody()
	//todo set apk path as log path
	return o.WriteAppend("log.log", o.body)
}

func (o *object) cleanMessageStyle() {
	indexByte := strings.IndexByte(o.msg, '[')
	if indexByte == 0 {
		//o.msg = strings.Replace(o.msg, "[", "[\n", 1) //特殊处理，fmt.sprint不知道怎么把这个加上了，应该是格式化了切片类型的原因
		//o.msg = strings.ReplaceAll(o.msg, "[", "") //pb klv 或者packed的时候怎么破?
		//o.msg = strings.ReplaceAll(o.msg, "]", "")
		b := []byte(o.msg)
		b = b[1 : len(b)-1]
		o.msg = string(b)
	}
}

const (
	colorFormat = "\x1b[1m\x1b[%dm%s\x1b[0m"
)

func (o *object) printColorBody() {
	ColorBody := ""
	switch o.tag {
	case tagHex:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiBlue, o.body)
	case tagHexDump:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiBlue, o.body)
	case tagJson:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiBlue, o.body)
	case tagStruct:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiBlue, o.body)
	case tagInfo:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiCyan, o.body)
	case tagTrace:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiMagenta, o.body)
	case tagError:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiRed, o.body)
	case tagWarning:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiYellow, o.body)
	case tagSuccess:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiGreen, o.body)
	}
	if o.debug {
		fmt.Println(ColorBody)
	}
}
func (o *object) level() string { return strings.ToUpper(o.tag)[0:4] }

const (
	tagHex     = `hex `
	tagHexDump = `dump`
	tagJson    = `json`
	tagStruct  = `struct`
	tagInfo    = `info`
	tagTrace   = `trace`
	tagError   = `error`
	tagWarning = `warning`
	tagSuccess = `success`
)
