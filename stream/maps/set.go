package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	"github.com/ddkwork/golibrary/mylog"
)

type Set[K comparable] struct {
	items StdMap[K, struct{}]
}

func (m *Set[K]) Clear() {
	m.items = nil
}

func (m *Set[K]) Len() int {
	return m.items.Len()
}

func (m *Set[K]) Range(f func(k K) bool) {
	mylog.CheckNil(m)
	if m.items == nil {
		return
	}
	mylog.CheckNil(m.items)
	for k := range m.items {
		if !f(k) {
			break
		}
	}
}

func (m *Set[K]) Has(k K) bool {
	return m.items.Has(k)
}

func (m *Set[K]) Delete(k K) {
	m.items.Delete(k)
}

func (m *Set[K]) Values() []K {
	return m.items.Keys()
}

func (m *Set[K]) Add(k ...K) SetI[K] {
	if m.items == nil {
		m.items = make(map[K]struct{})
	}
	for _, i := range k {
		m.items.Set(i, struct{}{})
	}
	return m
}

func (m *Set[K]) Merge(in SetI[K]) {
	mylog.CheckNil(m)
	if in == nil {
		return
	}
	mylog.CheckNil(in)
	if m.items == nil {
		m.items = make(map[K]struct{}, in.Len())
	}
	in.Range(func(k K) bool {
		m.items[k] = struct{}{}
		return true
	})
}

func (m *Set[K]) Equal(m2 SetI[K]) bool {
	if m.Len() != m2.Len() {
		return false
	}
	ret := true
	m2.Range(func(k K) bool {
		if !m.Has(k) {
			ret = false
			return false
		}
		return true
	})
	return ret
}

func (m *Set[K]) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	mylog.Check(enc.Encode(m.Values()))
	return b.Bytes(), nil
}

func (m *Set[K]) UnmarshalBinary(data []byte) (err error) {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	var v []K
	mylog.Check(dec.Decode(&v))
	for _, v2 := range v {
		m.Add(v2)
	}
	return
}

func (m *Set[K]) MarshalJSON() (out []byte, err error) {
	return json.Marshal(m.Values())
}

func (m *Set[K]) UnmarshalJSON(in []byte) (err error) {
	var v []K
	mylog.Check(json.Unmarshal(in, &v))
	for _, v2 := range v {
		m.Add(v2)
	}
	return
}

func (m *Set[K]) String() string {
	ret := "{"
	for i, v := range m.Values() {
		ret += fmt.Sprintf("%#v", v)
		if i < m.Len()-1 {
			ret += ","
		}
	}
	ret += "}"
	return ret
}
