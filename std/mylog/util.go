package mylog

import (
	"fmt"
)

func sprint(msg ...any) string {
	if len(msg) == 0 {
		return ""
	}
	data := fmt.Sprint(msg...)
	if data == "" {
		return `""`
	}
	switch {
	case data[0] == '[' && data[len(data)-1] == ']':
		data = data[1 : len(data)-1]
	}
	return data
}
