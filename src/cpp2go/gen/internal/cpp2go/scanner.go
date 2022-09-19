package cpp2go

type (
	Scanner interface {
		Translate() (ok bool)
	}
	scanner struct{}
)

func NewScanner() Scanner { return &scanner{} }
func (s *scanner) Translate() (ok bool) {
	return true
}

