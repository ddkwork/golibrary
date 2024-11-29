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

	"gioui.org/app"
	_ "gioui.org/app/permission/storage"
	"golang.org/x/exp/constraints"
)

type (
	log struct {
		kind           kind
		title          string
		message        string
		body           string
		debug          bool
		isHttp         bool
		callBack       func()
		isShortHexdump bool
	}
)

func (l *log) SetCallBack(callBack func()) {
	l.callBack = callBack
}

func SetCallBack(callBack func()) {
	l.SetCallBack(callBack)
}

func LogPath() (path string) {
	return filepath.Join(DataDir(), "log.log")
}

func DataDir() string {
	if IsAndroid() {
		dir := Check2(app.DataDir())

		return dir
	}
	if IsTermux() {
		return "/data/data/com.termux/files/usr" // todo choose another dir
	}
	return "."
}

func LogFileBody() string {
	return string(Check2(os.ReadFile(LogPath())))
}

func New() *log {
	return &log{
		kind:     0,
		title:    "",
		message:  "",
		body:     "",
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
	CheckIgnore(os.Truncate(LogPath(), io.SeekStart))
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
	GitProxy(true)
	FormatAllFiles()
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

func Reason() (reason string) { return l.Reason() }
func HexDump[K keyType, V constraints.Unsigned | []byte | *bytes.Buffer](title K, buf V) {
	key := formatKey(title)
	switch v := any(buf).(type) {
	case []byte:
		l.hexDump(key, v)
	case *bytes.Buffer:
		l.hexDump(key, v.Bytes())
	default:
		panic("unsupported type")
	}
}

func Todo(body string) {
	Warning("TODO", body)
}

type keyType interface{ string | constraints.Integer }

func formatKey[K keyType](title K) (key string) {
	switch k := any(title).(type) {
	case string:
		key = k
	default:
		key = fmt.Sprintf("%d", k)
	}
	return key
}

func Hex[K keyType, V constraints.Unsigned](title K, v V) string {
	return l.Hex(formatKey(title), FormatInteger(v))
}
func Info[K keyType](title K, msg ...any)     { l.Info(formatKey(title), msg...) }
func Trace[K keyType](title K, msg ...any)    { l.Trace(formatKey(title), msg...) }
func Warning[K keyType](title K, msg ...any)  { l.Warning(formatKey(title), msg...) }
func MarshalJson[K keyType](title K, msg any) { l.MarshalJson(formatKey(title), msg) }
func Json[K keyType](title K, msg ...any)     { l.Json(formatKey(title), msg...) }
func Success[K keyType](title K, msg ...any)  { l.Success(formatKey(title), msg...) }
func Struct[K keyType](title K, msg any) {
	switch t := any(title).(type) {
	case string:
		if t == "" {
			l.Struct(reflect.TypeOf(msg).Name(), msg)
			return
		}
	}
	l.Struct(formatKey(title), msg)
}
func SetDebug(debug bool)                             { l.debug = debug }
func Request(Request *http.Request, body bool)        { l.Request(Request, body) }
func Response(Response *http.Response, body bool)     { l.Response(Response, body) }
func DumpRequest(req *http.Request, body bool) string { return l.DumpRequest(req, body) }
func DumpResponse(resp *http.Response, body bool) string {
	return l.DumpResponse(resp, body)
}

func Body() string    { return l.body }
func Message() string { return l.message }
