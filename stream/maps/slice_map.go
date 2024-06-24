package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
)

type SliceMap[K comparable, V any] struct {
	items StdMap[K, V]
	order []K
	lessF func(key1, key2 K, val1, val2 V) bool
}

func (m *SliceMap[K, V]) SetSortFunc(f func(key1, key2 K, val1, val2 V) bool) {
	if m == nil {
		panic("cannot set a sort function on a nil SliceMap")
	}
	m.lessF = f
	if f != nil && len(m.order) > 0 {
		sort.Slice(m.order, func(i, j int) bool {
			return f(m.order[i], m.order[j], m.items[m.order[i]], m.items[m.order[j]])
		})
	}
}

func (m *SliceMap[K, V]) Set(key K, val V) {
	var ok bool
	var oldVal V

	if m == nil {
		panic("cannot set a value on a nil SliceMap")
	}

	if m.items == nil {
		m.items = make(map[K]V)
	}

	_, ok = m.items[key]
	if m.lessF != nil {
		if ok {

			loc := sort.Search(len(m.items), func(n int) bool {
				return !m.lessF(m.order[n], key, m.items[m.order[n]], oldVal)
			})
			m.order = append(m.order[:loc], m.order[loc+1:]...)
		}

		loc := sort.Search(len(m.order), func(n int) bool {
			return m.lessF(key, m.order[n], val, m.items[m.order[n]])
		})

		m.order = append(m.order, key)
		copy(m.order[loc+1:], m.order[loc:])
		m.order[loc] = key
	} else {
		if !ok {
			m.order = append(m.order, key)
		}
	}
	m.items[key] = val
}

func (m *SliceMap[K, V]) SetAt(index int, key K, val V) {
	if m == nil {
		panic("cannot set a value on a nil SliceMap")
	}

	if m.lessF != nil {
		panic("cannot use SetAt if you are also using a sort function")
	}

	if index >= len(m.order) {

		m.Set(key, val)
		return
	}

	var ok bool
	var emptyKey K

	if _, ok = m.items[key]; ok {
		m.Delete(key)
	}
	if index <= -len(m.items) {
		index = 0
	}
	if index < 0 {
		index = len(m.items) + index
	}

	m.order = append(m.order, emptyKey)
	copy(m.order[index+1:], m.order[index:])
	m.order[index] = key

	m.items[key] = val
}

func (m *SliceMap[K, V]) Delete(key K) {
	if m == nil {
		return
	}

	if _, ok := m.items[key]; ok {
		if m.lessF != nil {
			oldVal := m.items[key]
			loc := sort.Search(len(m.items), func(n int) bool {
				return !m.lessF(m.order[n], key, m.items[m.order[n]], oldVal)
			})
			m.order = append(m.order[:loc], m.order[loc+1:]...)
		} else {
			for i, v := range m.order {
				if v == key {
					m.order = append(m.order[:i], m.order[i+1:]...)
					break
				}
			}
		}
		delete(m.items, key)
	}
}

func (m *SliceMap[K, V]) Get(key K) (val V) {
	if m == nil {
		return
	}
	return m.items.Get(key)
}

func (m *SliceMap[K, V]) Load(key K) (val V) {
	if m == nil {
		return
	}
	return m.items.Load(key)
}

func (m *SliceMap[K, V]) Has(key K) (ok bool) {
	if m == nil {
		return
	}
	return m.items.Has(key)
}

func (m *SliceMap[K, V]) GetAt(position int) (val V) {
	if m == nil {
		return
	}
	if position < len(m.order) && position >= 0 {
		val, _ = m.items[m.order[position]]
	}
	return
}

func (m *SliceMap[K, V]) GetKeyAt(position int) (key K) {
	if m == nil {
		return
	}
	if position < len(m.order) && position >= 0 {
		key = m.order[position]
	}
	return
}

func (m *SliceMap[K, V]) Values() (vals []V) {
	if m == nil {
		return
	}
	return m.items.Values()
}

func (m *SliceMap[K, V]) Keys() (keys []K) {
	if m == nil {
		return
	}
	return m.items.Keys()
}

func (m *SliceMap[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return m.items.Len()
}

func (m *SliceMap[K, V]) MarshalBinary() (data []byte, err error) {
	if m == nil {
		return
	}
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	mylog.Check(encoder.Encode(map[K]V(m.items)))
	if err == nil {
		mylog.Check(encoder.Encode(m.order))
	}
	data = buf.Bytes()
	return
}

func (m *SliceMap[K, V]) UnmarshalBinary(data []byte) (err error) {
	var items map[K]V
	var order []K

	if m == nil {
		panic("cannot Unmarshal into a nil SliceMap")
	}

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	mylog.Check(dec.Decode(&items))
	mylog.Check(dec.Decode(&order))

	if err == nil {
		m.items = items
		m.order = order
	}
	return err
}

func (m *SliceMap[K, V]) MarshalJSON() (data []byte, err error) {
	if m == nil {
		return
	}
	return m.items.MarshalJSON()
}

func (m *SliceMap[K, V]) UnmarshalJSON(data []byte) (err error) {
	var items map[K]V

	if m == nil {
		panic("cannot unmarshall into a nil SliceMap")
	}
	if mylog.Check(json.Unmarshal(data, &items)); err == nil {
		m.items = items

		m.order = make([]K, len(m.items))
		i := 0
		for k := range m.items {
			m.order[i] = k
			i++
		}
	}
	return
}

func (m *SliceMap[K, V]) Merge(in MapI[K, V]) {
	in.Range(func(k K, v V) bool {
		m.Set(k, v)
		return true
	})
}

func (m *SliceMap[K, V]) Range(f func(key K, value V) bool) {
	if m != nil && m.items != nil {
		for _, k := range m.order {
			if !f(k, m.items[k]) {
				break
			}
		}
	}
}

func (m *SliceMap[K, V]) Equal(m2 MapI[K, V]) bool {
	if m == nil {
		return m2 == nil || m2.Len() == 0
	}
	return m.items.Equal(m2)
}

func (m *SliceMap[K, V]) Clear() {
	if m == nil {
		return
	}
	m.items = nil
	m.order = nil
}

func (m *SliceMap[K, V]) String() string {
	var s string

	if m == nil {
		return s
	}

	s = "{"
	m.Range(func(k K, v V) bool {
		s += fmt.Sprintf(`%#v:%#v,`, k, v)
		return true
	})
	s = strings.TrimRight(s, ",")
	s += "}"
	return s
}

type Equaler interface {
	Equal(a any) bool
}

func equalValues(a, b any) bool {
	if e, ok := a.(Equaler); ok {
		return e.Equal(b)
	}

	return a == b
}
