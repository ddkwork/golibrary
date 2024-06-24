package gen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
)

const goPackageDocURL = "https://protobuf.dev/reference/go/go-generated#package"

func (opts Options) Run(f func(*Plugin) error) {
}

type Plugin struct {
	Request string

	Files       []*File
	FilesByPath map[string]*File

	SupportedFeatures uint64

	annotateCode bool
	pathType     pathType
	module       string
	genFiles     []*GeneratedFile
	opts         Options
	err          error
}

type Options struct {
	ParamFunc func(name, value string) error

	ImportRewriteFunc func(GoImportPath) GoImportPath
}

func (opts Options) New(req string) (*Plugin, error) {
	gen := &Plugin{
		Request:     req,
		FilesByPath: make(map[string]*File),

		opts: opts,
	}

	return gen, nil
}

func (gen *Plugin) Error(err error) {
	if gen.err == nil {
		gen.err = err
	}
}

type File struct {
	GoDescriptorIdent GoIdent
	GoPackageName     GoPackageName
	GoImportPath      GoImportPath

	Enums      []*Enum
	Messages   []*Message
	Extensions []*Extension
	Services   []*Service

	Generate bool

	GeneratedFilenamePrefix string

	location Location
}

func splitImportPathAndPackageName(s string) (GoImportPath, GoPackageName) {
	if i := strings.Index(s, ";"); i >= 0 {
		return GoImportPath(s[:i]), GoPackageName(s[i+1:])
	}
	return GoImportPath(s), ""
}

type Enum struct {
	Desc string

	GoIdent GoIdent

	Values []*EnumValue

	Location Location
	Comments CommentSet
}

func newEnum(gen *Plugin, f *File, parent *Message, desc string) *Enum {
	enum := &Enum{
		Desc: desc,
	}

	return enum
}

type EnumValue struct {
	Desc string

	GoIdent GoIdent

	Parent *Enum

	Location Location
	Comments CommentSet
}

func newEnumValue(gen *Plugin, f *File, message *Message, enum *Enum, desc string) *EnumValue {
	if message != nil {
	}
	name := ""
	return &EnumValue{
		Desc:    desc,
		GoIdent: f.GoImportPath.Ident(name),
		Parent:  enum,
	}
}

type Message struct {
	Desc string

	GoIdent GoIdent

	Fields []*Field
	Oneofs []*Oneof

	Enums      []*Enum
	Messages   []*Message
	Extensions []*Extension

	Location Location
	Comments CommentSet
}

func newMessage(gen *Plugin, f *File, parent *Message, desc string) *Message {
	return nil
}

type Field struct {
	Desc     string
	GoName   string
	GoIdent  GoIdent
	Parent   *Message
	Oneof    *Oneof
	Extendee *Message
	Enum     *Enum
	Message  *Message
	Location Location
	Comments CommentSet
}

func newField(gen *Plugin, f *File, message *Message) *Field {
	field := &Field{
		GoIdent: GoIdent{
			GoImportPath: f.GoImportPath,
		},
		Parent: message,
	}
	return field
}

type Oneof struct {
	GoName   string
	GoIdent  GoIdent
	Parent   *Message
	Fields   []*Field
	Location Location
	Comments CommentSet
}

type Extension = Field

type Service struct {
	GoName string

	Methods []*Method

	Location Location
	Comments CommentSet
}

type Method struct {
	Desc string

	GoName string

	Parent *Service

	Input  *Message
	Output *Message

	Location Location
	Comments CommentSet
}

func newMethod(gen *Plugin, f *File, service *Service, desc string) *Method {
	method := &Method{
		Desc:   desc,
		Parent: service,
	}
	return method
}

type GeneratedFile struct {
	gen              *Plugin
	skip             bool
	filename         string
	goImportPath     GoImportPath
	buf              bytes.Buffer
	packageNames     map[GoImportPath]GoPackageName
	usedPackageNames map[GoPackageName]bool
	manualImports    map[GoImportPath]bool
	annotations      map[string][]Annotation
}

