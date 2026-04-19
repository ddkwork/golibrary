package mylog_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/ddkwork/golibrary/std/mylog"
)

func BenchmarkName(b *testing.B) {
	for b.Loop() {
		mylog.Call(xx)
	}
}

func TestFix111111111111111111111111111111111111111111111111111111111LongKey(t *testing.T) {
	mylog.Info("测试超长的key")
}

func TestFix(t *testing.T) {
	mylog.Info("tttttttttttttttttttttttt")
	mylog.Trace("kkkkkkkkkkkkkkkkkkkk")
	mylog.HexDump([]byte{
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
}

func TestPrettyFormat2(test *testing.T) {
	key := "Path"
	value := `C:\Program Files\Eclipse Adoptium\jdk-23.0.2.7-hotspot\bin;D:\todo\ewdk\dist\sdk\bin\Hostx64\x64;E:\Program Files\Microsoft Visual Studio\2022\BuildTools\VC\Tools\MSVC\14.44.35207\bin\HostX64\x64;E:\Program Files\Windows Kits\10\bin\10.0.28000.0\x64;C:\Windows\system32;C:\Windows;C:\Windows\System32\Wbem;C:\Windows\System32\WindowsPowerShell\v1.0\;C:\Windows\System32\OpenSSH\;C:\Program Files\Go\bin;C:\Program Files\Git\cmd;C:\Program Files\CMake\bin;C:\TDM-GCC-64\bin;C:\Program Files\LLVM\bin;C:\Program Files\CodeArts Agent\bin`
	m := map[string]string{key: value}
	type Struct struct {
		Path  string
		Value string
	}
	obj := Struct{
		Path:  key,
		Value: value,
	}
	mylog.Struct(m)
	mylog.Struct(obj)
}

func TestLog(t *testing.T) {
	mylog.Call(func() {
		xx()
		safeAppLife()
	})
}

func safeAppLife() {
	mylog.Check(errors.New("this is a err value"))
	mylog.Check(errors.New("this is a err value"))
	mylog.Check2(os.Open(""))
	mylog.Check2(os.Stdout.Write(nil))
	mylog.Check2(os.Stdout.Write(nil))
}

func xx() {
	mylog.HexDump([]byte{
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
	mylog.Hex(uint32(0x888))
	mylog.Info("tttttttttttttttttttttttt")
	mylog.Trace("kkkkkkkkkkkkkkkkkkkk")
	mylog.Warning("mmmmmmmmmmmmmmmmmmmmmmmmmmmkkkkkkkkkkkkkkkkkkkkkkkkkkk")
	mylog.Json(`{"manifestVersion":"1.1","engineVersion":"3.3.2180.8236","info":{"id":"VisualStudio/17.3.0+32804.467","buildBranch":"d17.3","buildVersion":"17.3.32804.467","localBuild":"build-lab","manifestName":"VisualStudio","manifestType":"installer","productDisplayVersion":"17.3.0","productLine":"Dev17","productLineVersion":"2022","productMilestone":"RTW","`)
	mylog.Success("vgoTest pass")
	mylog.Struct(
		struct {
			A int
			B string
			C []byte
			D []byte
		}{
			A: 89,
			B: "jhjsbdd",
			C: []byte{0x11, 0x22, 0x33, 0x44},
			D: []byte{
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
			},
		})
	mylog.Warning()
	mylog.Hex(uint32(0x72B8))
	mylog.Hex(uint32(0x72B8))
	mylog.HexDump([]byte{0x11, 0x22, 0x33, 0x44})
	mylog.HexDump([]byte{0x11, 0x22, 0x33, 0x11, 0x22, 0x33, 0x11, 0x22, 0x33, 0x11, 0x22, 0x33, 0x11, 0x22, 0x33, 0x11, 0x22, 0x33, 0x11, 0x22, 0x33, 0x11, 0x22, 0x33, 0x44})
}

func TestAutoFillKey(t *testing.T) {
	mylog.Info("auto filled key test")
}

func TestEmptyValuePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for empty value, but didn't get one")
		} else {
			if !strings.Contains(r.(string), "log value cannot be empty") {
				t.Errorf("Expected 'log value cannot be empty' panic, got: %v", r)
			}
		}
	}()
	mylog.Info()
}

func TestValidLog(t *testing.T) {
	mylog.Info("valid_value")
	mylog.Warning("warn_value")
	mylog.Success("success_value")
}

func TestFormatSyntaxPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for format syntax, but didn't get one")
		} else {
			if !strings.Contains(r.(string), "log value cannot contain format syntax") {
				t.Errorf("Expected 'log value cannot contain format syntax' panic, got: %v", r)
			}
		}
	}()
	value := "value is " + "%s" // nolint:staticcheck
	mylog.Info(value)
}
