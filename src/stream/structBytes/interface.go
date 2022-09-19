package structBytes

import (
	"bytes"

	"github.com/ddkwork/golibrary/src/stream/structBytes/goBinary"
)

type (
	Interface interface {
		StructBytes() []byte
		StructToBytes(obj any) bool
		BytesToStruct(StructBytes []byte, obj any) bool
		goBinary.Interface
	}
	object struct {
		*bytes.Buffer
		goBinary goBinary.Interface
	}
)

func New() Interface {
	return &object{
		Buffer:   nil,
		goBinary: goBinary.New(),
	}
}

func (o *object) StructBytes() []byte        { return o.Bytes() }
func (o *object) StructToBytes(obj any) bool { return o.Write(obj) }
func (o *object) BytesToStruct(StructBytes []byte, obj any) bool {
	return o.Read(StructBytes, obj)
}
func (o *object) Encode(obj any) bool             { return o.goBinary.Encode(obj) }
func (o *object) Decode(buf []byte, obj any) bool { return o.goBinary.Decode(buf, obj) }
