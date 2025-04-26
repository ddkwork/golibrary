package fakeError

import (
	"fmt"
	"github.com/ddkwork/golibrary/mylog"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func RemoveComments(file *ast.File) {
	mylog.Todo("https://github.com/Greyh4t/nocomment")
	newComments := make([]*ast.CommentGroup, 0) // can not be nil,why
	for _, group := range file.Comments {
		var newGroup ast.CommentGroup
		for _, comment := range group.List {
			if strings.HasPrefix(comment.Text, "//go:") {
				newGroup.List = append(newGroup.List, comment)
			}
		}
		if len(newGroup.List) > 0 {
			newComments = append(newComments, &newGroup)
		}
	}
	file.Comments = newComments
}

var Skips = []string{
	`vendor`,
	`/gioview/`,
	`/gio/`,
	`/gio-cmd/`,
	`/gio-example/`,
	`/gio-x/`,
	`/toolbox/`,
	`/unison/`,
	`/ux/patch/`,
}

func FakeError(path string, removeComments ...bool) {
	if path == "" {
		path = "."
	}
	mylog.Check(filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		abs := mylog.Check2(filepath.Abs(path))
		abs = filepath.ToSlash(abs)
		if filepath.Ext(abs) == ".go" {
			if filepath.Base(abs) == "SkipCheckBase.go" {
				return nil
			}
			for _, skip := range Skips {
				if strings.Contains(abs, skip) {
					// Warning("skip", abs)
					return nil
				}
			}
			if !strings.HasSuffix(path, ".go") {
				mylog.Check(fmt.Errorf("not a go file: %s", path))
			}
			fileSet := token.NewFileSet()
			file := mylog.Check2(parser.ParseFile(fileSet, path, nil, parser.ParseComments))
			if len(removeComments) > 0 {
				RemoveComments(file)
			}
			mylog.Call(func() {
				s := fakeError(fileSet, file, string(mylog.Check2(os.ReadFile(path))))
				mylog.WriteGoFileWithDiff(path, []byte(s))
			})
		}
		return err
	}))
}

func fakeErrorTest(text string) string {
	ret := ""
	mylog.Call(func() {
		fileSet := token.NewFileSet()
		file := mylog.Check2(parser.ParseFile(fileSet, "", text, parser.ParseComments))
		ret = fakeError(fileSet, file, text)
	})
	return ret
}

func fakeError(fileSet *token.FileSet, file *ast.File, text string) string {
	Replaces := make([]Edit, 0)

	fnCall := func(n int) string {
		switch n {
		case 1:
			return "mylog.Check"
		default:
			return "mylog.Check" + strconv.Itoa(n)
		}
	}

	fnHandleAssign := func(x *ast.AssignStmt, e *Edit) {
		if len(x.Rhs) > 1 {
			return
		}
		last := len(x.Lhs) - 1
		lastIdent, ok := x.Lhs[last].(*ast.Ident)
		if !ok {
			return
		}
		if lastIdent.Name != "err" && lastIdent.Name != "_" {
			return
		}
		if lastIdent.Name == "_" {
			lastReturnType, ok := GetLastReturnType(x)
			if !ok {
				return
			}
			if lastReturnType != "error" {
				return
			}
		}
		left := ""
		for i, v := range x.Lhs {
			c := getNodeCode(v, fileSet, text)
			if c == "err" || c == "_" {
				continue
			}
			if i < len(x.Lhs)-1 {
				left += c + ","
			}
		}
		//只有一个变量需要删除 ,
		left = strings.TrimRight(left, ",")

		right := ""
		tk := x.Tok.String()
		if left == "" {
			tk = ""
		}
		for i, v := range x.Rhs {
			right += getNodeCode(v, fileSet, text)
			if i < len(x.Rhs)-1 { //todo bug
				right += ","
			}
		}

		ee := Edit{
			Start: x.Pos(),
			End:   x.End(),
			Line:  fileSet.Position(x.Pos()).Line,
			New:   left + tk + fnCall(len(x.Lhs)) + "(" + right + ")",
			edge:  edge(x),
		}
		if e != nil { //AssignStmt in IfStmt
			ee = Edit{
				Start: e.Start,
				End:   e.End,
				Line:  e.Line,
				New:   left + tk + fnCall(len(x.Lhs)) + "(" + right + ")",
				edge:  e.edge,
			}
		}
		Replaces = append(Replaces, ee)
	}

	skipAssign := false
	for n := range ast.Preorder(file) {
		switch x := n.(type) {
		case *ast.IfStmt: //if err := backendConn.Close(); err != nil {
			for ifStmt := range findIfErrNotNil(n) {
				isOneWorkCode := false //if 块内部语句只有1句，没有其他业务逻辑,则直接替换为mylog.Check(业务逻辑)
				for i, stmt := range ifStmt.Body.List {
					if i > 1 {
						panic("if 块内部语句超过1句")
					}
					if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
						c := getNodeCode(exprStmt, fileSet, text)
						if strings.TrimSpace(c) == "panic(err)" {
							isOneWorkCode = true
							break
						}
						if strings.HasPrefix(c, "log.") && strings.HasSuffix(c, "(err)") {
							isOneWorkCode = true
							break
						}
						if strings.Contains(c, "continue") { //todo bug
							isOneWorkCode = false
							break
						}
					}
				}
				if strings.HasPrefix(getNodeCode(ifStmt, fileSet, text), "if err != nil {") && isOneWorkCode {
					Replaces = append(Replaces, Edit{
						Start: ifStmt.Pos(),
						End:   ifStmt.End(),
						Line:  fileSet.Position(ifStmt.Pos()).Line,
						New:   "",
						edge:  edge(ifStmt),
					})
					skipAssign = true
					break
				}
				if ifStmt.Init == nil {
					break
				}
				if stmt, ok := ifStmt.Init.(*ast.AssignStmt); ok {
					if isOneWorkCode {
						fnHandleAssign(stmt, &Edit{
							Start: ifStmt.Pos(),
							End:   ifStmt.End(),
							Line:  fileSet.Position(ifStmt.Pos()).Line,
							New:   "",
							edge:  edge(ifStmt) + " # " + edge(stmt),
						})
						skipAssign = true
						break
					}
				}
			}
		case *ast.AssignStmt: //backendConn, err := net.DialTimeout("tcp", forwardTarget, 5*time.Second)
			if skipAssign {
				skipAssign = false
				continue
			}
			fnHandleAssign(x, nil)
		}
	}
	mylog.Struct(Replaces)
	return Apply(text, Replaces)
}

