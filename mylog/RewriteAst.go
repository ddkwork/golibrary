package mylog

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/tools/go/ast/astutil"
)

type handle struct {
	fileSet  *token.FileSet
	root     *ast.File
	path     string
	lines    []string
	lineInfo string
	line     int
	mod      string

	file *token.File
}

func newHandle(path string, noComments bool) *handle {
	if !strings.HasSuffix(path, ".go") {
		Check(fmt.Errorf("not a go file: %s", path))
	}
	fileSet := token.NewFileSet()
	h := &handle{
		fileSet:  fileSet,
		root:     Check2(parser.ParseFile(fileSet, path, nil, parser.ParseComments)),
		path:     path,
		lines:    nil,
		lineInfo: "",
		line:     0,
		mod:      "github.com/ddkwork/golibrary/mylog",
		file:     nil,
	}
	if noComments {
		h.removeComments()
	}
	buf := Check2(os.ReadFile(path))
	var lines []string
	scanner := bufio.NewReaderSize(bytes.NewReader(buf), 100*bufio.MaxScanTokenSize)
	for {
		line, _, e := scanner.ReadLine()
		if CheckEof(e) {
			break
		}
		fuck := ""
		if e != nil {
			fuck = path + ": " + e.Error() // todo not work
		}
		if fuck != "" {
			Check(fuck)
		}
		lines = append(lines, string(line))
	}
	h.lines = lines
	return h
}

func newCodeHandle(code string, noComments bool) *handle {
	path := "testHandle.go"
	fileSet := token.NewFileSet()
	h := &handle{
		fileSet:  fileSet,
		root:     Check2(parser.ParseFile(fileSet, path, code, parser.ParseComments)),
		path:     "testHandle.go",
		lines:    nil,
		lineInfo: "",
		line:     0,
		mod:      "github.com/ddkwork/golibrary/mylog",
		file:     nil,
	}
	if noComments {
		h.removeComments()
	}
	lines := []string{""}
	scanner := bufio.NewScanner(strings.NewReader(code))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	Check(scanner.Err())
	h.lines = lines
	return h
}

