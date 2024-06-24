package orderedmap

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/ddkwork/golibrary/stream"

	"github.com/ddkwork/golibrary/mylog"
)

func New[K comparable, V any]() (m *OrderedMap[K, V]) {
	return &OrderedMap[K, V]{
		m: map[K]*list.Element{},
		l: list.New(),
	}
}

type pair[K comparable, V any] struct {
	Key   K
	Value V
}

type OrderedMap[K comparable, V any] struct {
	lock sync.RWMutex
	m    map[K]*list.Element
	l    *list.List
}

func (m *OrderedMap[K, V]) Map() map[K]V {
	r := map[K]V{}
	for k, e := range m.m {
		r[k] = e.Value.(pair[K, V]).Value
	}
	return r
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	_, found := m.m[key]
	return found
}

func (m *OrderedMap[K, V]) Keys() (keys []K) {
	for _, k := range m.List() {
		keys = append(keys, k.Key)
	}
	return
}

func (m *OrderedMap[K, V]) Values() (values []V) {
	for _, k := range m.List() {
		values = append(values, k.Value)
	}
	return
}

func (m *OrderedMap[K, V]) List() []pair[K, V] {
	l := make([]pair[K, V], 0, m.l.Len())
	for v := m.l.Front(); v != nil; v = v.Next() {
		l = append(l, v.Value.(pair[K, V]))
	}
	return l
}

func (m *OrderedMap[K, V]) Update(k K, v V) {
	if _, ok := m.m[k]; ok {
		m.m[k].Value = pair[K, V]{
			Key:   k,
			Value: v,
		}
	}
}

func (m *OrderedMap[K, V]) Set(k K, v V) {
	if _, ok := m.m[k]; ok {
		return
	}
	e := m.l.PushBack(pair[K, V]{
		Key:   k,
		Value: v,
	})
	m.m[k] = e
}

func (m *OrderedMap[K, V]) Get(k K) (v V, ok bool) {
	e, ok := m.m[k]
	if ok {
		return e.Value.(pair[K, V]).Value, true
	}
	return v, false
}

func (m *OrderedMap[K, V]) Front() (pair[K, V], bool) {
	elem := m.l.Front()
	if elem == nil {
		return pair[K, V]{}, false
	}
	return elem.Value.(pair[K, V]), true
}

func (m *OrderedMap[K, V]) Back() (pair[K, V], bool) {
	elem := m.l.Back()
	if elem == nil {
		return pair[K, V]{}, false
	}
	return elem.Value.(pair[K, V]), true
}

func (m *OrderedMap[K, V]) Prev(k K) (pair[K, V], bool) {
	elem, ok := m.m[k]
	if !ok {
		return pair[K, V]{}, false
	}
	elem = elem.Prev()
	if elem == nil {
		return pair[K, V]{}, false
	}
	return elem.Value.(pair[K, V]), true
}

func (m *OrderedMap[K, V]) Next(k K) (pair[K, V], bool) {
	elem, ok := m.m[k]
	if !ok {
		return pair[K, V]{}, false
	}
	elem = elem.Next()
	if elem == nil {
		return pair[K, V]{}, false
	}
	return elem.Value.(pair[K, V]), true
}

func (m *OrderedMap[K, V]) Delete(k K) {
	if e, ok := m.m[k]; ok {
		delete(m.m, k)
		m.l.Remove(e)
	}
}

func (m *OrderedMap[K, V]) Reset() {
	m.m = map[K]*list.Element{}
	m.l = list.New()
}

func (m *OrderedMap[K, V]) Len() int {
	return len(m.m)
}

func (m *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	kvs := m.List()
	l := len(kvs)
	if l == 0 {
		return []byte("{}"), nil
	}
	buf := bytes.NewBuffer(nil)
	buf.WriteString("{")
	for i, kv := range kvs {
		kBytes := mylog.Check2(json.Marshal(kv.Key))
		buf.Write(kBytes)
		buf.WriteString(":")
		vBytes := mylog.Check2(json.Marshal(kv.Value))
		buf.Write(vBytes)
		if i < l-1 {
			buf.Write([]byte(","))
		}
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (m *OrderedMap[K, V]) UnmarshalJSON(b []byte) error {
	tmp := map[string]V{}
	mylog.Check(json.Unmarshal(b, &tmp))
	objectKeys := mylog.Check2(m.objectKeys(b))
	m.Reset()
	for _, objectKey := range objectKeys {
		var k K
		mylog.Check(json.Unmarshal([]byte(fmt.Sprintf(`"%v"`, objectKey)), &k))
		m.Set(k, tmp[objectKey])
	}
	return nil
}

func (m *OrderedMap[K, V]) objectKeys(b []byte) ([]string, error) {
	d := json.NewDecoder(bytes.NewReader(b))
	t := mylog.Check2(d.Token())
	if t != json.Delim('{') {
		return nil, errors.New("expected start of object")
	}
	var keys []string
	for {
		t := mylog.Check2(d.Token())
		if t == json.Delim('}') {
			return keys, nil
		}
		keys = append(keys, t.(string))
		if e := m.skipValue(d); e != nil {
			return nil, e
		}
	}
}

func (m *OrderedMap[K, V]) skipValue(d *json.Decoder) error {
	t := mylog.Check2(d.Token())
	switch t {
	case json.Delim('['), json.Delim('{'):
		for {
			if e := m.skipValue(d); e != nil {
				if errors.Is(e, end) {
					break
				}
				return e
			}
		}
	case json.Delim(']'), json.Delim('}'):
		return end
	}
	return nil
}

var end = errors.New("invalid end of array or object")

func (m *OrderedMap[K, V]) InsertAfter(key K, value V) {
	if m.Contains(key) {
		mylog.Check("key already exists")
	}
	m.l.InsertAfter(pair[K, V]{Key: key, Value: value}, m.m[key])
}

func (m *OrderedMap[K, V]) InsertBefore(key K, value V) {
	if m.Contains(key) {
		mylog.Check("key already exists")
	}
	m.l.InsertBefore(pair[K, V]{Key: key, Value: value}, m.m[key])
}

func (m *OrderedMap[K, V]) CopyFrom(from *OrderedMap[K, V]) {
	for _, kv := range from.List() {
		m.Set(kv.Key, kv.Value)
	}
}

func (m *OrderedMap[K, V]) ApplyTo(to *OrderedMap[K, V]) {
	for _, p := range m.List() {
		to.Set(p.Key, p.Value)
	}
}

func (m *OrderedMap[K, V]) String() string {
	return fmt.Sprintf("%v", m.Map())
}

func (m *OrderedMap[K, V]) GoString() string {
	g := stream.NewGeneratedFile()
	g.P("var m = orderedmap.New[string, string]()")
	g.P()

	g.P("func init() {")

	for _, p := range m.List() {
		g.P("m.Set(" +
			strconv.Quote(fmt.Sprint(p.Key)) +
			", " +
			strconv.Quote(fmt.Sprint(p.Value)) +
			")")
	}
	g.P("}")
	return g.Buffer.String()
	return fmt.Sprintf("%#v", m.Map())
}
