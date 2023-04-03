package gocollections

import (
	"fmt"
	"testing"
	"time"
)

var gs string

func BenchmarkTimeExpiredList(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = fmt.Sprint("hello")
	}
	gs = s
}

func BenchmarkSlice(b *testing.B) {
	var list []string
	for i := 0; i < b.N; i++ {
		list = append(list, fmt.Sprintf("test-%d", i))
	}
}

func BenchmarkTimeExpiredList_Add(b *testing.B) {
	tlist := NewTimeExpiredList[string](1 * time.Second)
	defer tlist.Discard()
	for i := 0; i < b.N; i++ {
		tlist.Add(fmt.Sprintf("test-%d", i))
	}
}
