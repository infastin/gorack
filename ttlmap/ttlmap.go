package ttlmap

import (
	"context"
	"sync"
	"time"
)

type element[K comparable, V any] struct {
	key     K
	value   V
	expires time.Time
	bucket  uint8
}

type bucket[K comparable, V any] struct {
	elems       map[K]*element[K, V]
	newestEntry time.Time
}

// Map is a map with expirable elements.
type Map[K comparable, V any] struct {
	elems             map[K]*element[K, V]
	mu                sync.Mutex
	ttl               time.Duration
	buckets           []bucket[K, V]
	nextCleanupBucket uint8
}

// New creates a new map with expirable elements.
// Whenever a key-value pair is inserted, a cleanup bucket is chosed
// and the inserted key-value pair is put into the bucket.
// With the interval of `ttl / numBuckets` a cleanup bucket is chosed and
// every expired element in that bucket is removed from the map (and from the bucket).
func New[K comparable, V any](ttl time.Duration, numBuckets uint8) *Map[K, V] {
	m := &Map[K, V]{
		elems:             make(map[K]*element[K, V]),
		mu:                sync.Mutex{},
		ttl:               ttl,
		buckets:           make([]bucket[K, V], numBuckets),
		nextCleanupBucket: 0,
	}
	for i := 0; i < len(m.buckets); i++ {
		m.buckets[i].elems = make(map[K]*element[K, V])
	}
	return m
}

// Put puts an element into the map.
func (m *Map[K, V]) Put(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	if elem, ok := m.elems[key]; ok {
		m.removeFromBucket(elem)
		elem.value = value
		elem.expires = now.Add(m.ttl)
		m.addToBucket(elem)
		return
	}
	elem := &element[K, V]{
		key:     key,
		value:   value,
		expires: now.Add(m.ttl),
	}
	m.elems[key] = elem
	m.addToBucket(elem)
}

// Get retrieves an element from the map under the specified key.
func (m *Map[K, V]) Get(key K) (value V, exists bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, exists := m.elems[key]
	if !exists {
		return value, false
	}
	if time.Now().After(elem.expires) {
		m.removeElement(elem)
		return value, false
	}
	return elem.value, true
}

// Has checks if an element under the specified key exists.
func (m *Map[K, V]) Has(key K) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, ok := m.elems[key]
	if !ok {
		return false
	}
	if time.Now().After(elem.expires) {
		m.removeElement(elem)
		return false
	}
	return true
}

// Upsert updates an existing element or inserts a new one using provided callback function.
func (m *Map[K, V]) Upsert(key K, cb func(exists bool, value V) V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	itm, ok := m.elems[key]
	if !ok {
		itm = &element[K, V]{key: key}
		m.elems[key] = itm
	} else {
		m.removeFromBucket(itm)
	}
	itm.value = cb(ok, itm.value)
	itm.expires = time.Now().Add(m.ttl)
	m.addToBucket(itm)
}

// Update updates an existing element using provided callback function.
func (m *Map[K, V]) Update(key K, cb func(value V) V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, ok := m.elems[key]
	if !ok {
		return false
	}
	m.removeFromBucket(elem)
	elem.value = cb(elem.value)
	elem.expires = time.Now().Add(m.ttl)
	m.addToBucket(elem)
	return true
}

// GetAndRemove removes an element from the map and returns it.
func (m *Map[K, V]) GetAndRemove(key K) (value V, exists bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, exists := m.elems[key]
	if !exists {
		return value, false
	}
	m.removeElement(elem)
	if time.Now().After(elem.expires) {
		return value, false
	}
	return elem.value, true
}

// Remove removes an element from the map.
func (m *Map[K, V]) Remove(key K) bool {
	_, ok := m.GetAndRemove(key)
	return ok
}

// Start starts a cleanup loop.
func (m *Map[K, V]) Start(ctx context.Context) (stop func()) {
	ctx, cancel := context.WithCancel(ctx)
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(m.ttl / time.Duration(len(m.buckets))):
				m.mu.Lock()
			}
			idx := m.nextCleanupBucket
			timeToExpire := time.Until(m.buckets[idx].newestEntry)
			if timeToExpire > 0 {
				m.mu.Unlock()
				select {
				case <-ctx.Done():
					return
				case <-time.After(timeToExpire):
					m.mu.Lock()
				}
			}
			for _, elem := range m.buckets[idx].elems {
				m.removeElement(elem)
			}
			m.nextCleanupBucket = (m.nextCleanupBucket + 1) % uint8(len(m.buckets))
			m.mu.Unlock()
		}
	}()
	return func() {
		cancel()
		<-doneCh
	}
}

// removeElement removes an element from the map.
// WARN: has to be called with lock!
func (m *Map[K, V]) removeElement(elem *element[K, V]) {
	delete(m.elems, elem.key)
	m.removeFromBucket(elem)
}

// addToBucket adds an element to the expire bucket so that it will be cleaned up when the time comes.
// WARN: has to be called with lock!
func (m *Map[K, V]) addToBucket(elem *element[K, V]) {
	numBuckets := uint8(len(m.buckets))
	idx := (numBuckets + m.nextCleanupBucket - 1) % numBuckets
	elem.bucket = idx
	m.buckets[idx].elems[elem.key] = elem
	if m.buckets[idx].newestEntry.Before(elem.expires) {
		m.buckets[idx].newestEntry = elem.expires
	}
}

// removeFromBucket removes an element from its corresponding expire bucket.
// WARN: has to be called with lock!
func (m *Map[K, V]) removeFromBucket(elem *element[K, V]) {
	delete(m.buckets[elem.bucket].elems, elem.key)
}
