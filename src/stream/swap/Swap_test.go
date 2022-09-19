package swap_test

import (
	"github.com/ddkwork/golibrary/src/stream"
	"github.com/ddkwork/golibrary/src/stream/swap"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName(t *testing.T) {
	s := swap.New()
	assert.Equal(t, byte(0x16), s.CutUint16(0x6613))
	assert.Equal(t, "ET5AA5Q3N2KTR8      ", s.SerialNumber("TEA55A3Q2NTK8R      "))
	assert.Equal(t, "TA1591503892      ", s.SerialNumber("AT5119058329      "))

	b := s.HexString("1122334455667788")
	assert.Equal(t, "8877665544332211", stream.NewBytes(b).HexString())
}
