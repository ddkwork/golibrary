package mydoc

import (
	"fmt"
	"github.com/ddkwork/golibrary/src/stream"
	"github.com/ddkwork/golibrary/src/stream/tool"
	"github.com/ddkwork/golibrary/src/tuitable"
)

// need close vcs "use non-modal command interface",
//
//go:generate go install github.com/traefik/yaegi/cmd/yaegi@v0.14.1
//go:generate yaegi extract github.com/aquasecurity/table
type (
	Doc interface {
		Append(info Row)
		Gen() (body string)
	}
	Row struct {
		Api      string
		Function string
		Note     string
		Chinese  string
		Todo     string
	}
	doc struct{ infos []Row }
)

func New() Doc                 { return &doc{infos: make([]Row, 0)} }
func (d *doc) Append(info Row) { d.infos = append(d.infos, info) }
func (d *doc) Gen() (body string) {
	t := tuitable.NewTable()
	t.SetHeaders("ID", "api", "function", "note", "chinese", "todo")
	for i, info := range d.infos {
		t.AddRow(fmt.Sprint(i+1), info.Api, info.Function, info.Note, info.Chinese, info.Todo)
	}
	body = t.Body()
	lines, ok := tool.File().ToLines(body)
	if !ok {
		panic(ok)
	}
	b := stream.New()
	b.WriteStringLn("//")
	for _, line := range lines {
		b.WriteStringLn("//  " + line)
	}
	b.WriteStringLn("//")
	body = b.String()
	return
}
