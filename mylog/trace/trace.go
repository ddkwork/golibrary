package main

import (
	"net/http"
	"os"
	"runtime/trace"

	"github.com/ddkwork/golibrary/mylog"
)

func main() {
	os.Setenv("GODEBUG", "traceallocfree=1")
	f := mylog.Check2(os.Create("./1.trace"))

	trace.Start(f)
	defer trace.Stop()
	toTrace()
}

func toTrace() {
	resp := mylog.Check2(http.Get("https://www.baidu.com"))
	_ = resp
}
