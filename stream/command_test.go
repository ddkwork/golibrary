package stream_test

import (
	"os"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

func Test_cmd(t *testing.T) {
	t.Skip()
	stream.RunCommand("ping www.baidu.com -t")
	return
	// println(stream.RunCommand("cc").Output.String())
	// return
	session := stream.RunCommand(" clang -fsyntax-only -nobuiltininc -emit-llvm -Xclang -fdump-record-layouts -Xclang -fdump-record-layouts-complete merged_headers.h")
	println(session.String())
	return
	stream.RunCommand("cd")
	stream.RunCommand("./a.sh")
	stream.RunCommand("echo 123444 > 1.log")
	mylog.Check(os.Remove("1.log"))
}
