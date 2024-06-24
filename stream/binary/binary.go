package binary

import (
	"errors"
	"io"
	"math"
	"reflect"
	"sync"

	"github.com/ddkwork/golibrary/mylog"
)

type ByteOrder interface {
	Uint16([]byte) uint16
	Uint32([]byte) uint32
	Uint64([]byte) uint64
	PutUint16(uint16) []byte
	PutUint32(uint32) []byte
	PutUint64(uint64) []byte
	String() string
}

var LittleEndian littleEndian

var BigEndian bigEndian

type littleEndian struct{}

func (littleEndian) Uint16(b []byte) uint16 {
	_ = b[1]
	return uint16(b[0]) | uint16(b[1])<<8
}

func (littleEndian) PutUint16(v uint16) (b []byte) {
	b = make([]byte, 2)
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	return
}

func (littleEndian) Uint32(b []byte) uint32 {
	_ = b[3]
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func (littleEndian) PutUint32(v uint32) (b []byte) {
	b = make([]byte, 4)
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	return
}

func (littleEndian) Uint64(b []byte) uint64 {
	_ = b[7]
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func (littleEndian) PutUint64(v uint64) (b []byte) {
	b = make([]byte, 8)
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
	return
}

func (littleEndian) String() string { return "LittleEndian" }

func (littleEndian) GoString() string { return "mybinary.LittleEndian" }

type bigEndian struct{}

func (bigEndian) Uint16(b []byte) uint16 {
	_ = b[1]
	return uint16(b[1]) | uint16(b[0])<<8
}

func (bigEndian) PutUint16(v uint16) (b []byte) {
	b = make([]byte, 2)
	b[0] = byte(v >> 8)
	b[1] = byte(v)
	return
}

func (bigEndian) Uint32(b []byte) uint32 {
	_ = b[3]
	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}

func (bigEndian) PutUint32(v uint32) (b []byte) {
	b = make([]byte, 4)
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
	return
}

func (bigEndian) Uint64(b []byte) uint64 {
	_ = b[7]
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}

func (bigEndian) PutUint64(v uint64) (b []byte) {
	b = make([]byte, 8)
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
	return
}

func (bigEndian) String() string { return "BigEndian" }

func (bigEndian) GoString() string { return "mybinary.BigEndian" }

func Read(r io.Reader, order ByteOrder, data any) error {
	if n := intDataSize(data); n != 0 {
		bs := make([]byte, n)
		mylog.Check2(io.ReadFull(r, bs))
		switch data := data.(type) {
		case *bool:
			*data = bs[0] != 0
		case *int8:
			*data = int8(bs[0])
		case *uint8:
			*data = bs[0]
		case *int16:
			*data = int16(order.Uint16(bs))
		case *uint16:
			*data = order.Uint16(bs)
		case *int32:
			*data = int32(order.Uint32(bs))
		case *uint32:
			*data = order.Uint32(bs)
		case *int64:
			*data = int64(order.Uint64(bs))
		case *uint64:
			*data = order.Uint64(bs)
		case *float32:
			*data = math.Float32frombits(order.Uint32(bs))
		case *float64:
			*data = math.Float64frombits(order.Uint64(bs))
		case []bool:
			for i, x := range bs {
				data[i] = x != 0
			}
		case []int8:
			for i, x := range bs {
				data[i] = int8(x)
			}
		case []uint8:
			copy(data, bs)
		case []int16:
			for i := range data {
				data[i] = int16(order.Uint16(bs[2*i:]))
			}
		case []uint16:
			for i := range data {
				data[i] = order.Uint16(bs[2*i:])
			}
		case []int32:
			for i := range data {
				data[i] = int32(order.Uint32(bs[4*i:]))
			}
		case []uint32:
			for i := range data {
				data[i] = order.Uint32(bs[4*i:])
			}
		case []int64:
			for i := range data {
				data[i] = int64(order.Uint64(bs[8*i:]))
			}
		case []uint64:
			for i := range data {
				data[i] = order.Uint64(bs[8*i:])
			}
		case []float32:
			for i := range data {
				data[i] = math.Float32frombits(order.Uint32(bs[4*i:]))
			}
		case []float64:
			for i := range data {
				data[i] = math.Float64frombits(order.Uint64(bs[8*i:]))
			}
		default:
			n = 0
		}
		if n != 0 {
			return nil
		}
	}

	v := reflect.ValueOf(data)
	size := -1
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
		size = dataSize(v)
	case reflect.Slice:
		size = dataSize(v)
	}
	if size < 0 {
		return errors.New("mybinary.Read: invalid type " + reflect.TypeOf(data).String())
	}
	d := &decoder{order: order, buf: make([]byte, size)}
	mylog.Check2(io.ReadFull(r, d.buf))
	d.value(v)
	return nil
}

func Write(w io.Writer, order ByteOrder, data any) error {
	if n := intDataSize(data); n != 0 {
		bs := make([]byte, n)
		switch v := data.(type) {
		case *bool:
			if *v {
				bs[0] = 1
			} else {
				bs[0] = 0
			}
		case bool:
			if v {
				bs[0] = 1
			} else {
				bs[0] = 0
			}
		case []bool:
			for i, x := range v {
				if x {
					bs[i] = 1
				} else {
					bs[i] = 0
				}
			}
		case *int8:
			bs[0] = byte(*v)
		case int8:
			bs[0] = byte(v)
		case []int8:
			for i, x := range v {
				bs[i] = byte(x)
			}
		case *uint8:
			bs[0] = *v
		case uint8:
			bs[0] = v
		case []uint8:
			bs = v
		case *int16:
			bs = order.PutUint16(uint16(*v))
		case int16:
			bs = order.PutUint16(uint16(v))
		case []int16:
			for i, x := range v {
				copy(bs[2*i:], order.PutUint16(uint16(x)))
			}
		case *uint16:
			bs = order.PutUint16(*v)
		case uint16:
			bs = order.PutUint16(v)
		case []uint16:
			for i, x := range v {
				copy(bs[2*i:], order.PutUint16(x))
			}
		case *int32:
			bs = order.PutUint32(uint32(*v))
		case int32:
			bs = order.PutUint32(uint32(v))
		case []int32:
			for i, x := range v {
				copy(bs[4*i:], order.PutUint32(uint32(x)))
			}
		case *uint32:
			bs = order.PutUint32(*v)
		case uint32:
			bs = order.PutUint32(v)
		case []uint32:
			for i, x := range v {
				copy(bs[4*i:], order.PutUint32(x))
			}
		case *int64:
			bs = order.PutUint64(uint64(*v))
		case int64:
			bs = order.PutUint64(uint64(v))
		case []int64:
			for i, x := range v {
				copy(bs[8*i:], order.PutUint64(uint64(x)))
			}
		case *uint64:
			bs = order.PutUint64(*v)
		case uint64:
			bs = order.PutUint64(v)
		case []uint64:
			for i, x := range v {
				copy(bs[8*i:], order.PutUint64(x))
			}
		case *float32:
			bs = order.PutUint32(math.Float32bits(*v))
		case float32:
			bs = order.PutUint32(math.Float32bits(v))
		case []float32:
			for i, x := range v {
				copy(bs[4*i:], order.PutUint32(math.Float32bits(x)))
			}
		case *float64:
			bs = order.PutUint64(math.Float64bits(*v))
		case float64:
			bs = order.PutUint64(math.Float64bits(v))
		case []float64:
			for i, x := range v {
				copy(bs[8*i:], order.PutUint64(math.Float64bits(x)))
			}
		}
		mylog.Check2(w.Write(bs))
		return nil
	}

	v := reflect.Indirect(reflect.ValueOf(data))
	size := dataSize(v)
	if size < 0 {
		return errors.New("mybinary.Write: invalid type " + reflect.TypeOf(data).String())
	}
	buf := make([]byte, size)
	e := &encoder{order: order, buf: buf}
	e.value(v)
	mylog.Check2(w.Write(buf))
	return nil
}

func Size(v any) int {
	return dataSize(reflect.Indirect(reflect.ValueOf(v)))
}

var structSize sync.Map

func dataSize(v reflect.Value) int {
	switch v.Kind() {
	case reflect.Slice:
		if s := sizeof(v.Type().Elem()); s >= 0 {
			return s * v.Len()
		}
		return -1

	case reflect.Struct:
		t := v.Type()
		if size, ok := structSize.Load(t); ok {
			return size.(int)
		}
		size := sizeof(t)
		structSize.Store(t, size)
		return size

	default:
		return sizeof(v.Type())
	}
}

func sizeof(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Array:
		if s := sizeof(t.Elem()); s >= 0 {
			return s * t.Len()
		}

	case reflect.Struct:
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			s := sizeof(t.Field(i).Type)
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return int(t.Size())
	}

	return -1
}

