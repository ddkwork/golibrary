package javaScript

import (
	_ "embed"
	"fmt"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterfaceJs(t *testing.T) {
	p := New()
	assert.True(t, p.Run(`
function add(){
 return 1+1+""
}
add()
`))
	mylog.Success("js result", p.Value().String())
	assert.Equal(t, p.Value().String(), "2")
	assert.Equal(t, p.Value().ToInteger(), int64(2))
}

func TestName(t *testing.T) {
	vm := goja.New()
	_, err := vm.RunString(SCRIPT)
	if err != nil {
		panic(err)
	}

	var fn func(string, string, string) string
	err = vm.ExportTo(vm.Get("encPassword"), &fn)
	if err != nil {
		panic(err)
	}

	publickey_mod := "cf40107b85f0a48c34a64fef862819a5a6c53f364edf9307047f1a34d3d762098b50b077e19cbfaaed84d189fa8148f3d552038b257490fb0c41de518b65dbe6fcbd9e32a6dfd07c3b221c826d6ef0e433f76faff1957d55de1f2d095cdd98f55d4354fd6e0156b23855817c84433baac45033b898dca4bb19ac02ab4c9da76a6d7eae0ffc3cede649d1273a9c2aea628607527f6fcb63a99a4c9ff7e50db618413ee1bdeb8a56e5104444b2553ec770b6dda002af77a7af5c726624aca9e4948c3b76724b4f39c620b07f4152bfa410a9caa4883787435301894a30b79281e2118bc22503e3b1e6a094939f7f2dec94386baca52b71b8785fc40fde242bfe5b"
	publickey_exp := "010001"
	password := "123456"

	enc := fn(publickey_mod, publickey_exp, password)
	fmt.Println(enc) // note, _this_ value in the function will be undefined.
}
