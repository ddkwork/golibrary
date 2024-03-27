package mylog_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/ddkwork/golibrary/mylog"
)

func TestTrue(t *testing.T) {
	// mylog.True(false)
	mylog.True(true)
}

// todo 似乎在携程中获取不到行号？但为什么单元测试又可以了
func TestGoRunTine(t *testing.T) {
	go d1()
	go d2()
	time.Sleep(time.Second)
}

func d1() {
	mylog.Warning("xx")
}

func d2() {
	mylog.Warning("oo")
}

func BenchmarkName(b *testing.B) {
	// mylog.SetDebug(false)
	for i := 0; i < b.N; i++ {
		if !mylog.Error(errors.New("this is a err message")) {
		}
		if !mylog.Error2(nil, errors.New("this is a err message")) {
		}
		if !mylog.Error2(os.Open("")) {
		}
		if !mylog.Error2(os.Stdout.Write(nil)) {
		}
		_, err := os.Stdout.Write(nil)
		if !mylog.Error(err) {
		}
		// r1, r2, errno := syscall.SyscallN(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

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
		mylog.Warning("key", "mmmmmmmmm")
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
		mylog.Hex("lost mid last byte", 0x72B8)
		mylog.Hex("firstEnd xor 0x72B8", 0x72B8)
		mylog.HexDump("key", []byte{0x11, 0x22, 0x33, 0x44})
	}
}

func BenchmarkSlog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slog.Error("err")
	}
}

func TestLog(t *testing.T) {
	// 错误处理的代码可读性设计
	if !mylog.Error(errors.New("this is a err message")) {
	}
	if !mylog.Error2(nil, errors.New("this is a err message")) {
	}
	if !mylog.Error2(os.Open("")) {
	}
	if !mylog.Error2(os.Stdout.Write(nil)) {
	}
	_, err := os.Stdout.Write(nil)
	if !mylog.Error(err) {
	}
	// r1, r2, errno := syscall.SyscallN(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

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
	mylog.Hex("杨开银lost last byte", 0x72B8)
	// mylog.Hex("杨开银lost mid last byte", 0x72B8)
	mylog.Hex("firstEnd xor 0x72B8", 0x72B8)
	mylog.HexDump("key", []byte{0x11, 0x22, 0x33, 0x44})
}
