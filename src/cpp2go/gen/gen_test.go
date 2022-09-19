package gen_test

import (
	"github.com/ddkwork/golibrary/src/cpp2go/gen"
	"path/filepath"
	"testing"
)

func TestName(t *testing.T) {
	a := []string{
		"enum",
		"struct",
		"define",
		"extern",
		"method",
		"scanner",
		"ast",
		"cpp2go",
	}
	g := gen.New()
	for _, s := range a {
		join := filepath.Join("internal/cpp2go", s+".go")
		g.SetFileName(join).SetPkgName("cpp2go")
		g.AppendInfos(gen.Info{
			InterfaceName: s,
			Methods:       nil,
		})
		g.AppendMethods(gen.Method{
			ApiName: "Translate",
			Body:    "",
		})
		if !g.Generate() {
			return
		}
	}
}
