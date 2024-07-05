package base

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
	"unicode"

	"github.com/alecthomas/chroma/v2"
	"github.com/richardwilkes/toolbox/eval"
	"golang.org/x/net/nettest"

	"github.com/snapcore/snapd/client"
	"github.com/snapcore/snapd/dirs"
	"github.com/snapcore/snapd/logger"
	"github.com/snapcore/snapd/osutil"
	"github.com/snapcore/snapd/polkit"

	"github.com/richardwilkes/gcs/v5/model/gurps"
	"github.com/richardwilkes/unison"

	"github.com/ddkwork/golibrary/mylog"
)

// NewTraitTableDockableFromFile loads a list of traits from a file and creates a new unison.Dockable for them.
func NewTraitTableDockableFromFile(filePath string) (unison.Dockable, error) {

	req := mylog.Check2(http.ReadRequest(br))
	code, err := handshaker.ReadHandshake(br, req)
	assert.Equal(t, err, ErrBadWebSocketVersion)
	assert.Equal(t, code, http.StatusBadRequest)

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

type testCtrlAndDataHandler struct {
	hybiFrameHandler
}

func (h *testCtrlAndDataHandler) WritePing(b []byte) (n int, err error) {
	h.hybiFrameHandler.conn.wio.Lock()
	defer h.hybiFrameHandler.conn.wio.Unlock()
	w := mylog.Check2(h.hybiFrameHandler.conn.frameWriterFactory.NewFrameWriter(PingFrame))
	n = mylog.Check2(w.Write(b))
	mylog.Check(w.Close())
	return n, nil
}

func ctrlAndDataServer(ws *Conn) {
	defer ws.Close()
	h := &testCtrlAndDataHandler{hybiFrameHandler: hybiFrameHandler{conn: ws}}
	ws.frameHandler = h

	go func() {
		for i := 0; ; i++ {
			var b []byte
			if i%2 != 0 { // with or without payload
				b = []byte(fmt.Sprintf("#%d-CONTROL-FRAME-FROM-SERVER", i))
			}
			mylog.Call(func() {
				mylog.Check2(h.WritePing(b))
				mylog.Check2(h.WritePong(b)) // unsolicited pong
			})
			time.Sleep(10 * time.Millisecond)
		}
	}()

	traits, err := gurps.NewTraitsFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	d := NewTraitTableDockable(filePath, traits)
	d.needsSaveAsPrompt = false
	return d, nil
}
func logic() {
	lines := make([]string, 0)
	newReader := bufio.NewReader(r)

	for {
		line, _, err := newReader.ReadLine()
		if mylog.CheckEof(err) {
			return lines
		}
		lines = append(lines, string(line))
	}
}

//go:embed SkipCheckBase.go
var newton string

// The evaluator operators and functions that will be used when calling NewEvaluator().
var (
	EvalOperators = eval.FixedOperators[DP](true)
	EvalFuncs     = eval.FixedFunctions[DP]()
)

// DebugVariableResolver produces debug output for the variable resolver when enabled.
var DebugVariableResolver = false

// NewEvaluator creates a new evaluator whose number type is an Int.
func NewEvaluator(resolver eval.VariableResolver) *eval.Evaluator {
	return &eval.Evaluator{
		Resolver:  resolver,
		Operators: EvalOperators,
		Functions: EvalFuncs,
	}
}

func skipCase() {
	b := make([]byte, maxControlFramePayloadLength)
	n, err := io.ReadFull(frame, b)
	if err != nil && err != io.EOF && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, err
	}
	io.Copy(io.Discard, frame)
	if frame.PayloadType() == PingFrame {
		mylog.Check2(handler.WritePong(b[:n]))
	}
}

//////////////////////////// closer demo /////////////////////////////////////////////
//todo  D:\workspace\workspace\branch\golibrary\mylog\awesomeProject4\close\m.go

