package caseconv_test

import (
	"testing"

	"github.com/ddkwork/golibrary/std/assert"
	"github.com/ddkwork/golibrary/std/stream/caseconv"
)

func Test_ToPascal(t *testing.T) {
	testCases := []testCase{
		{"", ""},
		{"test", "Test"},
		{"test string", "TestString"},
		{"Test String", "TestString"},
		{"TestV2", "TestV2"},
		{"version 1.2.10", "Version1210"},
		{"version 1.21.0", "Version1210"},
		{"LÅNGSTRUMP", "Långstrump"},
		{"PippiLÅNGSTRUMP", "PippiLångstrump"},
		{"PippilÅNGSTRUMP", "PippilÅngstrump"},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedOutput, caseconv.ToPascal(tc.input))
	}
}

type testCase struct {
	input          string
	expectedOutput string
}
