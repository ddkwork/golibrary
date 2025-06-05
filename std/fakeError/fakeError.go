package fakeError

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/ddkwork/golibrary/std/mylog"
)

func Walk(path string, removeComments ...bool) {
	if path == "" {
		path = "."
	}
	mylog.Check(filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info != nil {
			if info.IsDir() && info.Name() == "vendor" {
				return nil
			}
		}
		abs := mylog.Check2(filepath.Abs(path))
		abs = filepath.ToSlash(abs)
		if filepath.Ext(abs) == ".go" {
			if !strings.HasSuffix(path, ".go") {
				mylog.Check(fmt.Errorf("not a go file: %s", path))
			}
			fileSet := token.NewFileSet()
			file := mylog.Check2(parser.ParseFile(fileSet, path, nil, parser.ParseComments))
			if len(removeComments) > 0 {
				RemoveComments(file)
			}
			s := handle(fileSet, file, mylog.Check2(os.ReadFile(path)))
			mylog.WriteGoFile(path, s)
		}
		return err
	}))
}

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

func handle[T string | []byte](fileSet *token.FileSet, file *ast.File, b T) string {
	text := string(b)
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
			lastReturnType, ok := getLastReturnType(x)
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
		// 只有一个返回值需要删除 ,
		left = strings.TrimRight(left, ",")

		right := ""
		tk := x.Tok.String()
		if left == "" {
			tk = ""
		}
		for i, v := range x.Rhs {
			right += getNodeCode(v, fileSet, text)
			if i < len(x.Rhs)-1 {
				right += ","
			}
		}

		ee := Edit{
			StartPos:   x.Pos(),
			EndPos:     x.End(),
			LineNumber: fileSet.Position(x.Pos()).Line,
			filePath:   fileSet.Position(x.Pos()).Filename + ":" + strconv.Itoa(fileSet.Position(x.Pos()).Line),
			NewContent: left + tk + fnCall(len(x.Lhs)) + "(" + right + ")",
			edge:       edge(x),
			isContinue: false,
		}
		if e != nil { // AssignStmt in IfStmt
			ee = Edit{
				StartPos:   e.StartPos,
				EndPos:     e.EndPos,
				LineNumber: e.LineNumber,
				filePath:   e.filePath,
				NewContent: left + tk + fnCall(len(x.Lhs)) + "(" + right + ")",
				edge:       e.edge,
				isContinue: false,
			}
		}
		Replaces = append(Replaces, ee)
	}

	skipAssign := false
	isContinue := false
	for n := range ast.Preorder(file) {
		switch x := n.(type) {
		case *ast.IfStmt: // if err := backendConn.Close(); err != nil {
			for ifStmt := range findIfErrNotNil(n) {
				isOneWorkCode := false // if 块内部语句只有1句，没有其他业务逻辑,则直接替换为mylog.Check(业务逻辑)
				for i, stmt := range ifStmt.Body.List {
					if len(ifStmt.Body.List) == 1 {
						isOneWorkCode = true
						break
					}
					if i > 1 {
						mylog.Warning("if 块内部语句超过1句:" + getNodeCode(stmt, fileSet, text) + " " + fileSet.Position(stmt.Pos()).Filename + ":" + strconv.Itoa(fileSet.Position(stmt.Pos()).Line))
						break
					}
					switch row := stmt.(type) {
					case *ast.BranchStmt:
						if getNodeCode(row, fileSet, text) == "continue" {
							isOneWorkCode = true
							isContinue = true
						}
					case *ast.ExprStmt:
						c := getNodeCode(row, fileSet, text)
						switch {
						case c == "panic(err)":
							isOneWorkCode = true

							// drivercodegen-master.zip
							//	if err := makeExeVcxprojFile(); err != nil {
							//		log.Println("[-] Failed to makeExeVcxprojFile....")
							//		return
							//	}
						// case strings.HasPrefix(c, "log.") && strings.HasSuffix(c, "(err)"):
						case strings.HasPrefix(c, "log."):
							isOneWorkCode = true
						}
					case *ast.ReturnStmt:
						c := getNodeCode(row, fileSet, text)
						switch {
						case c == "return nil, nil, err":
							isOneWorkCode = true
						case strings.Contains(c, ", err"):
							isOneWorkCode = true
						case strings.HasSuffix(c, ", err)"): //		return nil, nil, fmt.Errorf("failed to get user info: %w", err)
							isOneWorkCode = true
						}
					}
				}
				if strings.HasPrefix(getNodeCode(ifStmt, fileSet, text), "if err != nil {") {
					b := `if err != nil {
					mylog.CheckIgnore(err)
					continue
				}`
					e := Edit{
						StartPos:   ifStmt.Pos(),
						EndPos:     ifStmt.End(),
						LineNumber: fileSet.Position(ifStmt.Pos()).Line,
						filePath:   fileSet.Position(ifStmt.Pos()).Filename + ":" + strconv.Itoa(fileSet.Position(ifStmt.Pos()).Line),
						NewContent: "",
						edge:       edge(ifStmt),
						isContinue: false,
					}
					if isContinue {
						e.NewContent = b
						isContinue = false
					}
					Replaces = append(Replaces, e)
					break
				}
				if ifStmt.Init == nil {
					break
				}
				if stmt, ok := ifStmt.Init.(*ast.AssignStmt); ok {
					if isOneWorkCode {
						fnHandleAssign(stmt, &Edit{
							StartPos:   ifStmt.Pos(),
							EndPos:     ifStmt.End(),
							LineNumber: fileSet.Position(ifStmt.Pos()).Line,
							filePath:   fileSet.Position(ifStmt.Pos()).Filename + ":" + strconv.Itoa(fileSet.Position(ifStmt.Pos()).Line),
							NewContent: "",
							edge:       edge(ifStmt) + " # " + edge(stmt),
							isContinue: false,
						})
						skipAssign = true
						break
					}
				}
			}
		case *ast.AssignStmt: // backendConn, err := net.DialTimeout("tcp", forwardTarget, 5*time.Second)
			if skipAssign {
				skipAssign = false
				continue
			}
			fnHandleAssign(x, nil)
		case *ast.DeferStmt: // todo 处理多个返回值，以及返回值检查,e.Obj.Decl.(*ast.FuncDecl)取不出来的返回类型的，nil
			mylog.Todo("defer func() { mylog.Check(backendConn.Close()) }()")
			break
			c := getNodeCode(x.Call, fileSet, text)
			switch x := x.Call.Fun.(type) {
			case *ast.FuncLit: // todo 改成检测是否实现io.Closer接口,io.NopCloser
				if strings.Contains(c, "Close") {
					Replaces = append(Replaces, Edit{
						StartPos:   x.Pos(),
						EndPos:     x.End(),
						LineNumber: fileSet.Position(x.Pos()).Line,
						filePath:   fileSet.Position(x.Pos()).Filename + ":" + strconv.Itoa(fileSet.Position(x.Pos()).Line),
						NewContent: "mylog.Check(" + getNodeCode(x, fileSet, text) + ")",
						edge:       edge(x),
						isContinue: false,
					})
				}
			case *ast.CallExpr:
				if strings.Contains(c, "Close") {
					Replaces = append(Replaces, Edit{
						StartPos:   x.Pos(),
						EndPos:     x.End(),
						LineNumber: fileSet.Position(x.Pos()).Line,
						filePath:   fileSet.Position(x.Pos()).Filename + ":" + strconv.Itoa(fileSet.Position(x.Pos()).Line),
						NewContent: "mylog.Check(" + getNodeCode(x, fileSet, text) + ")",
						edge:       edge(x),
						isContinue: false,
					})
				}
			}
		}
	}
	return Apply(text, Replaces)
}

