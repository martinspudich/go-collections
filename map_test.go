package gocollections

import (
	"reflect"
	"slices"
	"testing"
)

func TestCopyMap(t *testing.T) {
	tests := []struct {
		name string
		src  map[string]any
		want map[string]any
	}{
		{
			name: "simple copy - value int",
			src:  map[string]any{"a": 1, "b": 2},
			want: map[string]any{"a": 1, "b": 2},
		},
		{
			name: "simple copy - value string",
			src:  map[string]any{"a": "1", "b": "2"},
			want: map[string]any{"a": "1", "b": "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := map[string]any{}
			CopyMap(got, tt.src)
			result := reflect.DeepEqual(got, tt.want)
			if result != true {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestMapMutex(t *testing.T) {
	m := NewMapMutex[string, int]()

	if !m.IsEmpty() {
		t.Error("expected empty map")
	}

	m.Put("a", 1)
	m.Put("b", 2)

	if m.Size() != 2 {
		t.Errorf("expected size 2, got %d", m.Size())
	}

	val, ok := m.Get("a")
	if !ok || val != 1 {
		t.Errorf("expected 1, got %v", val)
	}

	if !m.ContainsKey("b") {
		t.Error("expected to contain key b")
	}

	if !m.ContainsValue(2) {
		t.Error("expected to contain value 2")
	}

	old := m.Put("a", 10)
	if old != 1 {
		t.Errorf("expected old value 1, got %v", old)
	}

	val, _ = m.Get("a")
	if val != 10 {
		t.Errorf("expected 10, got %v", val)
	}

	removed, ok := m.Remove("b")
	if !ok || removed != 2 {
		t.Errorf("expected removed value 2, got %v", removed)
	}

	if m.Size() != 1 {
		t.Errorf("expected size 1, got %d", m.Size())
	}

	keys := m.KeySet()
	if !slices.Contains(keys, "a") || len(keys) != 1 {
		t.Errorf("unexpected KeySet: %v", keys)
	}

	values := m.Values()
	if !slices.Contains(values, 10) || len(values) != 1 {
		t.Errorf("unexpected Values: %v", values)
	}

	m.Clear()
	if !m.IsEmpty() {
		t.Error("expected empty map after clear")
	}
}

func TestMapMutex_PutAll(t *testing.T) {
	m1 := NewMapMutex[string, int]()
	m1.Put("a", 1)
	m1.Put("b", 2)

	m2 := NewMapMutex[string, int]()
	m2.Put("c", 3)
	m2.PutAll(m1)

	if m2.Size() != 3 {
		t.Errorf("expected size 3, got %d", m2.Size())
	}

	for k, v := range m1.All() {
		val, ok := m2.Get(k)
		if !ok || val != v {
			t.Errorf("expected key %s to have value %v in m2", k, v)
		}
	}
}

func TestMapMutex_Iterators(t *testing.T) {
	m := NewMapMutex[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)

	// All()
	count := 0
	for k, v := range m.All() {
		if k == "a" && v != 1 {
			t.Errorf("wrong value for a: %v", v)
		}
		count++
	}
	if count != 2 {
		t.Errorf("expected 2 iterations, got %d", count)
	}

	// Keys()
	count = 0
	for k := range m.Keys() {
		if k != "a" && k != "b" {
			t.Errorf("unexpected key: %s", k)
		}
		count++
	}
	if count != 2 {
		t.Errorf("expected 2 keys, got %d", count)
	}

	// ValueSeq()
	count = 0
	for v := range m.ValueSeq() {
		if v != 1 && v != 2 {
			t.Errorf("unexpected value: %d", v)
		}
		count++
	}
	if count != 2 {
		t.Errorf("expected 2 values, got %d", count)
	}
}
