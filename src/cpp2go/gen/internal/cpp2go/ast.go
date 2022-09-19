package cpp2go

type (
	Ast interface {
		Translate() (ok bool)
	}
	ast struct{}
)

func NewAst() Ast { return &ast{} }
func (a *ast) Translate() (ok bool) {
	return true
}

