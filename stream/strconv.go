package stream

import (
	"fmt"
	"strconv"

	"github.com/ddkwork/golibrary/mylog"
)

func ParseFloat(s string) float64 {
	float, err := strconv.ParseFloat(s, 64)
	if !mylog.Error(err) {
		return 0
	}
	return float
}

func Float64ToString(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

func Float64Cut(value float64, bits int) (float64, error) {
	return strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(bits)+"f", value), 64)
}

func ParseInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if !mylog.Error(err) {
		return 0
	}
	return i
}

func ParseUint(s string) uint64 {
	i, err := strconv.ParseUint(s, 10, 64)
	if !mylog.Error(err) {
		return 0
	}
	return i
}

func Atoi(s string) int {
	atoi, err := strconv.Atoi(s)
	if !mylog.Error(err) {
		return 0
	}
	return atoi
}