type coder struct {
	order  ByteOrder
	buf    []byte
	offset int
}

type (
	decoder coder
	encoder coder
)

func (d *decoder) bool() bool {
	x := d.buf[d.offset]
	d.offset++
	return x != 0
}

func (e *encoder) bool(x bool) {
	if x {
		e.buf[e.offset] = 1
	} else {
		e.buf[e.offset] = 0
	}
	e.offset++
}

func (d *decoder) uint8() uint8 {
	x := d.buf[d.offset]
	d.offset++
	return x
}

func (e *encoder) uint8(x uint8) {
	e.buf[e.offset] = x
	e.offset++
}

func (d *decoder) uint16() uint16 {
	x := d.order.Uint16(d.buf[d.offset : d.offset+2])
	d.offset += 2
	return x
}

func (e *encoder) uint16(x uint16) {
	copy(e.buf[e.offset:e.offset+2], e.order.PutUint16(x))
	e.offset += 2
}

func (d *decoder) uint32() uint32 {
	x := d.order.Uint32(d.buf[d.offset : d.offset+4])
	d.offset += 4
	return x
}

func (e *encoder) uint32(x uint32) {
	copy(e.buf[e.offset:e.offset+4], e.order.PutUint32(x))
	e.offset += 4
}

func (d *decoder) uint64() uint64 {
	x := d.order.Uint64(d.buf[d.offset : d.offset+8])
	d.offset += 8
	return x
}

