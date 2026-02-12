package mylog

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"

	"github.com/ddkwork/golibrary/std/mylog/pretty"
)

func (l *log) hexDump(title string, dump string) {
	*l = log{
		kind: hexDumpKind,
		row: keyValue{
			key:   title,
			value: dump,
		},
		debug:    l.debug,
		callBack: l.callBack,
	}
	l.printAndWrite()
}

func (l *log) hex(title, msg string) string {
	*l = log{
		callBack: l.callBack,
		kind:     hexKind,
		row: keyValue{
			key:   title,
			value: msg,
		},
		debug: l.debug,
	}
	l.printAndWrite()
	return l.row.Value()
}

func (l *log) Info(title string, msg ...any) {
	*l = log{
		callBack: l.callBack,
		kind:     infoKind,
		row: keyValue{
			key:   title,
			value: sprint(msg...),
		},
		debug: l.debug,
	}
	l.printAndWrite()
}

func (l *log) Trace(title string, msg ...any) {
	*l = log{
		callBack: l.callBack,
		kind:     traceKind,
		row: keyValue{
			key:   title,
			value: sprint(msg...),
		},
		debug: l.debug,
	}
	l.printAndWrite()
}

func (l *log) Warning(title string, msg ...any) {
	*l = log{
		callBack: l.callBack,
		kind:     warningKind,
		row: keyValue{
			key:   title,
			value: sprint(msg...),
		},
		debug: l.debug,
	}
	l.printAndWrite()
}

func (l *log) Json(title string, msg ...any) {
	*l = log{
		callBack: l.callBack,
		kind:     jsonKind,
		row: keyValue{
			key:   title,
			value: sprint(msg...),
		},
		debug: l.debug,
	}
	l.printAndWrite()
}

func (l *log) Success(title string, msg ...any) {
	*l = log{
		callBack: l.callBack,
		kind:     successKind,
		row: keyValue{
			key:   title,
			value: sprint(msg...),
		},
		debug: l.debug,
	}
	l.printAndWrite()
}

const (
	inHttp  = ">>>-----------------------request--------------------------->>>\n"
	outHttp = "<<<----------------------response---------------------------<<<\n"
	endHttp = ""
	// endHttp = "----------------------------------------------------------------------------------------------------------------------------------------\n"
)

