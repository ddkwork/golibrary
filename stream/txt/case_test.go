package txt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream/txt"
)

func TestToCamelCase(t *testing.T) {
	assert.Equal(t, "SnakeCase", txt.ToCamelCase("snake__case"))
	assert.Equal(t, "CamelCase", txt.ToCamelCase("CamelCase"))
}

func TestToCamelCaseWithExceptions(t *testing.T) {
	mylog.Call(func() {
		assert.Equal(t, "ID", txt.ToCamelCaseWithExceptions("id", txt.StdAllCaps))
		assert.Equal(t, "世界ID", txt.ToCamelCaseWithExceptions("世界_id", txt.StdAllCaps))
		assert.Equal(t, "OneID", txt.ToCamelCaseWithExceptions("one_id", txt.StdAllCaps))
		assert.Equal(t, "IDOne", txt.ToCamelCaseWithExceptions("id_one", txt.StdAllCaps))
		assert.Equal(t, "OneIDTwo", txt.ToCamelCaseWithExceptions("one_id_two", txt.StdAllCaps))
		assert.Equal(t, "OneIDTwoID", txt.ToCamelCaseWithExceptions("one_id_two_id", txt.StdAllCaps))
		assert.Equal(t, "OneIDID", txt.ToCamelCaseWithExceptions("one_id_id", txt.StdAllCaps))
		assert.Equal(t, "Orchid", txt.ToCamelCaseWithExceptions("orchid", txt.StdAllCaps))
		assert.Equal(t, "OneURLTwo", txt.ToCamelCaseWithExceptions("one_url_two", txt.StdAllCaps))
		assert.Equal(t, "URLID", txt.ToCamelCaseWithExceptions("url_id", txt.StdAllCaps))
	})
}

func TestToSnakeCase(t *testing.T) {
	assert.Equal(t, "snake_case", txt.ToSnakeCase("snake_case"))
	assert.Equal(t, "camel_case", txt.ToSnakeCase("CamelCase"))
}
