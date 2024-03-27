package stream

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/des"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/rivo/uniseg"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safeType"
)

func (o *Stream) ReaderGzip() *Stream {
	b := ReaderGzip(o.Bytes()).Bytes()
	o.Reset()
	if !mylog.Error2(o.Write(b)) {
		return o // todo write err info ?
	}
	return o
}

func ReaderGzip[T safeType.Type](data T) *safeType.Data {
	reader, err := gzip.NewReader(bytes.NewReader(safeType.New(data).Bytes()))
	if !mylog.Error(err) {
		return nil
	}
	b := make([]byte, 1024*2)
	n, err := reader.Read(b)
	if !mylog.Error(err) {
		return nil
	}
	return safeType.New(b[:n])
}

func (o *Stream) SplitBytes(size int) (blocks [][]byte) {
	blocks = make([][]byte, 0)
	quantity := o.Len() / size
	remainder := o.Len() % size
	for i := range quantity {
		blocks = append(blocks, o.Bytes()[i*size:(i+1)*size])
		if remainder != 0 {
			blocks = append(blocks, o.Bytes()[i*size:i*size+remainder])
		}
	}
	return
}

func (o *Stream) SplitString(size int) (blocks []string) {
	blocks = make([]string, 0)
	splitBytes := o.SplitBytes(size)
	for _, splitByte := range splitBytes {
		blocks = append(blocks, string(splitByte))
	}
	return
}

func (o *Stream) WritePackageName() {
	o.WriteStringLn("package " + GetPackageName())
}

func GetPackageName() (pkgName string) {
	absPath, err := filepath.Abs(".")
	if !mylog.Error(err) {
		return
	}
	return filepath.Base(absPath)
}

func (o *Stream) CutString(left, right string) (cut string, found bool) {
	_, after, found := strings.Cut(o.String(), left)
	if !found {
		return
	}
	before, _, f := strings.Cut(after, right)
	if !f {
		return
	}
	before = strings.ReplaceAll(before, "\\", "")
	before = strings.ReplaceAll(before, "/", "")
	return before, true
}
func (o *Stream) CutWithIndex(x, y int) []byte { return o.Bytes()[x:y] }
func (o *Stream) NewLine() *Stream             { o.WriteString("\n"); return o }
func (o *Stream) QuoteWith(s string) *Stream   { o.WriteString(s); return o }
func (o *Stream) WriteBinary(order binary.ByteOrder, data any) bool {
	return mylog.Error(binary.Write(o, order, data))
}

func (o *Stream) ReadBinary(order binary.ByteOrder) (data any) {
	if !mylog.Error(binary.Read(o, order, data)) {
		return nil
	}
	return o
}

func (o *Stream) WriteBytesLn(b []byte) *Stream {
	o.Write(b)
	o.NewLine()
	return o
}

func (o *Stream) ReplaceAll(old, new string) *Stream {
	s := strings.ReplaceAll(o.String(), old, new)
	o.Reset()
	o.WriteString(s)
	return o
}

func (o *Stream) WriteStringLn(s string) *Stream {
	o.WriteString(s)
	o.NewLine()
	return o
}
func (o *Stream) Quote()          { o.WriteByte('"') }
func (o *Stream) ObjectBegin()    { o.WriteByte('{') }
func (o *Stream) ObjectEnd()      { o.WriteByte('}') }
func (o *Stream) SliceBegin()     { o.WriteByte('[') }
func (o *Stream) SliceEnd()       { o.WriteByte(']') }
func (o *Stream) Indent(deep int) { o.WriteString(strings.Repeat(" ", deep)) }

func (o *Stream) CheckDesBlockSize() bool {
	switch o.Len() {
	case 0:
		return mylog.Error("buffer len == 0")
	default:
		if o.Len()%des.BlockSize != 0 {
			return mylog.Error("Len()%des.BlockSize != 0")
		}
	}
	return true
}

func Concat[S ~[]E, E any](slices_ ...S) S { return slices.Concat(slices_...) }
func (o *Stream) MergeByte(streams ...[]byte) []byte {
	//return slices.Concat( streams...)
	for _, b := range streams {
		o.Write(b)
	}
	return o.Bytes()
}

func (o *Stream) BigNumXorWithAlign(arg1, arg2 []byte, align int) (xorStream []byte) {
	xor := new(big.Int).Xor(new(big.Int).SetBytes(arg1), new(big.Int).SetBytes(arg2))
	alignBuf := make([]byte, align-len(xor.Bytes()))
	switch len(xor.Bytes()) {
	case 0:
		xorStream = alignBuf
	case align:
		xorStream = xor.Bytes()
	default:
		o.AppendByteSlice(xorStream, alignBuf, xor.Bytes()) // todo test
	}
	return o.Bytes()
}

func (o *Stream) InsertStringWithSplit(size int, insert string) string {
	blocks := o.SplitString(size)
	o.Reset()
	for i, block := range blocks {
		o.WriteString(block)
		if i < len(blocks)-1 {
			o.WriteString(insert)
		}
	}
	return o.String()
}