func (l *log) Request(Request *http.Request, body bool) {
	*l = log{
		callBack: l.callBack,
		kind:     jsonKind,
		row: keyValue{
			key:   "",
			value: l.DumpRequest(Request, body),
		},
		debug: l.debug,
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
	b := strings.TrimSuffix(string(dumpRequest), "\n")
	if Request.Header.Get("Content-Type") == "application/json" {
		buf := strings.Builder{}
		for s := range strings.Lines(b) {
			if strings.TrimSpace(s) == "" {
				continue
			}
			if strings.HasPrefix(s, "{") {
				s = JsonIndent([]byte(s))
				s += "\n"
			}
			buf.WriteString(s)
		}
		b = buf.String()
	}
	s := inHttp + Request.URL.String() + "\n"
	s += b
	s += endHttp
	return s
}

func (l *log) Response(Response *http.Response, body bool) {
	*l = log{
		callBack: l.callBack,
		kind:     jsonKind,
		row: keyValue{
			key:   "",
			value: l.DumpResponse(Response, body),
		},
		debug: l.debug,
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
	b := strings.TrimSuffix(string(dumpResponse), "\n")
	if Response.Header.Get("Content-Type") == "application/json" {
		buf := strings.Builder{}
		for s := range strings.Lines(b) {
			if strings.TrimSpace(s) == "" {
				continue
			}
			if strings.HasPrefix(s, "{") {
				s = JsonIndent([]byte(s))
			}
			buf.WriteString(s)
		}
		b = buf.String()
	}
	s := outHttp + b
	s += "\n" + endHttp
	return s
}

func (l *log) Struct(title string, msg any) {
	*l = log{
		callBack: l.callBack,
		kind:     structKind,
		row: keyValue{
			key:   title,
			value: pretty.Format(msg),
		},
		debug: l.debug,
	}
	l.printAndWrite()
}

func (l *log) MarshalJson(title string, msg any) {
	indent := Check2(json.MarshalIndent(msg, "", " "))
	l.Json(title, string(indent))
}

// 不带行号的日志方法

func (l *log) hexDumpNoCaller(title string, dump string) {
	*l = log{
		kind:       hexDumpKind,
		row:        keyValue{key: title, value: dump},
		debug:      l.debug,
		callBack:   l.callBack,
		showCaller: false,
	}
	l.printAndWrite()
}

func (l *log) hexNoCaller(title, msg string) string {
	*l = log{
		callBack:   l.callBack,
		kind:       hexKind,
		row:        keyValue{key: title, value: msg},
		debug:      l.debug,
		showCaller: false,
	}
	l.printAndWrite()
	return l.row.Value()
}

func (l *log) InfoNoCaller(title string, msg ...any) {
	*l = log{
		callBack:   l.callBack,
		kind:       infoKind,
		row:        keyValue{key: title, value: sprint(msg...)},
		debug:      l.debug,
		showCaller: false,
	}
	l.printAndWrite()
}

func (l *log) TraceNoCaller(title string, msg ...any) {
	*l = log{
		callBack:   l.callBack,
		kind:       traceKind,
		row:        keyValue{key: title, value: sprint(msg...)},
		debug:      l.debug,
		showCaller: false,
	}
	l.printAndWrite()
}

func (l *log) WarningNoCaller(title string, msg ...any) {
	*l = log{
		callBack:   l.callBack,
		kind:       warningKind,
		row:        keyValue{key: title, value: sprint(msg...)},
		debug:      l.debug,
		showCaller: false,
	}
	l.printAndWrite()
}

func (l *log) JsonNoCaller(title string, msg ...any) {
	*l = log{
		callBack:   l.callBack,
		kind:       jsonKind,
		row:        keyValue{key: title, value: sprint(msg...)},
		debug:      l.debug,
		showCaller: false,
	}
	l.printAndWrite()
}

func (l *log) SuccessNoCaller(title string, msg ...any) {
	*l = log{
		callBack:   l.callBack,
		kind:       successKind,
		row:        keyValue{key: title, value: sprint(msg...)},
		debug:      l.debug,
		showCaller: false,
	}
	l.printAndWrite()
}

func (l *log) RequestNoCaller(Request *http.Request, body bool) {
	*l = log{
		callBack:   l.callBack,
		kind:       jsonKind,
		row:        keyValue{key: "", value: l.DumpRequest(Request, body)},
		debug:      l.debug,
		showCaller: false,
	}
	l.printAndWrite()
}

func (l *log) ResponseNoCaller(Response *http.Response, body bool) {
	*l = log{
		callBack:   l.callBack,
		kind:       jsonKind,
		row:        keyValue{key: "", value: l.DumpResponse(Response, body)},
		debug:      l.debug,
		showCaller: false,
	}
	l.printAndWrite()
}

func (l *log) StructNoCaller(title string, msg any) {
	*l = log{
		callBack:   l.callBack,
		kind:       structKind,
		row:        keyValue{key: title, value: pretty.Format(msg)},
		debug:      l.debug,
		showCaller: false,
	}
	l.printAndWrite()
}

func (l *log) MarshalJsonNoCaller(title string, msg any) {
	indent := Check2(json.MarshalIndent(msg, "", " "))
	l.JsonNoCaller(title, string(indent))
}

var lock sync.RWMutex

func (l *log) printAndWrite() {
	lock.Lock()
	defer lock.Unlock()

	v := l.row.Value()
	var end string
	if l.showCaller {
		end = " //" + caller()
	}
	switch l.kind {
	case hexDumpKind:
		isLongHexdump := strings.Contains(l.row.value, "\n")
		if isLongHexdump {
			v = "\n" + v
			if l.showCaller {
				end = "\n" + end
			}
		}
		v += end
	case jsonKind, structKind:
		v = "\n" + v
		if l.showCaller {
			end = "\n" + end
		}
		v += end
	default:
		v += end
	}
	l.row = keyValue{
		key:   GetTimeNowString() + l.kind.String() + l.textIndent(l.row.key, false),
		value: v,
	}
	s := l.row.key + separate + l.row.value
	s += "\n"
	l.printColorBody(s)
	if l.callBack != nil {
		l.callBack(s)
	}
	if IsAndroid() {
		println("android log is not support yet")
		return
	}
	WriteAppend(logPath(), s)
}
