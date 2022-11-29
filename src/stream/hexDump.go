package stream

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ddkwork/golibrary/mylog"
	"strings"
)

const (
	address    = "00000000  "
	sep        = "|"
	newLine    = "\n"
	addressLen = len(address)
)

func IsGoStyle(s string) {} //00000000  7e + |
func IsNotGoStyle(s string) {
	//没有分隔符，地址后面是一个空格，go风格的两个空格
}
func HasSep(s string) bool { return strings.Contains(s, sep) }

func HasAddress(s string) bool {
	//00000000  7e 15
	//08A73200 57 61 72
	switch {
	case len(s) < len("00000000"):
		return false
	case strings.Contains(s, address):
		return true
	case s[len("00000000")+1] == ' ':
		return true
	}
	return false
}

func NewHexDump(hexdump string) (buf []byte) {
	defer func() {
		s := New()
		//s.WriteStringLn("buf:=" + fmt.Sprintf("%#v", buf))
		cut := `[]byte`
		cxx := fmt.Sprintf("%#v", buf)
		cxx = cxx[len(cut):]
		s.WriteStringLn("char buf[] = " + cxx + ";")
		mylog.Json("gen c++ code", s.String())
		mylog.HexDump("recovery go buffer", buf)
	}()
	hexdump = strings.TrimSuffix(hexdump, newLine)
	switch {
	case !HasAddress(hexdump) && !strings.Contains(hexdump, sep): //没有地址和分隔符
		hexdump = strings.ReplaceAll(hexdump, " ", "")
		decodeString, err := hex.DecodeString(hexdump)
		if !mylog.Error(err) {
			return
		}
		buf = decodeString
		return
	case strings.Contains(hexdump, sep): //go风格
		split := strings.Split(hexdump, newLine)
		noAddres := make([]string, len(split))
		hexString := new(bytes.Buffer)
		for i, s := range split {
			if s == "" {
				continue
			}
			noAddres[i] = s[addressLen:strings.Index(s, sep)]
			noAddres[i] = strings.ReplaceAll(noAddres[i], " ", "")
			hexString.WriteString(noAddres[i])
		}
		decodeString, err := hex.DecodeString(hexString.String())
		if !mylog.Error(err) {
			return
		}
		buf = decodeString
		return
	case !strings.Contains(hexdump, sep): //非go风格
		panic("x64dbg copy")
	}
	return
}
