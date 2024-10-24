package mylog

import (
	"bytes"
	"fmt"
	"gioui.org/app"
	_ "gioui.org/app/permission/storage"
	"golang.org/x/exp/constraints"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type (
	object struct {
		kind     kind
		title    string
		message  string
		body     string
		debug    bool
		isHttp   bool
		callBack func()
	}
)

func (o *object) SetCallBack(callBack func()) {
	o.callBack = callBack
}

func SetCallBack(callBack func()) {
	defaultObject.SetCallBack(callBack)
}

func LogPath() (path string) {
	return filepath.Join(DataDir(), "log.log")
}

func DataDir() string {
	if IsAndroid() {
		dir, err := app.DataDir()
		if err != nil {
			panic(err)
		}
		return dir
	}
	if IsTermux() {
		return "/data/data/com.termux/files/usr" //todo choose another dir
	}
	return "."
}

func LogFileBody() string {
	return string(Check2(os.ReadFile(LogPath())))
}

func New() *object {
	return &object{
		kind:     0,
		title:    "",
		message:  "",
		body:     "",
		debug:    true,
		isHttp:   false,
		callBack: nil,
	}
}

var defaultObject = New()

func init() {
	if IsAndroid() || IsTermux() {
		return
	}
	CheckIgnore(os.Truncate(LogPath(), io.SeekStart))
	Trace("--------- title ---------", "------------------ info ------------------") //android not work,why?
	if IsWindows() {
		/*
				cmd/link: enable ASLR by default on Windows
			ASLR（Address Space Layout Randomization）是一种安全技术，用于防止缓冲区溢出攻击。它通过随机化内存中关键数据结构的位置，
			使得攻击者难以预测和利用这些位置来进行攻击。在Windows操作系统中，启用ASLR可以提高系统的安全性。

					// loadDll堆栈溢出
					// go build -buildmode=exe
					// go env -w GOFLAGS="-buildmode=exe"
					// https://github.com/golang/go/issues/42593
		*/
		cmd := exec.Command("go", "env", "-w", "GOFLAGS=-buildmode=exe")
		Check2(cmd.CombinedOutput())
	}
}

var (
	GithubWorkspace = os.Getenv("GITHUB_WORKSPACE")
	IsAction        = GithubWorkspace != ""
)

func ChdirToGithubWorkspace() {
	if IsAction {
		Check(os.Chdir(GithubWorkspace))
	}
	Info("GITHUB_WORKSPACE", Check2(os.Getwd()))
}

func Reason() (reason string)        { return defaultObject.Reason() }
func HexDump(title string, b []byte) { defaultObject.hexDump(title, b) }
func HexInteger[T constraints.Integer](msg T) string {
	return Hex("", msg) + "|" + fmt.Sprintf("%d", msg)
}

func Todo(body string) {
	Warning("TODO", body)
}

// todo key use cmp.ordering 支持自动格式化int，因为时候传入index作为key
func Hex[T constraints.Integer | []byte | *bytes.Buffer](title string, msg T) string {
	return defaultObject.Hex(title, msg)
}
func Info(title string, msg ...any)                   { defaultObject.Info(title, msg...) }
func Trace(title string, msg ...any)                  { defaultObject.Trace(title, msg...) }
func Warning(title string, msg ...any)                { defaultObject.Warning(title, msg...) }
func MarshalJson(title string, msg any)               { defaultObject.MarshalJson(title, msg) }
func Json(title string, msg ...any)                   { defaultObject.Json(title, msg...) }
func Success(title string, msg ...any)                { defaultObject.Success(title, msg...) }
func Struct(msg any)                                  { defaultObject.Struct(msg) }
func SetDebug(debug bool)                             { defaultObject.debug = debug }
func Request(Request *http.Request, body bool)        { defaultObject.Request(Request, body) }
func Response(Response *http.Response, body bool)     { defaultObject.Response(Response, body) }
func DumpRequest(req *http.Request, body bool) string { return defaultObject.DumpRequest(req, body) }
func DumpResponse(resp *http.Response, body bool) string {
	return defaultObject.DumpResponse(resp, body)
}

func Body() string    { return defaultObject.body }
func Message() string { return defaultObject.message }
