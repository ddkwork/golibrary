package cmd_test

import (
	"os"
	"testing"

	"github.com/ddkwork/golibrary/stream/cmd"
)

func Test_cmd(t *testing.T) {
	// cmd.SetDir("/tmp")
	// cmd.Run("ls -l && du -h")
	// cmd.Run("cd D:\\clone")
	cmd.Run("cd")
	cmd.Run("./a.sh")
	cmd.Run("echo 123444 > 1.log")
	os.Remove("1.log")
	// cmd.Run("ping www.baidu.com")
	// cmd.Run("ping www.baidu.com -t ")
}
