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
			if i%2 != 0 {
				b = []byte(fmt.Sprintf("#%d-CONTROL-FRAME-FROM-SERVER", i))
			}
			mylog.Call(func() {
				mylog.Check2(h.WritePing(b))
				mylog.Check2(h.WritePong(b))
			})
			time.Sleep(10 * time.Millisecond)
		}
	}()

	traits := mylog.Check2(gurps.NewTraitsFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath)))

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

var (
	EvalOperators = eval.FixedOperators[DP](true)
	EvalFuncs     = eval.FixedFunctions[DP]()
)

var DebugVariableResolver = false

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

func closer() {
	f := mylog.Check2(os.Open(""))

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
	buf := make([]byte, 64)
	go checkWrite(t, w, []byte("hello, world"), c)
	n := mylog.Check2(r.Read(buf))

	<-c
	r.Close()
	w.Close()

	if mylog.Check(w.Close()); err != nil {
		t.Errorf("w.Close: %v", err)
	}
}

func delayClose(t *testing.T, cl closer, ch chan int, tt pipeTest) {
	time.Sleep(1 * time.Millisecond)

	if tt.closeWithError {
		mylog.Check(cl.CloseWithError(tt.err))
	} else {
		mylog.Check(cl.Close())
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
		conn := mylog.Check2(ln.Accept())

		go func() {
			defer close(errc)
			defer conn.Close()

			f := mylog.Check2(os.Open(newton))

			defer f.Close()

			sbytes := mylog.Check2(io.Copy(conn, f))

			if sbytes != newtonLen {
				errc <- fmt.Errorf("sent %d bytes; expected %d", sbytes, newtonLen)
				return
			}
		}()
	}(ln)

	c := mylog.Check2(Dial("tcp", ln.Addr().String()))

	defer c.Close()

	h := sha256.New()
	rbytes := mylog.Check2(io.Copy(h, c))

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

func (r *flateReadWrapper) Read(p []byte) (int, error) {
	nextRune, _, err = t.input.ReadRune()
	nextRuneType = t.classifier.ClassifyRune(nextRune)

	if err == io.EOF {
		nextRuneType = eofRuneClass
		err = nil
	}

	conn := mylog.Check2(upgrader.Upgrade(w, r, nil))

	if r.fr == nil {
		return 0, io.ErrClosedPipe
	}
	n, err := r.fr.Read(p)
	if err == io.EOF {
		r.Close()
	}

	tok := iter()
	if tok == chroma.EOF {
		err = errors.New(chroma.EOF.String())
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
		err := xa.readFromFile(r)
		if err == io.EOF {
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
		if allow := mylog.Check2(strconv.ParseBool(allowHeader)); err != nil {
			logger.Noticef("error parsing %s header: %s", client.AllowInteractionHeader, err)
		} else if allow {
			flags |= polkit.CheckAllowInteraction
		}
	}

	switch authorized := mylog.Check2(polkitCheckAuthorization(ucred.Pid, ucred.Uid, action, nil, flags)); err {
	case nil:
		if authorized {
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
	if file := mylog.Check2(a.fileFinder.WindowFile()); err == nil {
		editor.Marks.ShiftDueToTextModification(file, startOfChange, lengthOfChange)
	}
}

func testRetType() {
	index, _ := mylog.Check3(rightWordBoundary())
}

func rightWordBoundary() (byteIndex, runeIndex int, err error) {
	return r.rightBoundary(unicode.IsSpace)
}

func (de *DateEdit) SetDate(date time.Time) error {
	stNew := de.timeToSystemTime(date)
	stOld := mylog.Check2(de.systemTime())

	if mylog.Check(de.setSystemTime(stNew)); err != nil {
		return err
	}

	de.dateChangedPublisher.Publish()

	return nil
}

func dataFieldFromPath(root reflect.Value, path string) (DataField, error) {
	parent, value := mylog.Check3(reflectValueFromPath(root, path))

	if i, ok := value.Interface().(DataField); ok {
		return i, nil
	}

	return &reflectField{parent: parent, value: value, key: path[strings.LastIndexByte(path, '.')+1:]}, nil
}

func getSnapDirOptions(snap string) (*dirs.SnapDirOptions, error) {
	var opts dirs.SnapDirOptions

	data := mylog.Check2(os.ReadFile(filepath.Join(dirs.SnapSeqDir, snap+".json")))
	if errors.Is(err, os.ErrNotExist) {
		return &opts, nil
	}

	var seq struct {
		MigratedToHiddenDir   bool `json:"migrated-hidden"`
		MigratedToExposedHome bool `json:"migrated-exposed-home"`
	}
	if mylog.Check(json.Unmarshal(data, &seq)); err != nil {
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

	host, port := mylog.Check3(net.SplitHostPort(""))

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
		defer recovery(nil)
		handler(e)
	}
}

func getFuncAndMethod() {
	src, _ := os.ReadFile("D:\\workspace\\workspace\\app\\widget\\TreeTable.go")

	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "D:\\workspace\\workspace\\app\\widget\\TreeTable.go", src, 0)

	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			body := fmt.Sprintf("%s\n", src[d.Pos()-1:d.End()-1])
			before, _, found := strings.Cut(body, "{")
			if found {
				if unicode.IsUpper(rune(before[5])) {
					println(before)
				}
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := typeSpec.Type.(*ast.FuncType); ok {
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

	mylog.Check(binary.Read(bytes.NewReader([]byte{0}), binary.BigEndian, &res))
	checkResult(t, "ReadBool", binary.BigEndian, err, res, false)
	res = false
	mylog.Check(binary.Read(bytes.NewReader([]byte{1}), binary.BigEndian, &res))
	checkResult(t, "ReadBool", binary.BigEndian, err, res, true)
	res = false
	mylog.Check(binary.Read(bytes.NewReader([]byte{2}), binary.BigEndian, &res))
	checkResult(t, "ReadBool", binary.BigEndian, err, res, true)
}

func checkResult(t *testing.T, dir string, order binary.ByteOrder, err error, have, want any) {
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
		mylog.Check(binary.Write(buf, binary.BigEndian, srcSlice.Interface()))

		dst := reflect.New(src.Type()).Elem()
		dstSlice := dst.Slice(0, dst.Len())
		mylog.Check(binary.Read(buf, binary.BigEndian, dstSlice.Interface()))

		if !reflect.DeepEqual(src.Interface(), dst.Interface()) {
			t.Fatal(src)
		}
	}
}

func ValidateXauthorityFile(path string) error {
	f := mylog.Check2(os.Open(path))

	defer f.Close()
	return ValidateXauthority(f)
}

func MockXauthority(cookies int) (string, error) {
	f := mylog.Check2(os.CreateTemp("", "xauth"))
	defer f.Close()
	for n := 0; n < cookies; n++ {
		data := []byte{
			0x01, 0x00,

			0x00, 0x04, 0x73, 0x6e, 0x61, 0x70,

			0x00, 0x01, 0xff,

			0x00, 0x05, 0x73, 0x6e, 0x61, 0x70, 0x64,

			0x00, 0x01, 0xff,
		}
		m := mylog.Check2(f.Write(data))
		if m != len(data) {
			return "", fmt.Errorf("Could write cookie")
		}
	}
	return f.Name(), nil
}

func fackX() {
	outf := mylog.Check2(osutil.NewAtomicFile(tempPath, 0644, 0, osutil.NoChown, osutil.NoChown))

	defer outf.Cancel()

	if mylog.Check(outf.CommitAs(targetName)); err != nil {
		return nil, fmt.Errorf("cannot commit file to assets cache: %v", err)
	}

	if _ := mylog.Check2(io.Copy(outf, tr)); err != nil {
		return nil, fmt.Errorf("cannot copy trusted asset to cache: %v", err)
	}

	u = mylog.Check2(s.markSuccessful(u))

	b2JSON := mylog.Check2(json.Marshal(b2))
	return bytes.Equal(b1JSON, b2JSON)

	if b == nil {
		return nil
	}
}
