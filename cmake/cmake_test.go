package cmake

import (
	"os"
	"testing"

	"github.com/ddkwork/golibrary/std/mylog"
)

func TestName(t *testing.T) {
	// ANDROID_HOME=D:\sdk
	println(os.Getenv("ANDROID_HOME"))
}

func TestInstallInfo(t *testing.T) {
	mylog.Struct(InstallInfo())
}
