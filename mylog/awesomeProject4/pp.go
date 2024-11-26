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
	input := `
package main

import "encoding/json"

func someFunction() {
	aspect, err := newAspect(aspectBundle, name, rules)
	err = json.Unmarshal(nil, nil)
	if err != nil {
	
	}

	if err := checkSchemaMismatch(bundle.Schema, rules); err != nil {
			return nil, err
		}
}
`

	fset := token.NewFileSet()
	f := mylog.Check2(parser.ParseFile(fset, "input.go", input, parser.ParseComments))

	astutil.Apply(f, func(c *astutil.Cursor) bool {
		switch as := c.Node().(type) {
		case *ast.IfStmt:

			return true
		case *ast.AssignStmt:
			if len(as.Rhs) > 1 {
				return true
			}

			last := len(as.Lhs) - 1
			ident, ok := as.Lhs[last].(*ast.Ident)
			if !ok {
				return true
			}

			if ident.Name != "err" && ident.Name != "_" {
				return true
			}

			rightBack, ok := as.Rhs[0].(*ast.CallExpr)
			if !ok {
				return true
			}

			var newRight *ast.CallExpr

			switch len(as.Lhs) {
			case 1:
				newRight = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "mylog"},
						Sel: &ast.Ident{Name: "Check"},
					},

					Args: []ast.Expr{rightBack},
				}
				c.Replace(&ast.ExprStmt{X: newRight})

			case 2:
				newRight = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "mylog"},
						Sel: &ast.Ident{Name: "Check2"},
					},
					Args: []ast.Expr{rightBack},
				}

			case 3:
				newRight = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "mylog"},
						Sel: &ast.Ident{Name: "Check3"},
					},
					Args: []ast.Expr{rightBack},
				}
			case 4:
				newRight = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "mylog"},
						Sel: &ast.Ident{Name: "Check4"},
					},
					Args: []ast.Expr{rightBack},
				}
			case 5:
				newRight = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "mylog"},
						Sel: &ast.Ident{Name: "Check5"},
					},
					Args: []ast.Expr{rightBack},
				}
			case 6:
				newRight = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "mylog"},
						Sel: &ast.Ident{Name: "Check6"},
					},
					Args: []ast.Expr{rightBack},
				}
			case 7:
				newRight = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "mylog"},
						Sel: &ast.Ident{Name: "Check7"},
					},
					Args: []ast.Expr{rightBack},
				}
			}
			as.Rhs[0] = newRight
			as.Lhs = as.Lhs[:last]

		}
		return true
	}, nil)

	var output strings.Builder
	mylog.Check(printer.Fprint(&output, fset, f))

	fmt.Println(output.String())
}
