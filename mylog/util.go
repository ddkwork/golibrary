package mylog

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/ddkwork/golibrary/stream/align"
	"mvdan.cc/gofumpt/format"
)

func sprint(msg_ ...any) string {
	msg := fmt.Sprint(msg_...)
	if msg == "" {
		return `""`
	}
	switch {
	case msg[0] == '[' && msg[len(msg)-1] == ']':
		msg = msg[1 : len(msg)-1]
	}
	return msg
}

func WriteGoFileWithDiff[T []byte](path string, data T) {
	source, e := format.Source(data, format.Options{})
	CheckIgnore(e)
	if e != nil {
		write(path, false, data)
		return
	}
	write(path, false, source)
	diff := Diff(path, Check2(os.ReadFile(path)), path, source)
	if diff != nil {
		println(string(diff))
	}
}

func (l *log) textIndent(src string, isLeftAlign bool) string {
	const hexDumpIndentLen = 26
	Separate := ` │ `
	spaceLen := hexDumpIndentLen - align.StringWidth[int](src)
	spaceStr := ``
	if spaceLen > 0 {
		spaceStr = strings.Repeat(" ", spaceLen)
	}
	if isLeftAlign {
		return src + spaceStr + Separate
	}
	return spaceStr + src + Separate
}

var (
	count int64
	lock  = new(sync.Mutex)
)

func main() {
	ch := make(chan struct{}, 2)
	go func() {
		for i := 0; i < 100000; i++ {
			lock.Lock()
			count++
			lock.Unlock()
		}
		ch <- struct{}{}
	}()

	go func() {
		for i := 0; i < 100000; i++ {
			lock.Lock()
			count--
			lock.Unlock()
		}
		ch <- struct{}{}
	}()

	<-ch
	<-ch
	close(ch)
	fmt.Println(count)
}
