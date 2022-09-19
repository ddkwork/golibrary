package myc2go

import (
	"github.com/ddkwork/golibrary/src/cpp2go/delete/myc2go/cppBlock"
	"github.com/ddkwork/golibrary/src/stream"
	"github.com/ddkwork/golibrary/src/stream/tool"
	"os"
	"strconv"
	"strings"
)

type (
	tmpInterface interface {
		removeComment(lines []string)
		read() []string
		write(lines []string)
		name() string
		finished()
		clean(b cppBlock.Lines)
		flush()
	}
	tmpObject struct{}
)

func newTmpObject() tmpInterface {
	return &tmpObject{}
}

func (t *tmpObject) clean(b cppBlock.Lines) {
	lines := t.read()
	for _, info := range b {
		for i2 := range lines {
			if info.Col == i2+1 {
				lines[i2] = ""
			}
		}
	}
	t.write(lines)
}

func (t *tmpObject) flush() { Check(os.Remove(t.name())) }
func (t *tmpObject) finished() {
	lines := t.read()
	for _, line := range lines {
		if line != "" {
			panic(t.name() + "not finished")
		}
	}
}
func (t *tmpObject) name() string { return "./tmp.bak" }
func (t *tmpObject) removeComment(lines []string) {
	for i, line := range lines {
		switch {
		case line == "#pragma once":
			lines[i] = ""
		case strings.Contains(line, `/*`) && strings.Contains(line, `*/`):
			start := strings.Index(line, `/*`)
			end := strings.LastIndex(line, `*/`)
			c := line[start-1 : end+2] //?
			dd := strconv.Quote(line)
			dd = dd //8-18 ? "        /* flags  */ 0, \\"
			line = strings.ReplaceAll(line, c, "")
			lines[i] = line //        /* flags  */ 0, \
		case strings.Contains(line, "//"):
			before, _, found := strings.Cut(line, "//")
			if !found {
				panic("// not found")
			}
			lines[i] = before
		}
	}
	t.write(lines)
}
func (t *tmpObject) read() []string {
	file, err := os.ReadFile(t.name())
	Check(err)
	lines, ok := tool.File().ToLines(file)
	if !ok {
		panic(t.name() + " ToLines")
	}
	return lines
}
func (t *tmpObject) write(lines []string) {
	if !tool.File().WriteTruncate(t.name(), stream.New().LinesToString(lines)) {
		panic(t.name())
	}
}
