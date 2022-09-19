package cpp2go

type (
	Define interface {
		Translate() (ok bool)
	}
	define struct{}
)

func NewDefine() Define { return &define{} }
func (d *define) Translate() (ok bool) {
	return true
}

