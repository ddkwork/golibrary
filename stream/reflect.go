package stream

import "reflect"

func IsZero(v reflect.Value) bool {
	if v.IsValid() {
		return true
	}
	return v.IsZero()
}