func closer() {
	f, err := os.Open("")
	if err != nil {

	}
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
	var err error
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

//////////////////////////////////////////////////////////////////////////////////////

func (r *flateReadWrapper) Read(p []byte) (int, error) {

	nextRune, _, err = t.input.ReadRune()
	nextRuneType = t.classifier.ClassifyRune(nextRune)

	if err == io.EOF {
		nextRuneType = eofRuneClass
		err = nil
	} else if err != nil {
		return nil, err
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}

	if r.fr == nil {
		return 0, io.ErrClosedPipe
	}
	n, err := r.fr.Read(p) //符合检测条件，判断下一行eof
	if err == io.EOF {     //然而某些神秘情况会隔几行注释才会到eof，这就操蛋
		// Preemptively place the reader back in the pool. This helps with
		// scenarios where the application does not call NextReader() soon after
		// this final read.
		r.Close()
	}

	tok := iter()
	if tok == chroma.EOF {
		err = errors.New(chroma.EOF.String()) //判断当前行eof
		break
	}

	if err := json.Unmarshal(nil, nil); err == io.EOF {
	}
	return n, err
}

func ValidateXauthority(r io.Reader) error {
	cookies := 0
	for {
		xa := &xauth{}
		err := xa.readFromFile(r) //检测到这里，判断下一行eof
		if err == io.EOF {        //不符合nil检测，这里安全
			break
		}
		cookies++
	}

	if cookies <= 0 {
		return fmt.Errorf("Xauthority file is invalid")
	}

	return nil
}

func checkPolkitActionImpl(r *http.Request, ucred *ucrednet, action string) *apiError {
	var flags polkit.CheckFlags
	allowHeader := r.Header.Get(client.AllowInteractionHeader)
	if allowHeader != "" {
		if allow, err := strconv.ParseBool(allowHeader); err != nil {
			logger.Noticef("error parsing %s header: %s", client.AllowInteractionHeader, err)
		} else if allow {
			flags |= polkit.CheckAllowInteraction
		}
	}
	// Pass both pid and uid from the peer ucred to avoid pid race
	switch authorized, err := polkitCheckAuthorization(ucred.Pid, ucred.Uid, action, nil, flags); err {
	case nil:
		if authorized {
			// polkit says user is authorised
			return nil
		}
	case polkit.ErrDismissed:
		return AuthCancelled("cancelled")
	default:
		logger.Noticef("polkit error: %s", err)
	}
	return Unauthorized("access denied")
}

func (a editableAdapter) shiftEditorItemsDueToTextModification(startOfChange, lengthOfChange int) {
	if file, err := a.fileFinder.WindowFile(); err == nil { //todo
		editor.Marks.ShiftDueToTextModification(file, startOfChange, lengthOfChange)
	}
}

func testRetType() {
	index, _, _ := rightWordBoundary()
}

func rightWordBoundary() (byteIndex, runeIndex int, err error) {
	return r.rightBoundary(unicode.IsSpace)
}

func (de *DateEdit) SetDate(date time.Time) error {
	stNew := de.timeToSystemTime(date)
	stOld, err := de.systemTime()
	if err != nil {
		return err
	} else if stNew == stOld || stNew != nil && stOld != nil && *stNew == *stOld {
		return nil
	}

	if err := de.setSystemTime(stNew); err != nil {
		return err
	}

	de.dateChangedPublisher.Publish()

	return nil
}

func dataFieldFromPath(root reflect.Value, path string) (DataField, error) {
	parent, value, err := reflectValueFromPath(root, path)
	if err != nil {
		return nil, err
	}

	// convert to DataField
	if i, ok := value.Interface().(DataField); ok {
		return i, nil
	}

	return &reflectField{parent: parent, value: value, key: path[strings.LastIndexByte(path, '.')+1:]}, nil
}

func getSnapDirOptions(snap string) (*dirs.SnapDirOptions, error) {
	var opts dirs.SnapDirOptions

	data, err := os.ReadFile(filepath.Join(dirs.SnapSeqDir, snap+".json"))
	if errors.Is(err, os.ErrNotExist) {
		return &opts, nil
	} else if err != nil {
		return nil, err
	}

	var seq struct {
		MigratedToHiddenDir   bool `json:"migrated-hidden"`
		MigratedToExposedHome bool `json:"migrated-exposed-home"`
	}
	if err := json.Unmarshal(data, &seq); err != nil {
		return nil, err
	}

	opts.HiddenSnapDataDir = seq.MigratedToHiddenDir
	opts.MigratedToExposedHome = seq.MigratedToExposedHome

	return &opts, nil
}

func tt() {

	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	GetLogicalDrives := kernel32.MustFindProc("GetLogicalDrives")
	n, _, _ := GetLogicalDrives.Call()
	s := strconv.FormatInt(int64(n), 2)

	host, port, err := net.SplitHostPort("")
	if err != nil {
		panic(err)
	}
	println(host, port)

	for _, s := range y {
		if c := m[s]; c > -8 {
			m[s] = c - 4
		}
	}

}

func recovery(handler recoveryHandler) {
	if recovered := recover(); recovered != nil && handler != nil {
		e, ok := recovered.(error)
		if !ok {
			e = fmt.Errorf("%+v", recovered)
		}
		defer recovery(nil) // nice design, avoid infinite loop
		handler(e)
	}
}

func getFuncAndMethod() {
	src, _ := os.ReadFile("D:\\workspace\\workspace\\app\\widget\\TreeTable.go")

	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "D:\\workspace\\workspace\\app\\widget\\TreeTable.go", src, 0)

	// 遍历AST并获取所有函数和方法签名

	// 遍历AST并获取所有函数和方法签名
	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl: // 检查是否是函数声明
			body := fmt.Sprintf("%s\n", src[d.Pos()-1:d.End()-1])
			before, _, found := strings.Cut(body, "{")
			if found {
				if unicode.IsUpper(rune(before[5])) {
					println(before)
				}
			}
		case *ast.GenDecl: // 检查是否是通用声明
			for _, spec := range d.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := typeSpec.Type.(*ast.FuncType); ok { // 检查是否是函数类型
						body := fmt.Sprintf("%s\n", src[typeSpec.Pos()-1:typeSpec.End()-1])
						before, _, found := strings.Cut(body, "{")
						if found {
							if unicode.IsUpper(rune(before[5])) {
								println(before)
							}
						}
					}
				}
			}
		}
	}

}
func TestReadBool(t *testing.T) {
	var res bool
	var err error
	err = binary.Read(bytes.NewReader([]byte{0}), binary.BigEndian, &res)
	checkResult(t, "ReadBool", binary.BigEndian, err, res, false)
	res = false
	err = binary.Read(bytes.NewReader([]byte{1}), binary.BigEndian, &res)
	checkResult(t, "ReadBool", binary.BigEndian, err, res, true)
	res = false
	err = binary.Read(bytes.NewReader([]byte{2}), binary.BigEndian, &res)
	checkResult(t, "ReadBool", binary.BigEndian, err, res, true)
}

