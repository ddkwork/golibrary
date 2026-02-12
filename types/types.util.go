package types

import (
	"reflect"
	"strings"
)

func ValueIsBytesType(v reflect.Value) bool {
	// i,ok := v.Interface().([]byte)
	// assert, b := reflect.TypeAssert[[]byte](v)
	return v.Type().Elem().Kind() == reflect.Uint8
}

// func isASCIILower(c byte) bool { return 'a' <= c && c <= 'z' }
// func isASCIIUpper(c byte) bool { return 'A' <= c && c <= 'Z' }
func isASCIIDigit(c byte) bool { return '0' <= c && c <= '9' }

// func IsASCIIAlpha(s string) bool {
//	for i := 0; i < len(s); i++ {
//		c := s[i] // 直接按字节获取
//		if !isASCIILower(c) && !isASCIIUpper(c) {
//			return false
//		}
//	}
//	return true
// }

func IsASCIIDigit(s string) bool {
	for i := range len(s) {
		if !isASCIIDigit(s[i]) {
			return false
		}
	}
	return len(s) > 0 // 确保字符串非空
}

// func IsAlphanumeric(s string) bool {
//	for i := 0; i < len(s); i++ {
//		c := s[i]
//		if !isASCIIDigit(c) && !isASCIILower(c) && !isASCIIUpper(c) {
//			return false
//		}
//	}
//	return len(s) > 0 // 确保字符串非空
// }

// func isOneByteInteger(n int) bool {
//	return n >= -128 && n <= 127 // 检查有符号整数
//	// return n >= 0 && n <= 255 // 可以用于无符号整数
// }

func IsIncludeLine(s string) bool {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "#") {
		return false
	}
	s = strings.TrimPrefix(s, "#")
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "include")
}
