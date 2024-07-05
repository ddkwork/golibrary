package binary

import (
	"errors"
	"io"

	"github.com/ddkwork/golibrary/mylog"
)

const (
	MaxVarintLen16 = 3
	MaxVarintLen32 = 5
	MaxVarintLen64 = 10
)

type (
	varintStream interface {
		Size() int
		SetSize(size int)
		Bytes() []byte
		SetBytes(bytes []byte)
		Value() uint64
		SetValue(value uint64)
		Ok() bool
		SetOk(ok bool)
	}
	_varintStream struct {
		size  int
		bytes []byte
		value uint64
		ok    bool
	}
)

func (v *_varintStream) Bytes() []byte {
	return v.bytes
}

func (v *_varintStream) SetBytes(bytes []byte) {
	v.bytes = bytes
}

func (v *_varintStream) Size() int {
	return v.size
}

func (v *_varintStream) SetSize(size int) {
	v.size = size
}

func (v *_varintStream) Value() uint64 {
	return v.value
}

func (v *_varintStream) SetValue(value uint64) {
	v.value = value
}

func (v *_varintStream) Ok() bool {
	return v.ok
}

func (v *_varintStream) SetOk(ok bool) {
	v.ok = ok
}

func New() varintStream {
	return &_varintStream{}
}

type (
	varint interface {
		PutUvarint(uint64) varintStream
		Uvarint([]byte) varintStream
		PutVarint(int64) varintStream
		Varint([]byte) varintStream
		ReadUvarint(io.ByteReader) varintStream
		ReadVarint(io.ByteReader) varintStream
	}
	_varint struct {
		s varintStream
	}
)

func (v *_varint) PutUvarint(value uint64) varintStream {
	i := 0
	bytes := make([]byte, MaxVarintLen64)
	for value >= 0x80 {
		bytes[i] = byte(value) | 0x80
		value >>= 7
		i++
	}
	bytes[i] = byte(value)
	v.s.SetSize(i + 1)
	v.s.SetBytes(bytes)
	return v.s
}

func (v *_varint) Uvarint(bytes []byte) varintStream {
	var (
		value uint64
		shift uint
	)
	for i, b := range bytes {
		if i == MaxVarintLen64 {
			return nil
		}
		if b < 0x80 {
			if i == MaxVarintLen64-1 && b > 1 {
				return nil
			}
			v.s.SetValue(value | uint64(b)<<shift)
			v.s.SetSize(i + 1)
			return v.s
		}
		value |= uint64(b&0x7f) << shift
		shift += 7
	}
	return nil
}

func (v *_varint) PutVarint(value int64) varintStream {
	unsignedValue := uint64(value) << 1
	if value < 0 {
		unsignedValue = ^unsignedValue
	}
	return v.PutUvarint(unsignedValue)
}

func (v *_varint) Varint(bytes []byte) varintStream {
	uvarint := v.Uvarint(bytes)
	value := int64(uvarint.Value() >> 1)
	if uvarint.Value()&1 != 0 {
		value = ^value
	}
	v.s.SetValue(uint64(value))
	v.s.SetSize(uvarint.Size())
	return v.s
}

func (v *_varint) ReadUvarint(reader io.ByteReader) varintStream {
	var value uint64
	var shift uint
	for i := 0; i < MaxVarintLen64; i++ {
		readByte := mylog.Check2(reader.ReadByte())
		if readByte < 0x80 {
			if i == MaxVarintLen64-1 && readByte > 1 {
				return nil
			}
			v.s.SetValue(value | uint64(readByte)<<shift)
			return v.s
		}
		value |= uint64(readByte&0x7f) << shift
		shift += 7
	}
	return nil
}

func (v *_varint) ReadVarint(reader io.ByteReader) varintStream {
	uvarint := v.ReadUvarint(reader)
	value := int64(uvarint.Value() >> 1)
	if uvarint.Value()&1 != 0 {
		value = ^value
	}
	v.s.SetValue(uint64(value))
	return v.s
}

func new_varint() varint {
	return &_varint{
		s: New(),
	}
}

func PutUvarint(buf []byte, x uint64) int {
	i := 0
	for x >= 0x80 {
		buf[i] = byte(x) | 0x80
		x >>= 7
		i++
	}
	buf[i] = byte(x)
	return i + 1
}

func Uvarint(buf []byte) (uint64, int) {
	var x uint64
	var s uint
	for i, b := range buf {
		if i == MaxVarintLen64 {
			return 0, -(i + 1)
		}
		if b < 0x80 {
			if i == MaxVarintLen64-1 && b > 1 {
				return 0, -(i + 1)
			}
			return x | uint64(b)<<s, i + 1
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return 0, 0
}

func PutVarint(buf []byte, x int64) int {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return PutUvarint(buf, ux)
}

func Varint(buf []byte) (int64, int) {
	ux, n := Uvarint(buf)
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, n
}

var overflow = errors.New("mybinary: varintStream overflows a 64-bit integer")

func ReadUvarint(r io.ByteReader) (uint64, error) {
	var x uint64
	var s uint
	for i := 0; i < MaxVarintLen64; i++ {
		b, e := r.ReadByte()
		if mylog.CheckEof(e) {
			return x, nil
		}
		if b < 0x80 {
			if i == MaxVarintLen64-1 && b > 1 {
				return x, overflow
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return x, overflow
}

func ReadVarint(r io.ByteReader) (int64, error) {
	ux := mylog.Check2(ReadUvarint(r))
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, nil
}
