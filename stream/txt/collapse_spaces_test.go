package txt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ddkwork/golibrary/stream/txt"
)

func TestCollapseSpaces(t *testing.T) {
	data := []string{
		"123", "123",
		" 123", "123",
		" 123 ", "123",
		"    abc  ", "abc",
		"  a b c   d", "a b c d",
		"", "",
		" ", "",
	}
	for i := 0; i < len(data); i += 2 {
		assert.Equal(t, data[i+1], txt.CollapseSpaces(data[i]))
	}
}
