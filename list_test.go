package gocollections

import (
	"reflect"
	"testing"
)

func TestListMutex(t *testing.T) {
	l := NewListMutex[int]()

	if !l.IsEmpty() {
		t.Error("List should be empty")
	}

	l.Add(1)
	l.Add(2)
	l.Add(3)

	if l.Size() != 3 {
		t.Errorf("Expected size 3, got %d", l.Size())
	}

	if l.IsEmpty() {
		t.Error("List should not be empty")
	}

	val, err := l.Get(1)
	if err != nil || val != 2 {
		t.Errorf("Expected 2, got %v (err: %v)", val, err)
	}

	old, err := l.Set(1, 20)
	if err != nil || old != 2 {
		t.Errorf("Expected old value 2, got %v (err: %v)", old, err)
	}

	val, _ = l.Get(1)
	if val != 20 {
		t.Errorf("Expected 20, got %v", val)
	}

	err = l.AddAt(1, 15)
	if err != nil {
		t.Errorf("AddAt failed: %v", err)
	}

	if l.Size() != 4 {
		t.Errorf("Expected size 4, got %d", l.Size())
	}

	val, _ = l.Get(1)
	if val != 15 {
		t.Errorf("Expected 15 at index 1, got %v", val)
	}

	removed, err := l.Remove(2)
	if err != nil || removed != 20 {
		t.Errorf("Expected removed 20, got %v (err: %v)", removed, err)
	}

	if l.Size() != 3 {
		t.Errorf("Expected size 3 after removal, got %d", l.Size())
	}

	if !l.Contains(15) {
		t.Error("List should contain 15")
	}

	if l.Contains(100) {
		t.Error("List should not contain 100")
	}

	all := l.GetAll()
	expected := []int{1, 15, 3}
	if !reflect.DeepEqual(all, expected) {
		t.Errorf("Expected %v, got %v", expected, all)
	}

	l.Clear()
	if l.Size() != 0 || !l.IsEmpty() {
		t.Error("List should be empty after Clear")
	}
}

func TestListMutex_Bounds(t *testing.T) {
	l := NewListMutex[any]()

	_, err := l.Get(0)
	if err != ErrIndexOutOfBounds {
		t.Errorf("Expected ErrIndexOutOfBounds, got %v", err)
	}

	err = l.AddAt(1, "test")
	if err != ErrIndexOutOfBounds {
		t.Errorf("Expected ErrIndexOutOfBounds, got %v", err)
	}

	_, err = l.Remove(0)
	if err != ErrIndexOutOfBounds {
		t.Errorf("Expected ErrIndexOutOfBounds, got %v", err)
	}

	l.Add("a")
	_, err = l.Set(1, "b")
	if err != ErrIndexOutOfBounds {
		t.Errorf("Expected ErrIndexOutOfBounds, got %v", err)
	}
}
