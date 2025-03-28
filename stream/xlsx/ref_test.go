package xlsx_test

import (
	"testing"

	"github.com/ddkwork/golibrary/assert"
	"github.com/ddkwork/golibrary/stream/xlsx"
)

func TestRef(t *testing.T) {
	for _, d := range []struct {
		Text string
		Col  int
		Row  int
	}{
		{"A1", 0, 0},
		{"Z9", 25, 8},
		{"AA1", 26, 0},
		{"AA99", 26, 98},
		{"ZZ100", 701, 99},
	} {
		ref := xlsx.ParseRef(d.Text)
		assert.Equal(t, d.Col, ref.Col)
		assert.Equal(t, d.Row, ref.Row)
		assert.Equal(t, d.Text, ref.String())
	}

	for r := 0; r < 100; r++ {
		for c := 0; c < 10000; c++ {
			in := xlsx.Ref{Row: r, Col: c}
			out := xlsx.ParseRef(in.String())
			assert.Equal(t, in, out)
		}
	}
}