func (o *Stream) InsertBytes(index int, insert []byte) []byte {
	start := o.Bytes()[:index]
	end := o.Bytes()[index:]
	o.Reset()
	o.Write(start)
	o.Write(insert)
	o.Write(end)
	return o.Bytes()
}

func (o *Stream) InsertByte(index int, ch byte) { o.InsertBytes(index, []byte{ch}) }

func (o *Stream) InsertRune(index int, r rune) {
	if uint32(r) < utf8.RuneSelf {
		o.InsertByte(index, byte(r))
		return
	}
	var buffer [4]byte
	n := utf8.EncodeRune(buffer[:], r)
	o.InsertBytes(index, buffer[:n])
}

func (o *Stream) InsertString(index int, s string) string {
	return string(o.InsertBytes(index, []byte(s)))
}

func (o *Stream) AppendByteSlice(bytesSlice ...[]byte) []byte {
	for _, slice := range bytesSlice {
		o.Write(slice)
	}
	return o.Bytes()
}

func (o *Stream) Contains(substr string) bool { return strings.Contains(o.String(), substr) }

// Lines returns the lines in the file x, including newlines.
// If the file does not end in a newline, one is supplied
// along with a warning about the missing newline.
func Lines(x []byte) []string {
	l := strings.SplitAfter(string(x), "\n")
	if l == nil {
		return nil
	}
	if l[len(l)-1] == "" {
		l = l[:len(l)-1]
	} else {
		// Treat last line as having a message about the missing newline attached,
		// using the same text as BSD/GNU diff (including the leading backslash).
		l[len(l)-1] += "\n\\ No newline at end of file\n"
	}
	return l
}

func (o *Stream) ToLines() (lines []string, ok bool) {
	lines = make([]string, 0)
	// open C:\Users\Admin\Downloads\Compressed\xx\Driver-SoulExtraction.sln: The system cannot find the file specified.
	if strings.Contains(o.String(), "The system cannot find the file specified") {
		return
	}
	newScanner := bufio.NewScanner(o.Buffer)
	for newScanner.Scan() {
		lines = append(lines, newScanner.Text())
	}
	ok = true
	return
	newReader := bufio.NewReader(o.Buffer)
	for {
		line, _, err := newReader.ReadLine()
		switch err {
		case io.EOF:
			return lines, true
		default:
			if !mylog.Error(err) {
				return
			}
		}
		lines = append(lines, string(line))
	}
}

// SplitFileByLines 按行数切割文件
func SplitFileByLines(filePath string, size int) {
	lines, ok := ReadFileToLines(filePath)
	if !ok {
		return
	}
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
		WriteTruncate(fmt.Sprint(i)+".txt", New("").LinesToString(lines[startIndex:endIndex]))
	}
}

// SerialNumber 这个才是swap功能，前后字节交换顺序
func (o *Stream) SerialNumber() string {
	to := bytes.Buffer{}
	for k, v := range o.Bytes() {
		if k%2 == 1 {
			to.WriteByte(v)
			to.WriteByte(o.Bytes()[k-1])
		}
	}
	return to.String()
}

func (o *Stream) CutUint16(u uint16) uint8 { // 6613-->16
	u >>= 4
	high4 := u & 0x0f << 4 // 60
	low4 := u & 0xf0 >> 4  // 01
	high4 |= low4          // 16
	return uint8(high4)
}

func CutUint16(u uint16) uint8 {
	ss := fmt.Sprintf("%x", u)
	// 6 6 1 3
	// 0 1 2 3
	out := string(ss[2]) + string(ss[1])
	decodeString, err := hex.DecodeString(out)
	if !mylog.Error(err) {
		return 0
	}
	if decodeString == nil {
		return 0
	}
	return decodeString[0]
}

// SwapBytes2HexString2 逆序，并不是交换 todo 重命名为 Reverse，如果输入参数复合8字节以下的字节对齐可以用二进制包，不过不通用
func SwapBytes2HexString2(src safeType.HexString) (dst string) {
	hexString := safeType.NewHexString(src)
	slices.Reverse(hexString.Bytes())
	return hex.EncodeToString(hexString.Bytes())
}

func SwapBytes2HexString(src []byte) (dst safeType.HexString) {
	return safeType.New(SwapBytes(src)).HexString()
}

func SwapBytes(src []byte) (dst []byte) {
	slices.Reverse(src)
	return src
	return
	// binary.BigEndian.PutUint64()三次才行，10字节的话，如果是更多字节那就不通用
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

func SlicesIndex(slice any, item any) int { // slices.Index
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
	// Create a new slice with the same type as the input slice.
	// The new slice has length + 1, and capacity + 1 to accommodate the new element.
	t := reflect.MakeSlice(s.Type(), s.Len()+1, s.Cap()+1)
	// Copy the elements before the insertion point to the new slice.
	reflect.Copy(t.Slice(0, index), s.Slice(0, index))
	// Insert the new element at the insertion point.
	t.Index(index).Set(reflect.ValueOf(value))
	// Copy the elements after the insertion point to the new slice.
	reflect.Copy(t.Slice(index+1, s.Len()+1), s.Slice(index, s.Len()))
	return t.Interface()
}
