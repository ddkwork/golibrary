package internal

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
)

const goPackageDocURL = "https://protobuf.dev/reference/go/go-generated#package"

func (opts Options) Run(f func(*Plugin) error) {
}

type Plugin struct {
	// 请求是原程序提供的代码操作请求 。
	Request string

	// 文件是要生成的文件集及其导入的所有文件集。 文件以表层顺序显示, 所以每个文件都出现在任何导入文件之前 。
	// Files appear in topological order, so each file appears before any
	// file that imports it.
	Files       []*File
	FilesByPath map[string]*File

	SupportedFeatures uint64

	// fileReg        *protoregistry.Files
	// enumsByName    map[protoreflect.FullName]*Enum
	// messagesByName map[protoreflect.FullName]*Message
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
		// fileReg:        new(protoregistry.Files),
		// enumsByName:    make(map[protoreflect.FullName]*Enum),
		// messagesByName: make(map[protoreflect.FullName]*Message),
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
	GoDescriptorIdent GoIdent       // 文件描述符的 Go 名称变量
	GoPackageName     GoPackageName // 此文件的 Go 软件包的名称
	GoImportPath      GoImportPath  // 此文件的 Go 软件包导入路径

	Enums      []*Enum      // 高层全体全体宣言
	Messages   []*Message   // 最高一级信息公告
	Extensions []*Extension // 最高一级延期声明
	Services   []*Service   // 最高一级服务

	Generate bool // 如果我们为此文件生成代码, 是否真实

	// 生成的FilenamePrefix 用于构建与此源文件相关的生成文件的文件名。 例如, 源文件“ dir/ foo. proto” 可能有“ dir/ foo” 的文件名前缀。 附加的“. pb.go” 产生一个“ dir/ foo. pb.go. ” 输出文件 。
	// files associated with this source file.
	//
	// For example, the source file "dir/foo.proto" might have a filename prefix
	// of "dir/foo". Appending ".pb.go" produces an output file of "dir/foo.pb.go".
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

	GoIdent GoIdent // Go 类型

	Values []*EnumValue // 价值申报

	Location Location   // 此enum 的位置
	Comments CommentSet // 与本集相关的评论
}

func newEnum(gen *Plugin, f *File, parent *Message, desc string) *Enum {
	enum := &Enum{
		Desc: desc,
	}
	//gen.enumsByName[desc.FullName()] = enum
	//for i, vds := 0, enum.Desc.Values(); i < vds.Len(); i++ {
	//	enum.Values = append(enum.Values, newEnumValue(gen, f, parent, enum, vds.Get(i)))
	//}
	return enum
}

type EnumValue struct {
	Desc string

	GoIdent GoIdent // 生成的 Go 声明名称

	Parent *Enum // 声明此值的千兆

	Location Location   // 此enum 值的位置
	Comments CommentSet // 与此宏值相关的注释
}

func newEnumValue(gen *Plugin, f *File, message *Message, enum *Enum, desc string) *EnumValue {
	if message != nil {
	}
	name := ""
	return &EnumValue{
		Desc:    desc,
		GoIdent: f.GoImportPath.Ident(name),
		Parent:  enum,
		// Comments: makeCommentSet(f.Desc.SourceLocations().ByDescriptor(desc)),
	}
}

type Message struct {
	Desc string

	GoIdent GoIdent // Go 类型

	Fields []*Field // 消息字段声明
	Oneofs []*Oneof // 声明信息之一

	Enums      []*Enum      // 嵌套的昆虫声明
	Messages   []*Message   // 嵌套信件声明
	Extensions []*Extension // 嵌套扩展声明

	Location Location   // 此信件地址
	Comments CommentSet // 与此信件相关的注释
}

func newMessage(gen *Plugin, f *File, parent *Message, desc string) *Message {
	return nil
}

