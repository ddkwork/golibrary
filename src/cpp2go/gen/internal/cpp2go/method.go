package cpp2go

type (
	Method interface {
		Translate() (ok bool)
	}
	method struct{}
)

func NewMethod() Method { return &method{} }
func (m *method) Translate() (ok bool) {
	return true
}

