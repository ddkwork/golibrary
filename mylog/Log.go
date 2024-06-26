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
		message: sprint(msg...),
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
	dumpRequest, e := httputil.DumpRequest(Request, body)
	if e != nil {
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
	dumpResponse, e := httputil.DumpResponse(Response, body)
	if e != nil {
		dumpResponse, e = httputil.DumpResponse(Response, false)
	}
	s := outHttp + strings.TrimSuffix(string(dumpResponse), "\n") + "\n"
	s += endHttp
	return s
}

func (o *object) Struct(msg any) {
	msg = reflect.Indirect(reflect.ValueOf(msg)).Interface()
	marshalIndent := Check2(json.MarshalIndent(msg, "", " "))
	body := string(marshalIndent)
	if body == "{}" || reflect.TypeOf(msg).Kind() == reflect.Slice {
		body = fmt.Sprintf("%#v", msg)
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
	indent := Check2(json.MarshalIndent(msg, "", " "))
	o.Json(title, string(indent))
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
	if IsTermux() {
		return
	}
	indentTitle := GetTimeNowString() + o.kind.String() + " " + o.textIndent(o.title, false)
	o.message = strings.TrimSuffix(o.message, "\n")
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
		o.body += "\n"
		WriteAppend(logFileName, o.body)
	}
}