func (gen *Plugin) NewGeneratedFile(filename string, goImportPath GoImportPath) *GeneratedFile {
	g := &GeneratedFile{
		gen:              gen,
		filename:         filename,
		goImportPath:     goImportPath,
		packageNames:     make(map[GoImportPath]GoPackageName),
		usedPackageNames: make(map[GoPackageName]bool),
		manualImports:    make(map[GoImportPath]bool),
		annotations:      make(map[string][]Annotation),
	}

	for _, s := range types.Universe.Names() {
		g.usedPackageNames[GoPackageName(s)] = true
	}

	gen.genFiles = append(gen.genFiles, g)
	return g
}

func (g *GeneratedFile) P(v ...interface{}) {
	for _, x := range v {
		switch x := x.(type) {
		case GoIdent:
			fmt.Fprint(&g.buf, g.QualifiedGoIdent(x))
		default:
			fmt.Fprint(&g.buf, x)
		}
	}
	fmt.Fprintln(&g.buf)
}

func (g *GeneratedFile) QualifiedGoIdent(ident GoIdent) string {
	if ident.GoImportPath == g.goImportPath {
		return ident.GoName
	}
	if packageName, ok := g.packageNames[ident.GoImportPath]; ok {
		return string(packageName) + "." + ident.GoName
	}
	packageName := cleanPackageName(path.Base(string(ident.GoImportPath)))

	return string(packageName) + "." + ident.GoName
}

func cleanPackageName(base string) string {
	return ""
}

func (g *GeneratedFile) Import(importPath GoImportPath) {
	g.manualImports[importPath] = true
}

func (g *GeneratedFile) Write(p []byte) (n int, err error) {
	return g.buf.Write(p)
}

func (g *GeneratedFile) Skip() {
	g.skip = true
}

func (g *GeneratedFile) Unskip() {
	g.skip = false
}

func (g *GeneratedFile) Annotate(symbol string, loc Location) {
	g.AnnotateSymbol(symbol, Annotation{Location: loc})
}

type Annotation struct {
	Location Location
}

func (g *GeneratedFile) AnnotateSymbol(symbol string, info Annotation) {
	g.annotations[symbol] = append(g.annotations[symbol], info)
}

