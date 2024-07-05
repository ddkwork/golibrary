package maps

import (
	"strings"
	"sync"

	"github.com/ddkwork/golibrary/mylog"
)

type SafeMap[K comparable, V any] struct {
	sync.RWMutex
	items StdMap[K, V]
}

func (m *SafeMap[K, V]) Clear() {
	if m.items == nil {
		return
	}
	mylog.CheckNil(m.items)
	m.Lock()
	m.items = nil
	m.Unlock()
}

func (m *SafeMap[K, V]) Set(k K, v V) {
	m.Lock()
	if m.items == nil {
		m.items = map[K]V{k: v}
	} else {
		m.items[k] = v
	}
	m.Unlock()
}

func (m *SafeMap[K, V]) Get(k K) (v V) {
	return m.Load(k)
}

func (m *SafeMap[K, V]) Has(k K) (exists bool) {
	_, ok := m.items[k]
	return ok
}

func (m *SafeMap[K, V]) Load(k K) (v V) {
	mylog.CheckNil(m.items)
	m.RLock()
	v, ok := m.items[k]
	mylog.Check(ok)
	mylog.CheckNil(v)
	m.RUnlock()
	return
}

func (m *SafeMap[K, V]) Delete(k K) {
	m.Lock()
	m.items.Delete(k)
	m.Unlock()
}

func (m *SafeMap[K, V]) Values() (v []V) {
	mylog.CheckNil(m.items)
	m.RLock()
	v = m.items.Values()
	m.RUnlock()
	return
}

func (m *SafeMap[K, V]) HasPrefix(prefix string) bool {
	for _, k := range m.Keys() {
		if strings.HasPrefix(any(k).(string), prefix) {
			return true
		}
	}
	return false
}

func (m *SafeMap[K, V]) Keys() (keys []K) {
	mylog.CheckNil(m.items)
	m.RLock()
	keys = m.items.Keys()
	m.RUnlock()
	return
}

func (m *SafeMap[K, V]) Len() (l int) {
	if m.items == nil {
		return
	}
	mylog.CheckNil(m.items)
	m.RLock()
	l = m.items.Len()
	m.RUnlock()
	return
}

func (m *SafeMap[K, V]) Range(f func(k K, v V) bool) {
	mylog.CheckNil(m)
	if m.items == nil {
		return
	}
	mylog.CheckNil(m.items)
	m.RLock()
	defer m.RUnlock()
	m.items.Range(f)
}

func (m *SafeMap[K, V]) Merge(in MapI[K, V]) {
	if m.items == nil {
		m.items = make(map[K]V, in.Len())
	}
	m.Lock()
	defer m.Unlock()
	m.items.Merge(in)
}

func (m *SafeMap[K, V]) Equal(m2 MapI[K, V]) bool {
	m.RLock()
	defer m.RUnlock()
	return m.items.Equal(m2)
}

func (m *SafeMap[K, V]) MarshalBinary() ([]byte, error) {
	m.RLock()
	defer m.RUnlock()
	return m.items.MarshalBinary()
}

func (m *SafeMap[K, V]) UnmarshalBinary(data []byte) (err error) {
	m.Lock()
	defer m.Unlock()
	return m.items.UnmarshalBinary(data)
}

func (m *SafeMap[K, V]) MarshalJSON() (out []byte, err error) {
	m.RLock()
	defer m.RUnlock()
	return m.items.MarshalJSON()
}

func (m *SafeMap[K, V]) UnmarshalJSON(in []byte) (err error) {
	m.Lock()
	defer m.Unlock()
	return m.items.UnmarshalJSON(in)
}

func (m *SafeMap[K, V]) String() string {
	m.RLock()
	defer m.RUnlock()
	return m.items.String()
}
