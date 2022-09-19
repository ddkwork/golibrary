package script

type (
	Interface interface {
		RunJava() (ok bool)
		RunGoCode() (ok bool) //todo https://github.com/search?q=goscript
	}
	object struct {
	}
)

func (o *object) RunJava() (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (o *object) RunGoCode() (ok bool) {
	//TODO implement me
	panic("implement me")
}

var Default = New()

func New() Interface { return &object{} }
