package golibrary

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"os"
)

func StaticCheck() {
	stream.RunCommand("go install honnef.co/go/tools/cmd/staticcheck@latest")
	// stream.RunCommand("go test -v ./...")
	stream.RunCommand("staticcheck ./...")
}

func UpdateSelf() {
	mylog.Check(os.Setenv("GOPROXY", "direct"))
	hash := stream.GetLastCommitHashLocal("D:\\workspace\\workspace\\golibrary")
	stream.RunCommand("go get -v -x github.com/ddkwork/golibrary@" + hash)
	stream.RunCommand("go mod tidy")
	//更新桌面的dep.txt
}
