package stream_test

import (
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/std/stream"
)

func TestCloneStreamForLogPkg(t *testing.T) {
	path := "safeStream.go"
	g := stream.NewGeneratedFile()
	g.Write(stream.NewBuffer(path).Bytes())
	g.ReplaceAll("package stream", "package mylog")
	g.ReplaceAll(strconv.Quote("github.com/ddkwork/golibrary/std/mylog"), "")
	g.ReplaceAll("mylog.", "")
	stream.WriteGoFile("../mylog/"+stream.BaseName(path)+"_gen.go", g.Buffer)
}
