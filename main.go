package gocollections

import (
	"errors"
	"sync"
	"time"
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

// Config struct is for configuration List or Map options.
type Config struct {
	// CleanJobInterval How often remove expired elements from collections. If it's too often, ex. 1 second and there
	// is too many elements, than it will cause performance issue.
	CleanJobInterval time.Duration
	// Size of expired element channel. If channel is full then last is removed before new is added.
	ExpiredElChanSize int
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
	Clear()
	Discard()
	Size() int
	ExpiredElChan() chan expiredElement[V]
}

type timeExpiredList[V any] struct {
	config        Config
	mu            sync.Mutex
	duration      time.Duration
	data          []expiredElement[V]
	dataString    []V
	expiredElChan chan expiredElement[V]
	quitChan      chan struct{}
}

// NewTimeExpiredList creates instance of TimeExpiredList interface. It runs goroutine for removing expired elements.
func NewTimeExpiredList[V any](duration time.Duration, configs ...Config) TimeExpiredList[V] {
	var config Config
	if len(configs) < 1 {
		// Default config if not provided
		config = Config{
			CleanJobInterval:  60 * time.Second,
			ExpiredElChanSize: 0,
		}
	} else {
		// Or use provided configuration
		config = configs[0]
	}

	tlist := &timeExpiredList[V]{
		config:        config,
		duration:      duration,
		data:          []expiredElement[V]{},
		dataString:    []V{},
		expiredElChan: make(chan expiredElement[V], config.ExpiredElChanSize),
		quitChan:      make(chan struct{}),
	}

	// Run goroutine for removing expired elements.
	go tlist.run()

	return tlist
}

// Add method add element to TimeExpiredList
func (l *timeExpiredList[V]) Add(value V) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.data = append(l.data, expiredElement[V]{expiredAt: time.Now().Add(l.duration), data: value})
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

// Clear method clears all elements from the list.
func (l *timeExpiredList[V]) Clear() {
	l.mu.Lock()
	l.mu.Unlock()
	l.data = []expiredElement[V]{}
}

// Discard method stops the goroutine for removing elements and discards data in internal slice.
func (l *timeExpiredList[V]) Discard() {
	l.quitChan <- struct{}{}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.data = nil
}

func (l *timeExpiredList[V]) ExpiredElChan() chan expiredElement[V] {
	return l.expiredElChan
}

// run method runs the goroutine for removing expired elements.
func (l *timeExpiredList[V]) run() {
	ticker := time.NewTicker(l.config.CleanJobInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.removeExpired()
		case <-l.quitChan:
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
			// If Element is not expired then add to new data slice.
			newData = append(newData, expiredElement[V]{data: val.data, expiredAt: val.expiredAt})
		} else {
			// If expired element channel is defined and size is bigger than 0, than send expired element to this channel.
			if cap(l.expiredElChan) > 0 {
				if len(l.expiredElChan) >= l.config.ExpiredElChanSize {
					// If expired element channel is full then remove first element.
					<-l.expiredElChan
				}
				// If Element is expired then add to expired channel.
				l.expiredElChan <- val
			}
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
	Clear()
	Discard()
	ExpiredElChan() chan expiredElement[V]
}

type timeExpiredMap[K comparable, V any] struct {
	config        Config
	mu            sync.Mutex
	duration      time.Duration           // default element duration
	data          map[K]expiredElement[V] // map of elements
	expiredElChan chan expiredElement[V]
	quitChan      chan struct{} // channel for indicating to end goroutines for removing expired elements
}

// NewTimeExpiredMap creates new TimeExpiredMap object.
func NewTimeExpiredMap[K comparable, V any](duration time.Duration, configs ...Config) TimeExpiredMap[K, V] {
	var config Config
	if len(configs) < 1 {
		// Default config if not provided
		config = Config{
			CleanJobInterval:  60 * time.Second,
			ExpiredElChanSize: 100,
		}
	} else {
		// Or use provided configuration
		config = configs[0]
	}

	tmap := &timeExpiredMap[K, V]{
		config:        config,
		duration:      duration,
		data:          make(map[K]expiredElement[V]),
		expiredElChan: make(chan expiredElement[V], config.ExpiredElChanSize),
		quitChan:      make(chan struct{}),
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

// Clear function clear all elements from map.
func (m *timeExpiredMap[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[K]expiredElement[V])
}

// Discard method stops the goroutine for removing elements and discards data in internal map.
func (m *timeExpiredMap[K, V]) Discard() {
	m.quitChan <- struct{}{}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = nil
}

func (m *timeExpiredMap[K, V]) ExpiredElChan() chan expiredElement[V] {
	return m.expiredElChan
}

// run method runs the goroutine for removing expired elements.
func (m *timeExpiredMap[K, V]) run() {
	ticker := time.NewTicker(m.config.CleanJobInterval)
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
	m.mu.Lock()
	defer m.mu.Unlock()
	for key, val := range m.data {
		if val.expiredAt.Before(time.Now()) {
			// If expired element channel is defined and size is bigger than 0, than send expired element to this channel.
			if cap(m.expiredElChan) > 0 {
				if len(m.expiredElChan) >= m.config.ExpiredElChanSize {
					// If expired element channel is full then remove first element.
					<-m.expiredElChan
				}
				// Send expired element to expired element channel.
				m.expiredElChan <- m.data[key]
			}
			// Delete element from map.
			delete(m.data, key)
		}
	}
}
