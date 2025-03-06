package mylog

import (
	"fmt"
	"runtime"
	"strings"
)

func Call(f func()) {
	callWithHandler(f, func(err error) { l.errorCall(err) })
}

func callWithHandler(f func(), errHandler func(err error)) {
	defer recovery(errHandler)
	f()
}

type recoveryHandler func(error)

func recovery(handler recoveryHandler) {
	if recovered := recover(); recovered != nil && handler != nil {
		e, ok := recovered.(error)
		if !ok {
			e = fmt.Errorf("%+v", recovered)
		}
		defer recovery(nil)
		handler(e)
	}
}

func (l *log) errorCall(err any) bool {
	if err == nil {
		return true
	}
	*l = log{
		callBack: l.callBack,
		kind:     errorKind,
		title:    "",
		message:  reason(err),
		row:      "",
		debug:    l.debug,
		isHttp:   false,
	}
	l.printAndWrite2()
	return false
}

func (l *log) printAndWrite2() {
	if IsTermux() {
		return
	}
	indentTitle := GetTimeNowString() + l.kind.String() + " " + l.textIndent(l.title, false)
	if l.message == "" {
		l.message = `""`
	}
	l.message = strings.TrimSuffix(l.message, "\n")
	stack := make([]string, 0)
	frames := runtime.CallersFrames(callStack())
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		skip := false
		for _, s := range RuntimePrefixesToFilter {
			if strings.HasPrefix(frame.Function, s) || strings.HasSuffix(frame.Function, ".func1") {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		s := fmt.Sprintf("%s+0x%x %s:%d", frame.Function, frame.PC-frame.Entry, frame.File, frame.Line)
		stack = append(stack, s+"\n")
	}

	builder := strings.Builder{}
	fnNewLine := func() {
		builder.WriteString("\n")
	}
	builder.WriteString(indentTitle + l.message)
	fnNewLine()
	const hexDumpIndentLen = 26 + 6
	for _, s := range stack {
		indent := l.textIndent("", false)
		indent = strings.Repeat(" ", hexDumpIndentLen) + indent

		builder.WriteString(indent)
		builder.WriteString(s)
	}
	l.message = builder.String()
	l.row = l.message
	l.printColorBody()
	l.row += "\n"
	if l.callBack != nil {
		l.callBack()
	}
	WriteAppend(logPath(), l.row)
}

var RuntimePrefixesToFilter = []string{
	"runtime.",
	"testing.",
	"github.com/ddkwork/golibrary/mylog.callWithHandler",
	"github.com/ddkwork/golibrary/mylog.Call",
}

func callStack() []uintptr {
	var pcs [512]uintptr
	n := runtime.Callers(6, pcs[:])
	cs := make([]uintptr, n)
	copy(cs, pcs[:n])
	return cs
}
