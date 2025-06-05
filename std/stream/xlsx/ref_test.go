package xlsx_test

import (
	"testing"

	"github.com/ddkwork/golibrary/std/assert"
	"github.com/ddkwork/golibrary/std/stream/xlsx"
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

	for r := range 100 {
		for c := range 10000 {
			in := xlsx.Ref{Row: r, Col: c}
			out := xlsx.ParseRef(in.String())
			assert.Equal(t, in, out)
		}
	}
}
