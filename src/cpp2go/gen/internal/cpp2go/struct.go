package cpp2go

type (
	Structtype interface {
		Translate() (ok bool)
	}
	structtype struct{}
)

func NewStructtype() Structtype { return &structtype{} }
func (s *structtype) Translate() (ok bool) {
	return true
}

