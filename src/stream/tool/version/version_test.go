package version

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterfaceVerSion(t *testing.T) {
	v := New()
	assert.True(t, v.VerSion("56196.439.0"))
	assert.Equal(t, uint64(56196), v.Major())
	assert.Equal(t, uint64(439), v.Minor())
	assert.Equal(t, uint64(0), v.Patch())
}
