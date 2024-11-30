package mylog

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"

	"github.com/ddkwork/golibrary/mylog/pretty"
)

func (l *log) hexDump(title string, dump string) {
	*l = log{
		kind:     hexDumpKind,
		title:    title,
		message:  dump,
		body:     "",
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
		body:     "",
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
		body:     "",
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
		body:     "",
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
		body:     "",
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
		body:     "",
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
		body:     "",
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
		body:     "",
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
		body:     "",
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
	}
	s := outHttp + strings.TrimSuffix(string(dumpResponse), "\n") + "\n"
	s += endHttp
	return s
}

func (l *log) Struct(title string, msg any) {
	*l = log{
		callBack: l.callBack,
		kind:     structKind,
		title:    title,
		message:  pretty.Format(msg),
		body:     "",
		debug:    l.debug,
		isHttp:   false,
	}
	l.printAndWrite()
}

func (l *log) MarshalJson(title string, msg any) {
	indent := Check2(json.MarshalIndent(msg, "", " "))
	l.Json(title, string(indent))
}

func (l *log) Reason() (reason string) {
	english2Chinese := map[string]string{
		"A certificate was explicitly revoked by its issuer.": "证书的颁发者明确吊销了证书。",
	}
	re := regexp.MustCompile(`(\.|\?|!)`)
	splitStr := re.ReplaceAllString(l.message, "$1\n")
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
	return trimTrailingEmptyLines(r.String())
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
	l.body = fn.key() + " " + fn.value()
	if l.isHttp {
		l.body = l.message
	}
	l.printColorBody()
	l.body += "\n"
	if l.callBack != nil {
		l.callBack()
	}
	WriteAppend(LogPath(), l.body)
}
