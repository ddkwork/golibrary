package file

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ddkwork/golibrary/mylog"
	mypath "github.com/ddkwork/golibrary/src/stream/tool/path"
	"github.com/hjson/hjson-go"
	"go/format"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type (
	Interface interface {
		ReadWriteTruncate(src, dst string) (ok bool)
		WriteTruncate(name string, data any) (ok bool)
		WriteAppend(name string, data any) (ok bool)
		WriteGoCode(name string, data any) (ok bool)
		WriteBinary(name string, data any) (ok bool)
		WriteJson(name string, Obj any) (ok bool)
		WriteHjson(name string, Obj any) (ok bool)
		ToLines(data any) (lines []string, ok bool)
		ReadToLines(path string) (lines []string, ok bool)
		GoCode() string
		Copy(source, destination string) (ok bool)
	}
	object struct{ goCode string }
)

func (o *object) ReadToLines(path string) (lines []string, ok bool) {
	file, err := os.ReadFile(path)
	if !mylog.Error(err) {
		return
	}
	return o.ToLines(file)
}

func New() Interface {
	return &object{
		goCode: "",
	}
}
func readFile(path string) (b []byte, ok bool) {
	file, err := os.ReadFile(path)
	if !mylog.Error(err) {
		return nil, false
	}
	return file, true
}
func (o *object) RaedToLines(path string) (lines []string, ok bool) {

	b, ok2 := readFile(path)
	if !ok2 {
		return
	}
	return o.ToLines(b)
}

func (o *object) Copy(source, destination string) (ok bool) {
	base := filepath.Base(source)
	return mylog.Error(filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		split := strings.Split(path, base)
		dst := filepath.Join(destination, base, split[1])
		switch {
		case info.IsDir():
			if !mypath.New().CreatDirectory(dst) {
				return err
			}
		default:
			if !o.ReadWriteTruncate(path, dst) {
				return err
			}
		}
		return err
	}))
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

func (o *object) ReadWriteTruncate(src, dst string) (ok bool) {
	b, ok2 := readFile(src)
	if !ok2 {
		return
	}
	return o.WriteTruncate(dst, b)
}
func (o *object) WriteTruncate(name string, data any) (ok bool) { return o.write(name, data, false) }
func (o *object) write(name string, data any, isAppend bool) (ok bool) {
	var (
		f   *os.File
		err error
	)
	if !mypath.New().CreatDirectory(filepath.Dir(name)) {
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
	if !mylog.Error2(f.Write(o.buffer(data).Bytes())) {
		return
	}
	if !mylog.Error2(f.WriteString("\n")) {
		return
	}
	return mylog.Error(f.Close())
}
func (o *object) WriteAppend(name string, data any) (ok bool) { return o.write(name, data, true) }
func (o *object) WriteGoCode(name string, data any) (ok bool) {
	b, err := format.Source(o.buffer(data).Bytes())
	if !mylog.Error(err) {
		return
	}
	return o.WriteTruncate(name, b)
}
func (o *object) WriteBinary(name string, data any) (ok bool) { return o.WriteTruncate(name, data) }
func (o *object) ToLines(data any) (lines []string, ok bool) {
	return toLines(o.buffer(data).String()), true
	newReader := bufio.NewReader(o.buffer(data))
	for {
		line, _, err := newReader.ReadLine()
		switch err {
		case io.EOF:
			return lines, true
		default:
			if !mylog.Error(err) {
				return
			}
		}
		lines = append(lines, string(line))
	}
}

func toLines(output string) []string {
	lines := strings.TrimSuffix(output, "\r\n")
	return strings.Split(lines, "\r\n")
}
func (o *object) WriteJson(name string, Obj any) (ok bool) {
	var oo any
	switch reflect.TypeOf(Obj).Kind() {
	case reflect.Struct:
		oo = &Obj
	case reflect.Ptr:
		oo = Obj
	}
	data, err := json.MarshalIndent(oo, " ", " ")
	if !mylog.Error(err) {
		return
	}
	return o.WriteTruncate(name, data)
}
func (o *object) WriteHjson(name string, Obj any) (ok bool) {
	data, err := hjson.MarshalWithOptions(Obj, hjson.EncoderOptions{
		Eol:            "",
		BracesSameLine: false,
		EmitRootBraces: false,
		QuoteAlways:    false,
		IndentBy:       " ",
		AllowMinusZero: false,
		UnknownAsNull:  false,
	})
	if !mylog.Error(err) {
		return
	}
	return o.WriteTruncate(name, data)
}
func (o *object) GoCode() string { return o.goCode }
