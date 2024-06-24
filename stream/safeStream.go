package stream

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/des"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/dc0d/caseconv"
	"github.com/rivo/uniseg"
	"mvdan.cc/gofumpt/format"

	"github.com/ddkwork/golibrary/mylog"
)

type (
	Type interface {
		string | HexString | HexDumpString | ~[]byte | ~*bytes.Buffer | *big.Int | *Buffer
	}
	Buffer struct{ *bytes.Buffer }
)

func NewBuffer[T Type](s T) *Buffer {
	switch s := any(s).(type) {
	case *Buffer:
		return s
	case []byte:
		return &Buffer{bytes.NewBuffer(s)}
	case *bytes.Buffer:
		return &Buffer{s}
	case string:
		if IsFilePath(s) {
			return &Buffer{Buffer: bytes.NewBuffer(mylog.Check2(os.ReadFile(s)))}
		}
		return &Buffer{Buffer: bytes.NewBufferString(s)}
	case HexString:
		return NewHexString(s)
	case HexDumpString:
		return NewHexDump(s)
	case *big.Int:
		return &Buffer{Buffer: bytes.NewBuffer(s.Bytes())}
	default:
		return &Buffer{Buffer: &bytes.Buffer{}}
	}
}

func NewHexString(s HexString) *Buffer {
	decodeString := mylog.Check2(hex.DecodeString(string(s)))
	return NewBuffer(decodeString)
}

type HexString string

func (b *Buffer) HexString() HexString      { return HexString(hex.EncodeToString(b.Bytes())) }
func (b *Buffer) HexStringUpper() HexString { return HexString(fmt.Sprintf("%#X", b.Bytes())[2:]) }

type HexDumpString string

func NewHexDump(hexdumpStr HexDumpString) (data *Buffer) {
	hexdump := string(hexdumpStr)
	defer func() {
		s := NewBuffer("")

		cut := `[]byte`
		cxx := fmt.Sprintf("%#v", data.Bytes())
		cxx = cxx[len(cut):]
		s.WriteString("char data[] = " + cxx + ";\n")
		mylog.Json("gen c++ code", s.String())
		mylog.HexDump("recovery go buffer", data.Bytes())
	}()
	hexdump = strings.TrimSuffix(hexdump, newLine)
	switch {
	case !hasAddress(hexdump) && !strings.Contains(hexdump, sep):
		hexdump = strings.ReplaceAll(hexdump, " ", "")
		decodeString := mylog.Check2(hex.DecodeString(hexdump))
		data = NewBuffer(decodeString)
	case strings.Contains(hexdump, sep):
		split := strings.Split(hexdump, newLine)
		noAddress := make([]string, len(split))
		hexString := new(bytes.Buffer)
		for i, s := range split {
			if s == "" {
				continue
			}
			noAddress[i] = s[addressLen:strings.Index(s, sep)]
			noAddress[i] = strings.ReplaceAll(noAddress[i], " ", "")
			hexString.WriteString(noAddress[i])
		}
		decodeString := mylog.Check2(hex.DecodeString(hexString.String()))
		data = NewBuffer(decodeString)
	default:
		split := strings.Split(hexdump, newLine)
		hexString := new(bytes.Buffer)
		for _, s := range split {
			if s == "" {
				continue
			}
			fields := strings.Split(s, " ")
			for j, field := range fields {
				if j > 0 && field == "" {
					fields = fields[1:j]
					break
				}
			}
			for _, field := range fields {
				hexString.WriteString(field)
			}
		}
		decodeString := mylog.Check2(hex.DecodeString(hexString.String()))
		data = NewBuffer(decodeString)
	}
	return
}

const (
	address    = "00000000  "
	sep        = "|"
	newLine    = "\n"
	addressLen = len(address)
)

func hasAddress(s string) bool {
	switch {
	case len(s) < len("00000000"):
		return false
	case strings.Contains(s, address):
		return true
	}
	return s[len("00000000")+1] == ' '
}