func checkResult(t *testing.T, dir string, order binary.ByteOrder, err error, have, want any) {
	if err != nil {
		t.Errorf("%v %v: %v", dir, order, err)
		return
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("%v %v:\n\thave %+v\n\twant %+v", dir, order, have, want)
	}
}

var intArrays = []any{
	&[100]int8{},
	&[100]int16{},
	&[100]int32{},
	&[100]int64{},
	&[100]uint8{},
	&[100]uint16{},
	&[100]uint32{},
	&[100]uint64{},
}

func TestSliceRoundTrip(t *testing.T) {
	buf := new(bytes.Buffer)
	for _, array := range intArrays {
		src := reflect.ValueOf(array).Elem()
		unsigned := false
		switch src.Index(0).Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			unsigned = true
		}
		for i := 0; i < src.Len(); i++ {
			if unsigned {
				src.Index(i).SetUint(uint64(i * 0x07654321))
			} else {
				src.Index(i).SetInt(int64(i * 0x07654321))
			}
		}
		buf.Reset()
		srcSlice := src.Slice(0, src.Len())
		err := binary.Write(buf, binary.BigEndian, srcSlice.Interface())
		if err != nil {
			t.Fatal(err)
		}
		dst := reflect.New(src.Type()).Elem()
		dstSlice := dst.Slice(0, dst.Len())
		err = binary.Read(buf, binary.BigEndian, dstSlice.Interface())
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(src.Interface(), dst.Interface()) {
			t.Fatal(src)
		}
	}
}

func ValidateXauthorityFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return ValidateXauthority(f)
}

func MockXauthority(cookies int) (string, error) {
	f, err := os.CreateTemp("", "xauth")
	defer f.Close()
	for n := 0; n < cookies; n++ {
		data := []byte{
			// Family
			0x01, 0x00,
			// Address
			0x00, 0x04, 0x73, 0x6e, 0x61, 0x70,
			// Number
			0x00, 0x01, 0xff,
			// Name
			0x00, 0x05, 0x73, 0x6e, 0x61, 0x70, 0x64,
			// Data
			0x00, 0x01, 0xff,
		}
		m, err := f.Write(data)
		if m != len(data) {
			return "", fmt.Errorf("Could write cookie")
		}
	}
	return f.Name(), nil
}

func fackX() {
	outf, err := osutil.NewAtomicFile(tempPath, 0644, 0, osutil.NoChown, osutil.NoChown)
	if err != nil {
		return nil, fmt.Errorf("cannot create temporary cache file: %v", err)
	}
	defer outf.Cancel()

	if err := outf.CommitAs(targetName); err != nil {
		return nil, fmt.Errorf("cannot commit file to assets cache: %v", err)
	}

	if _, err := io.Copy(outf, tr); err != nil {
		return nil, fmt.Errorf("cannot copy trusted asset to cache: %v", err)
	}

	if err != nil {
		// all internal errors at this point
		panic(err)
	}

	u, err = s.markSuccessful(u)
	if err != nil {
		return fmt.Errorf(errPrefix, err)
	}

	b2JSON, err := json.Marshal(b2)
	return bytes.Equal(b1JSON, b2JSON)

	if b == nil {
		return nil
	}

	if err != nil {
		return err
	}

}
