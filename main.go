package gocollections

import (
	"errors"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")

// TimeExpiredMap implementation of this interface is a map with elements which are expire base on expiration duration.
// Implementation of this map is running goroutine which removes expired element. To stop this goroutine call Discard()
// method when this map is not needed any more.
type TimeExpiredMap interface {
	Add(key string, object any)
	AddWithDuration(key string, data any, duration time.Duration)
	Get(key string) (any, error)
	Del(key string) error
	Contains(key string) bool
	Size() int
	Discard()
}

type expiredElement struct {
	data      any
	expiredAt time.Time
}

type timeExpiredMap struct {
	duration time.Duration             // default element duration
	data     map[string]expiredElement // map of elements
	quitChan chan struct{}             // channel for indicating to end goroutines for removing expired elements
}

// NewTimeExpiredMap creates new TimeExpiredMap object.
func NewTimeExpiredMap(duration time.Duration) TimeExpiredMap {
	timeExpiredMap := &timeExpiredMap{
		duration: duration,
		data:     make(map[string]expiredElement),
		quitChan: make(chan struct{}),
	}

	go timeExpiredMap.run()

	return timeExpiredMap
}

// Add method adds element to the map with key.
func (m *timeExpiredMap) Add(key string, data any) {
	m.data[key] = expiredElement{expiredAt: time.Now().Add(m.duration), data: data}
}

// AddWithDuration adds element to the map with key. It will set custom duration time of the element in the internal map.
func (m *timeExpiredMap) AddWithDuration(key string, data any, duration time.Duration) {
	m.data[key] = expiredElement{expiredAt: time.Now().Add(duration), data: data}
}

// Get method returns element by key.
func (m *timeExpiredMap) Get(key string) (any, error) {
	if !m.Contains(key) {
		return nil, ErrKeyNotFound
	}
	return m.data[key].data, nil
}

// Del method removes element from map.
func (m *timeExpiredMap) Del(key string) error {
	if !m.Contains(key) {
		return ErrKeyNotFound
	}
	delete(m.data, key)
	return nil
}

// Contains method returns true if key is in the map. Else return false.
func (m *timeExpiredMap) Contains(key string) bool {
	_, found := m.data[key]
	return found
}

// Size method returns size of the map.
func (m *timeExpiredMap) Size() int {
	return len(m.data)
}

// Discard method stops the goroutine for removing elements and discards data in internal map.
func (m *timeExpiredMap) Discard() {
	m.quitChan <- struct{}{}
	m.data = nil
}

// run method runs the goroutine for removing expired elements.
func (m *timeExpiredMap) run() {
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
func (m *timeExpiredMap) removeExpired() {
	for key, val := range m.data {
		if val.expiredAt.Before(time.Now()) {
			delete(m.data, key)
		}
	}
}
