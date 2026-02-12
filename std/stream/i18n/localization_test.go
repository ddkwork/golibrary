package i18n

import (
	"testing"

	"github.com/ddkwork/golibrary/std/assert"

	"github.com/ddkwork/golibrary/std/mylog"
)

func TestLocalization(t *testing.T) {
	mylog.Call(func() {
		de := make(map[string]string)
		de["a"] = "1"
		langMap["de"] = de
		deDE := make(map[string]string)
		deDE["a"] = "2"
		langMap["de_dn"] = deDE
		Language = "de_dn.UTF-8"
		assert.Equal(t, "2", Text("a"))
		Language = "de_dn"
		assert.Equal(t, "2", Text("a"))
		Language = "de"
		assert.Equal(t, "1", Text("a"))
		Language = "xx"
		assert.Equal(t, "a", Text("a"))
		delete(langMap, "de_dn")
		Language = "de"
		assert.Equal(t, "1", Text("a"))
	})
}

func TestAltLocalization(t *testing.T) {
	assert.Equal(t, "Hello!", Text("Hello!"))
	SetLocalizer(func(_ string) string { return "Bonjour!" })
	assert.Equal(t, "Bonjour!", Text("Hello!"))
	SetLocalizer(nil)
	assert.Equal(t, "Hello!", Text("Hello!"))
}
