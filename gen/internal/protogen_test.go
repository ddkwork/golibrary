package internal

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
	// g.P("out := ", protoimplPackage.Ident("TypeBuilder"), "{")
	// g.P("File: ", protoimplPackage.Ident("DescBuilder"), "{")
	// g.P("GoPackagePath: ", reflectPackage.Ident("TypeOf"), "(x{}).PkgPath(),")
	// g.P("NumEnums: ", len(f.allEnums), ",")
	// g.P("NumMessages: ", len(f.allMessages), ",")
	// g.P("NumExtensions: ", len(f.allExtensions), ",")
	// g.P("NumServices: ", len(f.Services), ",")
	g.P("},")
}
