package gen

import "testing"

//go:generate go install github.com/vektra/mockery/v2@latest
//go:generate mockery --all --with-expecter --inpackage

func TestName(t *testing.T) {
	//m := NewMockInterface(t)
	//m.EXPECT()
	//m.EXPECT().FileAction().RunAndReturn(func() {
	//})
}
