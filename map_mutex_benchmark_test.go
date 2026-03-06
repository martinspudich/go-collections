package gocollections

import (
	"fmt"
	"testing"
)

func BenchmarkMapMutex_Put(b *testing.B) {
	m := NewMapMutex[string, string]()
	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("key-%d", i)
		m.Put(key, "value")
	}
}

func BenchmarkMapMutex_Get(b *testing.B) {
	m := NewMapMutex[string, string]()
	for i := range 1000 {
		m.Put(fmt.Sprintf("key-%d", i), "value")
	}
	b.ResetTimer()
	for i := range b.N {
		_, _ = m.Get(fmt.Sprintf("key-%d", i%1000))
	}
}

func BenchmarkMapMutex_Remove(b *testing.B) {
	for range b.N {
		b.StopTimer()
		m := NewMapMutex[string, string]()
		for j := range 100 {
			m.Put(fmt.Sprintf("key-%d", j), "value")
		}
		b.StartTimer()
		_, _ = m.Remove("key-50")
	}
}

func BenchmarkMapMutex_ContainsKey(b *testing.B) {
	m := NewMapMutex[string, string]()
	for i := range 1000 {
		m.Put(fmt.Sprintf("key-%d", i), "value")
	}
	b.ResetTimer()
	for i := range b.N {
		m.ContainsKey(fmt.Sprintf("key-%d", i%1000))
	}
}