func (e *encoder) uint64(x uint64) {
	copy(e.buf[e.offset:e.offset+8], e.order.PutUint64(x))
	e.offset += 8
}

func (d *decoder) int8() int8 { return int8(d.uint8()) }

func (e *encoder) int8(x int8) { e.uint8(uint8(x)) }

func (d *decoder) int16() int16 { return int16(d.uint16()) }

func (e *encoder) int16(x int16) { e.uint16(uint16(x)) }

func (d *decoder) int32() int32 { return int32(d.uint32()) }

func (e *encoder) int32(x int32) { e.uint32(uint32(x)) }

func (d *decoder) int64() int64 { return int64(d.uint64()) }

func (e *encoder) int64(x int64) { e.uint64(uint64(x)) }

func (d *decoder) value(v reflect.Value) {
	switch v.Kind() {
	case reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			d.value(v.Index(i))
		}

	case reflect.Struct:
		t := v.Type()
		l := v.NumField()
		for i := 0; i < l; i++ {
			if v := v.Field(i); v.CanSet() || t.Field(i).Name != "_" {
				d.value(v)
			} else {
				d.skip(v)
			}
		}

	case reflect.Slice:
		l := v.Len()
		for i := 0; i < l; i++ {
			d.value(v.Index(i))
		}

	case reflect.Bool:
		v.SetBool(d.bool())

	case reflect.Int8:
		v.SetInt(int64(d.int8()))
	case reflect.Int16:
		v.SetInt(int64(d.int16()))
	case reflect.Int32:
		v.SetInt(int64(d.int32()))
	case reflect.Int64:
		v.SetInt(d.int64())

	case reflect.Uint8:
		v.SetUint(uint64(d.uint8()))
	case reflect.Uint16:
		v.SetUint(uint64(d.uint16()))
	case reflect.Uint32:
		v.SetUint(uint64(d.uint32()))
	case reflect.Uint64:
		v.SetUint(d.uint64())

	case reflect.Float32:
		v.SetFloat(float64(math.Float32frombits(d.uint32())))
	case reflect.Float64:
		v.SetFloat(math.Float64frombits(d.uint64()))

	case reflect.Complex64:
		v.SetComplex(complex128(complex(
			math.Float32frombits(d.uint32()),
			math.Float32frombits(d.uint32()),
		)))
	case reflect.Complex128:
		v.SetComplex(complex(
			math.Float64frombits(d.uint64()),
			math.Float64frombits(d.uint64()),
		))
	}
}

