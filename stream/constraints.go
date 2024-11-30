package stream

import (
	"fmt"
	"reflect"
	"strings"
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

type Ordered interface {
	Integer | Float | ~string
}

func FormatInteger[T Integer](data T) string {
	return FormatIntegerHex0x(data) + " (" + fmt.Sprintf("%d", reflect.ValueOf(data).Interface()) + ")"
	return FormatIntegerHex0x(data) + " # " + fmt.Sprintf("%d", reflect.ValueOf(data).Interface())
}

func FormatIntegerHex0x[T Integer](data T) string {
	return "0x" + FormatIntegerHex(data)
}

func FormatIntegerHex[T Integer](data T) string {
	format := ""
	switch any(data).(type) {
	case int, int64, uint, uint64, uintptr:
		format = "%016X"
	case int8, uint8:
		format = "%02X"
	case int16, uint16:
		format = "%04X"
	case int32, uint32:
		format = "%08X"
	default:
		panic("unsupported type ---> " + reflect.TypeOf(data).Name())
	}
	return fmt.Sprintf(format, data)
}

func Integer2Bool[T Integer](value T) bool {
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

func IsIncludeLine(s string) bool {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "#") {
		return false
	}
	s = strings.TrimPrefix(s, "#")
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "include")
}

func ValueIsBytesType(v reflect.Value) bool {
	return v.Type().Elem().Kind() == reflect.Uint8
}

func isASCIILower(c byte) bool { return 'a' <= c && c <= 'z' }
func isASCIIUpper(c byte) bool { return 'A' <= c && c <= 'Z' }
func isASCIIDigit(c byte) bool { return '0' <= c && c <= '9' }

func IsASCIIAlpha(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i] // 直接按字节获取
		if !isASCIILower(c) && !isASCIIUpper(c) {
			return false
		}
	}
	return true
}

func IsASCIIDigit(s string) bool {
	for i := 0; i < len(s); i++ {
		if !isASCIIDigit(s[i]) {
			return false
		}
	}
	return len(s) > 0 // 确保字符串非空
}

func IsAlphanumeric(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !isASCIIDigit(c) && !isASCIILower(c) && !isASCIIUpper(c) {
			return false
		}
	}
	return len(s) > 0 // 确保字符串非空
}

func isOneByteInteger(n int) bool {
	return n >= -128 && n <= 127 // 检查有符号整数
	// return n >= 0 && n <= 255 // 可以用于无符号整数
}
