package safemap

import (
	"container/list"
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"sync"
)

const Ordered = true

type SafeMap[K comparable, V any] struct {
	sync.RWMutex
	m        map[K]V
	ordered  bool
	keys     *list.List
	keyIndex map[K]*list.Element
}

//thx github.com/hitsumitomo/safemap

type api[K comparable, V any] interface {
	New(ordered ...bool) (m *SafeMap[K, V])                                      //实例化
	NewOrdered(seq iter.Seq2[K, V]) (m *SafeMap[K, V])                           //实例化有序，std map的代码有语法检查，这个是实例化的时候检查，实例化语法间接性差不多
	NewStringer(ordered ...bool) (m *SafeMap[string, string])                    //从字符串实例化
	NewStringerKeys(keys []string, ordered ...bool) (m *SafeMap[string, string]) //从字符串切片实例化
	Has(key K) (exists bool)                                                     //是否存在
	Get(key K) (value V, exist bool)                                             //获取
	Update(key K, value V)                                                       //更新
	Delete(key K)                                                                //删除
	Remove(key K)                                                                //移除key
	Set(key K, value V) (actual V, exist bool)                                   //设置，如果存在则不更新
	GetAndDelete(key K) (value V, exist bool)                                    //获取后删除
	removeKey(key K)                                                             //移除key
	Range(f func(k K, v V) bool)                                                 //遍历，todo 回调内执行删除会死锁,vt调试器bind的时候需要
	Reset()                                                                      //清空
	Len() int                                                                    //大小
	Empty() bool                                                                 //大小为0
	Keys() []K                                                                   //键列表
	Map() map[K]V                                                                //原始map
	CopyFromMap(data map[K]V)                                                    //从map复制
	MarshalJSON() (data []byte, err error)                                       //
	UnmarshalJSON(data []byte) (err error)                                       //
	String() string                                                              //
}

func New[K comparable, V any](ordered ...bool) (m *SafeMap[K, V]) {
	sm := &SafeMap[K, V]{
		m:        make(map[K]V),
		keys:     list.New(),
		keyIndex: make(map[K]*list.Element),
	}
	if len(ordered) > 0 && ordered[0] {
		sm.ordered = true
	}
	return sm
}

func (s *SafeMap[K, V]) init(ordered ...bool) *SafeMap[K, V] {
	*s = SafeMap[K, V]{
		RWMutex:  sync.RWMutex{},
		m:        make(map[K]V),
		ordered:  false,
		keys:     list.New(),
		keyIndex: make(map[K]*list.Element),
	}
	s.keys.Init()
	if len(ordered) > 0 && ordered[0] {
		s.ordered = true
	}
	return s
}
func (s *SafeMap[K, V]) Collect(seq iter.Seq2[K, V]) *SafeMap[K, V] {
	s.checkInit()
	s.Lock()
	defer s.Unlock()
	for k, v := range seq {
		_, exist := s.Set(k, v)
		if exist {
			panic("duplicate key")
		}
	}
	return s
}
func (s *SafeMap[K, V]) Reset() {
	s.checkInit()
	s.Lock()
	defer s.Unlock()
	s.m = make(map[K]V)
	s.keys.Init()
	s.keyIndex = make(map[K]*list.Element)
}

func NewOrdered[K comparable, V any](seq iter.Seq2[K, V]) (m *SafeMap[K, V]) {
	m = New[K, V](true)
	for k, v := range seq {
		_, exist := m.Set(k, v)
		if exist {
			panic("duplicate key")
		}
	}
	return
}

func NewStringer(ordered ...bool) (m *SafeMap[string, string]) {
	return New[string, string](ordered...)
}

func NewStringerKeys(keys []string, ordered ...bool) (m *SafeMap[string, string]) {
	m = NewStringer(ordered...)
	for _, key := range keys {
		m.Set(key, key)
	}
	return
}

func (s *SafeMap[K, V]) Has(key K) (exists bool) {
	s.checkInit()
	s.RLock()
	defer s.RUnlock()
	_, exists = s.m[key]
	return exists
}

func (s *SafeMap[K, V]) Get(key K) (value V, exist bool) {
	s.checkInit()
	s.RLock()
	defer s.RUnlock()
	value, exist = s.m[key]
	return value, exist
}

func (s *SafeMap[K, V]) Update(key K, value V) {
	s.checkInit()
	s.Lock()
	defer s.Unlock()
	if elem, exists := s.keyIndex[key]; !exists {
		s.keyIndex[key] = s.keys.PushBack(key)
	} else {
		elem.Value = key
	}
	s.m[key] = value
}

func (s *SafeMap[K, V]) checkInit() {
	if s.m == nil { //new(safemap.SafeMap[K, V])这种方式实例化代码简洁
		s.init()
	}
}

