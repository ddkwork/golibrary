package gen

import (
	"bytes"
	"testing"
)

func TestName(t *testing.T) {
	g := GeneratedFile{
		gen:              nil,
		skip:             false,
		filename:         "",
		goImportPath:     "",
		buf:              bytes.Buffer{},
		packageNames:     nil,
		usedPackageNames: nil,
		manualImports:    nil,
		annotations:      nil,
	}
	g.P("type x struct{}")

	g.P("},")
}
