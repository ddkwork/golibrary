package stream

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"

	"github.com/ddkwork/golibrary/std/assert"
)

var (
	seed          = []byte{1, 3, 2, 0}
	expectedTable = [][]byte{
		{1, 3, 2, 0},
		{1, 3, 0, 2},
		{1, 2, 3, 0},
		{1, 2, 0, 3},
		{1, 0, 3, 2},
		{1, 0, 2, 3},
		{3, 1, 2, 0},
		{3, 1, 0, 2},
		{3, 2, 1, 0},
		{3, 2, 0, 1},
		{3, 0, 1, 2},
		{3, 0, 2, 1},
		{2, 1, 3, 0},
		{2, 1, 0, 3},
		{2, 3, 1, 0},
		{2, 3, 0, 1},
		{2, 0, 1, 3},
		{2, 0, 3, 1},
		{0, 1, 3, 2},
		{0, 1, 2, 3},
		{0, 3, 1, 2},
		{0, 3, 2, 1},
		{0, 2, 1, 3},
		{0, 2, 3, 1},
	}
)

type table struct {
	data []uint32
}

func (t *table) FromByteSlice(slice [][]byte) {
	var uint32s []uint32
	for _, bytes := range slice { // permutations
		// 只处理长度为4的字节数组
		if len(bytes) == 4 {
			uint32s = append(uint32s, binary.LittleEndian.Uint32(bytes))
		}
	}
	if reflect.DeepEqual(uint32s, t.data) {
		t.data = uint32s
		return
	}
	panic("invalid data")
}

func newTable() *table {
	return &table{
		data: []uint32{
			0x00020301,
			0x02000301,
			0x00030201,
			0x03000201,
			0x02030001,
			0x03020001,
			0x00020103,
			0x02000103,
			0x00010203,
			0x01000203,
			0x02010003,
			0x01020003,
			0x00030102,
			0x03000102,
			0x00010302,
			0x01000302,
			0x03010002,
			0x01030002,
			0x02030100,
			0x03020100,
			0x02010300,
			0x01020300,
			0x03010200,
			0x01030200,
		},
	}
}

func (t *table) GoString() string {
	g := NewGeneratedFile()
	g.P("var tableBuf = []byte{")
	for _, u := range t.data {
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, u)
		g.P(fmt.Sprintf("%d,%d,%d,%d,", b[0], b[1], b[2], b[3]))
	}
	g.P("}")
	return g.String()
}

func Test_permuteBacktrack(t *testing.T) {
	assert.Equal(t, expectedTable, Permute(seed))
	newTable().FromByteSlice(Permute(seed))
}
