package syncMap

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// Map is a concurrent Map with amortized-constant-time loads, stores, and deletes.
// It is safe for multiple goroutines to call a Map's methods concurrently.
//
// It is optimized for use in concurrent loops with keys that are
// stable over time, and either few steady-state stores, or stores
// localized to one goroutine per key.
//
// For use cases that do not share these attributes, it will likely have
// comparable or worse performance and worse type safety than an ordinary
// Map paired with a read-write mutex.
//
// The zero Map is valid and empty.
//
// A Map must not be copied after first use.
type (
	Interface interface {
		Load(key any) (value any, ok bool)
		Store(key, value any)
		LoadOrStore(key, value any) (actual any, loaded bool)
		Delete(key any) (res bool)
		Range(f func(key, value any) bool)
	}
	object struct {
		counter *int64
		mu      sync.Mutex

		// read contains the portion of the Map's contents that are safe for
		// concurrent access (with or without mu held).
		//
		// The read field itself is always safe to load, but must only be stored with
		// mu held.
		//
		// Entries stored in read may be updated concurrently without mu, but updating
		// a previously-expunged entry requires that the entry be copied to the dirty
		// Map and unexpunged with mu held.
		read atomic.Value // readOnly

		// dirty contains the portion of the Map's contents that require mu to be
		// held. To ensure that the dirty Map can be promoted to the read Map quickly,
		// it also includes all of the non-expunged entries in the read Map.
		//
		// Expunged entries are not stored in the dirty Map. An expunged entry in the
		// clean Map must be unexpunged and added to the dirty Map before a new value
		// can be stored to it.
		//
		// If the dirty Map is nil, the next write to the Map will initialize it by
		// making a shallow copy of the clean Map, omitting stale entries.
		dirty map[any]*entry

		// misses counts the number of loads since the read Map was last updated that
		// needed to lock mu to determine whether the key was present.
		//
		// Once enough misses have occurred to cover the cost of copying the dirty
		// Map, the dirty Map will be promoted to the read Map (in the unamended
		// state) and the next store to the Map will make a new dirty copy.
		misses int
	}
)

func New() Interface {
	return &object{}
}

// readOnly is an immutable struct stored atomically in the Map.read field.
type readOnly struct {
	m       map[any]*entry
	amended bool // true if the dirty Map contains some key not in m.
}

// expunged is an arbitrary pointer that marks entries which have been deleted
// from the dirty Map.
var expunged = unsafe.Pointer(new(any))

// An entry is a slot in the Map corresponding to a particular key.
type entry struct {
	// p points to the any value stored for the entry.
	//
	// If p == nil, the entry has been deleted and m.dirty == nil.
	//
	// If p == expunged, the entry has been deleted, m.dirty != nil, and the entry
	// is missing from m.dirty.
	//
	// Otherwise, the entry is valid and recorded in m.read.m[key] and, if m.dirty
	// != nil, in m.dirty[key].
	//
	// An entry can be deleted by atomic replacement with nil: when m.dirty is
	// next created, it will atomically replace nil with expunged and leave
	// m.dirty[key] unset.
	//
	// An entry's associated value can be updated by atomic replacement, provided
	// p != expunged. If p == expunged, an entry's associated value can be updated
	// only after first setting m.dirty[key] = e so that lookups using the dirty
	// Map find the entry.
	p unsafe.Pointer // *any
}

func newEntry(i any) *entry {
	return &entry{p: unsafe.Pointer(&i)}
}

// Load returns the value stored in the Map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the Map.
func (o *object) Load(key any) (value any, ok bool) {
	read, _ := o.read.Load().(readOnly)
	e, ok := read.m[key]
	if !ok && read.amended {
		o.mu.Lock()
		// Avoid reporting a spurious miss if o.dirty got promoted while we were
		// blocked on o.mu. (If further loads of the same key will not miss, it's
		// not worth copying the dirty Map for this key.)
		read, _ = o.read.Load().(readOnly)
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = o.dirty[key]
			// Regardless of whether the entry was present, record a miss: this key
			// will take the slow path until the dirty Map is promoted to the read
			// Map.
			o.missLocked()
		}
		o.mu.Unlock()
	}
	if !ok {
		return nil, false
	}
	return e.load()
}

func (e *entry) load() (value any, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expunged {
		return nil, false
	}
	return *(*any)(p), true
}

// Store sets the value for a key.
func (o *object) Store(key, value any) {
	read, _ := o.read.Load().(readOnly)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	o.mu.Lock()
	if o.counter == nil {
		v := new(int64)
		*v = 0
		o.counter = v
	}
	read, _ = o.read.Load().(readOnly)
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			// The entry was previously expunged, which implies that there is a
			// non-nil dirty Map and this entry is not in it.
			o.dirty[key] = e
		}
		e.storeLocked(&value)
	} else if e, ok := o.dirty[key]; ok {
		e.storeLocked(&value)
	} else {
		atomic.AddInt64(o.counter, 1)
		if !read.amended {
			// We're adding the first new key to the dirty Map.
			// Make sure it is allocated and mark the read-only Map as incomplete.
			o.dirtyLocked()
			o.read.Store(readOnly{m: read.m, amended: true})
		}
		o.dirty[key] = newEntry(value)
	}
	o.mu.Unlock()
}