func WriteGoFile[T Type](name string, data T) {
	s := NewBuffer(data)
	source, e := format.Source(s.Bytes(), format.Options{})
	mylog.CheckIgnore(e)
	if e != nil {
		write(name, false, s.Bytes())
		return
	}
	write(name, false, source)
}

func WriteAppend[T Type](name string, data T)     { write(name, true, data) }
func WriteTruncate[T Type](name string, data T)   { write(name, false, data) }
func WriteBinaryFile[T Type](name string, data T) { write(name, false, data) }

func write[T Type](name string, isAppend bool, data T) {
	mylog.Call(func() {
		if !CreatDirectory(filepath.Dir(name)) {
			mylog.Check(fmt.Errorf("create directory failed: %s", filepath.Dir(name)))
		}
		fnCreateFile := func() (*os.File, error) {
			if isAppend {
				return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
			}
			return os.Create(name)
		}
		f, e := fnCreateFile()
		mylog.CheckIgnore(e)
		if e != nil {
			switch {
			case os.IsExist(e):
				mylog.Check(f.Close())
				write(name, isAppend, data)
			case os.IsNotExist(e):
				mylog.Check(f.Close())
				write(name, isAppend, data)
			case os.IsPermission(e):
				return
			}
		}
		s := NewBuffer(data)
		mylog.Check2(f.Write(s.Bytes()))
		mylog.Check(f.Close())
	})
}

func CreatDirectory(dir string) bool {
	dir = FixFilePath(dir)
	info, e := os.Stat(dir)
	mylog.CheckIgnore(e)
	if e == nil {
		if info.IsDir() {
			return true
		}
		mylog.Check("path exists but is not a directory " + dir)
	}
	mylog.Warning("", "如果第一次看到这个错误，则说明当前目录下没有这个目录，请手动检查目录结构，如果第二次运行还出现权限错误则需要检查代码和系统问题")
	switch {
	case os.IsExist(e):
		return info.IsDir()
	case os.IsNotExist(e):
		mylog.Check(os.MkdirAll(dir, os.ModePerm))
		return true
	default:
		mylog.Check(e)
	}
	return false
}

func (b *Buffer) ReaderGzip() *Buffer {
	buf := ReaderGzip(b.Bytes()).Bytes()
	b.Reset()
	mylog.Check2(b.Write(buf))
	return b
}

func ReaderGzip[T Type](data T) *Buffer {
	reader := mylog.Check2(gzip.NewReader(bytes.NewReader(NewBuffer(data).Bytes())))
	b := make([]byte, 1024*2)
	n := mylog.Check2(reader.Read(b))
	return NewBuffer(b[:n])
}

func (b *Buffer) SplitBytes(size int) (blocks [][]byte) {
	blocks = make([][]byte, 0)
	quantity := b.Len() / size
	remainder := b.Len() % size
	for i := range quantity {
		blocks = append(blocks, b.Bytes()[i*size:(i+1)*size])
		if remainder != 0 {
			blocks = append(blocks, b.Bytes()[i*size:i*size+remainder])
		}
	}
	return
}

func (b *Buffer) SplitString(size int) (blocks []string) {
	blocks = make([]string, 0)
	splitBytes := b.SplitBytes(size)
	for _, splitByte := range splitBytes {
		blocks = append(blocks, string(splitByte))
	}
	return
}

func (b *Buffer) WritePackageName() {
	b.WriteStringLn("package " + GetPackageName())
}

func GetPackageName() (pkgName string) {
	return filepath.Base(mylog.Check2(filepath.Abs(".")))
}

func (b *Buffer) CutString(left, right string) (cut string, found bool) {
	_, after, found := strings.Cut(b.String(), left)
	mylog.Check(found)
	before, _, f := strings.Cut(after, right)
	mylog.Check(f)
	before = strings.ReplaceAll(before, "\\", "")
	before = strings.ReplaceAll(before, "/", "")
	return before, true
}
func (b *Buffer) CutWithIndex(x, y int) []byte { return b.Bytes()[x:y] }
func (b *Buffer) NewLine() *Buffer             { b.WriteString("\n"); return b }
func (b *Buffer) QuoteWith(s string) *Buffer   { b.WriteString(s); return b }
func (b *Buffer) WriteBinary(order binary.ByteOrder, data any) {
	mylog.Check(binary.Write(b, order, data))
}

