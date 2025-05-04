package deepcopy

import (
	"testing"
	"time"
	"unsafe"

	"github.com/ddkwork/golibrary/assert"
	"github.com/google/go-cmp/cmp"
)

func TestSimple(t *testing.T) {
	Strings := []string{"a", "b", "c"}
	cpyS := Clone(Strings)
	assert.Equal(t, cpyS, Strings)

	Bools := []bool{true, true, false, false}
	cpyB := Clone(Bools)
	assert.Equal(t, cpyB, Bools)

	Bytes := []byte("hello")
	cpyBt := Clone(Bytes)
	assert.Equal(t, cpyBt, Bytes)

	Ints := []int{42}
	cpyI := Clone(Ints)
	assert.Equal(t, cpyI, Ints)

	Uints := []uint{1, 2, 3, 4, 5}
	cpyU := Clone(Uints)
	assert.Equal(t, cpyU, Uints)

	Float32s := []float32{3.14}
	cpyF := Clone(Float32s)
	assert.Equal(t, cpyF, Float32s)

	Interfaces := []any{"a", 42, true, 4.32}
	cpyIf := Clone(Interfaces)
	assert.Equal(t, cpyIf, Interfaces)
}

type Basics struct {
	String      string
	Strings     []string
	StringArr   [4]string
	Bool        bool
	Bools       []bool
	Byte        byte
	Bytes       []byte
	Int         int
	Ints        []int
	Int8        int8
	Int8s       []int8
	Int16       int16
	Int16s      []int16
	Int32       int32
	Int32s      []int32
	Int64       int64
	Int64s      []int64
	Uint        uint
	Uints       []uint
	Uint8       uint8
	Uint8s      []uint8
	Uint16      uint16
	Uint16s     []uint16
	Uint32      uint32
	Uint32s     []uint32
	Uint64      uint64
	Uint64s     []uint64
	Float32     float32
	Float32s    []float32
	Float64     float64
	Float64s    []float64
	Complex64   complex64
	Complex64s  []complex64
	Complex128  complex128
	Complex128s []complex128
	Interface   any
	Interfaces  []any
}

func TestMostTypes(t *testing.T) {
	test := Basics{
		String:      "kimchi",
		Strings:     []string{"uni", "ika"},
		StringArr:   [4]string{"malort", "barenjager", "fernet", "salmiakki"},
		Bool:        true,
		Bools:       []bool{true, false, true},
		Byte:        'z',
		Bytes:       []byte("abc"),
		Int:         42,
		Ints:        []int{0, 1, 3, 4},
		Int8:        8,
		Int8s:       []int8{8, 9, 10},
		Int16:       16,
		Int16s:      []int16{16, 17, 18, 19},
		Int32:       32,
		Int32s:      []int32{32, 33},
		Int64:       64,
		Int64s:      []int64{64},
		Uint:        420,
		Uints:       []uint{11, 12, 13},
		Uint8:       81,
		Uint8s:      []uint8{81, 82},
		Uint16:      160,
		Uint16s:     []uint16{160, 161, 162, 163, 164},
		Uint32:      320,
		Uint32s:     []uint32{320, 321},
		Uint64:      640,
		Uint64s:     []uint64{6400, 6401, 6402, 6403},
		Float32:     32.32,
		Float32s:    []float32{32.32, 33},
		Float64:     64.1,
		Float64s:    []float64{64, 65, 66},
		Complex64:   complex64(-64 + 12i),
		Complex64s:  []complex64{complex64(-65 + 11i), complex64(66 + 10i)},
		Complex128:  -128 + 12i,
		Complex128s: []complex128{-128 + 11i, 129 + 10i},
		Interfaces:  []any{42, true, "pan-galactic"},
	}

	cpy := Clone(test)
	assert.Equal(t, cpy, test)
}

func TestComplexSlices(t *testing.T) {
	orig3Int := [][][]int{{{1, 2, 3}, {11, 22, 33}}, {{7, 8, 9}, {66, 77, 88, 99}}}
	cpyI := Clone(orig3Int)
	assert.False(t, unsafe.SliceData(orig3Int) == unsafe.SliceData(cpyI))
	assert.Equal(t, cpyI, orig3Int)
}

type A struct {
	Int    int
	String string
	UintSl []uint
	NilSl  []string
	Map    map[string]int
	MapB   map[string]*B
	SliceB []B
	B
	T time.Time
}

type B struct {
	Vals []string
}

