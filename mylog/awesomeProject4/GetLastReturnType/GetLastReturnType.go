package main

import (
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/ddkwork/golibrary/mylog"
)

func main() {
	fs := token.NewFileSet()
	node := mylog.Check2(parser.ParseFile(fs, "source.go", `
package main

import "unicode"

func testRetType() {
    index, _ := rightWordBoundary()
}

func rightWordBoundary() (byteIndex, runeIndex int) {
    return 0, 0
}
`, parser.ParseComments))

	var testRetTypeFunc *ast.FuncDecl
	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "testRetType" {
			continue
		}
		testRetTypeFunc = fn
		break
	}

	astutil.Apply(testRetTypeFunc, func(c *astutil.Cursor) bool {
		assignStmt, ok := c.Node().(*ast.AssignStmt)
		if !ok {
			return true
		}

		for _, left := range assignStmt.Lhs {

			ident, ok := left.(*ast.Ident)
			if !ok || ident.Name != "_" {
				continue
			}

			lastReturnType := GetLastReturnType(assignStmt)
			println(lastReturnType)
		}

		return true
	}, nil)
}

func GetLastReturnType(assignStmt *ast.AssignStmt) string {
	expr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	if ok {
		funDecl, ok := expr.Fun.(*ast.Ident).Obj.Decl.(*ast.FuncDecl)
		if !ok {
			panic("Cannot get last retun type")
		}
		if results := funDecl.Type.Results; results != nil && len(results.List) > 0 {
			lastResult := results.List[len(results.List)-1]
			if ident, ok := lastResult.Type.(*ast.Ident); ok {
				return ident.Name
			}
		}
	}
	panic("Cannot get last retun type")
}