func (b *Buffer) ReadBinary(order binary.ByteOrder) (data any) {
	mylog.Check(binary.Read(b, order, data))
	return b
}

func (b *Buffer) WriteBytesLn(buf []byte) *Buffer {
	b.Write(buf)
	b.NewLine()
	return b
}

func (b *Buffer) ReplaceAll(old, new string) *Buffer {
	s := strings.ReplaceAll(b.String(), old, new)
	b.Reset()
	b.WriteString(s)
	return b
}

func (b *Buffer) TrimSuffix(suffix string) *Buffer {
	s := strings.TrimSuffix(b.String(), suffix)
	b.Reset()
	b.WriteString(s)
	return b
}

func (b *Buffer) TrimPrefix(prefix string) *Buffer {
	s := strings.TrimPrefix(b.String(), prefix)
	b.Reset()
	b.WriteString(s)
	return b
}

func (b *Buffer) TrimSpace() *Buffer {
	s := strings.TrimSpace(b.String())
	b.Reset()
	b.WriteString(s)
	return b
}

func (b *Buffer) WriteStringLn(s string) *Buffer {
	b.WriteString(s)
	b.NewLine()
	return b
}
func (b *Buffer) Quote()          { b.WriteByte('"') }
func (b *Buffer) ObjectBegin()    { b.WriteByte('{') }
func (b *Buffer) ObjectEnd()      { b.WriteByte('}') }
func (b *Buffer) SliceBegin()     { b.WriteByte('[') }
func (b *Buffer) SliceEnd()       { b.WriteByte(']') }
func (b *Buffer) Indent(deep int) { b.WriteString(strings.Repeat(" ", deep)) }

func (b *Buffer) CheckDesBlockSize() {
	mylog.Check(b.Len())
	mylog.Check(b.Len()%des.BlockSize == 0)
}

func Concat[S ~[]E, E any](slices_ ...S) S { return slices.Concat(slices_...) }
func (b *Buffer) MergeByte(streams ...[]byte) []byte {
	return slices.Concat(b.Bytes(), slices.Concat(streams...))
	for _, s := range streams {
		b.Write(s)
	}
	return b.Bytes()
}

func (b *Buffer) BigNumXorWithAlign(arg1, arg2 []byte, align int) (xorStream []byte) {
	xor := new(big.Int).Xor(new(big.Int).SetBytes(arg1), new(big.Int).SetBytes(arg2))
	alignBuf := make([]byte, align-len(xor.Bytes()))
	switch len(xor.Bytes()) {
	case 0:
		xorStream = alignBuf
	case align:
		xorStream = xor.Bytes()
	default:
		b.AppendByteSlice(xorStream, alignBuf, xor.Bytes())
	}
	return b.Bytes()
}

func (b *Buffer) InsertStringWithSplit(size int, insert string) string {
	blocks := b.SplitString(size)
	b.Reset()
	for i, block := range blocks {
		b.WriteString(block)
		if i < len(blocks)-1 {
			b.WriteString(insert)
		}
	}
	return b.String()
}

func (b *Buffer) InsertBytes(index int, insert []byte) []byte {
	start := b.Bytes()[:index]
	end := b.Bytes()[index:]
	b.Reset()
	b.Write(start)
	b.Write(insert)
	b.Write(end)
	return b.Bytes()
}

func (b *Buffer) InsertByte(index int, ch byte) { b.InsertBytes(index, []byte{ch}) }

func (b *Buffer) InsertRune(index int, r rune) {
	if uint32(r) < utf8.RuneSelf {
		b.InsertByte(index, byte(r))
		return
	}
	var buffer [4]byte
	n := utf8.EncodeRune(buffer[:], r)
	b.InsertBytes(index, buffer[:n])
}

