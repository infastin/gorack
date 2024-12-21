package ttlmap

import (
	"context"
	"sync"
	"time"
)

type item[K comparable, V any] struct {
	key     K
	value   V
	expires time.Time
	bucket  uint8
}

type bucket[K comparable, V any] struct {
	items       map[K]*item[K, V]
	newestEntry time.Time
}

type Map[K comparable, V any] struct {
	items             map[K]*item[K, V]
	mu                sync.Mutex
	ttl               time.Duration
	buckets           []bucket[K, V]
	nextCleanupBucket uint8
}

// Creates a map with expirable items.
func New[K comparable, V any](ttl time.Duration, numBuckets uint8) *Map[K, V] {
	m := &Map[K, V]{
		items:             make(map[K]*item[K, V]),
		mu:                sync.Mutex{},
		ttl:               ttl,
		buckets:           make([]bucket[K, V], numBuckets),
		nextCleanupBucket: 0,
	}
	for i := 0; i < len(m.buckets); i++ {
		m.buckets[i].items = make(map[K]*item[K, V])
	}
	return m
}

// Puts a value to the map.
func (m *Map[K, V]) Put(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	if item, ok := m.items[key]; ok {
		m.removeFromBucket(item)
		item.value = value
		item.expires = now.Add(m.ttl)
		m.addToBucket(item)
		return
	}
	item := &item[K, V]{
		key:     key,
		value:   value,
		expires: now.Add(m.ttl),
	}
	m.items[key] = item
	m.addToBucket(item)
}

// Looks up a key's value from the map.
func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, ok := m.items[key]
	if !ok {
		return value, false
	}
	if time.Now().After(item.expires) {
		m.removeItem(item)
		return value, false
	}
	return item.value, true
}

func (m *Map[K, V]) Has(key K) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, ok := m.items[key]
	if !ok {
		return false
	}
	if time.Now().After(item.expires) {
		m.removeItem(item)
		return false
	}
	return true
}

// Updates existing element or inserts a new one using provided callback function.
func (m *Map[K, V]) Upsert(key K, cb func(exists bool, value V) V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	itm, ok := m.items[key]
	if !ok {
		itm = &item[K, V]{key: key}
		m.items[key] = itm
	} else {
		m.removeFromBucket(itm)
	}
	itm.value = cb(ok, itm.value)
	itm.expires = time.Now().Add(m.ttl)
	m.addToBucket(itm)
}

// Updates existing item or inserts a new one using provided callback function.
func (m *Map[K, V]) Update(key K, cb func(value V) V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, ok := m.items[key]
	if !ok {
		return false
	}
	m.removeFromBucket(item)
	item.value = cb(item.value)
	item.expires = time.Now().Add(m.ttl)
	m.addToBucket(item)
	return true
}

// Looks up a key's value from the map and removes it.
func (m *Map[K, V]) GetAndRemove(key K) (value V, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, ok := m.items[key]
	if !ok {
		return value, false
	}
	m.removeItem(item)
	if time.Now().After(item.expires) {
		return value, false
	}
	return item.value, true
}

func (m *Map[K, V]) Remove(key K) bool {
	_, ok := m.GetAndRemove(key)
	return ok
}

// Starts a cleanup loop.
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
			for _, item := range m.buckets[idx].items {
				m.removeItem(item)
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

// Removes a given item from the map.
// NOTE: Has to be called with lock!
func (m *Map[K, V]) removeItem(item *item[K, V]) {
	delete(m.items, item.key)
	m.removeFromBucket(item)
}

// Adds item to expire bucket so that it will be cleaned up when the time comes.
// NOTE: Has to be called with lock!
func (m *Map[K, V]) addToBucket(item *item[K, V]) {
	numBuckets := uint8(len(m.buckets))
	idx := (numBuckets + m.nextCleanupBucket - 1) % numBuckets
	item.bucket = idx
	m.buckets[idx].items[item.key] = item
	if m.buckets[idx].newestEntry.Before(item.expires) {
		m.buckets[idx].newestEntry = item.expires
	}
}

// Removes item from its corresponding expire bucket.
// NOTE: Has to be called with lock!
func (m *Map[K, V]) removeFromBucket(item *item[K, V]) {
	delete(m.buckets[item.bucket].items, item.key)
}
