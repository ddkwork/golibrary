package strconv

import (
	"github.com/ddkwork/golibrary/mylog"
	"strconv"
)

type (
	Interface interface {
		ParseUint(s string, base int, bitSize int) bool
		FormatUint(i uint64, base int) string
		Value() uint64
	}
	object struct {
		value uint64
		error error
	}
)

func (o *object) FormatUint(i uint64, base int) string { return strconv.FormatUint(i, base) }

func (o *object) ParseUint(s string, base int, bitSize int) bool {
	o.value, o.error = strconv.ParseUint(s, base, bitSize)
	return mylog.Error(o.error)
}

func (o *object) Value() uint64 {
	return o.value
}

var Default = New()

func New() Interface { return &object{} }
