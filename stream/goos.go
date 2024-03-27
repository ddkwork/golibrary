package stream

import (
	"os"
	"runtime"

	"github.com/ddkwork/golibrary/mylog"
)

func IsTermux() bool  { return mylog.Error2(os.Stat("/data/data/com.termux/files/usr")) }
func IsAndroid() bool { return runtime.GOOS == `android` }
func IsWindows() bool { return runtime.GOOS == `windows` }
func IsLinux() bool   { return runtime.GOOS == `linux` }
