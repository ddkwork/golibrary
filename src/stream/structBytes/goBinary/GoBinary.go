package goBinary

import (
	"bytes"
	"encoding/gob"
	"github.com/ddkwork/golibrary/mylog"
)

type (
	Interface interface {
		Encode(obj any) bool
		Decode(buf []byte, obj any) bool
	}
	object struct {
		bytes.Buffer
		err error
	}
)

func New() Interface {
	return &object{}
}

func (o *object) Encode(obj any) bool {
	enc := gob.NewEncoder(&o.Buffer)
	return mylog.Error(enc.Encode(obj))
}

func (o *object) Decode(buf []byte, obj any) bool {
	if !mylog.Error2(o.Write(buf)) {
		return false
	}
	dec := gob.NewDecoder(&o.Buffer)
	return mylog.Error(dec.Decode(obj))
}
