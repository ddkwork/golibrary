package mylog

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRewriteAst(t *testing.T) {
	b := Check2(os.ReadFile("testdata/SkipCheckBase.go"))
	write("testdata/out/tmp.go", false, b)
	formatAllFiles(true, "testdata/out")
}

func TestRewriteAst2(t *testing.T) {
	t.Skip()
	formatAllFiles(false, "D:\\giotest")
}

func Test_handle_findEof(t *testing.T) {
	code := `
package main
func main() {
	req := mylog.Check2(http.ReadRequest(br))
	code, err := handshaker.ReadHandshake(br, req)
	cmp.Equal(t, err, ErrBadWebSocketVersion)
	cmp.Equal(t, code, http.StatusBadRequest)
}
`
	code = code[1:]
	h := newCodeHandle(code, false)
	hasEof := h.findEof("")
	cmp.Equal(t, hasEof)
}

func Test_handle_findEof2(t *testing.T) {
	code := `
package main
func main() {
	for {
		var count Count
		err := JSON.Receive(ws, &count)
		if mylog.CheckEof(err) {
			break
		}
		count.N++
		count.S = strings.Repeat(count.S, count.N)
		mylog.Check(JSON.Send(ws, count))
	}
}
`
	code = code[1:]
	h := newCodeHandle(code, false)
	hasEof := h.findEof("")
	cmp.Equal(t, hasEof)
}