func getNodeCode(astNode ast.Node, f *token.FileSet, code string) string {
	c := code[f.Position(astNode.Pos()).Offset:f.Position(astNode.End()).Offset]
	c = strings.TrimSpace(c)
	//mylog.Json("code dump", c)
	return c
}

func findIfErrNotNil(n ast.Node) iter.Seq[*ast.IfStmt] {
	return func(yield func(*ast.IfStmt) bool) {
		if stmt, ok := n.(*ast.IfStmt); ok {
			if b, ok := stmt.Cond.(*ast.BinaryExpr); ok {
				if b.Op == token.NEQ {
					if x, ok := b.X.(*ast.Ident); ok && x.Name == "err" {
						if y, ok := b.Y.(*ast.Ident); ok && y.Name == "nil" {
							if !yield(stmt) {
								return
							}
						}
					}
				}
			}
		}
	}
}

func GetLastReturnType(assignStmt *ast.AssignStmt) (lastReturnType string, b bool) {
	if expr, ok := assignStmt.Rhs[0].(*ast.CallExpr); ok {
		switch e := expr.Fun.(type) {
		case *ast.Ident:
			if e.Obj == nil {
				return
			}
			funDecl, ok := e.Obj.Decl.(*ast.FuncDecl)
			if !ok {
				return
			}
			if results := funDecl.Type.Results; results != nil && len(results.List) > 0 {
				lastResult := results.List[len(results.List)-1]
				if ident, ok := lastResult.Type.(*ast.Ident); ok {
					return ident.Name, true
				}
			}
		}
	}
	return
}

type Edit struct {
	Start, End token.Pos
	Line       int
	New        string
	edge       string
}

func Apply(text string, replaces []Edit) string {
	if len(replaces) == 0 {
		return text
	}
	// 按起始位置从大到小排序,即从后往前替换，避免处理过程中坐标变化
	sort.Slice(replaces, func(i, j int) bool {
		return replaces[i].Start > replaces[j].Start
	})
	for _, r := range replaces {
		if r.Start > r.End {
			panic("起始位置大于终止位置")
		}
		// 左+新内容+右
		left := text[:r.Start-1] //Start的要删除的第一个字符发偏移，需要cut掉,这个通过单元测试了，不要改
		right := text[r.End-1:]  //不这样连续替换后没有换行，两个mycheck在一行导致语法错误
		text = left + r.New + right
	}
	text = strings.ReplaceAll(text, `var err error`, "")
	text = strings.ReplaceAll(text, `import (`, `import (
	"github.com/ddkwork/golibrary/mylog"`)
	return string(mylog.Check2(format.Source([]byte(text))))
}

