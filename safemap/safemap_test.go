package safemap

import (
	"strconv"
	"sync"
	"testing"
)

func TestName(t *testing.T) {
	m := NewOrdered[string, int](func(yield func(string, int) bool) {
		yield("a", 1)
		yield("b", 2)
		yield("c", 3)
		yield("d", 4)
		yield("e", 5)
		yield("f", 6)
		yield("g", 7)
		yield("h", 8)
		yield("i", 9)
		yield("j", 10)
		yield("k", 11)
	})
	for k := range m.Range() {
		if k == "c" {
			m.Delete(k)
		}
		// println(strconv.Quote(k), v)
	}

	println(m.Empty())
	for k, v := range m.Range() {
		println(strconv.Quote(k), v)
	}

	//for k, v := range m.All() {
	//	println(strconv.Quote(k), v)
	//}
}

func TestSafeMap_StoreAndLoad(t *testing.T) {
	sm := New[int, string]()
	sm.Update(1, "one")

	value, ok := sm.Get(1)
	if !ok || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}
}

func TestSafeMap_Exists(t *testing.T) {
	sm := New[int, string]()
	sm.Update(1, "one")

	if !sm.Has(1) {
		t.Errorf("expected key 1 to exist")
	}

	if sm.Has(2) {
		t.Errorf("expected key 2 to not exist")
	}
}

func TestSafeMap_Delete(t *testing.T) {
	sm := New[int, string]()
	sm.Update(1, "one")
	sm.Delete(1)

	if sm.Has(1) {
		t.Errorf("expected key 1 to be deleted")
	}
}

func TestSafeMap_LoadAndDelete(t *testing.T) {
	sm := New[int, string]()
	sm.Update(1, "one")

	value, ok := sm.GetAndDelete(1)
	if !ok || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}

	if sm.Has(1) {
		t.Errorf("expected key 1 to be deleted")
	}
}

func TestSafeMap_LoadOrStore(t *testing.T) {
	sm := New[int, string]()
	value, loaded := sm.Set(1, "one")
	if loaded || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}

	value, loaded = sm.Set(1, "two")
	if !loaded || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}
}

func TestSafeMap_Range(t *testing.T) {
	sm := New[int, string]()
	sm.Update(1, "one")
	sm.Update(2, "two")

	keys := make(map[int]bool)
	for k := range sm.Range() {
		keys[k] = true
	}

	if len(keys) != 2 || !keys[1] || !keys[2] {
		t.Errorf("expected keys 1 and 2, got '%v'", keys)
	}
}

func TestSafeMap_Clear(t *testing.T) {
	sm := New[int, string]()
	sm.Update(1, "one")
	sm.Reset()

	if sm.Len() != 0 {
		t.Errorf("expected length 0, got '%v'", sm.Len())
	}
}

func TestSafeMap_StoreAndLoad_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Update(1, "one")

	value, ok := sm.Get(1)
	if !ok || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}
}

func TestSafeMap_Exists_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Update(1, "one")

	if !sm.Has(1) {
		t.Errorf("expected key 1 to exist")
	}

	if sm.Has(2) {
		t.Errorf("expected key 2 to not exist")
	}
}

func TestSafeMap_Delete_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Update(1, "one")
	sm.Delete(1)

	if sm.Has(1) {
		t.Errorf("expected key 1 to be deleted")
	}
}

func TestSafeMap_LoadAndDelete_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Update(1, "one")

	value, ok := sm.GetAndDelete(1)
	if !ok || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}

	if sm.Has(1) {
		t.Errorf("expected key 1 to be deleted")
	}
}

func TestSafeMap_LoadOrStore_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	value, loaded := sm.Set(1, "one")
	if loaded || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}

	value, loaded = sm.Set(1, "two")
	if !loaded || value != "one" {
		t.Errorf("expected 'one', got '%v'", value)
	}
}

func TestSafeMap_Range_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Update(1, "one")
	sm.Update(2, "two")

	keys := make(map[int]bool)
	for k := range sm.Range() {
		keys[k] = true
	}

	if len(keys) != 2 || !keys[1] || !keys[2] {
		t.Errorf("expected keys 1 and 2, got '%v'", keys)
	}
}

func TestSafeMap_Clear_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Update(1, "one")
	sm.Reset()

	if sm.Len() != 0 {
		t.Errorf("expected length 0, got '%v'", sm.Len())
	}
}

func TestSafeMap_Keys_Ordered(t *testing.T) {
	sm := New[int, string](Ordered)
	sm.Update(1, "one")
	sm.Update(2, "two")

	keys := sm.Keys()
	if len(keys) != 2 || (keys[0] != 1 && keys[1] != 2) {
		t.Errorf("expected keys 1 and 2, got '%v'", keys)
	}
}

func BenchmarkSafeMap_Store(b *testing.B) {
	sm := New[int, string]()
	for i := 0; b.Loop(); i++ {
		sm.Update(i, "value")
	}
}

func BenchmarkSafeMap_Load(b *testing.B) {
	sm := New[int, string]()
	sm.Update(1, "value")

	for b.Loop() {
		sm.Get(1)
	}
}

func BenchmarkSafeMap_Delete(b *testing.B) {
	sm := New[int, string]()
	for i := 0; b.Loop(); i++ {
		sm.Update(i, "value")
		sm.Delete(i)
	}
}

func BenchmarkSafeMap_Ordered_Store(b *testing.B) {
	sm := New[int, string](Ordered)
	for i := 0; b.Loop(); i++ {
		sm.Update(i, "value")
	}
}

func BenchmarkSafeMap_Ordered_Load(b *testing.B) {
	sm := New[int, string](Ordered)
	sm.Update(1, "value")

	for b.Loop() {
		sm.Get(1)
	}
}

func BenchmarkSafeMap_Ordered_Delete(b *testing.B) {
	sm := New[int, string](Ordered)
	for i := 0; b.Loop(); i++ {
		sm.Update(i, "value")
		sm.Delete(i)
	}
}

func BenchmarkSyncMap_Store(b *testing.B) {
	var sm sync.Map
	for i := 0; b.Loop(); i++ {
		sm.Store(i, "value")
	}
}

func BenchmarkSyncMap_Load(b *testing.B) {
	var sm sync.Map
	sm.Store(1, "value")

	for b.Loop() {
		sm.Load(1)
	}
}

func BenchmarkSyncMap_Delete(b *testing.B) {
	var sm sync.Map
	for i := 0; b.Loop(); i++ {
		sm.Store(i, "value")
		sm.Delete(i)
	}
}
