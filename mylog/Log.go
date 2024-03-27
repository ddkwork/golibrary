package mylog

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"reflect"
	"regexp"
	"strings"
)

// 泛型的类型约束
// 1 约束普通函数
// 2 约束对象方法，在实例化对象的时候约束，对象内应该构造一个能处理所有被约束类型的实现
// 3 泛型方法,实例化的时候传入泛型形参，每次实例化的时候就不用指定类型了，方法的形参也会自动推导类型，这样才通用，否则一个泛型对象只能接收一种类型的参数不方便
func reason(err any) string {
	switch err.(type) {
	case error:
		return err.(error).Error()
	case string:
		return err.(string)
	}
	return ""
}

// Error2 最佳的go语言错误处理调用规则：
// 1.所有正常日志错误日志均包含行号
// 2.一个err返回直接调用Error(err any)方法，一句完事
// 3.两个err返回，如果要忽略除了错误之外的返回值，比如写文件后返回的长度，用不到，这种情况直接一行代码调用Error2(_ any, err error)
// 4.两个以上返回值+err的情况一般我们是需要处理返回值的，像windows syscall之类的返回值，这种情况调用Error(err any)，写成两行代码
// 更多错误处理代码可读性的示例请看单元测试
func (o *object) Error2(_ any, err error) bool { return o.Error(err) } // error type only

// string or error type,string类型说白了就是用户总费用错误类型，和errors.newObject()一样
func (o *object) Error(err any) bool {
	if err == nil { // string 也是any，只不过它的值不可能是nil，这就巧妙的处理了自定义错误类型，说到底官方的实现最终还是返回字符串类型的方法
		return true
	}
	*o = object{
		kind:    errorKind,
		title:   "",
		message: reason(err),
		body:    "",
		debug:   o.debug,
		isHttp:  false,
		w:       o.w,
	}
	o.printAndWrite()
	return false
}

func (o *object) hexDump(title string, b []byte) {
	if len(b) > 257 {
		o.Warning("big data", len(b))
		b = b[:257]
	}
	*o = object{
		kind:    hexDumpKind,
		title:   title,
		message: hex.Dump(b),
		body:    "",
		debug:   o.debug,
		isHttp:  false,
		w:       o.w,
	}
	o.printAndWrite()
}

func (o *object) Hex(title string, msg any) string {
	*o = object{
		kind:    hexKind,
		title:   title,
		message: fmt.Sprintf("%#x", msg),
		body:    "",
		debug:   o.debug,
		isHttp:  false,
		w:       o.w,
	}
	o.printAndWrite()
	return o.message
}

func (o *object) Info(title string, msg ...any) {
	*o = object{
		kind:    infoKind,
		title:   title,
		message: sprint(msg...),
		body:    "",
		debug:   o.debug,
		isHttp:  false,
		w:       o.w,
	}
	o.printAndWrite()
}

func (o *object) Trace(title string, msg ...any) {
	*o = object{
		kind:    traceKind,
		title:   title,
		message: sprint(msg...),
		body:    "",
		debug:   o.debug,
		isHttp:  false,
		w:       o.w,
	}
	o.printAndWrite()
}

func (o *object) Warning(title string, msg ...any) {
	*o = object{
		kind:    warningKind,
		title:   title,
		message: sprint(msg...),
		body:    "",
		debug:   o.debug,
		isHttp:  false,
		w:       o.w,
	}
	o.printAndWrite()
}

func (o *object) Json(title string, msg ...any) {
	*o = object{
		kind:    jsonKind,
		title:   title,
		message: sprint(msg...), //[]
		body:    "",
		debug:   o.debug,
		isHttp:  false,
		w:       o.w,
	}
	o.printAndWrite()
}

func (o *object) Success(title string, msg ...any) {
	*o = object{
		kind:    successKind,
		title:   title,
		message: sprint(msg...),
		body:    "",
		debug:   o.debug,
		isHttp:  false,
		w:       o.w,
	}
	o.printAndWrite()
}

const (
	inHttp  = ">>>--------------->>>--------------->>>--------------->>>\n"
	outHttp = "<<<---------------<<<---------------<<<---------------<<<\n"
	endHttp = "----------------------------------------------------------------------------------------------------------------------------------------\n"
)

func (o *object) Request(Request *http.Request, body bool) {
	*o = object{
		kind:    jsonKind,
		title:   "",
		message: o.DumpRequest(Request, body),
		body:    "",
		debug:   o.debug,
		isHttp:  true,
		w:       o.w,
	}
	o.printAndWrite()
}

func (o *object) DumpRequest(Request *http.Request, body bool) string {
	if Request == nil {
		return ""
	}
	dumpRequest, err := httputil.DumpRequest(Request, body)
	if err != nil {
		return o.DumpRequest(Request, false)
	}
	s := inHttp + Request.URL.String() + "\n"
	s += strings.TrimSuffix(string(dumpRequest), "\n")
	s += endHttp
	return s
}

func (o *object) Response(Response *http.Response, body bool) {
	*o = object{
		kind:    jsonKind,
		title:   "",
		message: o.DumpResponse(Response, body),
		body:    "",
		debug:   o.debug,
		isHttp:  true,
		w:       o.w,
	}
	o.printAndWrite()
}

func (o *object) DumpResponse(Response *http.Response, body bool) string {
	if Response == nil {
		return ""
	}
	dumpResponse, err := httputil.DumpResponse(Response, body)
	if err != nil {
		dumpResponse, err = httputil.DumpResponse(Response, false)
	}
	s := outHttp + strings.TrimSuffix(string(dumpResponse), "\n") + "\n"
	s += endHttp
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
		body = fmt.Sprintf("%#v", msg) // not export
	}
	*o = object{
		kind:    structKind,
		title:   "",
		message: body,
		body:    "",
		debug:   o.debug,
		isHttp:  false,
		w:       o.w,
	}
	o.printAndWrite()
}

func (o *object) MarshalJson(title string, msg any) {
	indent, err := json.MarshalIndent(msg, "", " ")
	if !o.Error(err) {
		return
	}
	o.Info(title, string(indent))
}

func (o *object) Reason() (reason string) {
	english2Chinese := map[string]string{
		"A certificate was explicitly revoked by its issuer.": "证书的颁发者明确吊销了证书。",
	}
	re := regexp.MustCompile(`(\.|\?|!)`)
	splitStr := re.ReplaceAllString(o.message, "$1\n")
	lines := strings.Split(splitStr, "\n")
	r := bytes.NewBuffer(nil)
	for _, line := range lines {
		trimSpace := strings.TrimSpace(line)
		r.WriteString(trimSpace)
		r.WriteString("\n")
		for english, chinese := range english2Chinese {
			if english == trimSpace {
				r.WriteString(chinese)
				r.WriteString("\n")
			}
		}
	}
	return r.String()
}

func (o *object) printAndWrite() {
	if isTermux() {
		return
	}
	indentTitle := o.getTimeNowString() + o.kind.String() + " " + o.textIndent(o.title, false)
	if o.message == "" {
		o.message = `""`
	}
	c := " //" + caller()
	switch o.kind {
	case jsonKind, hexDumpKind, structKind:
		indentTitle += c + "\n"
		o.body = indentTitle + o.message
	default:
		o.body = indentTitle + o.message + c
	}
	if o.isHttp {
		o.body = o.message
	}
	o.printColorBody()
	if !IsAndroid() {
		o.writeAppend(logFileName, o.body) // todo set apk path as log path
	}
}
