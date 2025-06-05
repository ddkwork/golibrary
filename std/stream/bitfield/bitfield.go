package bitfield

type BitField []byte

func New(n int) BitField {
	n = 1 + ((n - 1) / 8)
	return make([]byte, n)
}

func NewFromUint32(n uint32) BitField {
	b := BitField(make([]byte, 4))
	b[0] = byte(n)
	b[1] = byte(n >> 8)
	b[2] = byte(n >> 16)
	b[3] = byte(n >> 24)
	return b
}

func NewFromUint64(n uint64) BitField {
	b := BitField(make([]byte, 8))
	b[0] = byte(n)
	b[1] = byte(n >> 8)
	b[2] = byte(n >> 16)
	b[3] = byte(n >> 24)
	b[4] = byte(n >> 32)
	b[5] = byte(n >> 40)
	b[6] = byte(n >> 48)
	b[7] = byte(n >> 56)
	return b
}

func (b BitField) Size() int     { return len(b) }
func (b BitField) Bytes() []byte { return b }

func (b BitField) SetX(index int, v int) BitField {
	field := New(v)
	for i := range b {
		if i+2 == index {
			for _, b2 := range field {
				b[i+2] = b2
			}
			return b
		}
	}
	return nil
}

func (b BitField) GetX(index int) BitField {
	for i := range b {
		if i+1 == len(b)-index {
			return b[i+1 : i+1+index]
		}
	}
	return nil
}

func (b BitField) Set(i uint32) {
	idx, offset := i/8, i%8
	b[idx] |= 1 << uint(offset)
}

func (b BitField) Clear(i uint32) {
	idx, offset := i/8, i%8
	b[idx] &= ^(1 << uint(offset))
}

func (b BitField) Flip(i uint32) {
	idx, offset := i/8, i%8
	b[idx] ^= 1 << uint(offset)
}

func (b BitField) Test(i uint32) bool {
	idx, offset := i/8, i%8
	return (b[idx] & (1 << uint(offset))) != 0
}

func (b BitField) ClearAll() {
	for idx := range b {
		b[idx] = 0
	}
}

func (b BitField) SetAll() {
	for idx := range b {
		b[idx] = 0xff
	}
}

func (b BitField) FlipAll() {
	for idx := range b {
		b[idx] = ^b[idx]
	}
}

func (b BitField) ANDMask(m BitField) {
	maxidx := len(m)
	for idx := range b {

		if idx > maxidx {
			b[idx] = 0
			continue
		}
		b[idx] &= m[idx]
	}
}

func (b BitField) ORMask(m BitField) {
	maxidx := len(m)
	for idx := range b {
		if idx > maxidx {
			break
		}
		b[idx] |= m[idx]
	}
}

func (b BitField) XORMask(m BitField) {
	maxidx := len(m)
	for idx := range b {
		if idx > maxidx {
			break
		}
		b[idx] ^= m[idx]
	}
}

func (b BitField) ToUint32() uint32 {
	var r uint32
	r |= uint32(b[0])
	r |= uint32(b[1]) << 8
	r |= uint32(b[2]) << 16
	r |= uint32(b[3]) << 24
	return r
}

func (b BitField) ToUint32Safe() uint32 {
	var r uint32
	for idx := range b {
		r |= uint32(b[idx]) << uint32(idx*8)
		if idx == 3 {
			break
		}
	}
	return r
}

func (b BitField) ToUint64() uint64 {
	var r uint64
	r |= uint64(b[0])
	r |= uint64(b[1]) << 8
	r |= uint64(b[2]) << 16
	r |= uint64(b[3]) << 24
	r |= uint64(b[4]) << 32
	r |= uint64(b[5]) << 40
	r |= uint64(b[6]) << 48
	r |= uint64(b[7]) << 56
	return r
}

func (b BitField) ToUint64Safe() uint64 {
	var r uint64
	for idx := range b {
		r |= uint64(b[idx]) << uint64(idx*8)
		if idx == 7 {
			break
		}
	}
	return r
}
