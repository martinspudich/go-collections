package gocollections

import (
	"errors"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")

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
	m.data[key] = expiredElement[V]{expiredAt: time.Now().Add(m.duration), data: data}
}

// AddWithDuration adds element to the map with key. It will set custom duration time of the element in the internal map.
func (m *timeExpiredMap[K, V]) AddWithDuration(key K, data V, duration time.Duration) {
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