func (b *Buffer) InsertString(index int, s string) string {
	return string(b.InsertBytes(index, []byte(s)))
}

func (b *Buffer) AppendByteSlice(bytesSlice ...[]byte) []byte {
	for _, slice := range bytesSlice {
		b.Write(slice)
	}
	return b.Bytes()
}

func (b *Buffer) Contains(substr string) bool      { return strings.Contains(b.String(), substr) }
func ReadFileToLines(path string) (lines []string) { return NewBuffer(path).ToLines() }

func Lines(x []byte) []string {
	l := strings.SplitAfter(string(x), "\n")
	mylog.CheckNil(l)
	if l[len(l)-1] == "" {
		l = l[:len(l)-1]
	} else {
		l[len(l)-1] += "\n\\ No newline at end of file\n"
	}
	return l
}

func (b *Buffer) ReplaceLine(index int, line string) *Buffer {
	lines := NewBuffer(b.String()).ToLines()
	lines[index] = line
	b.Reset()
	b.WriteString(b.LinesToString(lines)) // todo not working
	return b
}

func (b *Buffer) LinesToString(lines []string) string {
	for _, line := range lines {
		b.WriteStringLn(line)
	}
	return b.String()
}

func (b *Buffer) ToLines() (lines []string) {
	lines = make([]string, 0)
	newScanner := bufio.NewScanner(b.Buffer)
	for newScanner.Scan() {
		lines = append(lines, newScanner.Text())
	}
	return
}

func SplitFileByLines(filePath string, size int) {
	lines := ReadFileToLines(filePath)
	if lines == nil {
		return
	}
	count := len(lines) / size
	div := len(lines) % size
	if div != 0 {
		count++
	}
	for i := range count {
		startIndex := i * size
		endIndex := (i + 1) * size
		if endIndex > len(lines) {
			endIndex = len(lines)
		}
		WriteTruncate(fmt.Sprint(i)+".txt", NewBuffer("").LinesToString(lines[startIndex:endIndex]))
	}
}

func ToLines[T string | []byte | *os.File | *bytes.Buffer](data T) (lines []string) {
	var r io.Reader
	switch data := any(data).(type) {
	case string:
		r = strings.NewReader(data)
		if IsFilePath(data) {
			b := mylog.Check2(os.ReadFile(data))
			r = bytes.NewReader(b)
		}
	case []byte:
		r = bytes.NewReader(data)
	case *os.File:
		r = data
	case *bytes.Buffer:
		r = data
	default:
		mylog.Check(fmt.Errorf("unsupported type %T", data))
	}

	lines = make([]string, 0)
	newReader := bufio.NewReader(r)

	for {
		line, _, err := newReader.ReadLine()
		if mylog.CheckEof(err) {
			return lines
		}
		lines = append(lines, string(line))
	}
}

func ReadLines(fullPath string) ([]string, error) {
	f := mylog.Check2(os.Open(fullPath))
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func WriteLines(lines []string, fullPath string) error {
	f := mylog.Check2(os.Create(fullPath))
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range lines {
		mylog.Check2(fmt.Fprintln(w, line))
	}

	return w.Flush()
}

func IsZero(v reflect.Value) bool {
	if v.IsValid() {
		return true
	}
	return v.IsZero()
}

func ReflectVisibleFields(object any) []reflect.StructField {
	fields := reflect.VisibleFields(reflect.TypeOf(object))
	var exportedFields []reflect.StructField
	for _, field := range fields {
		if field.Tag.Get("table") == "_" || field.Tag.Get("json") == "_" {
			continue
		}
		if !unicode.IsUpper(rune(field.Name[0])) {
			mylog.Trace("field name is not exported: ", field.Name)
			continue
		}
		exportedFields = append(exportedFields, field)
	}
	return exportedFields
}

type Pool[T any] struct{ pool sync.Pool }

func NewPool[T any](fn func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return fn()
			},
		},
	}
}

