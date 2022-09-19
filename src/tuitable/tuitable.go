package tuitable

import (
	"github.com/aquasecurity/table"
	"strings"
	"time"
)

type Table struct {
	*table.Table
	builder *strings.Builder
}

func NewTable() *Table {
	b := new(strings.Builder)
	t := table.New(b)
	t.SetDividers(table.Dividers{
		ALL: "╪",
		NES: "╠",
		NSW: "╣",
		NEW: "╧",
		ESW: "╤",
		NE:  "╚",
		NW:  "╝",
		SW:  "╗",
		ES:  "╔",
		EW:  "═",
		NS:  "║",
	})
	return &Table{
		Table:   t,
		builder: b,
	}
}

func (t *Table) Body() string {
	t.render()
	return t.builder.String()
}
func (t *Table) render() *Table {
	time.Sleep(time.Second)
	t.Render()
	return t
}


