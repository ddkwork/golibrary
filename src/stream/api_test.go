package stream_test

import (
	"encoding/hex"
	"github.com/ddkwork/golibrary/src/stream"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName(t *testing.T) {
	s := stream.NewHexString("1111")
	println(s.HexString())
	println(hex.Dump(s.Bytes()))

	s = stream.NewBytes([]byte{0xff, 11, 22, 33})
	println(s.HexString())
	println(hex.Dump(s.Bytes()))

	buffer := stream.New()
	b1 := []byte{1, 2, 3}
	b2 := []byte{4, 5, 6}
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6}, buffer.Merge(b1, b2).Bytes())
}
