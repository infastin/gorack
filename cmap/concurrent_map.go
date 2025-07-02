package cmap

import (
	"encoding/json"
	"fmt"
	"hash/maphash"
	"iter"
	"maps"
	"reflect"
	"sync"
	"unsafe"
)

type options struct {
	shardCount   int
	shardingFunc any
}

// ShardingFunc is a function for sharding a map.
type ShardingFunc[K comparable] func(key K) uint64

// Option is used to configure concurrent map.
type Option func(*options)

// WithShardCount allows to set the number of shards in a map.
func WithShardCount(n int) Option {
	return func(o *options) {
		o.shardCount = n
	}
}

// WithShardingFunc allows to set the sharding function of a map.
func WithShardingFunc[K comparable](fn ShardingFunc[K]) Option {
	return func(o *options) {
		o.shardingFunc = fn
	}
}

// ConcurrentMap is a thread-safe map.
// To avoid lock bottlenecks this map is dived to several map shards.
type ConcurrentMap[K comparable, V any] struct {
	shards   []*shard[K, V]
	sharding ShardingFunc[K]
}

type shard[K comparable, V any] struct {
	items map[K]V
	mu    sync.RWMutex
}

// New creates a new concurrent map.
func New[K comparable, V any](opts ...Option) ConcurrentMap[K, V] {
	seed := maphash.MakeSeed()
	var shardingFunc ShardingFunc[K]
	if reflect.TypeFor[K]().Kind() == reflect.String {
		shardingFunc = func(key K) uint64 {
			return maphash.String(seed, *(*string)(unsafe.Pointer(&key)))
		}
	} else {
		shardingFunc = func(key K) uint64 {
			return maphash.Comparable(seed, key)
		}
	}
	options := options{
		shardCount:   32,
		shardingFunc: shardingFunc,
	}
	for _, opt := range opts {
		opt(&options)
	}
	shardingFunc, ok := options.shardingFunc.(ShardingFunc[K])
	if !ok {
		panic(fmt.Sprintf("cmap: invalid sharding function: expected %v, got %T",
			reflect.TypeFor[ShardingFunc[K]](), options.shardingFunc))
	}
	m := ConcurrentMap[K, V]{
		shards:   make([]*shard[K, V], options.shardCount),
		sharding: shardingFunc,
	}
	for i := range m.shards {
		m.shards[i] = &shard[K, V]{
			items: make(map[K]V),
			mu:    sync.RWMutex{},
		}
	}
	return m
}

// getShard returns shard under the specified key.
func (m ConcurrentMap[K, V]) getShard(key K) *shard[K, V] {
	return m.shards[uint(m.sharding(key))%uint(len(m.shards))]
}

// Set sets the given value under the specified key.
func (m ConcurrentMap[K, V]) Set(key K, value V) {
	// Get map shard.
	shard := m.getShard(key)
	shard.mu.Lock()
	shard.items[key] = value
	shard.mu.Unlock()
}

// UpsertCb is a callback to return a new element to be inserted into the map.
// It is called while lock is held, therefore it MUST NOT
// try to access other keys in the same map, as it can lead to deadlock.
type UpsertCb[V any] func(exist bool, valueInMap, newValue V) V

// Upsert updates an existing element or inserts a new one using UpsertCb.
// Returns the updated/inserted element.
func (m ConcurrentMap[K, V]) Upsert(key K, value V, cb UpsertCb[V]) (res V) {
	shard := m.getShard(key)
	shard.mu.Lock()
	v, ok := shard.items[key]
	res = cb(ok, v, value)
	shard.items[key] = res
	shard.mu.Unlock()
	return res
}

// UpdateCb is a callback to update an element in the map.
// It is called while lock is held, therefore it MUST NOT
// try to access other keys in same map, as it can lead to deadlock.
type UpdateCb[V any] func(valueInMap, newValue V) V

// Update updates an existing element using UpdateCb.
// If the element doesn't exist, returns false.
// Otherwise returns the updated element and true.
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

// SetIfAbsent sets the given value under the specified key
// if no value was associated with it.
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

// Get retrieves an element from the map under the specified key.
func (m ConcurrentMap[K, V]) Get(key K) (V, bool) {
	// Get shard
	shard := m.getShard(key)
	shard.mu.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.mu.RUnlock()
	return val, ok
}

// Count returns the number of elements within the map.
func (m ConcurrentMap[K, V]) Count() int {
	count := 0
	for _, shard := range m.shards {
		shard.mu.RLock()
		count += len(shard.items)
		shard.mu.RUnlock()
	}
	return count
}

// Has checks if an item under the specified key exists.
func (m ConcurrentMap[K, V]) Has(key K) bool {
	// Get shard.
	shard := m.getShard(key)
	shard.mu.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.mu.RUnlock()
	return ok
}

// Remove removes an element from the map.
func (m ConcurrentMap[K, V]) Remove(key K) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.mu.Lock()
	delete(shard.items, key)
	shard.mu.Unlock()
}

// RemoveCb is a callback to remove an element from the map.
// It is called while lock is held.
// If it returns true, the element will be removed from the map.
type RemoveCb[K any, V any] func(key K, v V, exists bool) bool

// RemoveCb removes an element from the map using cb.
// Returns the value returned by cb.
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

// RemoveFunc is a callback to remove elements in the map.
// Lock is held for all calls for a given shard
// therefore callback sees consistent view of a shard,
// but not across the shards.
// If it returns true, the element will be removed from the map.
type RemoveFunc[K any, V any] func(key K, v V) bool

// RemoveFunc removes any element from the map for which fn returns true.
func (m ConcurrentMap[K, V]) RemoveFunc(fn RemoveFunc[K, V]) {
	for _, shard := range m.shards {
		shard.mu.Lock()
		for key, value := range shard.items {
			if fn(key, value) {
				delete(shard.items, key)
			}
		}
		shard.mu.Unlock()
	}
}

// Pop removes an element from the map and returns it.
func (m ConcurrentMap[K, V]) Pop(key K) (v V, exists bool) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.mu.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.mu.Unlock()
	return v, exists
}

// IsEmpty checks if the map is empty.
func (m ConcurrentMap[K, V]) IsEmpty() bool {
	return m.Count() == 0
}

// IterCb is an iterator callback called for every element in the map.
// RLock is held for all calls for a given shard
// therefore callback sees consistent view of a shard,
// but not across the shards.
type IterCb[K comparable, V any] func(key K, v V) bool

// Iter is a callback based iterator, cheapest way to read all elements in a map.
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

// Seq is handy go1.23 iterator based on Iter() method.
func (m ConcurrentMap[K, V]) Seq() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.Iter(yield)
	}
}

// Clear removes all items from the map.
func (m ConcurrentMap[K, V]) Clear() {
	for _, shard := range m.shards {
		shard.mu.Lock()
		clear(shard.items)
		shard.mu.Unlock()
	}
}

// Items returns all items in the map.
func (m ConcurrentMap[K, V]) Items() map[K]V {
	return maps.Collect(m.Seq())
}

// Keys returns all keys in the map.
func (m ConcurrentMap[K, V]) Keys() []K {
	keys := make([]K, 0)
	for key := range m.Seq() {
		keys = append(keys, key)
	}
	return keys
}

// MarshalJSON encodes the map into a json object.
func (m ConcurrentMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Items())
}

// UnmarshalJSON decodes a json object into the map.
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
