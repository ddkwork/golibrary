package txt

import (
	"strings"
)

var (
	romanValues = []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	romanText   = []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
)

func RomanNumerals(value int) string {
	var buffer strings.Builder
	for value > 0 {
		for i, v := range romanValues {
			if value >= v {
				buffer.WriteString(romanText[i])
				value -= v
				break
			}
		}
	}
	return buffer.String()
}
