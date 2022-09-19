package cpp2go

type (
	Extern interface {
		Translate() (ok bool)
	}
	extern struct{}
)

func NewExtern() Extern { return &extern{} }
func (e *extern) Translate() (ok bool) {
	return true
}

