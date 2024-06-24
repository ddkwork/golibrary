package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/ddkwork/golibrary/mylog"
)

type SafeSliceMap[K comparable, V any] struct {
	sync.RWMutex
	items StdMap[K, V]
	order []K
	lessF func(key1, key2 K, val1, val2 V) bool
}

func (m *SafeSliceMap[K, V]) SetSortFunc(f func(key1, key2 K, val1, val2 V) bool) {
	m.Lock()
	defer m.Unlock()

	m.lessF = f
	if f != nil && len(m.order) > 0 {
		sort.Slice(m.order, func(i, j int) bool {
			return f(m.order[i], m.order[j], m.items[m.order[i]], m.items[m.order[j]])
		})
	}
}

func (m *SafeSliceMap[K, V]) Set(key K, val V) {
	var ok bool
	var oldVal V

	m.Lock()

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
	m.Unlock()
}

func (m *SafeSliceMap[K, V]) SetAt(index int, key K, val V) {
	if m.lessF != nil {
		panic("cannot use SetAt if you are also using a sort function")
	}

	if index >= len(m.order) {
		m.Set(key, val)
		return
	}

	var emptyKey K

	if m.Has(key) {
		m.Delete(key)
	}
	m.Lock()
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
	m.Unlock()
}

func (m *SafeSliceMap[K, V]) Delete(key K) {
	m.Lock()
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
	m.Unlock()
}

func (m *SafeSliceMap[K, V]) Get(key K) (val V) {
	m.RLock()
	defer m.RUnlock()
	return m.items.Get(key)
}

func (m *SafeSliceMap[K, V]) Load(key K) (val V) {
	m.RLock()
	defer m.RUnlock()
	return m.items.Load(key)
}

func (m *SafeSliceMap[K, V]) Has(key K) (ok bool) {
	m.RLock()
	defer m.RUnlock()
	return m.items.Has(key)
}

func (m *SafeSliceMap[K, V]) GetAt(position int) (val V) {
	m.RLock()
	defer m.RUnlock()
	if position < len(m.order) && position >= 0 {
		val, _ = m.items[m.order[position]]
	}
	return
}

func (m *SafeSliceMap[K, V]) GetKeyAt(position int) (key K) {
	m.RLock()
	defer m.RUnlock()
	if position < len(m.order) && position >= 0 {
		key = m.order[position]
	}
	return
}

func (m *SafeSliceMap[K, V]) Values() (vals []V) {
	m.RLock()
	defer m.RUnlock()
	return m.items.Values()
}

func (m *SafeSliceMap[K, V]) Keys() (keys []K) {
	m.RLock()
	defer m.RUnlock()
	return m.items.Keys()
}

func (m *SafeSliceMap[K, V]) Len() int {
	m.RLock()
	defer m.RUnlock()
	return m.items.Len()
}

func (m *SafeSliceMap[K, V]) MarshalBinary() (data []byte, err error) {
	m.RLock()
	defer m.RUnlock()

	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	mylog.Check(encoder.Encode(map[K]V(m.items)))
	if err == nil {
		mylog.Check(encoder.Encode(m.order))
	}
	data = buf.Bytes()
	return
}

func (m *SafeSliceMap[K, V]) UnmarshalBinary(data []byte) (err error) {
	var items map[K]V
	var order []K

	m.Lock()
	defer m.Unlock()

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if mylog.Check(dec.Decode(&items)); err == nil {
		mylog.Check(dec.Decode(&order))
	}

	if err == nil {
		m.items = items
		m.order = order
	}
	return err
}

func (m *SafeSliceMap[K, V]) MarshalJSON() (data []byte, err error) {
	m.RLock()
	defer m.RUnlock()

	return m.items.MarshalJSON()
}

func (m *SafeSliceMap[K, V]) UnmarshalJSON(data []byte) (err error) {
	var items map[K]V

	m.Lock()
	defer m.Unlock()

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

func (m *SafeSliceMap[K, V]) Merge(in MapI[K, V]) {
	in.Range(func(k K, v V) bool {
		m.Set(k, v)
		return true
	})
}

func (m *SafeSliceMap[K, V]) Range(f func(key K, value V) bool) {
	mylog.CheckNil(m)
	if m.items == nil {
		return
	}
	mylog.CheckNil(m.items)
	m.RLock()
	defer m.RUnlock()
	for _, k := range m.order {
		if !f(k, m.items[k]) {
			break
		}
	}
}

func (m *SafeSliceMap[K, V]) Equal(m2 MapI[K, V]) bool {
	m.RLock()
	defer m.RUnlock()
	return m.items.Equal(m2)
}

func (m *SafeSliceMap[K, V]) Clear() {
	m.Lock()
	m.items = nil
	m.order = nil
	m.Unlock()
}

func (m *SafeSliceMap[K, V]) String() string {
	var s string

	s = "{"

	m.Range(func(k K, v V) bool {
		s += fmt.Sprintf(`%#v:%#v,`, k, v)
		return true
	})
	s = strings.TrimRight(s, ",")
	s += "}"
	return s
}
