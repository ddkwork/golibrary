package golibrary

import (
	"github.com/ddkwork/golibrary/stream"
)

func StaticCheck() {
	stream.RunCommand("go install honnef.co/go/tools/cmd/staticcheck@latest")
	// stream.RunCommand("go test -v ./...")
	stream.RunCommand("staticcheck ./...")
}
