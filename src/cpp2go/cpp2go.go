package cpp2go

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/goplus/c2go"
	"github.com/goplus/c2go/cl"
	"github.com/goplus/c2go/clang/preprocessor"
	"github.com/goplus/c2go/clang/types"
	"github.com/goplus/gox"
	gotypes "go/types"
	"path/filepath"
	_ "unsafe"
)

// ╔════╤══════════════════════════════════════╤═════════════════════════════════════════════════╤══════╤═════════╤══════╗
// ║ ID ║                 api                  ║                    function                     ║ note ║ chinese ║ todo ║
// ╠════╪══════════════════════════════════════╪═════════════════════════════════════════════════╪══════╪═════════╪══════╣
// ║ 1  ║ TranslateCFile(path, pkg string)     ║ convert c to go                                 ║      ║         ║      ║
// ╠════╪══════════════════════════════════════╪═════════════════════════════════════════════════╪══════╪═════════╪══════╣
// ║ 2  ║ Translate(root string) (ok bool)     ║ Translate cpp or c to go, not full part convert ║      ║         ║      ║
// ╠════╪══════════════════════════════════════╪═════════════════════════════════════════════════╪══════╪═════════╪══════╣
// ║ 3  ║ RemoveComment(root string) (ok bool) ║ remove comment                                  ║      ║         ║      ║
// ╚════╧══════════════════════════════════════╧═════════════════════════════════════════════════╧══════╧═════════╧══════╝
type (
	Interface interface {
		TranslateCFile(path, pkg string)
		Scanner
	}
	object struct{ s Scanner }
)

const (
	flagShort = 1 << iota
	flagLong
	flagLongLong
	flagUnsigned
	flagSigned
	flagComplex
	flagStructOrUnion
)

var (
	//go:linkname intTypes parser.intTypes
	//go:linkname Long types.Long
	Long = types.Int
	//go:linkname Ulong types.Ulong
	Ulong    = types.Uint
	intTypes = [...]gotypes.Type{
		0:                                      types.Int,
		flagShort:                              gotypes.Typ[gotypes.Int16],
		flagLong:                               types.Int, // ctypes.Long,
		flagLong | flagLongLong:                gotypes.Typ[gotypes.Int64],
		flagUnsigned:                           types.Uint,
		flagShort | flagUnsigned:               gotypes.Typ[gotypes.Uint16],
		flagLong | flagUnsigned:                types.Uint, // ctypes.Ulong,
		flagLong | flagLongLong | flagUnsigned: gotypes.Typ[gotypes.Uint64],
		flagShort | flagLong | flagLongLong | flagUnsigned: nil,
	}
)

func init() {
	Long = types.Int
	Ulong = types.Uint
}

func New() Interface                              { return &object{s: newScanner()} }
func (o *object) Translate(root string) (ok bool) { return o.s.Translate(root) }
func (o *object) TranslateCFile(path, pkg string) {
	mylog.Struct(intTypes)
	abs, err := filepath.Abs(path)
	if !mylog.Error(err) {
		return
	}
	cl.SetDebug(cl.DbgFlagAll)
	preprocessor.SetDebug(preprocessor.DbgFlagAll)
	gox.SetDebug(gox.DbgFlagInstruction) // | gox.DbgFlagMatch)
	c2go.Run(pkg, abs, c2go.FlagDumpJson, nil)
}
func (o *object) RemoveComment(root string) (ok bool) { return o.s.RemoveComment(root) }
