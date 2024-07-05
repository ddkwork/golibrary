package mylog

import (
	"io"
	"reflect"
	"strings"
)

func reason(err any) string {
	switch err.(type) {
	case error:
		return err.(error).Error()
	case string:
		return strings.TrimSuffix(err.(string), "\n")
	}
	return ""
}

type errorX interface {
	Error() string
	Zero() bool
}

var errType = reflect.TypeFor[error]

func Check2[T any](ret T, err error) (r1 T) {
	check(newErrorMock(err))
	CheckNil(ret)
	return ret
}

func Check2Ignore[T any](ret T, err error) (r1 T) {
	defaultObject.errorCall(err)
	return ret
}

func CheckIgnore(err any) {
	defaultObject.errorCall(err)
}

func Check2Bool[T any](ret T, ok bool) (r1 T) {
	check(ok)
	CheckNil(ret)
	return ret
}

func Check3[T1 any, T2 any](ret1 T1, ret2 T2, err error) (r1 T1, r2 T2) {
	check(newErrorMock(err))
	CheckNil(ret1)
	CheckNil(ret2)
	return ret1, ret2
}

func Check3Bool[T1 any, T2 any](ret1 T1, ret2 T2, ok bool) (r1 T1, r2 T2) {
	check(ok)
	CheckNil(ret1)
	CheckNil(ret2)
	return ret1, ret2
}

func Check4[T1 any, T2 any, T3 any](ret1 T1, ret2 T2, ret3 T3, err error) (r1 T1, r2 T2, r3 T3) {
	check(newErrorMock(err))
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	return ret1, ret2, ret3
}

func Check4Bool[T1 any, T2 any, T3 any](ret1 T1, ret2 T2, ret3 T3, ok bool) (r1 T1, r2 T2, r3 T3) {
	check(ok)
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	return ret1, ret2, ret3
}

func Check5[T1 any, T2 any, T3 any, T4 any](ret1 T1, ret2 T2, ret3 T3, ret4 T4, err error) (r1 T1, r2 T2, r3 T3, r4 T4) {
	check(newErrorMock(err))
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	CheckNil(ret4)
	return ret1, ret2, ret3, ret4
}

func Check5Bool[T1 any, T2 any, T3 any, T4 any](ret1 T1, ret2 T2, ret3 T3, ret4 T4, ok bool) (r1 T1, r2 T2, r3 T3, r4 T4) {
	check(ok)
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	CheckNil(ret4)
	return ret1, ret2, ret3, ret4
}

func Check6[T1 any, T2 any, T3 any, T4 any, T5 any](ret1 T1, ret2 T2, ret3 T3, ret4 T4, ret5 T5, err error) (r1 T1, r2 T2, r3 T3, r4 T4, r5 T5) {
	check(newErrorMock(err))
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	CheckNil(ret4)
	CheckNil(ret5)
	return ret1, ret2, ret3, ret4, ret5
}

func Check6Bool[T1 any, T2 any, T3 any, T4 any, T5 any](ret1 T1, ret2 T2, ret3 T3, ret4 T4, ret5 T5, ok bool) (r1 T1, r2 T2, r3 T3, r4 T4, r5 T5) {
	check(ok)
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	CheckNil(ret4)
	CheckNil(ret5)
	return ret1, ret2, ret3, ret4, ret5
}

func Check7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](ret1 T1, ret2 T2, ret3 T3, ret4 T4, ret5 T5, ret6 T6, err error) (r1 T1, r2 T2, r3 T3, r4 T4, r5 T5, r6 T6) {
	check(newErrorMock(err))
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	CheckNil(ret4)
	CheckNil(ret5)
	CheckNil(ret6)
	return ret1, ret2, ret3, ret4, ret5, ret6
}

func Check7Bool[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](ret1 T1, ret2 T2, ret3 T3, ret4 T4, ret5 T5, ret6 T6, ok bool) (r1 T1, r2 T2, r3 T3, r4 T4, r5 T5, r6 T6) {
	check(ok)
	CheckNil(ret1)
	CheckNil(ret2)
	CheckNil(ret3)
	CheckNil(ret4)
	CheckNil(ret5)
	CheckNil(ret6)
	return ret1, ret2, ret3, ret4, ret5, ret6
}

type errorMock struct{ err error }

func newErrorMock(err error) *errorMock { return &errorMock{err: err} }
func (e *errorMock) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func Check[T any](result T) {
	switch result := any(result).(type) {
	case bool:
		check(result)
	case string:
		check(result)
	case error:
		check(newErrorMock(result))
	default:
		if isNil(result) {
			break
		}
	}
}

func CheckEof[T any](result T) (isEof bool) {
	switch result := any(result).(type) {
	case bool:
		check(result)
	case string:
		check(result)
	case error:
		return checkEof(newErrorMock(result))
	default:
		if isNil(result) {
			break
		}
	}
	return
}

func checkEof[T *errorMock | ~string | ~bool](result T) (isEof bool) {
	switch r := any(result).(type) {
	case bool:
		if !r {
			panic("return value must be true,but you got false")
		}
	case string:
		panic(r)
	case *errorMock:
		s := r.Error()
		switch s {
		case "":
		default:
			if s == io.EOF.Error() {
				return true
			}
			panic(s)
		}
	}
	return
}

func check[T *errorMock | ~string | ~bool](result T) {
	switch r := any(result).(type) {
	case bool:
		if !r {
			panic("return value must be true,but you got false")
		}
	case string:
		panic(r)
	case *errorMock:
		s := r.Error()
		switch s {
		case "":
		case "The operation completed successfully.":
		case io.EOF.Error():
			Success("detected EOF", "you should call CheckEof()")
		default:
			panic(s)
		}
	}
}

func CheckNil(ptr any) {
	if isNil(ptr) {
		panic("detected nil pointer")
	}
	switch r := ptr.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if r == 0 {
		}
		if r == -1 {
			panic("return -1")
		}
	case string:
		if r == "" {
		}
	case []byte:
		if r == nil {
			panic("detected nil []byte")
		}
		if len(r) == 0 {
		}
	}
}

func isNil(i any) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return reflect.ValueOf(i).IsNil()
	default:
		return false
	}
}
