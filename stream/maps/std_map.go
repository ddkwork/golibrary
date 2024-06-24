package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
)

type StdMap[K comparable, V any] map[K]V

func NewStdMap[K comparable, V any](sources ...map[K]V) StdMap[K, V] {
	m := StdMap[K, V]{}
	for _, i := range sources {
		m.Merge(Cast(i))
	}
	return m
}

func Cast[M ~map[K]V, K comparable, V any](m M) StdMap[K, V] {
	return StdMap[K, V](m)
}

func (m StdMap[K, V]) Clear() {
	for k := range m {
		delete(m, k)
	}
}

func (m StdMap[K, V]) Len() int {
	return len(m)
}

func (m StdMap[K, V]) Merge(in MapI[K, V]) {
	mylog.CheckNil(m)
	in.Range(func(k K, v V) bool {
		m[k] = v
		return true
	})
}

func (m StdMap[K, V]) Range(f func(k K, v V) bool) {
	for k, v := range m {
		if !f(k, v) {
			break
		}
	}
}

func (m StdMap[K, V]) Load(k K) (v V) {
	mylog.CheckNil(m)
	v, ok := m[k]
	mylog.Check(ok)
	mylog.CheckNil(v)
	return
}

func (m StdMap[K, V]) Get(k K) (v V) {
	return m.Load(k)
}

func (m StdMap[K, V]) Has(k K) (exists bool) {
	_, ok := m[k]
	return ok
}

func (m StdMap[K, V]) Set(k K, v V) {
	mylog.CheckNil(m)
	m[k] = v
}

func (m StdMap[K, V]) Delete(k K) {
	delete(m, k)
}

func (m StdMap[K, V]) Keys() (keys []K) {
	if m.Len() == 0 {
		return
	}

	keys = make([]K, m.Len())

	var i int
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func (m StdMap[K, V]) Values() (values []V) {
	if m.Len() == 0 {
		return
	}
	values = make([]V, m.Len())
	var i int
	for _, v := range m {
		values[i] = v
		i++
	}
	return values
}

func (m StdMap[K, V]) Equal(m2 MapI[K, V]) bool {
	if m.Len() != m2.Len() {
		return false
	}
	ret := true
	m2.Range(func(k K, v V) bool {
		if v2, ok := m[k]; !ok || !equalValues(v, v2) {
			ret = false
			return false
		}
		return true
	})
	return ret
}

func (m StdMap[K, V]) String() string {
	s := fmt.Sprintf("%#v", m)
	loc := strings.IndexRune(s, '{')
	return s[loc:]
}

func (m StdMap[K, V]) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer

	enc := gob.NewEncoder(&b)
	mylog.Check(enc.Encode(map[K]V(m)))
	return b.Bytes(), nil
}

func (m *StdMap[K, V]) UnmarshalBinary(data []byte) (err error) {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	var v map[K]V
	mylog.Check(dec.Decode(&v))
	*m = v
	return
}

func (m StdMap[K, V]) MarshalJSON() (out []byte, err error) {
	v := map[K]V(m)
	return json.Marshal(v)
}

func (m *StdMap[K, V]) UnmarshalJSON(in []byte) (err error) {
	var v map[K]V
	mylog.Check(json.Unmarshal(in, &v))
	*m = v
	return
}
