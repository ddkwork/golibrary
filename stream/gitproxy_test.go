package stream_test

import (
	"testing"

	"github.com/ddkwork/golibrary/stream"
)

func TestSetGitProxy(t *testing.T) {
	stream.GitProxy(true)
}

func TestUnSetGitProxy(t *testing.T) {
	stream.GitProxy(false)
}
