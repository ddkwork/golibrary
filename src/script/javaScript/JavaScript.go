package javaScript

import (
	_ "embed"
	"github.com/ddkwork/golibrary/mylog"

	"github.com/dop251/goja"
)

//go:embed steamPasswordEnc.js
var SCRIPT string

type (
	Interface interface { //todo add fn
		Value() goja.Value
		Run(src string) bool
	}
	object struct {
		value goja.Value
		err   error
	}
)

func New() Interface {
	return &object{
		value: nil,
		err:   nil,
	}
}

func (o *object) Value() goja.Value { return o.value }
func (o *object) Run(src string) bool {
	o.value, o.err = goja.New().RunString(src)
	return mylog.Error(o.err)
}
