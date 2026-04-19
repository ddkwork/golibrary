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
	"regexp"
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
	deferProcessedIfStmts := make(map[token.Pos]bool)

	fnCall := func(n int) string {
		switch n {
		case 1:
			return "mylog.Check"
		default:
			return "mylog.Check" + strconv.Itoa(n)
		}
	}

	fnHandleAssign := func(x *ast.AssignStmt, e *Edit) {
		if deferProcessedIfStmts[x.Pos()] {
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
		_, isCallExpr := x.Rhs[0].(*ast.CallExpr)
		if lastIdent.Name == "_" && !isCallExpr {
			return
		}
		if lastIdent.Name == "_" {
			if expr, ok := x.Rhs[0].(*ast.CallExpr); ok {
				switch fun := expr.Fun.(type) {
				case *ast.SelectorExpr:
					if !returnsError(fun) {
						return
					}
				case *ast.Ident:
					lastReturnType, ok := getLastReturnType(x)
					if !ok {
						return
					}
					if lastReturnType != "error" {
						return
					}
				}
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

		var right strings.Builder
		tk := x.Tok.String()
		if left == "" {
			tk = ""
		}
		for i, v := range x.Rhs {
			right.WriteString(getNodeCode(v, fileSet, text))
			if i < len(x.Rhs)-1 {
				right.WriteString(",")
			}
		}

		ee := Edit{
			StartPos:   x.Pos(),
			EndPos:     x.End(),
			LineNumber: fileSet.Position(x.Pos()).Line,
			filePath:   fileSet.Position(x.Pos()).Filename + ":" + strconv.Itoa(fileSet.Position(x.Pos()).Line),
			NewContent: left + tk + fnCall(len(x.Lhs)) + "(" + right.String() + ")",
			edge:       edge(x),
			isContinue: false,
		}
		if e != nil { // AssignStmt in IfStmt
			ee = Edit{
				StartPos:   e.StartPos,
				EndPos:     e.EndPos,
				LineNumber: e.LineNumber,
				filePath:   e.filePath,
				NewContent: left + tk + fnCall(len(x.Lhs)) + "(" + right.String() + ")",
				edge:       e.edge,
				isContinue: false,
			}
		}
		Replaces = append(Replaces, ee)
	}

	skipAssign := false
	isContinue := false

	getIfElseChainEnd := func(ifStmt *ast.IfStmt) token.Pos {
		current := ifStmt
		for {
			if current.Else == nil {
				return current.Body.End()
			}
			if elseIf, ok := current.Else.(*ast.IfStmt); ok {
				current = elseIf
			} else if block, ok := current.Else.(*ast.BlockStmt); ok {
				return block.End()
			} else {
				return current.Body.End()
			}
		}
	}

	fnHandleIfElseChain := func(rootIfStmt *ast.IfStmt) {
		var checkStmts []string
		var remainingIfContent string
		var finalElseStmts []string

		currentIf := rootIfStmt

		for currentIf != nil {
			if assignStmt, ok := currentIf.Init.(*ast.AssignStmt); ok {
				last := len(assignStmt.Lhs) - 1
				if lastIdent, ok := assignStmt.Lhs[last].(*ast.Ident); ok && lastIdent.Name == "err" {
					left := ""
					for i, v := range assignStmt.Lhs {
						c := getNodeCode(v, fileSet, text)
						if c == "err" || c == "_" {
							continue
						}
						if i < len(assignStmt.Lhs)-1 {
							left += c + ","
						}
					}
					left = strings.TrimRight(left, ",")

					var right strings.Builder
					for i, v := range assignStmt.Rhs {
						right.WriteString(getNodeCode(v, fileSet, text))
						if i < len(assignStmt.Rhs)-1 {
							right.WriteString(",")
						}
					}

					tk := assignStmt.Tok.String()
					if left == "" {
						tk = ""
					}

					checkStmts = append(checkStmts, left+tk+fnCall(len(assignStmt.Lhs))+"("+right.String()+")")
					deferProcessedIfStmts[currentIf.Pos()] = true
					deferProcessedIfStmts[assignStmt.Pos()] = true

					if elseStmt, ok := currentIf.Else.(*ast.IfStmt); ok {
						currentIf = elseStmt
					} else if blockStmt, ok := currentIf.Else.(*ast.BlockStmt); ok {
						for _, stmt := range blockStmt.List {
							finalElseStmts = append(finalElseStmts, getNodeCode(stmt, fileSet, text))
						}
						break
					} else {
						break
					}
				} else {
					if elseStmt, ok := currentIf.Else.(*ast.IfStmt); ok {
						currentIf = elseStmt
					} else if blockStmt, ok := currentIf.Else.(*ast.BlockStmt); ok {
						for _, stmt := range blockStmt.List {
							finalElseStmts = append(finalElseStmts, getNodeCode(stmt, fileSet, text))
						}
						break
					} else {
						break
					}
				}
			} else {
				cond := getNodeCode(currentIf.Cond, fileSet, text)
				var bodyStmts []string
				for _, stmt := range currentIf.Body.List {
					c := getNodeCode(stmt, fileSet, text)
					if strings.HasPrefix(c, "log.") {
						arg := strings.TrimPrefix(c, "log.Fatal(")
						arg = strings.TrimSuffix(arg, ")")
						bodyStmts = append(bodyStmts, "mylog.Check("+arg+")")
					} else {
						bodyStmts = append(bodyStmts, c)
					}
				}
				remainingIfContent = "if " + cond + " {\n"
				for _, s := range bodyStmts {
					remainingIfContent += "\t\t" + s + "\n"
				}
				remainingIfContent += "\t}"
				if blockStmt, ok := currentIf.Else.(*ast.BlockStmt); ok {
					for _, stmt := range blockStmt.List {
						finalElseStmts = append(finalElseStmts, getNodeCode(stmt, fileSet, text))
					}
				}
				break
			}
		}

		var newContent strings.Builder
		for _, s := range checkStmts {
			newContent.WriteString("\t" + s + "\n")
		}
		if remainingIfContent != "" {
			newContent.WriteString("\t" + remainingIfContent + "\n")
		}
		for _, s := range finalElseStmts {
			newContent.WriteString("\tmylog.Check(" + s + ")\n")
		}

		endPos := getIfElseChainEnd(rootIfStmt)
		Replaces = append(Replaces, Edit{
			StartPos:   rootIfStmt.Pos(),
			EndPos:     endPos,
			NewContent: newContent.String(),
		})
	}

	for n := range ast.Preorder(file) {
		switch x := n.(type) {
		case *ast.IfStmt: // if err := backendConn.Close(); err != nil {
			if deferProcessedIfStmts[x.Pos()] {
				continue
			}

			if x.Else != nil && x.Init != nil {
				if assignStmt, ok := x.Init.(*ast.AssignStmt); ok {
					last := len(assignStmt.Lhs) - 1
					if lastIdent, ok := assignStmt.Lhs[last].(*ast.Ident); ok && lastIdent.Name == "err" {
						fnHandleIfElseChain(x)
						continue
					}
				}
			}

			for ifStmt := range findIfErrNotNil(n, deferProcessedIfStmts) {
				isOneWorkCode := false // if 块内部语句只有1句，没有其他业务逻辑,则直接替换为mylog.Check(业务逻辑)
				for i, stmt := range ifStmt.Body.List {
					if len(ifStmt.Body.List) == 1 {
						isOneWorkCode = true
						break
					}
					if i > 1 {
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
						case strings.HasPrefix(c, "os.Exit(1)"):
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
					hasBusinessLogic := false
					for _, stmt := range ifStmt.Body.List {
						switch stmt.(type) {
						case *ast.AssignStmt, *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt:
							hasBusinessLogic = true
						}
					}
					if hasBusinessLogic {
						break
					}
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
					// 处理 if err != nil { ... } 的情况（没有 Init）
					// 检查 if 块中是否只有简单的错误处理（log/return/panic）
					isSimpleErrorCheck := false
					if len(ifStmt.Body.List) == 1 {
						switch row := ifStmt.Body.List[0].(type) {
						case *ast.ExprStmt:
							c := getNodeCode(row, fileSet, text)
							if strings.HasPrefix(c, "log.") || c == "panic(err)" {
								isSimpleErrorCheck = true
							}
						case *ast.ReturnStmt:
							c := getNodeCode(row, fileSet, text)
							if strings.Contains(c, ", err") || strings.HasSuffix(c, ", err)") || c == "return err" {
								isSimpleErrorCheck = true
							}
						}
					}
					if isSimpleErrorCheck {
						// 删除整个 if 语句
						Replaces = append(Replaces, Edit{
							StartPos:   ifStmt.Pos(),
							EndPos:     ifStmt.End(),
							LineNumber: fileSet.Position(ifStmt.Pos()).Line,
							filePath:   fileSet.Position(ifStmt.Pos()).Filename + ":" + strconv.Itoa(fileSet.Position(ifStmt.Pos()).Line),
							NewContent: "",
							edge:       edge(ifStmt),
							isContinue: false,
						})
						break
					}
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
			if deferProcessedIfStmts[x.Pos()] {
				continue
			}
			if skipAssign {
				skipAssign = false
				continue
			}
			fnHandleAssign(x, nil)
		case *ast.DeferStmt:
			c := getNodeCode(x.Call, fileSet, text)
			switch fun := x.Call.Fun.(type) {
			case *ast.FuncLit:
				for _, stmt := range fun.Body.List {
					if ifStmt, ok := stmt.(*ast.IfStmt); ok {
						if ifStmt.Init == nil {
							isSimpleErrorCheck := false
							if len(ifStmt.Body.List) == 1 {
								switch row := ifStmt.Body.List[0].(type) {
								case *ast.ExprStmt:
									c := getNodeCode(row, fileSet, text)
									if strings.HasPrefix(c, "log.") || c == "panic(err)" {
										isSimpleErrorCheck = true
									}
								case *ast.ReturnStmt:
									c := getNodeCode(row, fileSet, text)
									if strings.Contains(c, ", err") || strings.HasSuffix(c, ", err)") || c == "return err" {
										isSimpleErrorCheck = true
									}
								}
							}
							if isSimpleErrorCheck {
								Replaces = append(Replaces, Edit{
									StartPos:   ifStmt.Pos(),
									EndPos:     ifStmt.End(),
									LineNumber: fileSet.Position(ifStmt.Pos()).Line,
									filePath:   fileSet.Position(ifStmt.Pos()).Filename + ":" + strconv.Itoa(fileSet.Position(ifStmt.Pos()).Line),
									NewContent: "",
									edge:       edge(ifStmt),
									isContinue: false,
								})
							}
						} else {
							if assignStmt, ok := ifStmt.Init.(*ast.AssignStmt); ok {
								last := len(assignStmt.Lhs) - 1
								if lastIdent, ok := assignStmt.Lhs[last].(*ast.Ident); ok {
									if lastIdent.Name == "err" {
										isSimpleErrorCheck := false
										if len(ifStmt.Body.List) == 1 {
											switch row := ifStmt.Body.List[0].(type) {
											case *ast.ExprStmt:
												c := getNodeCode(row, fileSet, text)
												if strings.HasPrefix(c, "log.") {
													isSimpleErrorCheck = true
												}
											}
										}
										if isSimpleErrorCheck {
											left := ""
											for i, v := range assignStmt.Lhs {
												c := getNodeCode(v, fileSet, text)
												if c == "err" || c == "_" {
													continue
												}
												if i < len(assignStmt.Lhs)-1 {
													left += c + ","
												}
											}
											left = strings.TrimRight(left, ",")

											var right strings.Builder
											tk := assignStmt.Tok.String()
											if left == "" {
												tk = ""
											}
											for i, v := range assignStmt.Rhs {
												right.WriteString(getNodeCode(v, fileSet, text))
												if i < len(assignStmt.Rhs)-1 {
													right.WriteString(",")
												}
											}

											newContent := left + tk + fnCall(len(assignStmt.Lhs)) + "(" + right.String() + ")"
											deferProcessedIfStmts[ifStmt.Pos()] = true
											deferProcessedIfStmts[assignStmt.Pos()] = true
											Replaces = append(Replaces, Edit{
												StartPos:   ifStmt.Pos(),
												EndPos:     ifStmt.End(),
												LineNumber: fileSet.Position(ifStmt.Pos()).Line,
												filePath:   fileSet.Position(ifStmt.Pos()).Filename + ":" + strconv.Itoa(fileSet.Position(ifStmt.Pos()).Line),
												NewContent: newContent,
												edge:       edge(ifStmt) + " # " + edge(assignStmt),
												isContinue: false,
											})
										}
									}
								}
							}
						}
					}
				}
			case *ast.CallExpr:
				if strings.Contains(c, "Close") {
					skip := false
					for pos := range deferProcessedIfStmts {
						if x.Pos() > pos {
							skip = true
							break
						}
					}
					if !skip {
						Replaces = append(Replaces, Edit{
							StartPos:   x.Pos(),
							EndPos:     x.End(),
							LineNumber: fileSet.Position(x.Pos()).Line,
							filePath:   fileSet.Position(x.Pos()).Filename + ":" + strconv.Itoa(fileSet.Position(x.Pos()).Line),
							NewContent: "mylog.Check(" + c + ")",
							edge:       edge(x),
							isContinue: false,
						})
					}
				}
			}
		}
	}
	return simplifyNestedChecks(Apply(text, Replaces))
}

func getNodeCode(astNode ast.Node, f *token.FileSet, code string) string {
	c := code[f.Position(astNode.Pos()).Offset:f.Position(astNode.End()).Offset]
	c = strings.TrimSpace(c)
	return c
}

// todo 这种是否应该删除 if err != nil && !errors.Is(err, fs.ErrNotExist) {
func findIfErrNotNil(n ast.Node, deferProcessedIfStmts map[token.Pos]bool) iter.Seq[*ast.IfStmt] {
	return func(yield func(*ast.IfStmt) bool) {
		if stmt, ok := n.(*ast.IfStmt); ok {
			if deferProcessedIfStmts[stmt.Pos()] {
				return
			}
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

func returnsError(expr *ast.SelectorExpr) bool {
	knownErrorReturningFuncs := map[string]bool{
		"os.ReadFile":                 true,
		"os.WriteFile":                true,
		"os.Mkdir":                    true,
		"os.MkdirAll":                 true,
		"os.Remove":                   true,
		"os.RemoveAll":                true,
		"os.Rename":                   true,
		"os.Open":                     true,
		"os.Create":                   true,
		"os.OpenFile":                 true,
		"os.Stat":                     true,
		"os.Lstat":                    true,
		"filepath.Glob":               true,
		"filepath.Walk":               true,
		"filepath.Abs":                true,
		"filepath.EvalSymlinks":       true,
		"io.Copy":                     true,
		"io.CopyN":                    true,
		"io.CopyBuffer":               true,
		"io.ReadAll":                  true,
		"io.ReadFull":                 true,
		"io.ReadAtLeast":              true,
		"io.WriteString":              true,
		"io/ioutil.ReadFile":          true,
		"io/ioutil.WriteFile":         true,
		"io/ioutil.ReadDir":           true,
		"io/ioutil.TempDir":           true,
		"io/ioutil.TempFile":          true,
		"exec.Command":                false,
		"exec.Command.Output":         true,
		"exec.Command.Run":            true,
		"exec.Command.CombinedOutput": true,
		"exec.Command.Start":          true,
		"exec.Command.Wait":           true,
		"net.Dial":                    true,
		"net.DialTimeout":             true,
		"net.Listen":                  true,
		"http.Get":                    true,
		"http.Post":                   true,
		"http.PostForm":               true,
		"json.Marshal":                true,
		"json.Unmarshal":              true,
		"xml.Marshal":                 true,
		"xml.Unmarshal":               true,
		"fmt.Sscanf":                  true,
		"fmt.Scanf":                   true,
		"fmt.Scan":                    true,
		"fmt.Scanln":                  true,
		"fmt.Fscanf":                  true,
		"fmt.Fscan":                   true,
		"fmt.Fscanln":                 true,
		"fmt.Sprint":                  false,
		"fmt.Sprintf":                 false,
		"fmt.Println":                 false,
		"fmt.Printf":                  false,
	}
	var pkg, name string
	if ident, ok := expr.X.(*ast.Ident); ok {
		pkg = ident.Name
		name = expr.Sel.Name
		fullName := pkg + "." + name
		if v, ok := knownErrorReturningFuncs[fullName]; ok {
			return v
		}
	}
	return false
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

func simplifyNestedChecks(text string) string {
	endsWithNewline := strings.HasSuffix(text, "\n")
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if !strings.Contains(line, "mylog.Check") {
			continue
		}
		if !hasDirectNestedCheck(line) {
			continue
		}
		innerArg := findInnermostCheckArg(line)
		if innerArg != "" {
			var indent strings.Builder
			for _, c := range line {
				if c == ' ' || c == '\t' {
					indent.WriteString(string(c))
				} else {
					break
				}
			}
			lines[i] = indent.String() + "mylog.Check(" + innerArg + ")"
		}
	}
	result := strings.Join(lines, "\n")
	if !endsWithNewline && strings.HasSuffix(result, "\n") {
		result = strings.TrimSuffix(result, "\n")
	}
	return result
}

func hasDirectNestedCheck(line string) bool {
	checkPattern := regexp.MustCompile(`mylog\.Check\d*\s*\(`)
	locs := checkPattern.FindAllStringIndex(line, -1)
	if len(locs) < 2 {
		return false
	}
	for _, loc := range locs {
		parenStart := loc[1] - 1
		for parenStart < len(line) && line[parenStart] != '(' {
			parenStart++
		}
		if parenStart >= len(line) {
			continue
		}
		depth := 1
		j := parenStart + 1
		for j < len(line) && depth > 0 {
			if line[j] == '(' {
				depth++
			} else if line[j] == ')' {
				depth--
			}
			j++
		}
		arg := line[parenStart+1 : j-1]
		arg = strings.TrimSpace(arg)
		if strings.HasPrefix(arg, "mylog.Check") {
			return true
		}
	}
	return false
}

func findInnermostCheckArg(line string) string {
	checkPattern := regexp.MustCompile(`mylog\.Check\d*`)
	locs := checkPattern.FindAllStringIndex(line, -1)
	if len(locs) == 0 {
		return ""
	}
	lastLoc := locs[len(locs)-1]
	lastCheckEnd := lastLoc[1]
	if lastCheckEnd >= len(line) || line[lastCheckEnd] != '(' {
		return ""
	}
	argStart := lastCheckEnd + 1
	depth := 1
	for i := argStart; i < len(line); i++ {
		if line[i] == '(' {
			depth++
		}
		if line[i] == ')' {
			depth--
			if depth == 0 {
				return line[argStart:i]
			}
		}
	}
	return ""
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
		replaces[i].filePath = " " + filepath.ToSlash(replaces[i].filePath) + " "
		if i > 0 && strings.Contains(r.NewContent, "continue") && replaces[i-1].NewContent != "" {
			replaces[i-1].isContinue = true
		}
	}
	sort.Slice(replaces, func(i, j int) bool {
		return replaces[i].StartPos > replaces[j].StartPos
	})
	for _, r := range replaces {
		if r.isContinue {
			continue
		}
		if r.StartPos > r.EndPos {
			continue
		}
		startPos := int(r.StartPos) - 1
		endPos := int(r.EndPos) - 1
		if startPos < 0 || endPos < 0 || startPos > len(text) || endPos > len(text) || startPos >= endPos {
			continue
		}
		text = text[:startPos] + r.NewContent + text[endPos:]
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
