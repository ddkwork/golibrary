package txt

import (
	"fmt"
	"unicode"

	"github.com/ddkwork/golibrary/mylog"
)

func DigitToValue(ch rune) (int, error) {
	if ch < '\U00010000' {
		r16 := uint16(ch)
		for _, one := range unicode.Digit.R16 {
			if one.Lo <= r16 && one.Hi >= r16 {
				return int(r16 - one.Lo), nil
			}
		}
	} else {
		r32 := uint32(ch)
		for _, one := range unicode.Digit.R32 {
			if one.Lo <= r32 && one.Hi >= r32 {
				return int(r32 - one.Lo), nil
			}
		}
	}
	mylog.Check(fmt.Sprintf("Not a digit: '%v'", ch))
	return 0, nil
}
