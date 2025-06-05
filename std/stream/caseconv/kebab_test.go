package caseconv_test

import (
	"testing"

	"github.com/ddkwork/golibrary/std/assert"
	"github.com/ddkwork/golibrary/std/stream/caseconv"
)

func Test_ToKebab(t *testing.T) {
	testCases := []testCase{
		{"", ""},
		{"test", "test"},
		{"test string", "test-string"},
		{"Test String", "test-string"},
		{"TestV2", "test-v2"},
		{"PippiLÅNGSTRUMP", "pippi-långstrump"},
		{"PippilÅNGSTRUMP", "pippil-ångstrump"},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedOutput, caseconv.ToKebab(tc.input))
	}
}
