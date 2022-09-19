package cpp2go

type (
	Cpp2go interface {
		Translate() (ok bool)
	}
	cpp2go struct{}
)

func NewCpp2go() Cpp2go { return &cpp2go{} }
func (c *cpp2go) Translate() (ok bool) {
	return true
}

