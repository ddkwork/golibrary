package indent

import (
	"bytes"
	"golang.org/x/text/width"
	"strings"
)

func (o *object) Left(src string) string {
	return o.textIndent(src, true)
}

func (o *object) Right(src string) string {
	return o.textIndent(src, false)
}

func (o *object) textIndent(src string, isIndentLeft bool) string {
	//https://blog.csdn.net/raoxiaoya/article/details/108982887
	const hexDumpIndentLen = 28
	Separate := ` | `
	spaceLen := hexDumpIndentLen - o.width(src)
	spaceStr := ``
	if spaceLen > 0 {
		spaceStr = strings.Repeat(" ", spaceLen)
	}
	if isIndentLeft {
		return spaceStr + src + Separate
	}
	return src + spaceStr + Separate
	//title = fmt.Sprintf("%28s | ", title) //"%-28s 加个负号就是左对齐,英文状态下28的长度刚好与hexdump对齐，试了utf8没用的
}

func (o *object) width(s string) (w int) {
	for _, r := range []rune(s) {
		switch width.LookupRune(r).Kind() {
		case width.EastAsianWide, width.EastAsianFullwidth:
			w += 2
		default:
			w++
		}
	}
	return
}

func (o *object) NewLine(src *bytes.Buffer) error {
	return src.WriteByte('\n')
}