// tryStore stores a value if the entry has not been expunged.
//
// If the entry is expunged, tryStore returns false and leaves the entry
// unchanged.
func (e *entry) tryStore(i *any) (ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == expunged {
		return
	}
	for {
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return
		}
		p = atomic.LoadPointer(&e.p)
		if p == expunged {
			return true
		}
	}
}

// unexpungeLocked ensures that the entry is not marked as expunged.
//
// If the entry was previously expunged, it must be added to the dirty Map
// before m.mu is unlocked.
func (e *entry) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expunged, nil)
}

// storeLocked unconditionally stores a value to the entry.
//
// The entry must be known not to be expunged.
func (e *entry) storeLocked(i *any) {
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (o *object) LoadOrStore(key, value any) (actual any, loaded bool) {
	// Avoid locking if it's a clean hit.
	read, _ := o.read.Load().(readOnly)
	if e, ok := read.m[key]; ok {
		actual, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual, loaded
		}
	}

	o.mu.Lock()
	read, _ = o.read.Load().(readOnly)
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			o.dirty[key] = e
		}
		actual, loaded, _ = e.tryLoadOrStore(value)
	} else if e, ok := o.dirty[key]; ok {
		actual, loaded, _ = e.tryLoadOrStore(value)
		o.missLocked()
	} else {
		if !read.amended {
			// We're adding the first new key to the dirty Map.
			// Make sure it is allocated and mark the read-only Map as incomplete.
			o.dirtyLocked()
			o.read.Store(readOnly{m: read.m, amended: true})
		}
		o.dirty[key] = newEntry(value)
		actual, loaded = value, false
	}
	o.mu.Unlock()

	return actual, loaded
}

// tryLoadOrStore atomically loads or stores a value if the entry is not
// expunged.
//
// If the entry is expunged, tryLoadOrStore leaves the entry unchanged and
// returns with ok==false.
func (e *entry) tryLoadOrStore(i any) (actual any, loaded, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == expunged {
		return nil, false, false
	}
	if p != nil {
		return *(*any)(p), true, true
	}

	// Copy the interface after the first load to make this method more amenable
	// to escape analysis: if we hit the "load" path or the entry is expunged, we
	// shouldn't bother heap-allocating.
	ic := i
	for {
		if atomic.CompareAndSwapPointer(&e.p, nil, unsafe.Pointer(&ic)) {
			return i, false, true
		}
		p = atomic.LoadPointer(&e.p)
		if p == expunged {
			return nil, false, false
		}
		if p != nil {
			return *(*any)(p), true, true
		}
	}
}

// Delete deletes the value for a key.
func (o *object) Delete(key any) (res bool) {
	read, _ := o.read.Load().(readOnly)
	e, ok := read.m[key]
	if !ok && read.amended {
		o.mu.Lock()
		read, _ = o.read.Load().(readOnly)
		e, ok = read.m[key]
		if !ok && read.amended {
			delete(o.dirty, key)
			atomic.AddInt64(o.counter, -1)
			res = true
		}
		o.mu.Unlock()
	}
	if ok {
		e.delete()
		res = true
		atomic.AddInt64(o.counter, -1)
	}
	return res
}

func (o *object) Length() *int64 {
	return o.counter
}

func (e *entry) delete() bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expunged {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return true
		}
	}
}

// Range calls f sequentially for each key and value present in the Map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently, Range may reflect any mapping for that key
// from any point during the Range call.
//
// Range may be O(N) with the number of elements in the Map even if f returns
// false after a constant number of calls.
func (o *object) Range(f func(key, value any) bool) {
	// We need to be able to iterate over all of the keys that were already
	// present at the start of the call to Range.
	// If read.amended is false, then read.o satisfies that property without
	// requiring us to hold o.mu for a long time.
	read, _ := o.read.Load().(readOnly)
	if read.amended {
		// o.dirty contains keys not in read.o. Fortunately, Range is already O(N)
		// (assuming the caller does not break out early), so a call to Range
		// amortizes an entire copy of the Map: we can promote the dirty copy
		// immediately!
		o.mu.Lock()
		read, _ = o.read.Load().(readOnly)
		if read.amended {
			read = readOnly{m: o.dirty}
			o.read.Store(read)
			o.dirty = nil
			o.misses = 0
		}
		o.mu.Unlock()
	}

	for k, e := range read.m {
		v, ok := e.load()
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}

func (o *object) missLocked() {
	o.misses++
	if o.misses < len(o.dirty) {
		return
	}
	o.read.Store(readOnly{m: o.dirty})
	o.dirty = nil
	o.misses = 0
}

func (o *object) dirtyLocked() {
	if o.dirty != nil {
		return
	}

	read, _ := o.read.Load().(readOnly)
	o.dirty = make(map[any]*entry, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			o.dirty[k] = e
		}
	}
}

func (e *entry) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expunged) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expunged
}