type Field struct {
	Desc     string
	GoName   string     // 例如, “ 字段Name”
	GoIdent  GoIdent    // 例如, “ 邮件Name_ fieldName”
	Parent   *Message   // 声明此字段的信息; 如果最高扩展
	Oneof    *Oneof     // 含有; 如果不是其中之一的一部分,则无
	Extendee *Message   // 扩展字段的扩展消息;其他无
	Enum     *Enum      // 输入字段的类型类型; 无效
	Message  *Message   // 用于信件或组字段或组字段的类型;否则为零
	Location Location   // 此字段位置
	Comments CommentSet // 与此字段相关的评论意见
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
	GoName   string     // 例如,“其中之一”
	GoIdent  GoIdent    // 例如, “ 信件Name_ one ofName”
	Parent   *Message   // 此一电文中声明该电文
	Fields   []*Field   // 字段的一部分
	Location Location   // 此位置
	Comments CommentSet // 与本
}

type Extension = Field

type Service struct {
	GoName string

	Methods []*Method // 工具方法声明

	Location Location   // 本服务所在地
	Comments CommentSet // 与此服务相关的注释
}

type Method struct {
	Desc string

	GoName string

	Parent *Service // 声明使用此方法的服务

	Input  *Message
	Output *Message

	Location Location   // 此方法的位置
	Comments CommentSet // 与此方法相关的注释
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

	// Go 中所有预先申报的身份识别特征都已使用 。
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
	//for i, orig := 1, packageName; g.usedPackageNames[packageName]; i++ {
	//	packageName = orig + GoPackageName(strconv.Itoa(i))
	//}
	//g.packageNames[ident.GoImportPath] = packageName
	//g.usedPackageNames[packageName] = true
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
	// 位置是元素的来源. proto 文件。
	Location Location
}

func (g *GeneratedFile) AnnotateSymbol(symbol string, info Annotation) {
	g.annotations[symbol] = append(g.annotations[symbol], info)
}

func (g *GeneratedFile) Content() ([]byte, error) {
	if !strings.HasSuffix(g.filename, ".go") {
		return g.buf.Bytes(), nil
	}

	// Reformat生成的代码。
	original := g.buf.Bytes()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", original, parser.ParseComments)
	if err != nil {
		// 用行号打印错误代码。 这不应该在实际中发生, 但是在改变生成的代码时, 它可以认为这是一个调试援助 。
		// This should never happen in practice, but it can while changing generated code
		// so consider this a debugging aid.
		var src bytes.Buffer
		s := bufio.NewScanner(bytes.NewReader(original))
		for line := 1; s.Scan(); line++ {
			fmt.Fprintf(&src, "%5d\t%s\n", line, s.Bytes())
		}
		return nil, fmt.Errorf("%v: unparsable Go source: %v\n%v", g.filename, err, src.String())
	}

	// 收集所有导入的分类列表 。
	var importPaths [][2]string
	rewriteImport := func(importPath string) string {
		if f := g.gen.opts.ImportRewriteFunc; f != nil {
			return string(f(GoImportPath(importPath)))
		}
		return importPath
	}
	for importPath := range g.packageNames {
		pkgName := string(g.packageNames[GoImportPath(importPath)])
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

	// 修改 AST 以包括一个新的导入块。
	if len(importPaths) > 0 {
		// 在软件包语句后面插入块,或在软件包语句末尾附加的软件包语句或可能的注释。
		// possible comment attached to the end of the package statement.
		pos := file.Package
		tokFile := fset.File(file.Package)
		pkgLine := tokFile.Line(file.Package)
		for _, c := range file.Comments {
			if tokFile.Line(c.Pos()) > pkgLine {
				break
			}
			pos = c.End()
		}

		// 构造导入块 。
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
	if err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(&out, fset, file); err != nil {
		return nil, fmt.Errorf("%v: can not reformat Go source: %v", g.filename, err)
	}
	return out.Bytes(), nil
}

type GeneratedCodeInfo struct{}

func (g *GeneratedFile) generatedCodeInfo(content []byte) (*GeneratedCodeInfo, error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "", content, 0)
	if err != nil {
		return nil, err
	}

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
