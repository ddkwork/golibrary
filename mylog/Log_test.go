package mylog_test

import (
	"errors"
	"os"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
)

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mylog.Call(xx)
	}
}

func TestLog(t *testing.T) {
	mylog.Call(func() {
		xx()
		safeAppLife()
	})
}

func safeAppLife() {
	mylog.Check(errors.New("this is a err message"))
	mylog.Check(errors.New("this is a err message"))
	mylog.Check2(os.Open(""))
	mylog.Check2(os.Stdout.Write(nil))
	mylog.Check2(os.Stdout.Write(nil))
}

func xx() {
	mylog.HexDump("buf", []byte{
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
	mylog.Hex("key", 0x888)
	mylog.Info("key", "tttttttttttttttttttttttt")
	mylog.Trace("key", "kkkkkkkkkkkkkkkkkkkk")
	mylog.Warning("key", "mmmmmmmmmmmmmmmmmmmmmmmmmmmkkkkkkkkkkkkkkkkkkkkkkkkkkk")
	mylog.Json("key", `{"manifestVersion":"1.1","engineVersion":"3.3.2180.8236","info":{"id":"VisualStudio/17.3.0+32804.467","buildBranch":"d17.3","buildVersion":"17.3.32804.467","localBuild":"build-lab","manifestName":"VisualStudio","manifestType":"installer","productDisplayVersion":"17.3.0","productLine":"Dev17","productLineVersion":"2022","productMilestone":"RTW","`)
	mylog.Success("key", "vgoTest pass")
	mylog.Struct(struct {
		A int
		B string
		C []byte
	}{
		A: 89,
		B: "jhjsbdd",
		C: []byte{0x11, 0x22, 0x33, 0x44},
	})
	mylog.Warning("key")
	mylog.Hex("lost last byte", 0x72B8)
	mylog.Hex("firstEnd xor 0x72B8", 0x72B8)
	mylog.HexDump("key", []byte{0x11, 0x22, 0x33, 0x44})
}
