package gocollections

import (
	"iter"
)

// Map interface based on Java Map interface.
type Map[K comparable, V any] interface {
	// Put associates the specified value with the specified key in this map.
	// Returns the previous value associated with key, or zero value if there was no mapping for key.
	Put(key K, value V) V
	// Get returns the value to which the specified key is mapped, and a boolean indicating if the key was found.
	Get(key K) (V, bool)
	// Remove removes the mapping for a key from this map if it is present.
	// Returns the value to which this map previously associated the key, and a boolean indicating if the key was found.
	Remove(key K) (V, bool)
	// ContainsKey returns true if this map contains a mapping for the specified key.
	ContainsKey(key K) bool
	// ContainsValue returns true if this map maps one or more keys to the specified value.
	ContainsValue(value V) bool
	// Size returns the number of key-value mappings in this map.
	Size() int
	// IsEmpty returns true if this map contains no key-value mappings.
	IsEmpty() bool
	// Clear removes all of the mappings from this map.
	Clear()
	// PutAll copies all of the mappings from the specified map to this map.
	PutAll(other Map[K, V])
	// KeySet returns a slice containing all of the keys contained in this map.
	KeySet() []K
	// Values returns a slice containing all of the values contained in this map.
	Values() []V
	// All returns an iterator for all key-value pairs in the map.
	All() iter.Seq2[K, V]
	// Keys returns an iterator for all keys in the map.
	Keys() iter.Seq[K]
	// ValueSeq returns an iterator for all values in the map.
	ValueSeq() iter.Seq[V]
}

// CopyMap function will creates hard copy of the source map to destination map.
func CopyMap[K comparable, V any](dst, src map[K]V) {
	for k, v := range src {
		dst[k] = v
	}
}
