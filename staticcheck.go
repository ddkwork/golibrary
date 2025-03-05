package golibrary

import (
	"github.com/ddkwork/golibrary/stream"
	"strings"
)

func StaticCheck() {
	stream.RunCommand("go install honnef.co/go/tools/cmd/staticcheck@latest")
	// stream.RunCommand("go test -v ./...")
	stream.RunCommand("staticcheck ./...")
}

func UpdateDependencies() {
	for s := range strings.Lines(` go get -x gioui.org@main
	 go get -x gioui.org/cmd@main
	 go get -x gioui.org/example@main
	 go get -x gioui.org/x@main
	 go get -x github.com/oligo/gvcode@main
	 go get -x github.com/ddkwork/golibrary@master
	 go get -x github.com/ddkwork/ux@master
	 go get -x github.com/google/go-cmp@master
	 go get -x github.com/ddkwork/app@master
	 go get -x github.com/ddkwork/toolbox@master
	 go get -x github.com/ddkwork/unison@master
	 go get -x github.com/ebitengine/purego@main
	 go get -x github.com/saferwall/pe@main
	 ::go get -u -x all
	 go mod tidy`) {
		if strings.HasPrefix(s, "::") {
			continue
		}
		stream.RunCommand(s)

	}
}