func (p *Pool[T]) Put(v T) { p.pool.Put(v) }
func (p *Pool[T]) Get() T  { return p.pool.Get().(T) }

func SerialNumber(b string) string {
	to := bytes.Buffer{}
	for k, v := range b {
		if k%2 == 1 {
			to.WriteByte(byte(v))
			to.WriteByte(b[k-1])
		}
	}
	return to.String()
}

func CutUint16(u uint16) uint8 {
	u >>= 4
	high4 := u & 0x0f << 4
	low4 := u & 0xf0 >> 4
	high4 |= low4
	return uint8(high4)
}

func CutUint16_(u uint16) uint8 {
	ss := fmt.Sprintf("%x", u)

	out := string(ss[2]) + string(ss[1])
	decodeString := mylog.Check2(hex.DecodeString(out))
	if decodeString == nil {
		return 0
	}
	return decodeString[0]
}

func SwapBytes2HexString2(src HexString) (dst string) {
	hexString := NewHexString(src)
	slices.Reverse(hexString.Bytes())
	return hex.EncodeToString(hexString.Bytes())
}

func SwapBytes2HexString(src []byte) (dst HexString) {
	return NewBuffer(SwapBytes(src)).HexString()
}

func SwapBytes(src []byte) (dst []byte) {
	slices.Reverse(src)
	return src
	return

	to := bytes.Buffer{}
	for i := range src {
		to.WriteByte(src[len(src)-i-1])
	}
	return to.Bytes()
}

func AlignString(s string, length int) (ss string) {
	width := uniseg.StringWidth(s)
	if width < length {
		repeat := strings.Repeat(" ", length-width)
		ss = s + repeat
		return ss
	}
	return s
}

func SubStrRunes(s string, length int) string {
	switch {
	case len(s) > length:
		rs := []rune(s)
		return string(rs[:length])
	case len(s) < length:
		repeat := strings.Repeat(" ", length-len(s))
		return s + repeat
		return fmt.Sprintf("%-*s", length, s)
	}
	return s
}

func SlicesIndex(slice any, item any) int {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("findIndex: not slice")
	}
	for i := range s.Len() {
		if reflect.DeepEqual(s.Index(i).Interface(), item) {
			return i
		}
	}
	return -1
}

func SlicesInsert(slice any, index int, value any) any {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("Insert: not slice")
	}

	t := reflect.MakeSlice(s.Type(), s.Len()+1, s.Cap()+1)

	reflect.Copy(t.Slice(0, index), s.Slice(0, index))

	t.Index(index).Set(reflect.ValueOf(value))

	reflect.Copy(t.Slice(index+1, s.Len()+1), s.Slice(index, s.Len()))
	return t.Interface()
}

func MarshalJSON(v any) []byte {
	indent := mylog.Check2(json.MarshalIndent(v, "", " "))
	return indent
}

func MarshalJsonToFile(v any, name string) {
	indent := mylog.Check2(json.MarshalIndent(v, "", " "))
	ext := filepath.Ext(name)
	if ext != ".json" {
		name += ".json"
	}
	WriteTruncate(name, indent)
}

func JsonIndent(b []byte) string {
	buffer := new(bytes.Buffer)
	mylog.Check(json.Indent(buffer, b, "", " "))
	return buffer.String()
}

func RandomAny[T any](slice []T) T {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	return slice[random.Intn(len(slice))]
}

func GenA2Z() (letters []string) {
	return generateLatinAlphabets()
}

func generateLatinAlphabets() []string {
	var alphabets []string
	for i := 'A'; i <= 'Z'; i++ {
		alphabets = append(alphabets, string(i))
	}
	return alphabets
}

func GetWindowsLogicalDrives() []string {
	var driveLetters []string
	for _, s := range GenA2Z() {
		s += ":\\"
		_, e := os.Stat(s)
		if e == nil {
			driveLetters = append(driveLetters, s)
		}
	}
	return driveLetters
}

func ToCamel(data string, isCommit bool) string {
	s := fmt.Sprintf("%-50s", caseconv.ToCamel(data))
	if isCommit {
		s += "//" + s
	}
	return s
}

