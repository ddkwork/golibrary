package golibrary

import (
	"github.com/ddkwork/golibrary/fakeError"
	"testing"
)

func TestName(t *testing.T) {
	fakeError.Walk("")
	StaticCheck()
}
