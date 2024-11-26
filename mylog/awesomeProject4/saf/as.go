package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/ddkwork/golibrary/mylog"
)

func main() {
	src := `
package main

import (
	"fmt"
	"go/parser"
	"go/token"
)

func main() {
	if err != nil {
		fmt.Println(err)
	}

	if err := checkSchemaMismatch(bundle.
Schema, rules); err != nil {
		return nil, err
	}

	if f, err := parser.ParseFile(fset, "input.go", input, parser.ParseComments); err != nil {
		log.Fatal(err)
	}


	fset := token.
		NewFileSet()
}`

	fset := token.NewFileSet()

	file := mylog.Check2(parser.ParseFile(fset, "demo.go", src, 0))

	mergeMultiline := func(cursor *astutil.Cursor) bool {
		if basicLit, ok := cursor.Node().(*ast.BasicLit); ok {

			pos := fset.Position(basicLit.Pos()).Offset
			end := fset.Position(basicLit.End()).Offset
			if strings.Count(src[pos:end], "\n") > 0 {

				merged := strings.ReplaceAll(src[pos:end], "\n", "")
				cursor.Replace(&ast.BasicLit{ValuePos: basicLit.ValuePos, Kind: basicLit.Kind, Value: merged})
			}
		}
		return true
	}

	deleteIfErrNotNil := func(cursor *astutil.Cursor) bool {
		if ifStmt, ok := cursor.Node().(*ast.IfStmt); ok {
			if ifStmt.Init == nil && ifStmt.Cond != nil && ifStmt.Body != nil {
				binaryExpr, isBinary := ifStmt.Cond.(*ast.BinaryExpr)
				if isBinary && binaryExpr.Op == token.NEQ {
					if ident, isIdent := binaryExpr.X.(*ast.Ident); isIdent && ident.Name == "err" {
						if basicLit, isBlank := binaryExpr.Y.(*ast.Ident); isBlank && basicLit.Name == "nil" {
							cursor.Delete()
						}
					}
				}
			}
		}
		return true
	}

	modifyIfErrNotNil := func(cursor *astutil.Cursor) bool {
		if ifStmt, ok := cursor.Node().(*ast.IfStmt); ok {
			if ifStmt.Init != nil && ifStmt.Cond != nil && ifStmt.Body != nil {
				exprStmt, ok := ifStmt.Init.(*ast.AssignStmt)
				if ok && len(exprStmt.Lhs) == 1 && len(exprStmt.Rhs) == 1 {
					expr := exprStmt.Rhs[0]
					newExpr := &ast.ExprStmt{
						X: &ast.CallExpr{
							Fun:  &ast.Ident{Name: "mylog.Check"},
							Args: []ast.Expr{expr},
						},
					}
					cursor.Replace(newExpr)
				}
			}
		}
		return true
	}

	modifyIfErrNotNilAssignment := func(cursor *astutil.Cursor) bool {
		if ifStmt, ok := cursor.Node().(*ast.IfStmt); ok {
			if ifStmt.Init != nil && ifStmt.Cond != nil && ifStmt.Body != nil {
				assignStmt, ok := ifStmt.Init.(*ast.AssignStmt)
				if ok && len(assignStmt.Lhs) == 2 {
					f := assignStmt.Lhs[0]
					newExpr := &ast.AssignStmt{
						Lhs: []ast.Expr{f},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{&ast.CallExpr{
							Fun:  &ast.Ident{Name: "mylog.Check"},
							Args: []ast.Expr{assignStmt.Rhs[0]},
						}},
					}
					cursor.Replace(newExpr)
				}
			}
		}
		return true
	}

	astutil.Apply(file, deleteIfErrNotNil, nil)
	astutil.Apply(file, modifyIfErrNotNil, nil)
	astutil.Apply(file, modifyIfErrNotNilAssignment, nil)

	astutil.Apply(file, nil, mergeMultiline)

	var output strings.Builder
	mylog.Check(printer.Fprint(&output, fset, file))

	fmt.Println(output.String())
}
