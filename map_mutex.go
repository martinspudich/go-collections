package gocollections

import (
	"iter"
	"maps"
	"reflect"
	"slices"
	"sync"
)

// MapMutex is a thread-safe implementation of the Map interface using sync.RWMutex.
type MapMutex[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// NewMapMutex creates a new instance of MapMutex.
func NewMapMutex[K comparable, V any]() *MapMutex[K, V] {
	return &MapMutex[K, V]{
		data: make(map[K]V),
	}
}

// Put associates the specified value with the specified key in this map.
// Returns the previous value associated with key, or zero value if there was no mapping for key.
func (m *MapMutex[K, V]) Put(key K, value V) V {
	m.mu.Lock()
	defer m.mu.Unlock()
	old := m.data[key]
	m.data[key] = value
	return old
}

// Get returns the value to which the specified key is mapped, and a boolean indicating if the key was found.
func (m *MapMutex[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}

// Remove removes the mapping for a key from this map if it is present.
// Returns the value to which this map previously associated the key, and a boolean indicating if the key was found.
func (m *MapMutex[K, V]) Remove(key K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	val, ok := m.data[key]
	if ok {
		delete(m.data, key)
	}
	return val, ok
}

// ContainsKey returns true if this map contains a mapping for the specified key.
func (m *MapMutex[K, V]) ContainsKey(key K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.data[key]
	return ok
}

// ContainsValue returns true if this map maps one or more keys to the specified value.
func (m *MapMutex[K, V]) ContainsValue(value V) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, v := range m.data {
		if reflect.DeepEqual(v, value) {
			return true
		}
	}
	return false
}

// Size returns the number of key-value mappings in this map.
func (m *MapMutex[K, V]) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// IsEmpty returns true if this map contains no key-value mappings.
func (m *MapMutex[K, V]) IsEmpty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data) == 0
}

// Clear removes all of the mappings from this map.
func (m *MapMutex[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	clear(m.data)
}

// PutAll copies all of the mappings from the specified map to this map.
func (m *MapMutex[K, V]) PutAll(other Map[K, V]) {
	for k, v := range other.All() {
		m.Put(k, v)
	}
}

// KeySet returns a slice containing all of the keys contained in this map.
func (m *MapMutex[K, V]) KeySet() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return slices.Collect(maps.Keys(m.data))
}

// Values returns a slice containing all of the values contained in this map.
func (m *MapMutex[K, V]) Values() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return slices.Collect(maps.Values(m.data))
}

// All returns an iterator for all key-value pairs in the map.
func (m *MapMutex[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()
		for k, v := range m.data {
			if !yield(k, v) {
				return
			}
		}
	}
}

// Keys returns an iterator for all keys in the map.
func (m *MapMutex[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()
		for k := range m.data {
			if !yield(k) {
				return
			}
		}
	}
}

// ValueSeq returns an iterator for all values in the map.
func (m *MapMutex[K, V]) ValueSeq() iter.Seq[V] {
	return func(yield func(V) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()
		for _, v := range m.data {
			if !yield(v) {
				return
			}
		}
	}
}
