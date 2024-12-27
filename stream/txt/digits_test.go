package txt_test

import (
	"testing"

	"github.com/ddkwork/golibrary/assert"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream/txt"
)

func TestDigitToValue(t *testing.T) {
	mylog.Call(func() {
		checkDigitToValue('5', 5, t)
		checkDigitToValue('٥', 5, t)
		checkDigitToValue('𑁯', 9, t)
		mylog.Check2(txt.DigitToValue('a'))
	})
}

func checkDigitToValue(ch rune, expected int, t *testing.T) {
	value := mylog.Check2(txt.DigitToValue(ch))
	assert.Equal(t, expected, value)
}
