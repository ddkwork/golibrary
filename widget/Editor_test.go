package widget

import (
	"path/filepath"
	"testing"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/ddkwork/golibrary/stream"
	"github.com/stretchr/testify/assert"

	"cogentcore.org/core/gi"
	"cogentcore.org/core/texteditor"
)

func TestNasm(t *testing.T) {
	s := stream.NewReadFile("vmintrin.asm")
	// assert.NoError(t, quick.Highlight(os.Stdout, s.String(), "Nasm", "html", "manni"))
	nasm := stream.New("")
	assert.NoError(t, quick.Highlight(nasm, s.String(), "Nasm", "svg", "manni"))
	stream.WriteTruncate("nasm.svg", nasm.String())
	return
	lexer := lexers.Analyse(s.String())
	assert.NotNil(t, lexer)
	println(lexer.Config().Name)
}

func TestSetEditorBuf(t *testing.T) {
	b := gi.NewBody("imageConvert")
	// editor := NewEditor(b, "log view")
	// SetEditorBuf(editor, "111")
	// SetEditorBuf(editor, []byte{1, 2, 3, 4})
	texteditor.NewEditor(b).SetBuf(texteditor.NewBuf().SetText([]byte("111")))
	b.AssertRender(t, filepath.Join("Editor", "basic"))
}