var AStruct = A{
	Int:    42,
	String: "Konichiwa",
	UintSl: []uint{0, 1, 2, 3},
	NilSl:  nil,
	Map:    map[string]int{"a": 1, "b": 2},
	MapB: map[string]*B{
		"hi":  {Vals: []string{"hello", "bonjour"}},
		"bye": {Vals: []string{"good-bye", "au revoir"}},
	},
	SliceB: []B{
		{Vals: []string{"Ciao", "Aloha"}},
	},
	B: B{Vals: []string{"42"}},
	T: time.Now(),
}

func TestStructA(t *testing.T) {
	cpy := Clone(AStruct)
	assert.Equal(t, cpy, AStruct)
}

type Unexported struct {
	A  string
	B  int
	aa string
	bb int
	cc []int
	dd map[string]string
}

func TestUnexportedFields(t *testing.T) {
	t.Skip("tree to json, Unexported fields are not supported")
	u := &Unexported{
		A:  "A",
		B:  42,
		aa: "aa",
		bb: 42,
		cc: []int{1, 2, 3},
		dd: map[string]string{"hello": "bonjour"},
	}
	cpy := Clone(u)
	assert.Equal(t, cpy, u, cmp.AllowUnexported())
}

type T struct {
	time.Time
}

func TestTimeCopy(t *testing.T) {
	tests := []struct {
		Y    int
		M    time.Month
		D    int
		h    int
		m    int
		s    int
		nsec int
		TZ   string
	}{
		{2016, time.July, 4, 23, 11, 33, 3000, "America/New_York"},
		{2015, time.October, 31, 9, 44, 23, 45935, "UTC"},
		{2014, time.May, 5, 22, 0o1, 50, 219300, "Europe/Prague"},
	}

	for i, test := range tests {
		l, e := time.LoadLocation(test.TZ)
		if e != nil {
			t.Errorf("%d: unexpected error: %s", i, e)
			continue
		}
		var x T
		x.Time = time.Date(test.Y, test.M, test.D, test.h, test.m, test.s, test.nsec, l)
		c := Clone(x)
		assert.Equal(t, c, x)
	}
}

func TestPointerToStruct(t *testing.T) {
	type Foo struct {
		Bar int
	}
	f := &Foo{Bar: 42}
	cpy := Clone(f)
	assert.Equal(t, f, cpy)
}

func TestIssue9(t *testing.T) {
	x := 42
	testA := map[string]*int{
		"a": nil,
		"b": &x,
	}
	copyA := Clone(testA)
	assert.Equal(t, testA, copyA)

	type Foo struct {
		Alpha string
	}

	type Bar struct {
		Beta  string
		Gamma int
		Delta *Foo
	}

	type Biz struct {
		Epsilon map[int]*Bar
	}

	testB := Biz{
		Epsilon: map[int]*Bar{
			0: {},
			1: {
				Beta:  "don't panic",
				Gamma: 42,
				Delta: nil,
			},
			2: {
				Beta:  "sudo make me a sandwich.",
				Gamma: 11,
				Delta: &Foo{
					Alpha: "okay.",
				},
			},
		},
	}

	copyB := Clone(testB)
	assert.Equal(t, testB, copyB)

	// testC := map[*Foo][]string{
	//	{Alpha: "Henry Dorsett Case"}: {
	//		"Cutter",
	//	},
	//	{Alpha: "Molly Millions"}: {
	//		"Rose Kolodny",
	//		"Cat Mother",
	//		"Steppin' Razor",
	//	},
	// }
	//
	// copyC := Clone(testC)
	// assert.Equal(t, testC, copyC) //差异在于 map 中的键的内存地址不同
	//
	// type Bizz struct {
	//	*Foo
	// }
	//
	// testD := map[Bizz]string{
	//	{&Foo{"Neuromancer"}}: "Rio",
	//	{&Foo{"Wintermute"}}:  "Berne",
	// }
	// copyD := Clone(testD)
	// assert.Equal(t, testD, copyD) //差异在于 map 中的键的内存地址不同
}

type I struct {
	A string
}

func (i *I) Clone() any {
	return &I{A: "custom copy"}
}

type NestI struct {
	I *I
}

func TestInterface(t *testing.T) {
	i := &I{A: "A"}
	copied := Clone(i)
	assert.Equal(t, "custom copy", copied.A)

	ni := &NestI{I: &I{A: "A"}}
	copiedNest := Clone(ni)
	assert.Equal(t, "custom copy", copiedNest.I.A)
}
