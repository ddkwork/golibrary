package structBytes

import (
	"bytes"
	"encoding/binary"
	"github.com/ddkwork/golibrary/mylog"
	"reflect"
)

func (o *object) Read(StructBytes []byte, obj any) bool {
	o.Buffer = bytes.NewBuffer(StructBytes)
	return o.ReadValue(reflect.ValueOf(obj), 0)
}

func (o *object) ReadValue(v reflect.Value, depth int) (ok bool) {
	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			v.Set(reflect.ValueOf(v.Type()))
		}
		fallthrough
	case reflect.Ptr:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if !o.ReadValue(v.Elem(), depth+1) {
			return
		}

	case reflect.Struct:
		l := v.NumField()
		for i := 0; i < l; i++ {
			if !o.ReadValue(v.Field(i), depth+1) {
				return
			}
		}

	case reflect.Slice:
		if v.IsNil() {
			return mylog.Error("切片必须在解码之前初始化为正确的长度")
		}
		fallthrough
	case reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			if !o.ReadValue(v.Index(i), depth+1) {
				return
			}
		}

	case reflect.Int:
		var i int32
		if !mylog.Error(binary.Read(o.Buffer, binary.BigEndian, &i)) {
			return
		}
		v.SetInt(int64(i))

	case reflect.Bool:
		boolValue := uint8(0)
		if !mylog.Error(binary.Read(o.Buffer, binary.BigEndian, &boolValue)) {
			return
		}
		v.SetBool(boolValue != 0)

	default:
		return mylog.Error(binary.Read(o.Buffer, binary.BigEndian, v.Addr().Interface()))
	}
	return true
}
