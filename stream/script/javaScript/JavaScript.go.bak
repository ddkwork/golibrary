package javaScript

import (
	_ "embed"

	"github.com/ddkwork/golibrary/mylog"

	"github.com/dop251/goja"
)

func Run(src string) goja.Value {
	return mylog.Check2(goja.New().RunString(src))
}
