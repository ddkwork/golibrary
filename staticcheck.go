package golibrary

import (
	"github.com/ddkwork/golibrary/stream"
)

func StaticCheck() {
	stream.RunCommand("go install honnef.co/go/tools/cmd/staticcheck@latest")
	// stream.RunCommand("go test -v ./...")
	stream.RunCommand("staticcheck ./...")
}

func UpdateSelf() {
	hash := stream.GetLastCommitHashLocal("D:\\workspace\\workspace\\golibrary")
	stream.RunCommand("go get -v -x github.com/ddkwork/golibrary@" + hash)
	stream.RunCommand("go mod tidy")
	//更新桌面的dep.txt
}
