package cmake

import (
	"os"
	"testing"

	"github.com/ddkwork/golibrary/std/mylog"
)

func TestName(t *testing.T) {
	println(os.Getenv("ANDROID_HOME"))
}

func TestInstallInfo(t *testing.T) {
	mylog.Success(Module())
}
