package mylog

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"

	"github.com/ddkwork/golibrary/mylog/pretty"
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
	inHttp  = ">>>--------------->>>--------------->>>--------------->>>\n"
	outHttp = "<<<---------------<<<---------------<<<---------------<<<\n"
	endHttp = "----------------------------------------------------------------------------------------------------------------------------------------\n"
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
	s := inHttp + Request.URL.String() + "\n"
	s += strings.TrimSuffix(string(dumpRequest), "\n")
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
	s := outHttp + strings.TrimSuffix(string(dumpResponse), "\n")
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

var lock sync.Mutex

func (l *log) printAndWrite() {
	lock.Lock()
	defer lock.Unlock()

	v := l.row.Value()
	end := " //" + caller()
	switch l.kind {
	case hexDumpKind:
		isLongHexdump := strings.Contains(l.row.value, "\n")
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
