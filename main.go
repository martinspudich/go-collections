package gocollections

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	// CleanJobInterval How often remove expired elements from collections. If it's too often, ex. 1 second and there
	// is too many elements, than it will cause performance issue.
	CleanJobInterval = 60 * time.Second
)

var (
	ErrKeyNotFound     = errors.New("key not found")
	ErrIndexOutOfBound = errors.New("index out of bound")
	ErrExpired         = errors.New("element expired") // When an element is present in the collection but the validity time expires.
)

type expiredElement[V any] struct {
	data      V
	expiredAt time.Time
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
	mu         sync.Mutex
	duration   time.Duration
	data       []expiredElement[V]
	dataString []V
	quitChan   chan struct{}
}

// NewTimeExpiredList creates instance of TimeExpiredList interface. It runs goroutine for removing expired elements.
func NewTimeExpiredList[V any](duration time.Duration) TimeExpiredList[V] {
	tlist := &timeExpiredList[V]{
		duration:   duration,
		data:       []expiredElement[V]{},
		dataString: []V{},
		quitChan:   make(chan struct{}),
	}

	// Run goroutine for removing expired elements.
	go tlist.run()

	return tlist
}

// Add method add element to TimeExpiredList
func (l *timeExpiredList[V]) Add(value V) {
	l.mu.Lock()
	defer l.mu.Unlock()
	//l.dataString = append(l.dataString, value)
	l.data = append(l.data, expiredElement[V]{expiredAt: time.Now().Add(l.duration), data: value})
	//l.data = append(l.data, expiredElement[V]{})
}

// Get returns element by index.
func (l *timeExpiredList[V]) Get(i int) (V, error) {
	var result V
	l.mu.Lock()
	defer l.mu.Unlock()
	if i < 0 && i >= len(l.data) {
		return result, ErrIndexOutOfBound
	}
	if l.data[i].expiredAt.Before(time.Now()) {
		return result, ErrExpired
	}
	result = l.data[i].data
	return result, nil
}

// GetAll returns TimeExpiredElements values in slice.
func (l *timeExpiredList[V]) GetAll() []V {
	var result []V
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, v := range l.data {
		if v.expiredAt.Before(time.Now()) {
			// skip element if expired.
			continue
		}
		result = append(result, v.data)
	}
	return result
}

// Del removes element by index.
func (l *timeExpiredList[V]) Del(i int) error {
	if i < 0 || i >= len(l.data) {
		return ErrIndexOutOfBound
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	l.data = append(l.data[:i], l.data[i+1:]...)
	return nil
}

// Size returns size of the list
func (l *timeExpiredList[V]) Size() int {
	var count = 0
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, e := range l.data {
		// Don't count if element already expired.
		if e.expiredAt.After(time.Now()) {
			count++
		}
	}
	return count
}

// Discard method stops the goroutine for removing elements and discards data in internal slice.
func (l *timeExpiredList[V]) Discard() {
	l.quitChan <- struct{}{}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.data = nil
}

// run method runs the goroutine for removing expired elements.
func (l *timeExpiredList[V]) run() {
	ticker := time.NewTicker(CleanJobInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("running remove expired")
			l.removeExpired()
		case <-l.quitChan:
			fmt.Println("exiting run")
			return
		}
	}
}

// removeExpired method removes expired elements in list.
func (l *timeExpiredList[V]) removeExpired() {
	var newData []expiredElement[V]
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, val := range l.data {
		if val.expiredAt.After(time.Now()) {
			//newData = append(newData, val)
			newData = append(newData, expiredElement[V]{data: val.data, expiredAt: val.expiredAt})
		}
	}
	l.data = newData
}

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

type timeExpiredMap[K comparable, V any] struct {
	mu       sync.Mutex
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
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = expiredElement[V]{expiredAt: time.Now().Add(m.duration), data: data}
}

// AddWithDuration adds element to the map with key. It will set custom duration time of the element in the internal map.
func (m *timeExpiredMap[K, V]) AddWithDuration(key K, data V, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = expiredElement[V]{expiredAt: time.Now().Add(duration), data: data}
}

// Get method returns element by key.
func (m *timeExpiredMap[K, V]) Get(key K) (V, error) {
	var result V
	if !m.Contains(key) {
		return result, ErrKeyNotFound
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data[key].expiredAt.Before(time.Now()) {
		return result, ErrExpired
	}
	return m.data[key].data, nil
}

// Del method removes element from map.
func (m *timeExpiredMap[K, V]) Del(key K) error {
	if !m.Contains(key) {
		return ErrKeyNotFound
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}

// Contains method returns true if key is in the map. Else return false.
func (m *timeExpiredMap[K, V]) Contains(key K) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, found := m.data[key]
	if e.expiredAt.Before(time.Now()) {
		// if element expire, then return false
		return false
	}
	return found
}

// Size method returns size of the map.
func (m *timeExpiredMap[K, V]) Size() int {
	var count = 0
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, d := range m.data {
		// Don't count if element already expired.
		if d.expiredAt.After(time.Now()) {
			count++
		}
	}
	return count
}

// Discard method stops the goroutine for removing elements and discards data in internal map.
func (m *timeExpiredMap[K, V]) Discard() {
	m.quitChan <- struct{}{}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = nil
}

// run method runs the goroutine for removing expired elements.
func (m *timeExpiredMap[K, V]) run() {
	ticker := time.NewTicker(CleanJobInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
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