func (g *GeneratedFile) Content() ([]byte, error) {
	if !strings.HasSuffix(g.filename, ".go") {
		return g.buf.Bytes(), nil
	}

	original := g.buf.Bytes()
	fset := token.NewFileSet()
	file, e := parser.ParseFile(fset, "", original, parser.ParseComments)
	if e != nil {

		var src bytes.Buffer
		s := bufio.NewScanner(bytes.NewReader(original))
		for line := 1; s.Scan(); line++ {
			fmt.Fprintf(&src, "%5d\t%s\n", line, s.Bytes())
		}
		return nil, fmt.Errorf("%v: unparsable Go source: %v\n%v", g.filename, e, src.String())
	}

	var importPaths [][2]string
	rewriteImport := func(importPath string) string {
		if f := g.gen.opts.ImportRewriteFunc; f != nil {
			return string(f(GoImportPath(importPath)))
		}
		return importPath
	}
	for importPath := range g.packageNames {
		pkgName := string(g.packageNames[importPath])
		pkgPath := rewriteImport(string(importPath))
		importPaths = append(importPaths, [2]string{pkgName, pkgPath})
	}
	for importPath := range g.manualImports {
		if _, ok := g.packageNames[importPath]; !ok {
			pkgPath := rewriteImport(string(importPath))
			importPaths = append(importPaths, [2]string{"_", pkgPath})
		}
	}
	sort.Slice(importPaths, func(i, j int) bool {
		return importPaths[i][1] < importPaths[j][1]
	})

	if len(importPaths) > 0 {

		pos := file.Package
		tokFile := fset.File(file.Package)
		pkgLine := tokFile.Line(file.Package)
		for _, c := range file.Comments {
			if tokFile.Line(c.Pos()) > pkgLine {
				break
			}
			pos = c.End()
		}

		impDecl := &ast.GenDecl{
			Tok:    token.IMPORT,
			TokPos: pos,
			Lparen: pos,
			Rparen: pos,
		}
		for _, importPath := range importPaths {
			impDecl.Specs = append(impDecl.Specs, &ast.ImportSpec{
				Name: &ast.Ident{
					Name:    importPath[0],
					NamePos: pos,
				},
				Path: &ast.BasicLit{
					Kind:     token.STRING,
					Value:    strconv.Quote(importPath[1]),
					ValuePos: pos,
				},
				EndPos: pos,
			})
		}
		file.Decls = append([]ast.Decl{impDecl}, file.Decls...)
	}

	var out bytes.Buffer
	mylog.Check((&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(&out, fset, file))
	return out.Bytes(), nil
}

type GeneratedCodeInfo struct{}

func (g *GeneratedFile) generatedCodeInfo(content []byte) (*GeneratedCodeInfo, error) {
	fset := token.NewFileSet()
	astFile := mylog.Check2(parser.ParseFile(fset, "", content, 0))

	seenAnnotations := make(map[string]bool)
	annotate := func(s string, ident *ast.Ident) {
		seenAnnotations[s] = true
		for range g.annotations[s] {
		}
	}
	for _, decl := range astFile.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					annotate(spec.Name.Name, spec.Name)
					switch st := spec.Type.(type) {
					case *ast.StructType:
						for _, field := range st.Fields.List {
							for _, name := range field.Names {
								annotate(spec.Name.Name+"."+name.Name, name)
							}
						}
					case *ast.InterfaceType:
						for _, field := range st.Methods.List {
							for _, name := range field.Names {
								annotate(spec.Name.Name+"."+name.Name, name)
							}
						}
					}
				case *ast.ValueSpec:
					for _, name := range spec.Names {
						annotate(name.Name, name)
					}
				}
			}
		case *ast.FuncDecl:
			if decl.Recv == nil {
				annotate(decl.Name.Name, decl.Name)
			} else {
				recv := decl.Recv.List[0].Type
				if s, ok := recv.(*ast.StarExpr); ok {
					recv = s.X
				}
				if id, ok := recv.(*ast.Ident); ok {
					annotate(id.Name+"."+decl.Name.Name, decl.Name)
				}
			}
		}
	}
	for a := range g.annotations {
		if !seenAnnotations[a] {
			return nil, fmt.Errorf("%v: no symbol matching annotation %q", g.filename, a)
		}
	}

	return nil, nil
}

type GoIdent struct {
	GoName       string
	GoImportPath GoImportPath
}

func (id GoIdent) String() string { return fmt.Sprintf("%q.%v", id.GoImportPath, id.GoName) }

type Descriptor string

func newGoIdent(f *File, d Descriptor) *GoIdent {
	return nil
}

type GoImportPath string

func (p GoImportPath) String() string { return strconv.Quote(string(p)) }

func (p GoImportPath) Ident(s string) GoIdent {
	return GoIdent{GoName: s, GoImportPath: p}
}

type GoPackageName string

type pathType int

const (
	pathTypeImport pathType = iota
	pathTypeSourceRelative
)

type Location struct {
	SourceFile string
	Path       string
}

type CommentSet struct {
	LeadingDetached []Comments
	Leading         Comments
	Trailing        Comments
}

func makeCommentSet(loc string) CommentSet {
	var leadingDetached []Comments
	return CommentSet{
		LeadingDetached: leadingDetached,
	}
}

type Comments string

func (c Comments) String() string {
	if c == "" {
		return ""
	}
	var b []byte
	for _, line := range strings.Split(strings.TrimSuffix(string(c), "\n"), "\n") {
		b = append(b, "//"...)
		b = append(b, line...)
		b = append(b, "\n"...)
	}
	return string(b)
}

type Types struct{}

type extensionRegistry struct {
	base  *Types
	local *Types
}

func newExtensionRegistry() *extensionRegistry {
	GlobalTypes := &Types{}
	return &extensionRegistry{
		base:  GlobalTypes,
		local: &Types{},
	}
}

type ExtensionType struct{}

func (e *extensionRegistry) FindExtensionByName(field string) (*ExtensionType, error) {
	return nil, nil
}

type FieldNumber string

type FullName string

func (e *extensionRegistry) FindExtensionByNumber(message FullName, field FieldNumber) (*ExtensionType, error) {
	return nil, nil
}
