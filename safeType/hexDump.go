package safeType

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
)

type HexDumpString string

func NewHexDump(hexdumpStr HexDumpString) (data *Data) {
	hexdump := string(hexdumpStr)
	defer func() {
		s := New("")
		// s.WriteStringLn("data:=" + fmt.Sprintf("%#v", data))
		cut := `[]byte`
		cxx := fmt.Sprintf("%#v", data.Bytes())
		cxx = cxx[len(cut):]
		s.WriteString("char data[] = " + cxx + ";\n")
		mylog.Json("gen c++ code", s.String())
		mylog.HexDump("recovery go buffer", data.Bytes())
	}()
	hexdump = strings.TrimSuffix(hexdump, newLine)
	switch {
	case !hasAddress(hexdump) && !strings.Contains(hexdump, sep): // 没有地址和分隔符
		hexdump = strings.ReplaceAll(hexdump, " ", "")
		decodeString, err := hex.DecodeString(hexdump)
		if !mylog.Error(err) {
			return
		}
		data = New(decodeString)
	case strings.Contains(hexdump, sep): // go风格
		split := strings.Split(hexdump, newLine)
		noAddress := make([]string, len(split))
		hexString := new(bytes.Buffer)
		for i, s := range split {
			if s == "" {
				continue
			}
			noAddress[i] = s[addressLen:strings.Index(s, sep)]
			noAddress[i] = strings.ReplaceAll(noAddress[i], " ", "")
			hexString.WriteString(noAddress[i])
		}
		decodeString, err := hex.DecodeString(hexString.String())
		if !mylog.Error(err) {
			return
		}
		data = New(decodeString)
	default: // 非go风格
		split := strings.Split(hexdump, newLine)
		hexString := new(bytes.Buffer)
		for _, s := range split {
			if s == "" {
				continue
			}
			fields := strings.Split(s, " ")
			for j, field := range fields {
				if j > 0 && field == "" {
					fields = fields[1:j]
					break
				}
			}
			for _, field := range fields {
				hexString.WriteString(field)
			}
		}
		decodeString, err := hex.DecodeString(hexString.String())
		if !mylog.Error(err) {
			return
		}
		data = New(decodeString)
	}
	return
}

const (
	address    = "00000000  "
	sep        = "|"
	newLine    = "\n"
	addressLen = len(address)
)

func hasAddress(s string) bool {
	// 00000000  7e 15
	// 08A73200 57 61 72
	switch {
	case len(s) < len("00000000"):
		return false
	case strings.Contains(s, address):
		return true
	}
	return s[len("00000000")+1] == ' '
}
