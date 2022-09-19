package platform

import "runtime"

type (
	Interface interface {
		IsAndroid() bool
	}
	object struct{}
)

func New() Interface { return &object{} }

func (o *object) IsAndroid() bool { return runtime.GOOS == "android" }
