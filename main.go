package gocollections

import (
	"errors"
	"sync"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")
var ErrIndexOutOfBound = errors.New("index out of bound")

/*
Time Expired Map
*/

// TimeExpiredMap implementation of this interface is a map with elements which are expire base on expiration duration.
// Implementation of this map is running goroutine which removes expired element. To stop this goroutine call Discard()
// method when this map is not needed any more.
type TimeExpiredMap[K comparable, V any] interface {
	Add(key K, object V)
	AddWithDuration(key K, data V, duration time.Duration)
	Get(key K) (V, error)
	Del(key K) error
	Contains(key K) bool
	Size() int
	Discard()
}

type expiredElement[V any] struct {
	data      V
	expiredAt time.Time
}

type timeExpiredMap[K comparable, V any] struct {
	sync.Mutex
	duration time.Duration           // default element duration
	data     map[K]expiredElement[V] // map of elements
	quitChan chan struct{}           // channel for indicating to end goroutines for removing expired elements
}

// NewTimeExpiredMap creates new TimeExpiredMap object.
func NewTimeExpiredMap[K comparable, V any](duration time.Duration) TimeExpiredMap[K, V] {
	tmap := &timeExpiredMap[K, V]{
		duration: duration,
		data:     make(map[K]expiredElement[V]),
		quitChan: make(chan struct{}),
	}

	go tmap.run()

	return tmap
}

// Add method adds element to the map with key.
func (m *timeExpiredMap[K, V]) Add(key K, data V) {
	m.Lock()
	defer m.Unlock()
	m.data[key] = expiredElement[V]{expiredAt: time.Now().Add(m.duration), data: data}
}

// AddWithDuration adds element to the map with key. It will set custom duration time of the element in the internal map.
func (m *timeExpiredMap[K, V]) AddWithDuration(key K, data V, duration time.Duration) {
	m.Lock()
	defer m.Unlock()
	m.data[key] = expiredElement[V]{expiredAt: time.Now().Add(duration), data: data}
}

// Get method returns element by key.
func (m *timeExpiredMap[K, V]) Get(key K) (V, error) {
	var result V
	if !m.Contains(key) {
		return result, ErrKeyNotFound
	}
	return m.data[key].data, nil
}

// Del method removes element from map.
func (m *timeExpiredMap[K, V]) Del(key K) error {
	if !m.Contains(key) {
		return ErrKeyNotFound
	}
	m.Lock()
	defer m.Unlock()
	delete(m.data, key)
	return nil
}

// Contains method returns true if key is in the map. Else return false.
func (m *timeExpiredMap[K, V]) Contains(key K) bool {
	_, found := m.data[key]
	return found
}

// Size method returns size of the map.
func (m *timeExpiredMap[K, V]) Size() int {
	return len(m.data)
}

// Discard method stops the goroutine for removing elements and discards data in internal map.
func (m *timeExpiredMap[K, V]) Discard() {
	m.quitChan <- struct{}{}
	m.Lock()
	defer m.Unlock()
	m.data = nil
}

// run method runs the goroutine for removing expired elements.
func (m *timeExpiredMap[K, V]) run() {
	for {
		select {
		case <-time.After(1 * time.Second):
			m.removeExpired()
		case <-m.quitChan:
			return
		}
	}
}

// removeExpired method removes expired elements.
func (m *timeExpiredMap[K, V]) removeExpired() {
	for key, val := range m.data {
		if val.expiredAt.Before(time.Now()) {
			delete(m.data, key)
		}
	}
}

/*
Time Expired List
*/

// TimeExpiredList is a list collection which values expires in time.
type TimeExpiredList[V any] interface {
	Add(value V)
	Get(index int) (V, error)
	GetAll() []V
	Del(i int) error
	Discard()
	Size() int
}

type timeExpiredList[V any] struct {
	sync.Mutex
	duration time.Duration
	data     []expiredElement[V]
	quitChan chan struct{}
}

// NewTimeExpiredList creates instance of TimeExpiredList interface. It runs goroutine for removing expired elements.
func NewTimeExpiredList[V any](duration time.Duration) TimeExpiredList[V] {
	tlist := &timeExpiredList[V]{
		duration: duration,
		data:     []expiredElement[V]{},
		quitChan: make(chan struct{}),
	}

	// Run goroutine for removing expired elements.
	go tlist.run()

	return tlist
}

// Add method add element to TimeExpiredList
func (l *timeExpiredList[V]) Add(value V) {
	l.Lock()
	defer l.Unlock()
	l.data = append(l.data, expiredElement[V]{expiredAt: time.Now().Add(l.duration), data: value})
}

// Get returns element by index.
func (l *timeExpiredList[V]) Get(i int) (V, error) {
	var result V
	if i < 0 && i >= len(l.data) {
		return result, ErrIndexOutOfBound
	}
	result = l.data[i].data
	return result, nil
}

// GetAll returns TimeExpiredElements values in slice.
func (l *timeExpiredList[V]) GetAll() []V {
	var result []V
	for _, v := range l.data {
		result = append(result, v.data)
	}
	return result
}

// Del removes element by index.
func (l *timeExpiredList[V]) Del(i int) error {
	if i < 0 && i >= len(l.data) {
		return ErrIndexOutOfBound
	}

	l.Lock()
	defer l.Unlock()
	l.data = append(l.data[:i], l.data[i+1:]...)
	return nil
}

// Size returns size of the list
func (l *timeExpiredList[V]) Size() int {
	return len(l.data)
}

// Discard method stops the goroutine for removing elements and discards data in internal slice.
func (l *timeExpiredList[V]) Discard() {
	l.quitChan <- struct{}{}
	l.Lock()
	defer l.Unlock()
	l.data = nil
}

// run method runs the goroutine for removing expired elements.
func (l *timeExpiredList[V]) run() {
	for {
		select {
		case <-time.After(1 * time.Second):
			l.removeExpired()
		case <-l.quitChan:
			return
		}
	}
}

// removeExpired method removes expired elements in list.
func (l *timeExpiredList[V]) removeExpired() {
	for i, val := range l.data {
		if val.expiredAt.Before(time.Now()) {
			_ = l.Del(i)
		}
	}
}
