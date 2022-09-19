package flot

import (
	"fmt"
	"strconv"
)

type (
	Interface interface {
		Float64ToString(f float64, cut int) string           //Float64转String
		Float64Cut(value float64, bits int) (float64, error) //Float64截取，bits为保留几位小数
	}
	object struct {
	}
)

func New() Interface { return &object{} }

func (o *object) Float64ToString(f float64, cut int) string {
	return strconv.FormatFloat(f, 'f', cut, 64)
}

func (o *object) Float64Cut(value float64, bits int) (float64, error) {
	return strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(bits)+"f", value), 64)
}
