package txt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ddkwork/golibrary/stream/txt"
)

func TestWrap(t *testing.T) {
	table := []struct {
		Prefix string
		Text   string
		Max    int
		Out    string
	}{
		{Prefix: "// ", Text: "short", Max: 78, Out: "// short"},
		{Prefix: "// ", Text: "some text that is longer", Max: 12, Out: "// some text\n// that is\n// longer"},
		{Prefix: "// ", Text: "some text\nwith embedded line feeds", Max: 16, Out: "// some text\n// with embedded\n// line feeds"},
		{Prefix: "", Text: "some text that is longer", Max: 12, Out: "some text\nthat is\nlonger"},
		{Prefix: "", Text: "some text that is longer", Max: 4, Out: "some\ntext\nthat\nis\nlonger"},
		{Prefix: "", Text: "some text that is longer, yep", Max: 4, Out: "some\ntext\nthat\nis\nlonger,\nyep"},
		{Prefix: "", Text: "some text\nwith embedded line feeds", Max: 16, Out: "some text\nwith embedded\nline feeds"},
	}
	for i, one := range table {
		assert.Equal(t, one.Out, txt.Wrap(one.Prefix, one.Text, one.Max), "#%d", i)
	}
}
