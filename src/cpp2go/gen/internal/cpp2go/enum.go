package cpp2go

type (
	Enum interface {
		Translate() (ok bool)
	}
	enum struct{}
)

func NewEnum() Enum { return &enum{} }
func (e *enum) Translate() (ok bool) {
	return true
}

