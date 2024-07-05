package maps

type Map[K comparable, V any] struct {
	items StdMap[K, V]
}

func (m *Map[K, V]) Clear() {
	m.items = nil
}

func (m Map[K, V]) Len() int {
	return m.items.Len()
}

func (m Map[K, V]) Range(f func(k K, v V) bool) {
	m.items.Range(f)
}

func (m Map[K, V]) Load(k K) V {
	return m.items.Load(k)
}

func (m Map[K, V]) Get(k K) V {
	return m.items.Get(k)
}

func (m Map[K, V]) Has(k K) bool {
	return m.items.Has(k)
}

func (m Map[K, V]) Delete(k K) {
	m.items.Delete(k)
}

func (m Map[K, V]) Keys() []K {
	return m.items.Keys()
}

func (m Map[K, V]) Values() []V {
	return m.items.Values()
}

func (m *Map[K, V]) Set(k K, v V) {
	if m.items == nil {
		m.items = map[K]V{k: v}
	} else {
		m.items.Set(k, v)
	}
}

func (m *Map[K, V]) Merge(in MapI[K, V]) {
	if m.items == nil {
		m.items = make(map[K]V, in.Len())
	}
	m.items.Merge(in)
}

func (m Map[K, V]) Equal(m2 MapI[K, V]) bool {
	return m.items.Equal(m2)
}

func (m Map[K, V]) MarshalBinary() ([]byte, error) {
	return m.items.MarshalBinary()
}

func (m *Map[K, V]) UnmarshalBinary(data []byte) (err error) {
	return m.items.UnmarshalBinary(data)
}

func (m Map[K, V]) MarshalJSON() (out []byte, err error) {
	return m.items.MarshalJSON()
}

func (m *Map[K, V]) UnmarshalJSON(in []byte) (err error) {
	return m.items.UnmarshalJSON(in)
}

func (m Map[K, V]) String() string {
	return m.items.String()
}
