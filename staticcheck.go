package golibrary

import (
	"github.com/ddkwork/golibrary/stream"
)

func StaticCheck() {
	stream.RunCommand("go install honnef.co/go/tools/cmd/staticcheck@master")
	// stream.RunCommand("go test -v ./...")
	stream.RunCommand("staticcheck ./...")
}

func UpdateSelf() {
	// mylog.Check(os.Setenv("GOPROXY", "direct"))//制定id的情况下，使用模块代理也是安全的
	hash := stream.GetLastCommitHashLocal("D:\\workspace\\workspace\\golibrary")
	stream.RunCommand("go get -v -x github.com/ddkwork/golibrary@" + hash)
	stream.RunCommand("go mod tidy")
	// 更新桌面的go.mod
}
