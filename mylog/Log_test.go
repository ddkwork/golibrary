package mylog_test

import (
	"errors"
	"github.com/ddkwork/golibrary/mylog"
	"testing"
)

func TestLog(t *testing.T) {
	l := mylog.Default
	if !l.Error(errors.New("this is a err msg")) {

	}
	l.Hex("hex value", 0x888)
	l.HexDump("buf", []byte{
		0x11, 0x22, 0x33, 0x44,
		0x11, 0x22, 0x33, 0x44,
		0x11, 0x22, 0x33, 0x44,
		0x11, 0x22, 0x33, 0x44,
		0x11, 0x22, 0x33, 0x44,
		0x11, 0x22, 0x33, 0x44,
		0x11, 0x22, 0x33, 0x44,
		0x11, 0x22, 0x33, 0x44,
		0x11, 0x22, 0x33, 0x44,
		0x11, 0x22, 0x33, 0x44,
	})
	l.Struct(struct {
		A int
		B string
		C []byte
	}{
		A: 89,
		B: "jhjsbdd",
		C: []byte{0x11, 0x22, 0x33, 0x44},
	})
	l.Info("infomation", "tttttttttttttttttttttttt")
	l.Trace("trace", "kkkkkkkkkkkkkkkkkkkk")
	l.Warning("warnning", "mmmmmmmmm")
	l.Warning("warnning")
	l.Success("Success", "vgoTest pass")
	l.Hex("firstEnd xor 0x72B8,机器码丢弃一个字节", 0x72B8)
	l.Hex("firstEnd xor 0x72B8", 0x72B8)
	l.Json("warnning", `{"manifestVersion":"1.1","engineVersion":"3.3.2180.8236","info":{"id":"VisualStudio/17.3.0+32804.467","buildBranch":"d17.3","buildVersion":"17.3.32804.467","localBuild":"build-lab","manifestName":"VisualStudio","manifestType":"installer","productDisplayVersion":"17.3.0","productLine":"Dev17","productLineVersion":"2022","productMilestone":"RTW","
`)
}
