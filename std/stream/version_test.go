package stream_test

import (
	"testing"

	"github.com/ddkwork/golibrary/std/assert"
	"github.com/ddkwork/golibrary/std/stream"
)

func TestInterfaceVerSion(t *testing.T) {
	v := stream.NewVersion("10.0.22631.2506")
	assert.Equal(t, uint64(10), v.Major)
	assert.Equal(t, uint64(0), v.Minor)
	assert.Equal(t, uint64(22631), v.Patch)
	assert.Equal(t, uint64(2506), v.Build)
	assert.Equal(t, "10.0.22631.2506", v.String())
}
