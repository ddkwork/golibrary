package cmd_test

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/src/stream/tool/cmd"
	"testing"
)

func Test_cmd(t *testing.T) {
	//cmd.SetDir("/tmp")
	run, err := cmd.Run("ls -l && du -h")
	if !mylog.Error(err) {
		return
	}
	mylog.Info(run)
}

func Test_cmd1(t *testing.T) {
	a, b := cmd.Run("./a.sh")
	t.Log("ok...", a, b)
}

func Test_cmd2(t *testing.T) {
	a, b := cmd.Run("echo 123444 > 1.log")
	t.Log("ok...", a, b)
}
