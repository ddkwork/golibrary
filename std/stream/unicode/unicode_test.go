package unicode_test

import (
	"testing"

	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream/unicode"
)

func TestName(t *testing.T) {
	mylog.Call(func() {
		u := unicode.New()
		key := "9c4077ce-81b6-4edf-8ade-a9ba7da278ba"
		u.FromString(key)
		mylog.HexDump("unicode", u.Bytes())
		u.ToString(u.Bytes())
		println(u.String())
	})
}
