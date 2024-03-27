package stream

import (
	"encoding/json"
	"path/filepath"

	"github.com/ddkwork/golibrary/mylog"
)

func MarshalJSON(v any, name string) ([]byte, error) {
	return json.MarshalIndent(v, "", " ")
}

func MarshalJsonToFile(v any, name string) bool {
	indent, err := json.MarshalIndent(v, "", " ")
	if !mylog.Error(err) {
		return false
	}
	ext := filepath.Ext(name)
	if ext != ".json" {
		name += ".json"
	}
	return WriteTruncate(name, indent)
}
