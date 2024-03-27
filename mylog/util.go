package mylog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func (o *object) write(name string, data any) (ok bool) {
	if !o.creatDirectory(filepath.Dir(name)) {
		return
	}
	if !o.Error2(o.w.Write(o.buffer(data).Bytes())) {
		return
	}
	if !o.Error2(o.w.WriteString("\n")) {
		return
	}
	return true
	// return o.Error(o.w.Close())
}

var l = new(sync.RWMutex)

func (o *object) writeAppend(name string, data any) (ok bool) {
	// l.RLock()
	// defer l.RUnlock()
	return o.write(name, data)
}

func (o *object) creatDirectory(dir string) bool {
	fnMakeDir := func() bool { return o.Error(os.MkdirAll(dir, os.ModePerm)) }
	info, err := os.Stat(dir)
	switch {
	case os.IsExist(err):
		return true
	case os.IsNotExist(err):
		return fnMakeDir()
	case err == nil:
		return info.IsDir()
	default:
		return o.Error(err)
	}
}

func (o *object) buffer(data any) *bytes.Buffer {
	switch data.(type) {
	case string:
		return bytes.NewBufferString(data.(string))
	case []byte:
		return bytes.NewBuffer(data.([]byte))
	}
	return bytes.NewBufferString("error file data type " + fmt.Sprintf("%t", data))
}

func (o *object) getTimeNowString() string { return time.Now().Format("[2006-01-02 15:04:05]") }

func IsAndroid() bool { return runtime.GOOS == `android` }
func IsWindows() bool { return runtime.GOOS == `windows` }
func IsLinux() bool   { return runtime.GOOS == `linux` }

func sprint(msg_ ...any) string {
	msg := fmt.Sprint(msg_...) // 遇到[]byte会加上[],
	if msg == "" {             // 只输入了标题，没贴输入格式化内容
		return ""
	}
	switch {
	case msg[0] == '[' && msg[len(msg)-1] == ']':
		msg = msg[1 : len(msg)-1]
	}
	return msg
}

func isTermux() bool {
	_, err := os.Stat("/data/data/com.termux/files/usr")
	return err == nil
}