func ToCamelUpper(s string, isCommit bool) string {
	camel := ToCamel(s, isCommit)
	camel = strings.TrimSpace(camel)
	return strings.ToUpper(string(camel[0])) + camel[1:]
}

func ToCamelToLower(s string, isCommit bool) string {
	camel := ToCamel(s, isCommit)
	camel = strings.TrimSpace(camel)
	return strings.ToLower(string(camel[0])) + camel[1:]
}

func CurrentDirName(path string) (currentDirName string) {
	if path == "" {
		path = mylog.Check2(os.Getwd())
	}
	split := strings.Split(path, "\\")
	if split == nil {
		return BaseName(filepath.Dir(path))
	}
	return split[len(split)-1]
}

func CopyDir(dst, src string) error {
	if !CreatDirectory(dst) {
		mylog.Check("CreatDirectory err")
	}
	entries := mylog.Check2(os.ReadDir(src))
	for _, entry := range entries {
		if entry.IsDir() {
			mylog.Check(CopyDir(filepath.Join(dst, entry.Name()), filepath.Join(src, entry.Name())))
		} else {
			mylog.Check(copyFile(filepath.Join(dst, entry.Name()), filepath.Join(src, entry.Name())))
		}
	}
	return nil
}

func copyFile(dst, src string) (err error) {
	s := mylog.Check2(os.Open(src))
	defer func() { mylog.Check(s.Close()) }()
	d := mylog.Check2(os.Create(dst))
	defer func() { mylog.Check(d.Close()) }()
	mylog.Check2(io.Copy(d, s))
	return nil
}

func CopyFile(path, dstPath string) {
	WriteTruncate(dstPath, NewBuffer(path).Bytes())
}

func MoveFile(src, dst string) {
	srcInfo := mylog.Check2(os.Stat(src))
	if !srcInfo.Mode().IsRegular() {
		mylog.Check(fmt.Sprintf("%s is not a regular file", src))
	}
	dstInfo := mylog.Check2(os.Stat(dst))
	if !dstInfo.Mode().IsRegular() {
		mylog.Check(fmt.Sprintf("%s is not a regular file", dst))
	}
	mylog.Check(os.SameFile(srcInfo, dstInfo))
	mylog.Check(os.Rename(src, dst))
	var in, out *os.File
	out = mylog.Check2(os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcInfo.Mode()))
	defer func() { mylog.Check(out.Close()) }()
	in = mylog.Check2(os.Open(src))
	mylog.Check2(io.Copy(out, in))
	mylog.Check(os.Remove(src))
}

func IsFilePathEx(path string) (ok bool) {
	return isFilePath(path, true)
}

func IsFilePath(path string) bool {
	return isFilePath(path, false)
}

func isFilePath(path string, debug bool) bool {
	pattern := []string{"://", "\n", "*", "?", "<", ">", "|"}
	for _, s := range pattern {
		if strings.Contains(path, s) {
			return false
		}
	}
	stat, e := os.Stat(path)
	if e != nil {
		if debug {
			mylog.CheckIgnore(e)
		}
		return false
	}
	// strings.Contains(path, "/") || strings.Contains(path, "\\")
	mode := stat.Mode()
	return !mode.IsDir() && mode.IsRegular()
}

func DirDepth(dirPath string) (depth int) {
	return strings.Count(dirPath, string(filepath.Separator))
}

func IsDirRoot(path string) bool {
	return !strings.Contains(path, string(filepath.Separator))
}

func IsDirEx(path string) (ok bool) {
	return isDir(path, true)
}

func IsDir(path string) bool {
	return isDir(path, false)
}

func isDir(path string, debug bool) bool {
	if strings.HasPrefix(path, ".") && IsDirRoot(path) {
		return true
	}
	fi, e := os.Stat(path)
	if e != nil {
		return false
	}
	return fi != nil && fi.IsDir()
}

func FixFilePath(path string) string {
	return strings.ReplaceAll(strings.ReplaceAll(path, "\\", "/"), "//", "/")
}