func (e *encoder) value(v reflect.Value) {
	switch v.Kind() {
	case reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			e.value(v.Index(i))
		}

	case reflect.Struct:
		t := v.Type()
		l := v.NumField()
		for i := 0; i < l; i++ {
			if v := v.Field(i); v.CanSet() || t.Field(i).Name != "_" {
				e.value(v)
			} else {
				e.skip(v)
			}
		}

	case reflect.Slice:
		l := v.Len()
		for i := 0; i < l; i++ {
			e.value(v.Index(i))
		}

	case reflect.Bool:
		e.bool(v.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v.Type().Kind() {
		case reflect.Int8:
			e.int8(int8(v.Int()))
		case reflect.Int16:
			e.int16(int16(v.Int()))
		case reflect.Int32:
			e.int32(int32(v.Int()))
		case reflect.Int64:
			e.int64(v.Int())
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch v.Type().Kind() {
		case reflect.Uint8:
			e.uint8(uint8(v.Uint()))
		case reflect.Uint16:
			e.uint16(uint16(v.Uint()))
		case reflect.Uint32:
			e.uint32(uint32(v.Uint()))
		case reflect.Uint64:
			e.uint64(v.Uint())
		}

	case reflect.Float32, reflect.Float64:
		switch v.Type().Kind() {
		case reflect.Float32:
			e.uint32(math.Float32bits(float32(v.Float())))
		case reflect.Float64:
			e.uint64(math.Float64bits(v.Float()))
		}

	case reflect.Complex64, reflect.Complex128:
		switch v.Type().Kind() {
		case reflect.Complex64:
			x := v.Complex()
			e.uint32(math.Float32bits(float32(real(x))))
			e.uint32(math.Float32bits(float32(imag(x))))
		case reflect.Complex128:
			x := v.Complex()
			e.uint64(math.Float64bits(real(x)))
			e.uint64(math.Float64bits(imag(x)))
		}
	}
}

func (d *decoder) skip(v reflect.Value) {
	d.offset += dataSize(v)
}

func (e *encoder) skip(v reflect.Value) {
	n := dataSize(v)
	zero := e.buf[e.offset : e.offset+n]
	for i := range zero {
		zero[i] = 0
	}
	e.offset += n
}

func intDataSize(data any) int {
	switch data := data.(type) {
	case bool, int8, uint8, *bool, *int8, *uint8:
		return 1
	case []bool:
		return len(data)
	case []int8:
		return len(data)
	case []uint8:
		return len(data)
	case int16, uint16, *int16, *uint16:
		return 2
	case []int16:
		return 2 * len(data)
	case []uint16:
		return 2 * len(data)
	case int32, uint32, *int32, *uint32:
		return 4
	case []int32:
		return 4 * len(data)
	case []uint32:
		return 4 * len(data)
	case int64, uint64, *int64, *uint64:
		return 8
	case []int64:
		return 8 * len(data)
	case []uint64:
		return 8 * len(data)
	case float32, *float32:
		return 4
	case float64, *float64:
		return 8
	case []float32:
		return 4 * len(data)
	case []float64:
		return 8 * len(data)
	}
	return 0
}
