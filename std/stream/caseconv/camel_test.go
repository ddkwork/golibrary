package caseconv_test

import (
	"testing"

	"github.com/ddkwork/golibrary/std/assert"
	"github.com/ddkwork/golibrary/std/stream/caseconv"
)

func Test_ToCamel(t *testing.T) {
	testCases := []testCase{
		{"", ""},
		{"test", "test"},
		{"test string", "testString"},
		{"Test String", "testString"},
		{"TestV2", "testV2"},
		{"_foo_bar_", "fooBar"},
		{"version 1.2.10", "version1210"},
		{"version 1.21.0", "version1210"},
		{"version 1.2.10", "version1210"},
		{"PippiLÅNGSTRUMP", "pippiLångstrump"},
		{"PippilÅNGSTRUMP", "pippilÅngstrump"},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedOutput, caseconv.ToCamel(tc.input))
	}
}
