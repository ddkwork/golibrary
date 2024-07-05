package maps

import (
	"encoding/gob"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeSliceMap_Mapi(t *testing.T) {
	runMapiTests[SafeSliceMap[string, int]](t, makeMapi[SafeSliceMap[string, int]])
}

func init() {
	gob.Register(new(SafeSliceMap[string, int]))
}

func TestSafeSliceMap_MapiWithSortFunction(t *testing.T) {
	runMapiTests[SafeSliceMap[string, int]](t,
		func(sources ...mapT) MapI[string, int] {
			m := new(SafeSliceMap[string, int])
			for _, s := range sources {
				m.Merge(s)
			}
			m.SetSortFunc(func(k1, k2 string, v1, v2 int) bool {
				return k1 < k2
			})
			return m
		},
	)
}

func ExampleSafeSliceMap_SetSortFunc() {
	m := new(SafeSliceMap[string, int])

	m.Set("b", 2)
	m.Set("a", 1)

	fmt.Println(m)

	m.SetSortFunc(func(k1, k2 string, v1, v2 int) bool {
		return k1 < k2
	})
	fmt.Println(m)
}

func ExampleSafeSliceMap_SetAt() {
	m := new(SafeSliceMap[string, int])
	m.Set("b", 2)
	m.Set("a", 1)
	m.SetAt(1, "c", 3)
	fmt.Println(m)
}

func ExampleSafeSliceMap_GetAt() {
	m := new(SafeSliceMap[string, int])
	m.Set("b", 2)
	m.Set("c", 3)
	m.Set("a", 1)
	v := m.GetAt(1)
	fmt.Print(v)
}

func ExampleSafeSliceMap_GetKeyAt() {
	m := new(SafeSliceMap[string, int])
	m.Set("b", 2)
	m.Set("c", 3)
	m.Set("a", 1)
	v := m.GetKeyAt(1)
	fmt.Print(v)
}

func TestSafeSliceMap_SetAt(t *testing.T) {
	m := new(SafeSliceMap[string, int])
	m.Set("b", 2)
	m.Set("a", 1)
	m.SetAt(5, "c", 3)
	assert.Equal(t, 3, m.GetAt(2))

	m.SetAt(-1, "d", 4)
	assert.Equal(t, 4, m.GetAt(2))

	m.SetAt(-7, "e", 5)
	assert.Equal(t, 5, m.GetAt(0))

	m.Set("e", 6)
	assert.Equal(t, 6, m.GetAt(0))

	m.Delete("e")
	m.Set("e", 6)
	assert.Equal(t, 6, m.GetAt(4))

	m.SetAt(3, "e", 6)
	assert.Equal(t, 6, m.GetAt(3))

	m.SetSortFunc(func(k1, k2 string, v1, v2 int) bool {
		return k1 < k2
	})
	assert.Panics(t, func() {
		m.SetAt(3, "f", 4)
	})
}

func TestSafeSliceMap_GetAt(t *testing.T) {
	m := new(SafeSliceMap[string, int])
	assert.Equal(t, 0, m.GetAt(0))
	assert.Equal(t, "", m.GetKeyAt(0))
}
