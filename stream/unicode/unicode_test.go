package unicode_test

import (
	"testing"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream/unicode"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	u := unicode.New()
	key := "9c4077ce-81b6-4edf-8ade-a9ba7da278ba"
	assert.True(t, u.FromString(key))
	mylog.HexDump("unicode", u.Bytes())
	u.ToString(u.Bytes())
	println(u.String())
}
