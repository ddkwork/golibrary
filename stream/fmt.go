package stream

import (
	"encoding/json"

	"github.com/ddkwork/golibrary/mylog"
)

func (o *Stream) JsonIndent() string {
	if !mylog.Error(json.Indent(o.Buffer, o.Bytes(), "", " ")) {
		return ""
	}
	return o.String()
}
