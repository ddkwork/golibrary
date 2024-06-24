package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/ddkwork/golibrary/mylog"
)

func main() {
	src := `
package main

import (
	"os"
)

func closer() {
	f := os.Open("")
	defer f.Close()

	r, w := io.Pipe()
	go func() {
		fmt.Fprint(w, "some io.Reader stream to be read\n")
		w.Close()
	}()
}

func TestPipe1(t *testing.T) {
	c := make(chan int)
	r, w := Pipe()
	var buf = make([]byte, 64)
	go checkWrite(t, w, []byte("hello, world"), c)
	n, err := r.Read(buf)
	if err != nil {
		t.Errorf("read: %v", err)
	} else if n != 12 || string(buf[0:12]) != "hello, world" {
		t.Errorf("bad read: got %q", buf[0:n])
	}
	<-c
	r.Close()
	w.Close()

if err = w.Close(); err != nil {
			t.Errorf("w.Close: %v", err)
		}

}

func delayClose(t *testing.T, cl closer, ch chan int, tt pipeTest) {
	time.Sleep(1 * time.Millisecond)
	
	if tt.closeWithError {
		err = cl.CloseWithError(tt.err)
	} else {
		err = cl.Close()
	}
	if err != nil {
		t.Errorf("delayClose: %v", err)
	}
	ch <- 0
}

type file struct {
	file  *os.File
	data  []byte
	atEOF bool
}

func (f *file) close() { f.file.Close() }

func TestPipe(t *testing.T) {
	nettest.TestConn(t, func() (c1, c2 net.Conn, stop func(), err error) {
		c1, c2 = net.Pipe()
		stop = func() {
			c1.Close()
			c2.Close()
		}
		return
	})
}

func TestSendfile(t *testing.T) {
	ln := newLocalListener(t, "tcp")
	defer ln.Close()

	errc := make(chan error, 1)
	go func(ln Listener) {
		// Wait for a connection.
		conn, err := ln.Accept()
		if err != nil {
			errc <- err
			close(errc)
			return
		}

		go func() {
			defer close(errc)
			defer conn.Close()

			f, err := os.Open(newton)
			if err != nil {
				errc <- err
				return
			}
			defer f.Close()

			// Return file data using io.Copy, which should use
			// sendFile if available.
			sbytes, err := io.Copy(conn, f)
			if err != nil {
				errc <- err
				return
			}

			if sbytes != newtonLen {
				errc <- fmt.Errorf("sent %d bytes; expected %d", sbytes, newtonLen)
				return
			}
		}()
	}(ln)

	// Connect to listener to retrieve file and verify digest matches
	// expected.
	c, err := Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	h := sha256.New()
	rbytes, err := io.Copy(h, c)
	if err != nil {
		t.Error(err)
	}

	if rbytes != newtonLen {
		t.Errorf("received %d bytes; expected %d", rbytes, newtonLen)
	}

	if res := hex.EncodeToString(h.Sum(nil)); res != newtonSHA256 {
		t.Error("retrieved data hash did not match")
	}

	for err := range errc {
		t.Error(err)
	}
}


`

	fset := token.NewFileSet()
	node := mylog.Check2(parser.ParseFile(fset, "example.go", src, parser.ParseComments))

	astutil.Apply(node, func(cr *astutil.Cursor) bool {
		if deferStmt, ok := cr.Node().(*ast.DeferStmt); ok {
			if callExpr, ok := deferStmt.Call.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := callExpr.X.(*ast.Ident); ok && ident.Name == "f" && callExpr.Sel.Name == "Close" {
					cr.Replace(&ast.DeferStmt{
						Call: &ast.CallExpr{
							Fun: &ast.FuncLit{
								Type: &ast.FuncType{},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										&ast.ExprStmt{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   &ast.Ident{Name: "mylog"},
													Sel: &ast.Ident{Name: "Check"},
												},
												Args: []ast.Expr{&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X:   &ast.Ident{Name: "f"},
														Sel: &ast.Ident{Name: "Close"},
													},
												}},
											},
										},
									},
								},
							},
						},
					})
				}
			}
		}
		return true
	}, nil)

	var buf bytes.Buffer
	format.Node(&buf, fset, node)
	fmt.Println(buf.String())
}
