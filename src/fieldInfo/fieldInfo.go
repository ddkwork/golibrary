package fieldInfo

import (
	"fmt"
	"github.com/ddkwork/golibrary/src/tuitable"
)

type (
	FieldInfo interface {
		Append(info Row)
		Gen() (body string)
	}
	Row struct {
		Number   string
		Deep     string
		Name     string
		Kind     string
		Size     string
		ValueFmt string //Hex|Decimal Bool String (HexDump Map)
		Value    any
		Stream   string
		Json     string
	}
	fieldInfo struct {
		infos []Row
	}
)

func New() FieldInfo { return &fieldInfo{} }
func (f *fieldInfo) Append(info Row) {
	f.infos = append(f.infos, info)
}
func (f *fieldInfo) Gen() (body string) {
	t := tuitable.NewTable()
	t.SetHeaders("ID", "number", "deep", "name", "kind", "size", "ValueFmt", "stream", "Json")
	for i, info := range f.infos {
		t.AddRow(fmt.Sprint(i+1), info.Number, info.Deep, info.Name, info.Kind, info.Size,
			info.ValueFmt, info.Stream, info.Json)
	}
	return t.Body()
}
