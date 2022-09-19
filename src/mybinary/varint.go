// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mybinary

// This file implements "varintStream" encoding of 64-bit integers.
// The encoding is:
// - unsigned integers are serialized 7 bits at a time, starting with the
//   least significant bits
// - the most significant bit (msb) in each output byte indicates if there
//   is a continuation byte (msb = 1)
// - signed integers are mapped to unsigned integers using "zig-zag"
//   encoding: Positive values x are written as 2*x + 0, negative values
//   are written as 2*(^x) + 1; that is, negative numbers are complemented
//   and whether to complement is encoded in bit 0.
//
// Design note:
// At most 10 bytes are needed for 64-bit values. The encoding could
// be more dense: a full 64-bit value needs an extra byte just to hold bit 63.
// Instead, the msb of the previous byte could be used to hold bit 63 since we
// know there can't be more than 64 bits. This is a trivial improvement and
// would reduce the maximum encoding length to 9 bytes. However, it breaks the
// invariant that the msb is always the "continuation bit" and thus makes the
// format incompatible with a varintStream encoding for larger numbers (say 128-bit).

import (
	"errors"
	"github.com/ddkwork/golibrary/mylog"

	"io"
)

// MaxVarintLenN is the maximum length of a varintStream-encoded N-bit integer.
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
			mylog.Error("i == MaxVarintLen64")
			return nil
		}
		if b < 0x80 {
			if i == MaxVarintLen64-1 && b > 1 {
				mylog.Error("i == MaxVarintLen64-1 && b > 1 ")
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
	v.s.SetValue(uint64(value)) //todo add int64 value field
	v.s.SetSize(uvarint.Size())
	return v.s
}

func (v *_varint) ReadUvarint(reader io.ByteReader) varintStream {
	var value uint64
	var shift uint
	for i := 0; i < MaxVarintLen64; i++ {
		readByte, err := reader.ReadByte()
		if !mylog.Error(err) {
			return nil
		}
		if readByte < 0x80 {
			if i == MaxVarintLen64-1 && readByte > 1 {
				mylog.Error("i == MaxVarintLen64-1 && readByte > 1")
				return nil
			}
			v.s.SetValue(value | uint64(readByte)<<shift) //todo set size
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
	v.s.SetValue(uint64(value)) //todo add int64 value field and set size
	return v.s
}

func new_varint() varint {
	return &_varint{
		s: New(),
	}
}

// PutUvarint encodes a uint64 into bytes and returns the number of bytes written.
// If the buffer is too small, PutUvarint will panic.
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

// Uvarint decodes a uint64 from bytes and returns that value and the
// number of bytes read (> 0). If an error occurred, the value is 0
// and the number of bytes n is <= 0 meaning:
//
//	n == 0: bytes too small
//	n  < 0: value larger than 64 bits (overflow)
//	        and -n is the number of bytes read
func Uvarint(buf []byte) (uint64, int) {
	var x uint64
	var s uint
	for i, b := range buf {
		if i == MaxVarintLen64 {
			// Catch byte reads past MaxVarintLen64.
			// See issue https://golang.org/issues/41185
			return 0, -(i + 1) // overflow
		}
		if b < 0x80 {
			if i == MaxVarintLen64-1 && b > 1 {
				return 0, -(i + 1) // overflow
			}
			return x | uint64(b)<<s, i + 1
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return 0, 0
}

// PutVarint encodes an int64 into bytes and returns the number of bytes written.
// If the buffer is too small, PutVarint will panic.
func PutVarint(buf []byte, x int64) int {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return PutUvarint(buf, ux)
}

// Varint decodes an int64 from bytes and returns that value and the
// number of bytes read (> 0). If an error occurred, the value is 0
// and the number of bytes n is <= 0 with the following meaning:
//
//	n == 0: bytes too small
//	n  < 0: value larger than 64 bits (overflow)
//	        and -n is the number of bytes read
func Varint(buf []byte) (int64, int) {
	ux, n := Uvarint(buf) // ok to continue in presence of error
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, n
}

var overflow = errors.New("mybinary: varintStream overflows a 64-bit integer")

// ReadUvarint reads an encoded unsigned integer from r and returns it as a uint64.
func ReadUvarint(r io.ByteReader) (uint64, error) {
	var x uint64
	var s uint
	for i := 0; i < MaxVarintLen64; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return x, err
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

// ReadVarint reads an encoded signed integer from r and returns it as an int64.
func ReadVarint(r io.ByteReader) (int64, error) {
	ux, err := ReadUvarint(r) // ok to continue in presence of error
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, err
}
