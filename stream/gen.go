package stream

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
)

type GeneratedFile struct{ *Buffer }

func NewGeneratedFile() *GeneratedFile {
	return &GeneratedFile{
		Buffer: NewBuffer(""),
	}
}

func (g *GeneratedFile) P(v ...any) {
	for _, x := range v {
		mylog.Check2(fmt.Fprint(g, x))
	}
	mylog.Check2(fmt.Fprintln(g))
}

func (g *GeneratedFile) Enum(kindName string, kinds, tooltip []string) {
	mylog.Call(func() {
		for i, kind := range kinds {
			kinds[i] = ToCamelUpper(kind, false)
		}
		kindNameUpper := ToCamelUpper(kindName, false)
		InvalidKind := "Invalid" + kindNameUpper + "Kind"

		g.P("package ", GetPackageName())
		g.P("")

		g.P("import (")
		g.P(strconv.Quote("golang.org/x/exp/constraints"))
		g.P(strconv.Quote("strings"))
		g.P(")")
		g.P("")

		g.P("// Code generated by GeneratedFile enum - DO NOT EDIT.")

		g.P("")
		g.P("type ", kindNameUpper, "Kind byte")
		g.P("")
		g.P("const (")
		for i, kind := range kinds {
			if i == 0 {
				g.P(" ", kind, "Kind", " ", kindNameUpper, "Kind = iota")
				continue
			}
			g.P(" "+kind, "Kind")
		}
		g.P(" ", InvalidKind)

		g.P(")")
		g.P("")

		g.P("func ConvertInteger2", kindNameUpper, "Kind[T constraints.Integer](v T)", kindNameUpper, "Kind  {")
		g.P("  return ", kindNameUpper, "Kind(v)")
		g.P("}")
		g.P("")

		g.P("func (k ", kindNameUpper, "Kind) AssertKind(kinds string) ", kindNameUpper, "Kind {")
		g.P(" for _, kind := range k.Kinds() {")
		// g.P("  if kinds == kind.String() {")
		g.P("  if strings.ToLower(kinds) == strings.ToLower(kind.String()) {")
		g.P("   return kind")
		g.P("  }")
		g.P(" }")
		g.P(" return " + InvalidKind)
		g.P("}")
		g.P("")

		g.P("func (k ", kindNameUpper, "Kind) String() string {")
		g.P(" switch k {")
		for _, kind := range kinds {
			g.P(" case ", kind, "Kind:")
			g.P("  return ", strconv.Quote(kind))
		}
		g.P(" default:")
		g.P("  return ", strconv.Quote(InvalidKind))
		g.P(" }")
		g.P("}")
		g.P("")

		if len(tooltip) > 0 {
			g.P("func (k ", kindNameUpper, "Kind) Tooltip() string {")
			g.P(" switch k {")
			for i, kind := range kinds {
				g.P(" case ", kind, "Kind:")
				g.P("  return ", strconv.Quote(tooltip[i]))
			}
			g.P(" default:")
			g.P("  return ", strconv.Quote(InvalidKind))
			g.P(" }")
			g.P("}")
			g.P("")
		}

		g.P("func (k ", kindNameUpper, "Kind) Keys() []string {")
		g.P(" return []string{")
		for _, kind := range kinds {
			g.P("  ", strconv.Quote(kind), ",")
		}
		g.P("  ", strconv.Quote(InvalidKind), ",")
		g.P(" }")
		g.P("}")
		g.P("")

		g.P("func (k ", kindNameUpper, "Kind) Kinds() []", kindNameUpper, "Kind {")
		g.P(" return []", kindNameUpper, "Kind{")
		for _, kind := range kinds {
			g.P("  ", kind, "Kind,")
		}
		g.P("  ", InvalidKind, ",")
		g.P(" }")
		g.P("}")
		g.P("")

		g.P("func (k ", kindNameUpper, "Kind) SvgFileName() string {")
		g.P(" switch k {")
		for _, kind := range kinds {
			g.P(" case ", kind, "Kind:")
			g.P("  return ", strconv.Quote(kind))
		}
		g.P("	default:")
		g.P("  return ", strconv.Quote(InvalidKind))
		g.P(" }")
		g.P("}")
		g.P("")
		WriteGoFile(kindName+"_enum_gen.go", g.Buffer)
		g.Reset()
	})
}

func (g *GeneratedFile) ReadTemplates(path, pkg string) {
	pkg = strings.TrimSuffix(pkg, "_test")
	mylog.Call(func() {
		s := NewBuffer("package " + pkg + "_test")
		s.NewLine()
		s.WriteStringLn(`
   import (
			"github.com/ddkwork/golibrary/stream"
		)`)
		s.NewLine()

		s.WriteStringLn("func generateIR(path string,callBack func(b*stream.Buffer)) {")
		s.WriteStringLn("g := stream.NewGeneratedFile()")
		lines := NewBuffer(path).ToLines()
		for _, line := range lines {
			needNewLine := false
			if line == "" {
				needNewLine = true
			}
			if strings.HasPrefix(line, "package") {
				packageName := "package " + GetPackageName()
				line = "g.P(" + strconv.Quote(packageName) + ")"
			} else {
				line = strings.ReplaceAll(line, "\t", " ")
				line = "g.P(" + strconv.Quote(line) + ")"

			}
			s.WriteStringLn(line)
			if needNewLine {
				s.NewLine()
			}
		}
		s.NewLine()
		s.WriteStringLn("callBack(g.Buffer)")
		// s.WriteStringLn("stream.WriteGoFile(" + strconv.Quote(BaseName(path)+"_gen.go") + ", g.Buffer())")
		s.WriteStringLn("stream.WriteGoFile(path, g.Buffer)")
		s.WriteStringLn("}")
		WriteGoFile("generateIR_gen_test.go", s.Bytes())
	})
}

func Unquote(line string) string {
	begin := strings.Index(line, `\"`)
	if begin < 0 {
		return line
	}
	split := strings.Split(line, `\"`)
	ss := NewBuffer("")
	for _, s := range split {
		if strings.Contains(s, `"`) {
			s = strings.ReplaceAll(s, `"`, ``)
			ss.WriteString(s)
		} else {
			ss.WriteString("strconv.Quote(")
			ss.WriteString(strconv.Quote(s))
			ss.WriteString(")")
		}
	}
	return ss.String()
}