func getNodeCode(astNode ast.Node, f *token.FileSet, code string) string {
	c := code[f.Position(astNode.Pos()).Offset:f.Position(astNode.End()).Offset]
	c = strings.TrimSpace(c)
	// mylog.Json("code dump", c)
	return c
}

// todo 这种是否应该删除 if err != nil && !errors.Is(err, fs.ErrNotExist) {
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

func getLastReturnType(assignStmt *ast.AssignStmt) (lastReturnType string, b bool) {
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
	StartPos, EndPos token.Pos
	LineNumber       int
	filePath         string
	NewContent       string
	edge             string
	isContinue       bool
}

// Apply 按起始位置从大到小排序,即从后往前替换，避免处理过程中坐标变化
// 单行:左+新内容+右
// 多行:前+新内容+后
// StartPos:要删除的第一个字符的偏移,这个通过单元测试了，不要改
// end:不这样连续替换后没有换行，两个mycheck在一行导致语法错误
func Apply(text string, replaces []Edit) string {
	if len(replaces) == 0 {
		return text
	}
	for i, r := range replaces {
		replaces[i].filePath = " " + filepath.ToSlash(replaces[i].filePath) + " " // 为了使行号可点击定位到文件
		if strings.Contains(r.NewContent, "continue") && replaces[i-1].NewContent != "" {
			replaces[i-1].isContinue = true
		}
	}
	mylog.Struct(replaces)
	sort.Slice(replaces, func(i, j int) bool {
		return replaces[i].StartPos > replaces[j].StartPos
	})
	for _, r := range replaces {
		if r.isContinue {
			continue
		}
		if r.StartPos > r.EndPos {
			panic("起始位置大于终止位置")
		}
		text = text[:r.StartPos-1] + r.NewContent + text[r.EndPos-1:]
	}
	text = strings.ReplaceAll(text, `var err error`, "")
	lib := "github.com/ddkwork/golibrary/std/mylog"
	if !strings.Contains(text, lib) {
		text = strings.Replace(text, `import (`, `import (
			"github.com/ddkwork/golibrary/std/mylog"`, 1)
	}
	return text
}

type edgeType interface {
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
		// *ast.Package |
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

// var (
// 	anyType    = reflect.TypeFor[any]()
// 	stringType = reflect.TypeFor[string]()
// 	bytesType  = reflect.TypeFor[[]byte]()
// 	byteType   = reflect.TypeFor[byte]()
// )

func edge[T edgeType](n T) string {
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
	// case *ast.Package:
	// 	return "Package"
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
