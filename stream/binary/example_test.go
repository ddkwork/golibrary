package binary_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/ddkwork/golibrary/mylog"
)

func ExampleWrite() {
	buf := new(bytes.Buffer)
	pi := math.Pi
	mylog.Check(binary.Write(buf, binary.LittleEndian, pi))
	fmt.Printf("% x", buf.Bytes())
}

func ExampleWrite_multi() {
	buf := new(bytes.Buffer)
	data := []any{
		uint16(61374),
		int8(-54),
		uint8(254),
	}
	for _, v := range data {
		mylog.Check(binary.Write(buf, binary.LittleEndian, v))
	}
	fmt.Printf("%x", buf.Bytes())
}

func ExampleRead() {
	var pi float64
	b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
	buf := bytes.NewReader(b)
	mylog.Check(binary.Read(buf, binary.LittleEndian, &pi))
	fmt.Print(pi)
}

func ExampleRead_multi() {
	b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40, 0xff, 0x01, 0x02, 0x03, 0xbe, 0xef}
	r := bytes.NewReader(b)

	var data struct {
		PI   float64
		Uate uint8
		Mine [3]byte
		Too  uint16
	}

	mylog.Check(binary.Read(r, binary.LittleEndian, &data))

	fmt.Println(data.PI)
	fmt.Println(data.Uate)
	fmt.Printf("% x\n", data.Mine)
	fmt.Println(data.Too)
}

func ExampleByteOrder_put() {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b[0:], 0x03e8)
	binary.LittleEndian.PutUint16(b[2:], 0x07d0)
	fmt.Printf("% x\n", b)
}

func ExampleByteOrder_get() {
	b := []byte{0xe8, 0x03, 0xd0, 0x07}
	x1 := binary.LittleEndian.Uint16(b[0:])
	x2 := binary.LittleEndian.Uint16(b[2:])
	fmt.Printf("%#04x %#04x\n", x1, x2)
}

func ExamplePutUvarint() {
	buf := make([]byte, binary.MaxVarintLen64)

	for _, x := range []uint64{1, 2, 127, 128, 255, 256} {
		n := binary.PutUvarint(buf, x)
		fmt.Printf("%x\n", buf[:n])
	}
}

func ExamplePutVarint() {
	buf := make([]byte, binary.MaxVarintLen64)

	for _, x := range []int64{-65, -64, -2, -1, 0, 1, 2, 63, 64} {
		n := binary.PutVarint(buf, x)
		fmt.Printf("%x\n", buf[:n])
	}
}

func ExampleUvarint() {
	inputs := [][]byte{
		{0x01},
		{0x02},
		{0x7f},
		{0x80, 0x01},
		{0xff, 0x01},
		{0x80, 0x02},
	}
	for _, b := range inputs {
		x, n := binary.Uvarint(b)
		if n != len(b) {
			fmt.Println("Uvarint did not consume all of in")
		}
		fmt.Println(x)
	}
}

func ExampleVarint() {
	inputs := [][]byte{
		{0x81, 0x01},
		{0x7f},
		{0x03},
		{0x01},
		{0x00},
		{0x02},
		{0x04},
		{0x7e},
		{0x80, 0x01},
	}
	for _, b := range inputs {
		x, n := binary.Varint(b)
		if n != len(b) {
			fmt.Println("Varint did not consume all of in")
		}
		fmt.Println(x)
	}
}
