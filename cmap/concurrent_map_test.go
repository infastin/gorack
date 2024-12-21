package cmap

import (
	"encoding/json"
	"sort"
	"strconv"
	"testing"
)

type Animal struct {
	name string
}

func TestMapCreation(t *testing.T) {
	m := New[string]()
	if m.shards == nil {
		t.Error("map is null.")
	}
	if m.Count() != 0 {
		t.Error("new map should be empty.")
	}
}

func TestInsert(t *testing.T) {
	m := New[Animal]()

	elephant := Animal{"elephant"}
	m.Set("elephant", elephant)

	monkey := Animal{"monkey"}
	m.Set("monkey", monkey)

	if m.Count() != 2 {
		t.Error("map should contain exactly two elements.")
	}
}

func TestInsertAbsent(t *testing.T) {
	m := New[Animal]()

	elephant := Animal{"elephant"}
	m.SetIfAbsent("elephant", elephant)

	monkey := Animal{"monkey"}
	if ok := m.SetIfAbsent("elephant", monkey); ok {
		t.Error("map set a new value even the entry is already present")
	}
}

func TestGet(t *testing.T) {
	m := New[Animal]()

	// Get a missing element.
	val, ok := m.Get("Money")
	if ok {
		t.Error("ok should be false when item is missing from map.")
	}
	if val != (Animal{}) {
		t.Error("Missing values should return as null.")
	}

	elephant := Animal{"elephant"}
	m.Set("elephant", elephant)

	// Retrieve inserted element.
	elephant, ok = m.Get("elephant")
	if !ok {
		t.Error("ok should be true for item stored within the map.")
	}
	if elephant.name != "elephant" {
		t.Error("item was modified.")
	}
}

func TestHas(t *testing.T) {
	m := New[Animal]()

	// Get a missing element.
	if m.Has("Money") {
		t.Error("element shouldn't exists")
	}

	elephant := Animal{"elephant"}
	m.Set("elephant", elephant)

	if !m.Has("elephant") {
		t.Error("element exists, expecting Has to return True.")
	}
}

func TestRemove(t *testing.T) {
	m := New[Animal]()

	monkey := Animal{"monkey"}
	m.Set("monkey", monkey)
	m.Remove("monkey")

	if m.Count() != 0 {
		t.Error("Expecting count to be zero once item was removed.")
	}

	temp, ok := m.Get("monkey")
	if ok != false {
		t.Error("Expecting ok to be false for missing items.")
	}
	if temp != (Animal{}) {
		t.Error("Expecting item to be nil after its removal.")
	}

	// Remove a none existing element.
	m.Remove("noone")
}

func TestRemoveCb(t *testing.T) {
	m := New[Animal]()

	monkey := Animal{"monkey"}
	m.Set("monkey", monkey)

	elephant := Animal{"elephant"}
	m.Set("elephant", elephant)

	var (
		mapKey   string
		mapVal   Animal
		wasFound bool
	)
	cb := func(key string, val Animal, exists bool) bool {
		mapKey = key
		mapVal = val
		wasFound = exists
		return val.name == "monkey"
	}

	// Monkey should be removed
	result := m.RemoveCb("monkey", cb)
	if !result {
		t.Errorf("Result was not true")
	}
	if mapKey != "monkey" {
		t.Error("Wrong key was provided to the callback")
	}
	if mapVal != monkey {
		t.Errorf("Wrong value was provided to the value")
	}
	if !wasFound {
		t.Errorf("Key was not found")
	}
	if m.Has("monkey") {
		t.Errorf("Key was not removed")
	}

	// Elephant should not be removed
	result = m.RemoveCb("elephant", cb)
	if result {
		t.Errorf("Result was true")
	}
	if mapKey != "elephant" {
		t.Error("Wrong key was provided to the callback")
	}
	if mapVal != elephant {
		t.Errorf("Wrong value was provided to the value")
	}
	if !wasFound {
		t.Errorf("Key was not found")
	}
	if !m.Has("elephant") {
		t.Errorf("Key was removed")
	}

	// Unset key should remain unset
	result = m.RemoveCb("horse", cb)
	if result {
		t.Errorf("Result was true")
	}
	if mapKey != "horse" {
		t.Error("Wrong key was provided to the callback")
	}
	if mapVal != (Animal{}) {
		t.Errorf("Wrong value was provided to the value")
	}
	if wasFound {
		t.Errorf("Key was found")
	}
	if m.Has("horse") {
		t.Errorf("Key was created")
	}
}

func TestPop(t *testing.T) {
	m := New[Animal]()

	monkey := Animal{"monkey"}
	m.Set("monkey", monkey)

	v, exists := m.Pop("monkey")
	if !exists || v != monkey {
		t.Error("Pop didn't find a monkey.")
	}

	v2, exists2 := m.Pop("monkey")
	if exists2 || v2 == monkey {
		t.Error("Pop keeps finding monkey")
	}

	if m.Count() != 0 {
		t.Error("Expecting count to be zero once item was Pop'ed.")
	}

	temp, ok := m.Get("monkey")
	if ok != false {
		t.Error("Expecting ok to be false for missing items.")
	}
	if temp != (Animal{}) {
		t.Error("Expecting item to be nil after its removal.")
	}
}

