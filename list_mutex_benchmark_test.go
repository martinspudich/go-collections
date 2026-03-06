package gocollections

import (
	"fmt"
	"testing"
)

func BenchmarkListMutex_Add(b *testing.B) {
	l := NewListMutex[string]()
	b.ResetTimer()
	for i := range b.N {
		l.Add(fmt.Sprintf("item-%d", i))
	}
}

func BenchmarkListMutex_Get(b *testing.B) {
	l := NewListMutex[string]()
	for i := range 1000 {
		l.Add(fmt.Sprintf("item-%d", i))
	}
	b.ResetTimer()
	for i := range b.N {
		_, _ = l.Get(i % 1000)
	}
}

func BenchmarkListMutex_Remove(b *testing.B) {
	// Note: Removing from the beginning of a slice is O(n), so we benchmark it accordingly.
	// To avoid an empty list, we add elements in each iteration or pre-fill.
	// However, standard benchmarks usually pre-fill.
	b.Run("RemoveFirst", func(b *testing.B) {
		for range b.N {
			b.StopTimer()
			l := NewListMutex[string]()
			for j := range 100 {
				l.Add(fmt.Sprintf("item-%d", j))
			}
			b.StartTimer()
			_, _ = l.Remove(0)
		}
	})
	b.Run("RemoveLast", func(b *testing.B) {
		for range b.N {
			b.StopTimer()
			l := NewListMutex[string]()
			for j := range 100 {
				l.Add(fmt.Sprintf("item-%d", j))
			}
			b.StartTimer()
			_, _ = l.Remove(99)
		}
	})
}

func BenchmarkListMutex_Contains(b *testing.B) {
	l := NewListMutex[string]()
	for i := range 1000 {
		l.Add(fmt.Sprintf("item-%d", i))
	}
	b.ResetTimer()
	for range b.N {
		l.Contains("item-500")
	}
}