func BaseName(path string) string {
	return TrimExtension(filepath.Base(mylog.Check2(filepath.Abs(path))))
}
func TrimExtension(path string) string { return path[:len(path)-len(filepath.Ext(path))] }

func JoinHomeDir(path string) (join string)  { return joinHome(path, true) }
func JoinHomeFile(path string) (join string) { return joinHome(path, false) }
func joinHome(path string, isDir bool) (join string) {
	join = filepath.Join(HomeDir(), path)
	if !IsFilePath(join) {
		switch isDir {
		case true:
			mylog.Check(os.MkdirAll(join, os.ModePerm))
		default:
			f := mylog.Check2(os.Create(join))
			mylog.Check(f.Close())
		}
	}
	return
}

func HomeDir() string {
	if u, e := user.Current(); e == nil {
		return u.HomeDir
	}
	if dir, e := os.UserHomeDir(); e == nil {
		return dir
	}
	return "."
}

func RunDir() string {
	return filepath.Base(mylog.Check2(os.Getwd()))
}

func ParseFloat(sizeStr string) (size float64) {
	return mylog.Check2(strconv.ParseFloat(sizeStr, 64))
}

func Float64ToString(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

func Float64Cut(value float64, bits int) (float64, error) {
	return strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(bits)+"f", value), 64)
}

func ParseInt(s string) int64 {
	return mylog.Check2(strconv.ParseInt(s, 10, 64))
}

func ParseUint(s string) uint64 {
	return mylog.Check2(strconv.ParseUint(s, 10, 64))
}

func Atoi(s string) int {
	return mylog.Check2(strconv.Atoi(s))
}

func IsTermux() bool {
	r0, e := os.Stat("/data/data/com.termux/files/usr")
	if e == nil {
		return r0.IsDir()
	}
	return false
}
func IsAix() bool       { return runtime.GOOS == "aix" }
func IsAndroid() bool   { return runtime.GOOS == "android" }
func IsDarwin() bool    { return runtime.GOOS == "darwin" }
func IsDragonfly() bool { return runtime.GOOS == "dragonfly" }
func IsFreebsd() bool   { return runtime.GOOS == "freebsd" }
func IsHurd() bool      { return runtime.GOOS == "hurd" }
func IsIllumos() bool   { return runtime.GOOS == "illumos" }
func IsIos() bool       { return runtime.GOOS == "ios" }
func IsJs() bool        { return runtime.GOOS == "js" }
func IsLinux() bool     { return runtime.GOOS == "linux" }
func IsNacl() bool      { return runtime.GOOS == "nacl" }
func IsNetbsd() bool    { return runtime.GOOS == "netbsd" }
func IsOpenbsd() bool   { return runtime.GOOS == "openbsd" }
func IsPlan9() bool     { return runtime.GOOS == "plan9" }
func IsSolaris() bool   { return runtime.GOOS == "solaris" }
func IsWasip1() bool    { return runtime.GOOS == "wasip1" }
func IsWindows() bool   { return runtime.GOOS == "windows" }
func IsZos() bool       { return runtime.GOOS == "zos" }

var knownOS = map[string]bool{
	"aix":       true,
	"android":   true,
	"darwin":    true,
	"dragonfly": true,
	"freebsd":   true,
	"hurd":      true,
	"illumos":   true,
	"ios":       true,
	"js":        true,
	"linux":     true,
	"nacl":      true,
	"netbsd":    true,
	"openbsd":   true,
	"plan9":     true,
	"solaris":   true,
	"wasip1":    true,
	"windows":   true,
	"zos":       true,
}

var unixOS = map[string]bool{
	"aix":       true,
	"android":   true,
	"darwin":    true,
	"dragonfly": true,
	"freebsd":   true,
	"hurd":      true,
	"illumos":   true,
	"ios":       true,
	"linux":     true,
	"netbsd":    true,
	"openbsd":   true,
	"solaris":   true,
}

