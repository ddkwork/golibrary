package gen

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/src/caseconv"

	"github.com/ddkwork/golibrary/src/stream"
	"github.com/ddkwork/golibrary/src/stream/tool"
	"go/format"
	"go/token"
	"strings"
)

type (
	Method struct {
		ApiName string
		Body    string
	}
	Info struct {
		InterfaceName string
		Methods       []Method
	}
	Object struct {
		fileName string
		pkgName  string
		infos    []Info
	}
)

func New() *Object {
	return &Object{
		fileName: "",
		pkgName:  "",
		infos:    make([]Info, 0),
	}
}
func (g *Object) SetFileName(fileName string) *Object {
	g.fileName = fileName
	return g
}
func (g *Object) SetPkgName(pkgName string) *Object {
	g.pkgName = pkgName
	return g
}
func (g *Object) AppendInfos(info Info) *Object {
	g.infos = append(g.infos, info)
	return g
}

var (
	indent   = " "
	ok       = "ok"
	trueType = "true"
	boolType = "bool"
	NewType  = "New"
)

func (g *Object) AppendMethods(method Method) *Object {
	method.ApiName = method.ApiName +
		token.LPAREN.String() +
		token.RPAREN.String() +
		token.LPAREN.String() +
		ok +
		indent +
		boolType +
		token.RPAREN.String()
	for i, info := range g.infos {
		info.Methods = append(info.Methods, method)
		g.infos[i] = info
	}
	return g
}
func (g *Object) Generate() (ok bool) {
	b := stream.New()
	for _, info := range g.infos {
		b.WriteStringLn(token.PACKAGE.String() +
			indent +
			g.pkgName)
		if token.IsKeyword(strings.ToLower(info.InterfaceName)) {
			info.InterfaceName += token.TYPE.String()
		}
		InterfaceName := caseconv.ToCamelUpper(info.InterfaceName, false)
		objectName := caseconv.ToCamelToLower(InterfaceName, false)
		receiverName := string(objectName[0])
		b.WriteStringLn(token.TYPE.String() +
			indent +
			token.LPAREN.String())
		b.WriteStringLn(InterfaceName +
			indent +
			token.INTERFACE.String() +
			token.LBRACE.String())
		for _, method := range info.Methods {
			b.WriteStringLn(method.ApiName)
		}
		b.WriteStringLn(token.RBRACE.String())
		b.WriteStringLn(objectName +
			indent +
			token.STRUCT.String() +
			token.LBRACE.String() +
			token.RBRACE.String())
		b.WriteStringLn(token.RPAREN.String())
		b.WriteStringLn(
			token.FUNC.String() +
				indent +
				NewType +
				InterfaceName +
				token.LPAREN.String() +
				token.RPAREN.String() +
				InterfaceName +
				token.LBRACE.String() +
				indent +
				token.RETURN.String() +
				indent +
				token.AND.String() +
				indent +
				objectName +
				token.LBRACE.String() +
				token.RBRACE.String() +
				indent +
				token.RBRACE.String())
		for _, method := range info.Methods {
			b.WriteStringLn(
				token.FUNC.String() +
					indent +
					token.LPAREN.String() +
					receiverName +
					indent +
					token.MUL.String() +
					objectName +
					token.RPAREN.String() +
					method.ApiName + token.LBRACE.String())
			if method.Body != "" {
				b.WriteStringLn(method.Body)
			}
			b.WriteStringLn(token.RETURN.String() +
				indent +
				trueType)
			b.WriteStringLn(token.RBRACE.String())
		}
		source, err := format.Source(b.Bytes())
		if !mylog.Error(err) {
			if !tool.File().WriteTruncate(g.fileName, b.String()) {
				return
			}
			panic(111)
		}
		if !tool.File().WriteTruncate(g.fileName, source) {
			return
		}
		b.Reset()
	}
	return true
}