type Type interface {
	*ast.ArrayType |
		*ast.AssignStmt |
		*ast.BadDecl |
		*ast.BadExpr |
		*ast.BadStmt |
		*ast.BasicLit |
		*ast.BinaryExpr |
		*ast.BlockStmt |
		*ast.BranchStmt |
		*ast.CallExpr |
		*ast.CaseClause |
		*ast.ChanType |
		*ast.CommClause |
		*ast.Comment |
		*ast.CommentGroup |
		*ast.CompositeLit |
		*ast.DeclStmt |
		*ast.DeferStmt |
		*ast.Ellipsis |
		*ast.EmptyStmt |
		*ast.ExprStmt |
		*ast.Field |
		*ast.FieldList |
		*ast.File |
		*ast.ForStmt |
		*ast.FuncDecl |
		*ast.FuncLit |
		*ast.FuncType |
		*ast.GenDecl |
		*ast.GoStmt |
		*ast.Ident |
		*ast.IfStmt |
		*ast.ImportSpec |
		*ast.IncDecStmt |
		*ast.IndexExpr |
		*ast.IndexListExpr |
		*ast.InterfaceType |
		*ast.KeyValueExpr |
		*ast.LabeledStmt |
		*ast.MapType |
		*ast.Package |
		*ast.ParenExpr |
		*ast.RangeStmt |
		*ast.ReturnStmt |
		*ast.SelectStmt |
		*ast.SelectorExpr |
		*ast.SendStmt |
		*ast.SliceExpr |
		*ast.StarExpr |
		*ast.StructType |
		*ast.SwitchStmt |
		*ast.TypeAssertExpr |
		*ast.TypeSpec |
		*ast.TypeSwitchStmt |
		*ast.UnaryExpr |
		*ast.ValueSpec
}

var (
	anyType    = reflect.TypeFor[any]()
	stringType = reflect.TypeFor[string]()
	bytesType  = reflect.TypeFor[[]byte]()
	byteType   = reflect.TypeFor[byte]()
)

func edge[T Type](n T) string {
	switch any(n).(type) {
	case *ast.ArrayType:
		return "ArrayType"
	case *ast.AssignStmt:
		return "AssignStmt"
	case *ast.BadDecl:
		return "BadDecl"
	case *ast.BadExpr:
		return "BadExpr"
	case *ast.BadStmt:
		return "BadStmt"
	case *ast.BasicLit:
		return "BasicLit"
	case *ast.BinaryExpr:
		return "BinaryExpr"
	case *ast.BlockStmt:
		return "BlockStmt"
	case *ast.BranchStmt:
		return "BranchStmt"
	case *ast.CallExpr:
		return "CallExpr"
	case *ast.CaseClause:
		return "CaseClause"
	case *ast.ChanType:
		return "ChanType"
	case *ast.CommClause:
		return "CommClause"
	case *ast.Comment:
		return "Comment"
	case *ast.CommentGroup:
		return "CommentGroup"
	case *ast.CompositeLit:
		return "CompositeLit"
	case *ast.DeclStmt:
		return "DeclStmt"
	case *ast.DeferStmt:
		return "DeferStmt"
	case *ast.Ellipsis:
		return "Ellipsis"
	case *ast.EmptyStmt:
		return "EmptyStmt"
	case *ast.ExprStmt:
		return "ExprStmt"
	case *ast.Field:
		return "Field"
	case *ast.FieldList:
		return "FieldList"
	case *ast.File:
		return "File"
	case *ast.ForStmt:
		return "ForStmt"
	case *ast.FuncDecl:
		return "FuncDecl"
	case *ast.FuncLit:
		return "FuncLit"
	case *ast.FuncType:
		return "FuncType"
	case *ast.GenDecl:
		return "GenDecl"
	case *ast.GoStmt:
		return "GoStmt"
	case *ast.Ident:
		return "Ident"
	case *ast.IfStmt:
		return "IfStmt"
	case *ast.ImportSpec:
		return "ImportSpec"
	case *ast.IncDecStmt:
		return "IncDecStmt"
	case *ast.IndexExpr:
		return "IndexExpr"
	case *ast.IndexListExpr:
		return "IndexListExpr"
	case *ast.InterfaceType:
		return "InterfaceType"
	case *ast.KeyValueExpr:
		return "KeyValueExpr"
	case *ast.LabeledStmt:
		return "LabeledStmt"
	case *ast.MapType:
		return "MapType"
	case *ast.Package:
		return "Package"
	case *ast.ParenExpr:
		return "ParenExpr"
	case *ast.RangeStmt:
		return "RangeStmt"
	case *ast.ReturnStmt:
		return "ReturnStmt"
	case *ast.SelectStmt:
		return "SelectStmt"
	case *ast.SelectorExpr:
		return "SelectorExpr"
	case *ast.SendStmt:
		return "SendStmt"
	case *ast.SliceExpr:
		return "SliceExpr"
	case *ast.StarExpr:
		return "StarExpr"
	case *ast.StructType:
		return "StructType"
	case *ast.SwitchStmt:
		return "SwitchStmt"
	case *ast.TypeAssertExpr:
		return "TypeAssertExpr"
	case *ast.TypeSpec:
		return "TypeSpec"
	case *ast.TypeSwitchStmt:
		return "TypeSwitchStmt"
	case *ast.UnaryExpr:
		return "UnaryExpr"
	case *ast.ValueSpec:
		return "ValueSpec"
	default:
		panic("unknown type")
	}
}
