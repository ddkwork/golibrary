package stream_test

import (
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/stream"
)

func TestName(t *testing.T) {
	cloneStreamForLogPkg("safeStream.go")
	cloneStreamForLogPkg("constraints.go")

	path := "constraints.go"
	g := stream.NewGeneratedFile()
	g.Write(stream.NewBuffer(path).Bytes())
	g.ReplaceAll("package stream", "package pretty")
	g.ReplaceAll(strconv.Quote("github.com/ddkwork/golibrary/mylog"), "")
	g.ReplaceAll("mylog.", "")
	stream.WriteTruncate("../mylog/pretty/"+stream.BaseName(path)+"_gen.go", g.Buffer)

	g = stream.NewGeneratedFile()
	g.Write(stream.NewBuffer(path).Bytes())
	g.ReplaceAll("package stream", "package assert")
	g.ReplaceAll(strconv.Quote("github.com/ddkwork/golibrary/mylog"), "")
	g.ReplaceAll("mylog.", "")
	stream.WriteTruncate("../assert/"+stream.BaseName(path)+"_gen.go", g.Buffer)
}

func cloneStreamForLogPkg(path string) {
	g := stream.NewGeneratedFile()
	g.Write(stream.NewBuffer(path).Bytes())
	g.ReplaceAll("package stream", "package mylog")
	g.ReplaceAll(strconv.Quote("github.com/ddkwork/golibrary/mylog"), "")
	g.ReplaceAll("mylog.", "")
	stream.WriteTruncate("../mylog/"+stream.BaseName(path)+"_gen.go", g.Buffer)
}
