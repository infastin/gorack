package cmap

import (
	"encoding/json"
	"hash/maphash"
	"iter"
	"sync"
)

type options[K comparable] struct {
	shardCount   int
	shardingFunc func(key K) uint64
}

// Used to configure concurrent map.
type Option[K comparable] func(*options[K])

// Allows to set the number of shards in a map.
func WithShardCount[K comparable](n int) Option[K] {
	return func(o *options[K]) {
		o.shardCount = n
	}
}

// A thread-safe map.
// To avoid lock bottlenecks this map is dived to several map shards.
type ConcurrentMap[K comparable, V any] struct {
	shards   []*shard[K, V]
	sharding func(key K) uint64
}

type shard[K comparable, V any] struct {
	items map[K]V
	mu    sync.RWMutex
}

// Creates a new concurrent map.
func New[V any](opts ...Option[string]) ConcurrentMap[string, V] {
	seed := maphash.MakeSeed()
	sharding := func(key string) uint64 {
		return maphash.String(seed, key)
	}
	options := options[string]{
		shardCount:   32,
		shardingFunc: sharding,
	}
	for _, opt := range opts {
		opt(&options)
	}
	m := ConcurrentMap[string, V]{
		shards:   make([]*shard[string, V], options.shardCount),
		sharding: options.shardingFunc,
	}
	for i := range m.shards {
		m.shards[i] = &shard[string, V]{
			items: make(map[string]V),
			mu:    sync.RWMutex{},
		}
	}
	return m
}

// Returns shard under given key.
func (m ConcurrentMap[K, V]) getShard(key K) *shard[K, V] {
	return m.shards[uint(m.sharding(key))%uint(len(m.shards))]
}

// Sets the given value under the specified key.
func (m ConcurrentMap[K, V]) Set(key K, value V) {
	// Get map shard.
	shard := m.getShard(key)
	shard.mu.Lock()
	shard.items[key] = value
	shard.mu.Unlock()
}

// Callback to return new element to be inserted into the map.
// It is called while lock is held, therefore it MUST NOT
// try to access other keys in same map, as it can lead to deadlock since
// Go sync.RWLock is not reentrant.
type UpsertCb[V any] func(exist bool, valueInMap, newValue V) V

// Updates existing element or inserts a new one using UpsertCb.
// Returns updated/inserted element.
func (m ConcurrentMap[K, V]) Upsert(key K, value V, cb UpsertCb[V]) (res V) {
	shard := m.getShard(key)
	shard.mu.Lock()
	v, ok := shard.items[key]
	res = cb(ok, v, value)
	shard.items[key] = res
	shard.mu.Unlock()
	return res
}

// Callback to update an element in the map.
// It is called while lock is held, therefore it MUST NOT
// try to access other keys in same map, as it can lead to deadlock since
// Go sync.RWLock is not reentrant.
type UpdateCb[V any] func(valueInMap, newValue V) V

// Updates an existing element using UpdateCb.
// If the element doesn't exist, returns false.
// Otherwise returns updated element and true.
func (m ConcurrentMap[K, V]) Update(key K, value V, cb UpdateCb[V]) (res V, updated bool) {
	shard := m.getShard(key)
	shard.mu.Lock()
	v, ok := shard.items[key]
	if !ok {
		shard.mu.Unlock()
		return res, false
	}
	res = cb(v, value)
	shard.items[key] = res
	shard.mu.Unlock()
	return res, true
}

// Sets the given value under the specified key if no value was associated with it.
func (m ConcurrentMap[K, V]) SetIfAbsent(key K, value V) bool {
	// Get map shard.
	shard := m.getShard(key)
	shard.mu.Lock()
	_, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}
	shard.mu.Unlock()
	return !ok
}

// Retrieves an element from map under given key.
func (m ConcurrentMap[K, V]) Get(key K) (V, bool) {
	// Get shard
	shard := m.getShard(key)
	shard.mu.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.mu.RUnlock()
	return val, ok
}

// Returns the number of elements within the map.
func (m ConcurrentMap[K, V]) Count() int {
	count := 0
	for _, shard := range m.shards {
		shard.mu.RLock()
		count += len(shard.items)
		shard.mu.RUnlock()
	}
	return count
}

// Looks up an item under specified key.
func (m ConcurrentMap[K, V]) Has(key K) bool {
	// Get shard.
	shard := m.getShard(key)
	shard.mu.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.mu.RUnlock()
	return ok
}

// Removes an element from the map.
func (m ConcurrentMap[K, V]) Remove(key K) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.mu.Lock()
	delete(shard.items, key)
	shard.mu.Unlock()
}

// RemoveCb is a callback executed in a map.RemoveCb() call, while Lock is held
// If returns true, the element will be removed from the map
type RemoveCb[K any, V any] func(key K, v V, exists bool) bool

// Locks the shard containing the key, retrieves its current value and calls the callback with those params.
// If callback returns true and element exists, it will remove it from the map.
// Returns the value returned by the callback (even if element was not present in the map).
func (m ConcurrentMap[K, V]) RemoveCb(key K, cb RemoveCb[K, V]) bool {
	// Try to get shard.
	shard := m.getShard(key)
	shard.mu.Lock()
	v, ok := shard.items[key]
	remove := cb(key, v, ok)
	if remove && ok {
		delete(shard.items, key)
	}
	shard.mu.Unlock()
	return remove
}

// Removes an element from the map and returns it.
func (m ConcurrentMap[K, V]) Pop(key K) (v V, exists bool) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.mu.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.mu.Unlock()
	return v, exists
}

// Checks if map is empty.
func (m ConcurrentMap[K, V]) IsEmpty() bool {
	return m.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel.
type Tuple[K comparable, V any] struct {
	Key K
	Val V
}

// Iterator callback called for every (key, value) pair found in maps.
// RLock is held for all calls for a given shard
// therefore callback sees consistent view of a shard,
// but not across the shards.
type IterCb[K comparable, V any] func(key K, v V) bool

// Callback based iterator, cheapest way to read all elements in a map.
func (m ConcurrentMap[K, V]) Iter(fn IterCb[K, V]) {
	for _, shard := range m.shards {
		shard.mu.RLock()
		for key, value := range shard.items {
			if !fn(key, value) {
				shard.mu.RUnlock()
				return
			}
		}
		shard.mu.RUnlock()
	}
}

// Handy go1.23 iterator based on Iter() method.
func (m ConcurrentMap[K, V]) Seq() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.Iter(yield)
	}
}

// Removes all items from map.
func (m ConcurrentMap[K, V]) Clear() {
	for _, shard := range m.shards {
		shard.mu.Lock()
		clear(shard.items)
		shard.mu.Unlock()
	}
}

// Returns all items as map[string]V.
func (m ConcurrentMap[K, V]) Items() map[K]V {
	items := make(map[K]V)
	for key, value := range m.Seq() {
		items[key] = value
	}
	return items
}

// Returns all keys as []string.
func (m ConcurrentMap[K, V]) Keys() []K {
	keys := make([]K, 0)
	for key := range m.Seq() {
		keys = append(keys, key)
	}
	return keys
}

// Reviles ConcurrentMap "private" variables to json marshal.
func (m ConcurrentMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Items())
}

// Reverse process of Marshal.
func (m *ConcurrentMap[K, V]) UnmarshalJSON(b []byte) (err error) {
	tmp := make(map[K]V)
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	for key, value := range tmp {
		m.Set(key, value)
	}
	return nil
}
