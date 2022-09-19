package mylog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func (o *object) write(name string, data any, isAppend bool) (ok bool) {
	var (
		f   *os.File
		err error
	)
	if !o.CreatDirectory(filepath.Dir(name)) {
		return
	}
	switch isAppend {
	case true:
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	case false:
		f, err = os.Create(name) //== os.Truncate(name, 0)
	}
	if err != nil {
		if os.IsNotExist(err) {
			return o.write(name, data, isAppend)
		}
	}
	if !o.Error2(f.Write(o.buffer(data).Bytes())) {
		return
	}
	if !o.Error2(f.WriteString("\n")) {
		return
	}
	return o.Error(f.Close())
}
func (o *object) WriteAppend(name string, data any) (ok bool) { return o.write(name, data, true) }
func (o *object) CreatDirectory(dir string) bool {
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
func (o *object) buffer(data any) *bytes.Buffer { //todo replaced as stream pkg
	switch data.(type) {
	case string:
		return bytes.NewBufferString(data.(string))
	case []byte:
		return bytes.NewBuffer(data.([]byte))
	}
	return bytes.NewBufferString("error file data type " + fmt.Sprintf("%t", data))
}
func (o *object) GetTimeNowString() string { return time.Now().Format("2006-01-02 15:04:05 ") }

const (
	android = `android`
	linux   = `linux`
	windows = `windows`
)

func IsAndroid() bool { return runtime.GOOS == android }
func IsWindows() bool { return runtime.GOOS == windows }
func IsLinux() bool   { return runtime.GOOS == linux }
