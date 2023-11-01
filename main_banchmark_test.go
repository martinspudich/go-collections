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
	b.ReportAllocs()
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

func BenchmarkTimeExpiredList_Expired(b *testing.B) {
	tlist := &timeExpiredList[string]{
		duration:   1 * time.Nanosecond,
		data:       []expiredElement[string]{},
		dataString: []string{},
		quitChan:   make(chan struct{}),
	}
	for i := 0; i < b.N; i++ {
		tlist.Add(fmt.Sprintf("test-%d", i))
	}
	tlist.removeExpired()
}

func BenchmarkTimeExpiredList_Size(b *testing.B) {
	tlist := NewTimeExpiredList[string](10 * time.Millisecond)
	defer tlist.Discard()
	for i := 0; i < b.N; i++ {
		tlist.Add(fmt.Sprintf("test-%d", i))
	}
	tlist.Size()
}

func BenchmarkTimeExpiredMap_Clear(b *testing.B) {
	tmap := NewTimeExpiredMap[string, string](5 * time.Second)
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("test-%d", i)
		tmap.Add(key, key)
	}
	tmap.Clear()
}