func (h *handle) removeComments() {
	Todo("https://github.com/Greyh4t/nocomment")
	newComments := make([]*ast.CommentGroup, 0) // can not be nil,why
	for _, group := range h.root.Comments {
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
	h.root.Comments = newComments
}

func FormatAllFiles()           { formatAllFiles(false, "") }
func FormatAllFilesNoComments() { formatAllFiles(true, "") }

var Skips = []string{
	`vendor`,
	`\gioview\`,
	`\gio\`,
	`\gio-cmd\`,
	`\gio-example\`,
	`\gio-x\`,
	`\toolbox\`,
	`\unison\`,
	`\ux\patch\`,
}

func formatAllFiles(noComments bool, path string) {
	if path == "" {
		path = "."
	}
	Check(filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		abs := Check2(filepath.Abs(path))
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
			newHandle(path, noComments).rewriteAst()
		}
		return err
	}))
}

func (h *handle) findEof(stmtType string) (hasEof bool) {
	if strings.Contains(h.lineInfo, "mylog.Check") {
		return false
	}
	eofString := []string{
		".EOF",
	}

	checked := false
	stop := false
	todoLineInfo := ""
	skips := []string{
		"CheckEof",
		"cmp.Equal",
		"io.EOF.Error()",
	}
	if h.line == 0 {
		Warning("for unit test", "because no ast parser")
		for i, s := range h.lines {
			if strings.Contains(s, "err ") {
				h.line = i
				break
			}
		}
	}

	for i, s := range h.lines[h.line:] {
		if stop {
			break
		}
		if strings.TrimSpace(s) == "}" {
			return
		}
		nextLines := h.lines[i+h.line:]
		const maxScan = 5
		for j, nextLine := range nextLines {
			if stop {
				break
			}
			if strings.TrimSpace(nextLine) == "}" {
				return
			}
			if j > maxScan {
				return
			}
			if i+h.line+j >= len(h.lines) {
				return
			}
			for _, eof := range eofString {
				if strings.Contains(nextLine, eof) {
					hasEof = true
					break
				}
			}
			for _, skip := range skips {
				if strings.Contains(nextLine, skip) {
					checked = true
					stop = true
					todoLineInfo = fmt.Sprintf(h.path+":%d ", i+j+h.line) + nextLine
					break
				}
			}
		}
	}
	if hasEof && !checked {
		Success("find eof in "+stmtType, "you should call CheckEof()  ", todoLineInfo)
		return true
	}
	return hasEof || checked
}

func (h *handle) rewriteAst() {
	needImport := false
	astutil.Apply(h.root, func(cursor *astutil.Cursor) bool {
		n := cursor.Node()
		if cursor.Node() == nil {
			return true
		}
		h.line = h.fileSet.Position(cursor.Node().Pos()).Line
		if h.line > len(h.lines)-1 {
			h.line = len(h.lines) - 1 // todo bug
			// Warning("line > len(lines) " + h.lineInfo)
		}
		h.lineInfo = fmt.Sprintf(h.path+":%d ", h.line) + h.lines[h.line]

		if h.line == 49 {
			// Todo("stop here for debug")
		}

		switch x := n.(type) {
		case *ast.IfStmt:
			if x.Else != nil {
				elseIf, ok := x.Else.(*ast.IfStmt)
				if ok {
					if isIfError(elseIf) {
						x.Else = nil
						Trace("else if err != nil", h.lineInfo)
						return true
					}
				}
			}

			if !isIfError(x) {
				return true
			}

			if x.Init == nil {
				Trace("if err != nil", h.lineInfo)
				cursor.Delete()
				return true
			}

			if !h.findEof("IfStmt") {
				return true
			}

			needImport = true

			if x.Init != nil && x.Cond != nil && x.Body != nil {
				exprStmt, ok := x.Init.(*ast.AssignStmt)
				if !ok {
					return true
				}
				if ok && len(exprStmt.Rhs) == 1 {
					switch len(exprStmt.Lhs) {
					case 1:
						expr := exprStmt.Rhs[0]
						newExpr := &ast.ExprStmt{
							X: &ast.CallExpr{
								Fun:  &ast.Ident{Name: "mylog.Check"},
								Args: []ast.Expr{expr},
							},
						}
						cursor.Replace(newExpr)
					case 2:
						l0, ok := exprStmt.Lhs[0].(*ast.Ident)
						if ok {
							if l0.Name == "_" {
								expr := exprStmt.Rhs[0]
								newExpr := &ast.ExprStmt{
									X: &ast.CallExpr{
										Fun:  &ast.Ident{Name: "mylog.Check2"},
										Args: []ast.Expr{expr},
									},
								}
								cursor.Replace(newExpr)
							}
						} else {
							last := len(exprStmt.Lhs) - 1
							newExpr := &ast.AssignStmt{
								Lhs: exprStmt.Lhs[:last],
								Tok: token.DEFINE,
								Rhs: []ast.Expr{&ast.CallExpr{
									Fun:  &ast.Ident{Name: "mylog.Check2"},
									Args: []ast.Expr{exprStmt.Rhs[0]},
								}},
							}
							cursor.Replace(newExpr)
						}

					case 3:
						last := len(exprStmt.Lhs) - 1
						newExpr := &ast.AssignStmt{
							Lhs: exprStmt.Lhs[:last],
							Tok: token.DEFINE,
							Rhs: []ast.Expr{&ast.CallExpr{
								Fun:  &ast.Ident{Name: "mylog.Check3"},
								Args: []ast.Expr{exprStmt.Rhs[0]},
							}},
						}
						cursor.Replace(newExpr)

					case 4:
						last := len(exprStmt.Lhs) - 1
						newExpr := &ast.AssignStmt{
							Lhs: exprStmt.Lhs[:last],
							Tok: token.DEFINE,
							Rhs: []ast.Expr{&ast.CallExpr{
								Fun:  &ast.Ident{Name: "mylog.Check4"},
								Args: []ast.Expr{exprStmt.Rhs[0]},
							}},
						}
						cursor.Replace(newExpr)

					case 5:
						last := len(exprStmt.Lhs) - 1
						newExpr := &ast.AssignStmt{
							Lhs: exprStmt.Lhs[:last],
							Tok: token.DEFINE,
							Rhs: []ast.Expr{&ast.CallExpr{
								Fun:  &ast.Ident{Name: "mylog.Check5"},
								Args: []ast.Expr{exprStmt.Rhs[0]},
							}},
						}
						cursor.Replace(newExpr)

					case 6:
						last := len(exprStmt.Lhs) - 1
						newExpr := &ast.AssignStmt{
							Lhs: exprStmt.Lhs[:last],
							Tok: token.DEFINE,
							Rhs: []ast.Expr{&ast.CallExpr{
								Fun:  &ast.Ident{Name: "mylog.Check6"},
								Args: []ast.Expr{exprStmt.Rhs[0]},
							}},
						}
						cursor.Replace(newExpr)

					case 7:
						last := len(exprStmt.Lhs) - 1
						newExpr := &ast.AssignStmt{
							Lhs: exprStmt.Lhs[:last],
							Tok: token.DEFINE,
							Rhs: []ast.Expr{&ast.CallExpr{
								Fun:  &ast.Ident{Name: "mylog.Check7"},
								Args: []ast.Expr{exprStmt.Rhs[0]},
							}},
						}
						cursor.Replace(newExpr)

					}
				}
			}
		case *ast.RangeStmt:
			// Todo("check loopVar and integer range " + lineInfo)
		case *ast.ForStmt:
			// Todo("check loopVar and integer range " + lineInfo)
		case *ast.DeferStmt:
			// Todo("check closer in defer " + lineInfo)

		case *ast.AssignStmt:
			if h.findEof("AssignStmt") {
				return true
			}
			if len(x.Rhs) > 1 {
				return true
			}

			last := len(x.Lhs) - 1
			lastIdent, ok := x.Lhs[last].(*ast.Ident)
			if !ok {
				return true
			}

			if lastIdent.Name != "err" && lastIdent.Name != "_" {
				return true
			}

			rightBack, ok := x.Rhs[0].(*ast.CallExpr)
			if !ok {
				return true
			}

			if lastIdent.Name == "_" {
				lastReturnType, ok := GetLastReturnType(x)
				if !ok {
					return true
				}

				if lastReturnType != "error" {
					return true
				}
			}

			needImport = true

			var newRight *ast.CallExpr

			switch len(x.Lhs) {
			case 1:
				newRight = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "mylog"},
						Sel: &ast.Ident{Name: "Check"},
					},
					Args: []ast.Expr{rightBack},
				}
				cursor.Replace(&ast.ExprStmt{X: newRight})

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
			x.Rhs[0] = newRight
			x.Lhs = x.Lhs[:last]
		}
		return true
	}, nil)

	if needImport {
		astutil.AddImport(h.fileSet, h.root, h.mod)
		Success("add import", h.path)
	}

	Call(func() {
		var buf bytes.Buffer
		Check(format.Node(&buf, h.fileSet, h.root))
		pattern := "var " + "err " + "error"
		s := strings.ReplaceAll(buf.String(), pattern, "")
		WriteGoFileWithDiff(h.path, []byte(s))
	})
}

func isIfError(stmt *ast.IfStmt) (isError bool) {
	binaryExpr, isBinary := stmt.Cond.(*ast.BinaryExpr)
	if !isBinary {
		return
	}
	if binaryExpr.Op == token.NEQ {
		if ident, isIdent := binaryExpr.X.(*ast.Ident); isIdent && ident.Name == "err" {
			if basicLit, isBlank := binaryExpr.Y.(*ast.Ident); isBlank && basicLit.Name == "nil" {
				return true
			}
		}
	}
	return
}

func (h *handle) debug(title string, n ast.Node) {
	Trace(title+" start", h.fileSet.Position(n.Pos()).String())
	Warning(title+" end", h.fileSet.Position(n.End()).String())
}

func (h *handle) Position(p token.Pos) token.Position {
	return h.file.PositionFor(p, false)
}

func (h *handle) removeLines(fromLine, toLine int) {
	for fromLine < toLine {
		h.file.MergeLine(fromLine)
		toLine--
	}
}

func (h *handle) removeLinesBetween(from, to token.Pos) {
	h.removeLines(h.Line(from)+1, h.Line(to))
}

func (h *handle) Line(p token.Pos) int {
	return h.Position(p).Line
}

func (h *handle) Offset(p token.Pos) int {
	return h.file.Offset(p)
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

func identEqual(expr ast.Expr, name string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == name
}

func mergeOk() {
	src := `
package main

import "fmt"

func main() {
	lines := strings.
Split(src, lineBreak)
}
`

	mergedSrc := mergeLines(src, ".", "\n")

	fmt.Println(mergedSrc)
}

func mergeLines(src, mergeToken, lineBreak string) string {
	lines := strings.Split(src, lineBreak)
	var mergedLines []string
	mergeBuffer := ""
	for _, line := range lines {
		if strings.HasSuffix(line, mergeToken) && mergeBuffer == "" {
			mergeBuffer = line
		} else {
			if mergeBuffer != "" {
				mergedLines = append(mergedLines, mergeBuffer+line)
				mergeBuffer = ""
			} else {
				mergedLines = append(mergedLines, line)
			}
		}
	}
	if mergeBuffer != "" {
		mergedLines = append(mergedLines, mergeBuffer)
	}
	return strings.Join(mergedLines, lineBreak)
}

func getAstLine() {
	src := ``
	fset := token.NewFileSet()
	node := Check2(parser.ParseFile(fset, "demo", src, 0))

	for _, group := range node.Decls {
		start := fset.Position(group.Pos())
		end := fset.Position(group.End())
		lines := strings.Split(src[start.Offset:end.Offset], "\n")
		for i, line := range lines {
			fmt.Printf("Line %d: %s\n", start.Line+i, line)
		}

	}
}

func getFuncAndMethod() {
	src := Check2(os.ReadFile("D:\\workspace\\workspace\\app\\widget\\TreeTable.go"))

	fset := token.NewFileSet()
	node := Check2(parser.ParseFile(fset, "D:\\workspace\\workspace\\app\\widget\\TreeTable.go", src, 0))

	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			body := fmt.Sprintf("%s\n", src[d.Pos()-1:d.End()-1])
			before, _, found := strings.Cut(body, "{")
			if found {
				if unicode.IsUpper(rune(before[5])) {
					println(before)
				}
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := typeSpec.Type.(*ast.FuncType); ok {
						body := fmt.Sprintf("%s\n", src[typeSpec.Pos()-1:typeSpec.End()-1])
						before, _, found := strings.Cut(body, "{")
						if found {
							if unicode.IsUpper(rune(before[5])) {
								println(before)
							}
						}
					}
				}
			}
		}
	}
}
