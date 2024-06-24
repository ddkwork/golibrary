package mylog

import (
	"fmt"
	"runtime"
	"strings"
)

func Call(f func()) {
	callWithHandler(f, func(err error) { defaultObject.errorCall(err) })
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

func (o *object) errorCall(err any) bool {
	if err == nil {
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
	o.printAndWrite2()
	return false
}

func (o *object) printAndWrite2() {
	if IsTermux() {
		return
	}
	indentTitle := GetTimeNowString() + o.kind.String() + " " + o.textIndent(o.title, false)
	if o.message == "" {
		o.message = `""`
	}
	o.message = strings.TrimSuffix(o.message, "\n")
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
	builder.WriteString(indentTitle + o.message)
	fnNewLine()
	const hexDumpIndentLen = 26 + 6
	for _, s := range stack {
		indent := o.textIndent("", false)
		indent = strings.Repeat(" ", hexDumpIndentLen) + indent

		builder.WriteString(indent)
		builder.WriteString(s)
	}
	o.message = builder.String()
	o.body = o.message
	o.printColorBody()
	if !IsAndroid() {
		o.body += "\n"
		WriteAppend(logFileName, o.body)
	}
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
