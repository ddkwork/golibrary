package net

type (
	Interface interface {
		Xxobject() xxInterface
	}
	object struct {
		xxobject xxInterface
	}
)

func (o *object) Xxobject() xxInterface {
	return o.xxobject
}

var Default = New()

func New() Interface {
	return &object{
		xxobject: xxNew(),
	}
}