func (s *SafeMap[K, V]) Set(key K, value V) (actual V, exist bool) {
	s.checkInit()
	s.Lock()
	defer s.Unlock()
	actual, exist = s.m[key]
	if exist {
		return actual, true //todo add log ?
	}
	s.m[key] = value
	if s.ordered {
		if elem, exists := s.keyIndex[key]; !exists {
			s.keyIndex[key] = s.keys.PushBack(key)
		} else {
			elem.Value = key
		}
	}
	return value, false
}

func (s *SafeMap[K, V]) Range(callback func(k K, v V) bool) {
	s.checkInit()
	if s.ordered {
		for e := s.keys.Front(); e != nil; e = e.Next() {
			k := e.Value.(K)
			if !callback(k, s.m[k]) {
				return
			}
		}
		return
	}
	for k, v := range maps.All(s.m) {
		if !callback(k, v) {
			return
		}
	}
}

func (s *SafeMap[K, V]) RangeKeys() iter.Seq[K] { //todo test 这种不支持有序遍历keys，使用Keys() 老方案遍历算了
	s.checkInit()
	return func(yield func(K) bool) {
		for k := range s.m {
			if !yield(k) {
				return
			}
		}
	}
}

func (s *SafeMap[K, V]) Keys() []K {
	s.checkInit()
	s.RLock()
	defer s.RUnlock()
	keys := make([]K, s.keys.Len())
	i := 0
	if s.ordered {
		for e := s.keys.Front(); e != nil; e = e.Next() {
			keys[i] = e.Value.(K)
			i++
		}
	} else {
		for k := range maps.Keys(s.m) {
			keys[i] = k
		}
	}
	return keys
}
func (s *SafeMap[K, V]) Values() []V {
	s.checkInit()
	s.RLock()
	defer s.RUnlock()

	values := make([]V, s.keys.Len())
	i := 0
	if s.ordered {
		for e := s.keys.Front(); e != nil; e = e.Next() {
			values[i] = s.m[e.Value.(K)]
			i++
		}
	} else {
		for v := range maps.Values(s.m) {
			values[i] = v
		}
	}
	return values
}
func (s *SafeMap[K, V]) All() iter.Seq2[K, V] { //todo 移除它，遍历就足够了
	s.checkInit()
	return func(yield func(k K, v V) bool) {
		s.Range(yield)
	}
}

func (s *SafeMap[K, V]) Remove(key K) {
	s.removeKey(key)
}
func (s *SafeMap[K, V]) Delete(key K) {
	s.removeKey(key)
}
func (s *SafeMap[K, V]) removeKey(key K) {
	s.checkInit()
	s.Lock()
	defer s.Unlock()

	if elem, exists := s.keyIndex[key]; exists {
		s.keys.Remove(elem)
		delete(s.keyIndex, key)
	}
	delete(s.m, key)
}

func (s *SafeMap[K, V]) GetAndDelete(key K) (value V, exist bool) {
	value, exist = s.m[key]
	if exist {
		s.removeKey(key)
	}
	return value, exist
}

func (s *SafeMap[K, V]) Len() int {
	s.checkInit()
	s.RLock()
	defer s.RUnlock()
	return len(s.m)
}

func (s *SafeMap[K, V]) Empty() bool {
	return s.Len() == 0
}

func (s *SafeMap[K, V]) Map() map[K]V {
	s.checkInit()
	s.RLock()
	defer s.RUnlock()
	return maps.Clone(s.m)
}
func (s *SafeMap[K, V]) CopyFromMap(data map[K]V) {
	s.checkInit()
	s.Lock()
	defer s.Unlock()
	s.m = maps.Clone(data) //todo add Clone method
	s.keys.Init()
	s.keyIndex = make(map[K]*list.Element)
	for k := range data {
		s.keyIndex[k] = s.keys.PushBack(k)
	}
}

func (s *SafeMap[K, V]) MarshalJSON() (data []byte, err error) {
	s.checkInit()
	s.RLock()
	defer s.RUnlock()
	return json.Marshal(s.m)
}

func (s *SafeMap[K, V]) UnmarshalJSON(data []byte) (err error) {
	s.checkInit()
	s.Lock()
	defer s.Unlock()
	var m map[K]V
	if err = json.Unmarshal(data, &m); err != nil {
		return err
	}
	s.m = m
	s.keys.Init()
	s.keyIndex = make(map[K]*list.Element)
	return nil
}
func (s *SafeMap[K, V]) String() string {
	s.checkInit()
	s.RLock()
	defer s.RUnlock()
	return fmt.Sprintf("%#v", s.m) //todo 使用结构体打印包格式化
}
