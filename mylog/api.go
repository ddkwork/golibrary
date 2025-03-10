package mylog

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
)

type (
	log struct {
		kind     kind   // type of log
		title    string // key
		message  string // value
		row      string // merge title and message,means key-value
		debug    bool   // display print and write log
		isHttp   bool   // todo
		callBack func() // for ux logView
	}
)

func (l *log) SetCallBack(callBack func()) {
	l.callBack = callBack
}

func SetCallBack(callBack func()) {
	l.SetCallBack(callBack)
}

func logPath() (path string) {
	return filepath.Join(dataDir(), "log.log")
}

func dataDir() string {
	switch {
	// case IsAndroid():
	// return Check2(app.dataDir())
	case IsTermux():
		return "/data/data/com.termux/files/usr" // todo choose another dir
	default: // windows,linux
		return "."
	}
}

func New() *log {
	return &log{
		kind:     0,
		title:    "",
		message:  "",
		row:      "",
		debug:    true,
		isHttp:   false,
		callBack: nil,
	}
}

var l = New()

func init() {
	if IsAndroid() || IsTermux() {
		return
	}
	if FileLineCountIsMoreThan(logPath(), 2) {
		CheckIgnore(os.Truncate(logPath(), io.SeekStart))
	}
	Trace("--------- title ---------", "------------------ info ------------------") // android not work,why?
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
		Check2(exec.Command("go", "env", "-w", "GOFLAGS=-buildmode=exe").CombinedOutput())
	}
	Check2(exec.Command("go", "env", "-w", "GOPROXY=https://goproxy.cn").CombinedOutput())
	// GitProxy(true)
	// FormatAllFiles()
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

func HexDump[K keyType, V []byte | *bytes.Buffer](title K, buf V) {
	l.hexDump(formatKey(title), DumpHex(buf))
}

func Todo(body string) {
	Warning("TODO", body)
}

type keyType interface{ string | Integer }

func formatKey[K keyType](title K) (key string) {
	switch k := any(title).(type) {
	case string:
		key = k
	default:
		key = fmt.Sprintf("%d", k)
	}
	return key
}

func Hex[K keyType, V Unsigned](title K, v V) string {
	return l.hex(formatKey(title), FormatInteger(v))
}
func Info[K keyType](title K, msg ...any)     { l.Info(formatKey(title), msg...) }
func Trace[K keyType](title K, msg ...any)    { l.Trace(formatKey(title), msg...) }
func Warning[K keyType](title K, msg ...any)  { l.Warning(formatKey(title), msg...) }
func MarshalJson[K keyType](title K, msg any) { l.MarshalJson(formatKey(title), msg) }
func Json[K keyType](title K, msg ...any)     { l.Json(formatKey(title), msg...) }
func Success[K keyType](title K, msg ...any)  { l.Success(formatKey(title), msg...) }
func Struct[K keyType](title K, object any) {
	if reflect.Indirect(reflect.ValueOf(object)).Kind() != reflect.Struct {
		// panic("object must be a struct")//打印单一类型的切片还是必要的
	}
	switch t := any(title).(type) {
	case string:
		if t == "" {
			l.Struct(reflect.TypeOf(object).Name(), object)
			return
		}
	}
	l.Struct(formatKey(title), object)
}
func SetDebug(debug bool)                             { l.debug = debug }
func Request(Request *http.Request, body bool)        { l.Request(Request, body) }
func Response(Response *http.Response, body bool)     { l.Response(Response, body) }
func DumpRequest(req *http.Request, body bool) string { return l.DumpRequest(req, body) }
func DumpResponse(resp *http.Response, body bool) string {
	return l.DumpResponse(resp, body)
}

func Row() string { return l.row }
