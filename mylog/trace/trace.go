package main

import (
	"net/http"
	"os"
	"runtime/trace"
)

func main() {
	os.Setenv("GODEBUG", "traceallocfree=1")
	f, err := os.Create("./1.trace")
	if err != nil {
		panic(err)
	}
	trace.Start(f)
	defer trace.Stop()
	toTrace()
}
func toTrace() {
	resp, err := http.Get("https://www.baidu.com")
	_ = resp
	_ = err
}
