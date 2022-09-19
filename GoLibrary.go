package golibrary

import (
	"reflect"
	"unsafe"
)

func Struct(s interface{}) []byte {
	//theme.LightTheme()
	//theme.FromJSON()
	//widget.NewButton()

	v := reflect.ValueOf(s)
	sz := int(v.Elem().Type().Size())
	return unsafe.Slice((*byte)(unsafe.Pointer(v.Pointer())), sz)
}