var knownArch = map[string]bool{
	"386":         true,
	"amd64":       true,
	"amd64p32":    true,
	"arm":         true,
	"armbe":       true,
	"arm64":       true,
	"arm64be":     true,
	"loong64":     true,
	"mips":        true,
	"mipsle":      true,
	"mips64":      true,
	"mips64le":    true,
	"mips64p32":   true,
	"mips64p32le": true,
	"ppc":         true,
	"ppc64":       true,
	"ppc64le":     true,
	"riscv":       true,
	"riscv64":     true,
	"s390":        true,
	"s390x":       true,
	"sparc":       true,
	"sparc64":     true,
	"wasm":        true,
}

const TimeLayout = "2006-01-02 15:04:05"

func FormatTime(t time.Time) string { return t.Format(TimeLayout) }
func UnFormatTime(s string) time.Time {
	parse := mylog.Check2(time.Parse(TimeLayout, s))
	return parse
}

func FormatDuration(d time.Duration) string { return d.String() }
func UnFormatDuration(s string) time.Duration {
	duration := mylog.Check2(time.ParseDuration(s))
	return duration
}
func GetTimeNowString() string { return time.Now().Format("2006-01-02 15:04:05 ") }

func GetTimeStamp13Bits() int64 { return time.Now().UnixNano() / 1000000 }

func GetTimeStamp() string { return strconv.FormatInt(time.Now().UnixNano()/1000000, 10) }

func GetDiffDays(dstTime string) string {
	t := mylog.Check2(time.Parse("2006-01-02", dstTime))
	now := t.Sub(time.Now())
	days := int(now.Hours() / 24)
	years := days / 365
	months := (days % 365) / 30
	remainingDays := (days % 365) % 30
	hours := int(now.Hours()) % 24
	minutes := int(now.Minutes()) % 60
	seconds := int(now.Seconds()) % 60

	s := NewBuffer("")
	s.WriteStringLn(fmt.Sprintf("相差天数 %d 天", days))
	s.WriteStringLn(fmt.Sprintf("相差年数 %d 年", years))
	s.WriteStringLn(fmt.Sprintf("相差月数 %d 月", months))
	s.WriteStringLn(fmt.Sprintf("相差时数 %d 时", hours))
	s.WriteStringLn(fmt.Sprintf("相差分数 %d 分", minutes))
	s.WriteStringLn(fmt.Sprintf("相差秒数 %d 秒", seconds))
	s.WriteStringLn(fmt.Sprintf("相差时间 %d 年 %d 月 %d 天 %d 时 %d 分 %d 秒",
		years, months, remainingDays, hours, minutes, seconds))
	return s.String()
}

func GetUserConfigDirs() (UserConfigDirs map[string]string) {
	UserConfigDirs = make(map[string]string)
	if runtime.GOOS == "windows" {
		dir := mylog.Check2(os.UserConfigDir())
		u := mylog.Check2(user.Current())
		UserConfigDirs[u.Username] = dir
	} else if IsTermux() {
		dir := mylog.Check2(os.UserConfigDir())
		u := mylog.Check2(user.Current())
		UserConfigDirs[u.Username] = dir
	} else {
		file := mylog.Check2(os.Open("/etc/passwd"))
		defer func() { mylog.Check(file.Close()) }()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.Split(line, ":")
			if len(parts) > 0 {
				username := parts[0]
				u, e := user.Lookup(username)
				mylog.CheckIgnore(e)
				if e != nil {
					continue
				}
				dir := u.HomeDir + "/.config"
				if strings.Contains(dir, "root") || strings.Contains(dir, "home") {
					UserConfigDirs[username] = dir
				}
			}
		}
	}
	return UserConfigDirs
}

var RegexpCenter = `(.+?)`

func RegexpWebBodyBlocks(tagName string) string {
	return `<` + tagName + `[^>]*?>[\w\W]*?<\/` + tagName + `>`
}

func IntegerToIP(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

var (
	RegexpIp     = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))`)
	RegexpIpPort = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))):([0-9]+)`)
)
