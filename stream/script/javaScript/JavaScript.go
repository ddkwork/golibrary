package javaScript

import (
	_ "embed"

	"github.com/ddkwork/golibrary/mylog"

	"github.com/dop251/goja"
)

//go:embed steamPasswordEnc.js
var SCRIPT string

func Run(src string) (goja.Value, bool) {
	value, err := goja.New().RunString(src)
	return value, mylog.Error(err)
}
