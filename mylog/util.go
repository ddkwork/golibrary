package mylog

import (
	"fmt"
	"go/format"
	"os"

	"github.com/google/go-cmp/cmp"
)

func sprint(msg ...any) string {
	data := fmt.Sprint(msg...)
	if data == "" {
		return `""`
	}
	switch {
	case data[0] == '[' && data[len(data)-1] == ']':
		data = data[1 : len(data)-1]
	}
	return data
}

func WriteGoFileWithDiff[T []byte](path string, data T) {
	source, e := format.Source(data)
	CheckIgnore(e)
	if e != nil {
		write(path, false, data)
		return
	}
	write(path, false, source)
	diff := cmp.Diff(Check2(os.ReadFile(path)), source)
	if diff != "" {
		println(diff)
	}
}
