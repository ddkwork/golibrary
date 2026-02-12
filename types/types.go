package types

import (
	"cmp"
)

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
	Signed | Unsigned
}

type Float interface {
	~float32 | ~float64
}

type Complex interface {
	~complex64 | ~complex128
}

type Ordered = cmp.Ordered

type Number interface {
	Integer | Float | Complex
}

// todo D:\workspace\workspace\ux\demo\airtable\sdk\cell.value.detect.go
func ParseNumber[T Number | ~string]() {}

func ParseBool[T Integer | ~string]() {
}

func Integer2Bool[T Integer](value T) bool { // todo多维表格那波转换移动到这里，支持字符串解析到bool
	switch v := any(value).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		if v == 1 {
			return true
		}
	}
	return false
}

func Bool2Integer[T Integer](b bool) T {
	var zero T
	switch any(zero).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		if b {
			zero = 1
		}
	}
	return zero
}
