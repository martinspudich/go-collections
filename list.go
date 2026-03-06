package gocollections

import "errors"

// ErrIndexOutOfBounds is returned when an operation is performed on an index that is out of range.
var (
	ErrIndexOutOfBounds = errors.New("index out of bounds")
)

// List interface defines the standard operations for a list-like collection.
type List[T any] interface {
	// Add appends the specified value to the end of the list.
	Add(value T)
	// AddAt inserts the specified value at the specified position in the list.
	// Returns ErrIndexOutOfBounds if the index is out of range.
	AddAt(index int, value T) error
	// Get returns the element at the specified position in the list.
	// Returns ErrIndexOutOfBounds if the index is out of range.
	Get(index int) (T, error)
	// Set replaces the element at the specified position in the list with the specified value.
	// Returns the old value and ErrIndexOutOfBounds if the index is out of range.
	Set(index int, value T) (T, error)
	// Remove removes the element at the specified position in the list.
	// Returns the removed value and ErrIndexOutOfBounds if the index is out of range.
	Remove(index int) (T, error)
	// Clear removes all elements from the list.
	Clear()
	// Size returns the number of elements in the list.
	Size() int
	// IsEmpty returns true if the list contains no elements.
	IsEmpty() bool
	// Contains returns true if the list contains the specified value.
	Contains(value T) bool
	// GetAll returns a slice containing all elements in the list.
	GetAll() []T
}
