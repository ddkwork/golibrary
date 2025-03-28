package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ddkwork/golibrary/stream"

	"github.com/ddkwork/golibrary/mylog"
)

func main() {
	outPath := "language.i18n"
	kv := make(map[string]string)
	fileSet := token.NewFileSet()
	root := mylog.Check2(filepath.Abs("."))
	mylog.Check(filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if !fi.IsDir() && filepath.Ext(path) == ".go" {
			fmt.Println(path)
			var file *ast.File
			if file = mylog.Check2(parser.ParseFile(fileSet, path, nil, 0)); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			const (
				LookForPackageState = iota
				LookForTextCallState
				LookForParameterState
			)
			state := LookForPackageState
			ast.Inspect(file, func(node ast.Node) bool {
				switch x := node.(type) {
				case *ast.Ident:
					switch state {
					case LookForPackageState:
						if x.Name == "i18n" {
							state = LookForTextCallState
						}
					case LookForTextCallState:
						if x.Name == "Text" {
							state = LookForParameterState
						} else {
							state = LookForPackageState
						}
					default:
						state = LookForPackageState
					}
				case *ast.BasicLit:
					if state == LookForParameterState {
						if x.Kind == token.STRING {
							var v string
							if v = mylog.Check2(strconv.Unquote(x.Value)); err != nil {
								fmt.Fprintln(os.Stderr, err)
							} else {
								kv[v] = v
							}
						}
					}
					state = LookForPackageState
				case nil:
				default:
					state = LookForPackageState
				}
				return true
			})
		}
		return nil
	}))

	keys := make([]string, 0, len(kv))
	for key := range kv {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return stream.NaturalLess(keys[i], keys[j], true)
	})
	out := mylog.Check2(os.OpenFile(outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644))

	fmt.Fprintf(out, `# Generated on %v
#
# Key-value pairs are defined as one or more lines prefixed with "k:" for the
# key, followed by one or more lines prefixed with "v:" for the value. These
# prefixes are then followed by a quoted string, using escaping rules for Go
# strings where needed. When two or more lines are present in a row, they will
# be concatenated together with an intervening \n character.
#
# Do NOT modify the 'k' values. They are the values as seen in the code.
#
# Replace the 'v' values with the appropriate translation.
`, time.Now().Format(time.RFC1123))
	for _, key := range keys {
		fmt.Fprintln(out)
		for s := range strings.Lines(key) {
			mylog.Check2(fmt.Fprintf(out, "k:%q\n", s))
			mylog.Check2(fmt.Fprintf(out, "v:%q\n", s))
		}
	}
	mylog.Check(out.Close())
}
