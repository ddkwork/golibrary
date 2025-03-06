package mylog

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/ddkwork/golibrary/mylog/pretty"
)

func (l *log) hexDump(title string, dump string) {
	*l = log{
		kind:     hexDumpKind,
		title:    title,
		message:  dump,
		row:      "",
		debug:    l.debug,
		isHttp:   false,
		callBack: l.callBack,
	}
	l.printAndWrite()
}

func (l *log) hex(title, msg string) string {
	*l = log{
		callBack: l.callBack,
		kind:     hexKind,
		title:    title,
		message:  msg,
		row:      "",
		debug:    l.debug,
		isHttp:   false,
	}
	l.printAndWrite()
	return l.message
}

func (l *log) Info(title string, msg ...any) {
	*l = log{
		callBack: l.callBack,
		kind:     infoKind,
		title:    title,
		message:  sprint(msg...),
		row:      "",
		debug:    l.debug,
		isHttp:   false,
	}
	l.printAndWrite()
}

func (l *log) Trace(title string, msg ...any) {
	*l = log{
		callBack: l.callBack,
		kind:     traceKind,
		title:    title,
		message:  sprint(msg...),
		row:      "",
		debug:    l.debug,
		isHttp:   false,
	}
	l.printAndWrite()
}

func (l *log) Warning(title string, msg ...any) {
	*l = log{
		callBack: l.callBack,
		kind:     warningKind,
		title:    title,
		message:  sprint(msg...),
		row:      "",
		debug:    l.debug,
		isHttp:   false,
	}
	l.printAndWrite()
}

func (l *log) Json(title string, msg ...any) {
	*l = log{
		callBack: l.callBack,
		kind:     jsonKind,
		title:    title,
		message:  sprint(msg...),
		row:      "",
		debug:    l.debug,
		isHttp:   false,
	}
	l.printAndWrite()
}

func (l *log) Success(title string, msg ...any) {
	*l = log{
		callBack: l.callBack,
		kind:     successKind,
		title:    title,
		message:  sprint(msg...),
		row:      "",
		debug:    l.debug,
		isHttp:   false,
	}
	l.printAndWrite()
}

const (
	inHttp  = ">>>--------------->>>--------------->>>--------------->>>\n"
	outHttp = "<<<---------------<<<---------------<<<---------------<<<\n"
	endHttp = "----------------------------------------------------------------------------------------------------------------------------------------\n"
)

func (l *log) Request(Request *http.Request, body bool) {
	*l = log{
		callBack: l.callBack,
		kind:     jsonKind,
		title:    "",
		message:  l.DumpRequest(Request, body),
		row:      "",
		debug:    l.debug,
		isHttp:   true,
	}
	l.printAndWrite()
}

func (l *log) DumpRequest(Request *http.Request, body bool) string {
	if Request == nil {
		return ""
	}
	dumpRequest, e := httputil.DumpRequest(Request, body)
	if e != nil {
		return l.DumpRequest(Request, false)
	}
	s := inHttp + Request.URL.String() + "\n"
	s += strings.TrimSuffix(string(dumpRequest), "\n")
	s += endHttp
	return s
}

func (l *log) Response(Response *http.Response, body bool) {
	*l = log{
		callBack: l.callBack,
		kind:     jsonKind,
		title:    "",
		message:  l.DumpResponse(Response, body),
		row:      "",
		debug:    l.debug,
		isHttp:   true,
	}
	l.printAndWrite()
}

func (l *log) DumpResponse(Response *http.Response, body bool) string {
	if Response == nil {
		return ""
	}
	dumpResponse, e := httputil.DumpResponse(Response, body)
	if e != nil {
		dumpResponse, e = httputil.DumpResponse(Response, false)
		if e != nil {
			panic(e)
		}
	}
	s := outHttp + strings.TrimSuffix(string(dumpResponse), "\n")
	s += "\n" + endHttp
	return s
}

func (l *log) Struct(title string, msg any) {
	*l = log{
		callBack: l.callBack,
		kind:     structKind,
		title:    title,
		message:  pretty.Format(msg),
		row:      "",
		debug:    l.debug,
		isHttp:   false,
	}
	l.printAndWrite()
}

func (l *log) MarshalJson(title string, msg any) {
	indent := Check2(json.MarshalIndent(msg, "", " "))
	l.Json(title, string(indent))
}

type keyValue struct {
	key   func() string
	value func() string // 重点是处理换行逻辑
}

func (l *log) printAndWrite() {
	fn := keyValue{
		key: func() string {
			return GetTimeNowString() + l.kind.String() + " " + l.textIndent(l.title, false)
		},
		value: func() string {
			v := strings.TrimSuffix(l.message, "\n")
			end := " //" + caller()
			switch l.kind {
			case hexDumpKind:
				isLongHexdump := strings.Contains(l.message, "\n")
				if isLongHexdump {
					v = "\n" + v
					end = "\n" + end
				}
				v += end
			case jsonKind, structKind:
				v = "\n" + v
				end = "\n" + end
				v += end
			default:
				v += end
			}
			return v
		},
	}
	l.row = fn.key() + " " + fn.value()
	if l.isHttp {
		l.row = l.message
	}
	l.printColorBody()
	l.row += "\n"
	if l.callBack != nil {
		l.callBack()
	}
	if IsAndroid() {
		println("android log is not support yet")
		return
	}
	WriteAppend(logPath(), l.row)
}
