package golibrary

import (
	"testing"

	"github.com/ddkwork/golibrary/std/fakeError"
	"github.com/ddkwork/golibrary/std/stream"
)

func TestName(t *testing.T) {
	fakeError.Walk(".", true)
	stream.Fmt(".")
	stream.Fix(".")
}
