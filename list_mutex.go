package gocollections

import (
	"reflect"
	"sync"
)

// ListMutex is a thread-safe implementation of the List interface using sync.RWMutex.
type ListMutex[T any] struct {
	mu    sync.RWMutex
	items []T
}

// NewListMutex creates a new instance of ListMutex.
func NewListMutex[T any]() *ListMutex[T] {
	return &ListMutex[T]{
		items: make([]T, 0),
	}
}

// Add appends the specified value to the end of the list.
func (l *ListMutex[T]) Add(value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.items = append(l.items, value)
}

// AddAt inserts the specified value at the specified position in the list.
// Returns ErrIndexOutOfBounds if the index is out of range.
func (l *ListMutex[T]) AddAt(index int, value T) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if index < 0 || index > len(l.items) {
		return ErrIndexOutOfBounds
	}
	var zero T
	l.items = append(l.items, zero)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = value
	return nil
}

// Get returns the element at the specified position in the list.
// Returns ErrIndexOutOfBounds if the index is out of range.
func (l *ListMutex[T]) Get(index int) (T, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if index < 0 || index >= len(l.items) {
		var zero T
		return zero, ErrIndexOutOfBounds
	}
	return l.items[index], nil
}

// Set replaces the element at the specified position in the list with the specified value.
// Returns the old value and ErrIndexOutOfBounds if the index is out of range.
func (l *ListMutex[T]) Set(index int, value T) (T, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if index < 0 || index >= len(l.items) {
		var zero T
		return zero, ErrIndexOutOfBounds
	}
	old := l.items[index]
	l.items[index] = value
	return old, nil
}

// Remove removes the element at the specified position in the list.
// Returns the removed value and ErrIndexOutOfBounds if the index is out of range.
func (l *ListMutex[T]) Remove(index int) (T, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if index < 0 || index >= len(l.items) {
		var zero T
		return zero, ErrIndexOutOfBounds
	}
	old := l.items[index]
	l.items = append(l.items[:index], l.items[index+1:]...)
	return old, nil
}

// Clear removes all elements from the list.
func (l *ListMutex[T]) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.items = make([]T, 0)
}

// Size returns the number of elements in the list.
func (l *ListMutex[T]) Size() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.items)
}

// IsEmpty returns true if the list contains no elements.
func (l *ListMutex[T]) IsEmpty() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.items) == 0
}

// Contains returns true if the list contains the specified value.
func (l *ListMutex[T]) Contains(value T) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, item := range l.items {
		if reflect.DeepEqual(item, value) {
			return true
		}
	}
	return false
}

// GetAll returns a slice containing all elements in the list.
func (l *ListMutex[T]) GetAll() []T {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return append([]T(nil), l.items...)
}
