package mylog

import (
	"fmt"
	"go/format"
	"os"
	"strings"
	"sync"

	"github.com/google/go-cmp/cmp"

	"github.com/ddkwork/golibrary/stream/align"
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