func TestCount(t *testing.T) {
	m := New[Animal]()
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	if m.Count() != 100 {
		t.Error("Expecting 100 element within map.")
	}
}

func TestIsEmpty(t *testing.T) {
	m := New[Animal]()
	if !m.IsEmpty() {
		t.Error("new map should be empty")
	}
	m.Set("elephant", Animal{"elephant"})
	if m.IsEmpty() {
		t.Error("map shouldn't be empty.")
	}
}

func TestSeq(t *testing.T) {
	m := New[Animal]()
	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	counter := 0
	// Iterate over elements.
	for _, val := range m.Seq() {
		if val == (Animal{}) {
			t.Error("Expecting an object.")
		}
		counter++
	}
	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestClear(t *testing.T) {
	m := New[Animal]()
	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	m.Clear()
	if m.Count() != 0 {
		t.Error("We should have 0 elements.")
	}
}

func TestIter(t *testing.T) {
	m := New[Animal]()
	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	counter := 0
	// Iterate over elements.
	m.Iter(func(key string, v Animal) bool {
		counter++
		return true
	})
	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestItems(t *testing.T) {
	m := New[Animal]()
	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	items := m.Items()
	if len(items) != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestConcurrent(t *testing.T) {
	m := New[int]()
	ch := make(chan int)
	const iterations = 1000
	var a [iterations]int
	// Using go routines insert 1000 ints into our map.
	go func() {
		for i := 0; i < iterations/2; i++ {
			// Add item to map.
			m.Set(strconv.Itoa(i), i)
			// Retrieve item from map.
			val, _ := m.Get(strconv.Itoa(i))
			// Write to channel inserted value.
			ch <- val
		} // Call go routine with current index.
	}()
	go func() {
		for i := iterations / 2; i < iterations; i++ {
			// Add item to map.
			m.Set(strconv.Itoa(i), i)
			// Retrieve item from map.
			val, _ := m.Get(strconv.Itoa(i))
			// Write to channel inserted value.
			ch <- val
		} // Call go routine with current index.
	}()
	// Wait for all go routines to finish.
	counter := 0
	for elem := range ch {
		a[counter] = elem
		counter++
		if counter == iterations {
			break
		}
	}
	// Sorts array, will make is simpler to verify all inserted values we're returned.
	sort.Ints(a[0:iterations])
	// Make sure map contains 1000 elements.
	if m.Count() != iterations {
		t.Error("Expecting 1000 elements.")
	}
	// Make sure all inserted values we're fetched from map.
	for i := 0; i < iterations; i++ {
		if i != a[i] {
			t.Error("missing value", i)
		}
	}
}

func TestJsonMarshal(t *testing.T) {
	expected := "{\"a\":1,\"b\":2}"

	m := New[int](WithShardCount[string](2))
	m.Set("a", 1)
	m.Set("b", 2)

	j, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
	}
	if string(j) != expected {
		t.Error("json", string(j), "differ from expected", expected)
	}
}

func TestKeys(t *testing.T) {
	m := New[Animal]()
	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	if len(m.Keys()) != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestUpsert(t *testing.T) {
	dolphin := Animal{"dolphin"}
	whale := Animal{"whale"}
	tiger := Animal{"tiger"}
	lion := Animal{"lion"}

	cb := func(exists bool, valueInMap Animal, newValue Animal) Animal {
		if !exists {
			return newValue
		}
		valueInMap.name += newValue.name
		return valueInMap
	}

	m := New[Animal]()
	m.Set("marine", dolphin)
	m.Upsert("marine", whale, cb)
	m.Upsert("predator", tiger, cb)
	m.Upsert("predator", lion, cb)

	if m.Count() != 2 {
		t.Error("map should contain exactly two elements.")
	}

	marineAnimals, ok := m.Get("marine")
	if marineAnimals.name != "dolphinwhale" || !ok {
		t.Error("Set, then Upsert failed")
	}

	predators, ok := m.Get("predator")
	if !ok || predators.name != "tigerlion" {
		t.Error("Upsert, then Upsert failed")
	}
}

func TestKeysWhenRemoving(t *testing.T) {
	m := New[Animal]()
	// Insert 100 elements.
	const total = 100
	for i := 0; i < total; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	// Remove 10 elements concurrently.
	num := 10
	for i := 0; i < num; i++ {
		go func(c *ConcurrentMap[string, Animal], n int) {
			c.Remove(strconv.Itoa(n))
		}(&m, i)
	}
	for _, k := range m.Keys() {
		if k == "" {
			t.Error("Empty keys returned")
		}
	}
}