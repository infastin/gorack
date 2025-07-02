package cmap

import (
	"iter"
	"slices"
	"sort"
	"strconv"
	"testing"
)

type Animal struct {
	name string
}

func TestMapCreation(t *testing.T) {
	m := New[string, string]()
	if m.shards == nil {
		t.Error("map is nil")
	}
	if m.Count() != 0 {
		t.Error("new map must be empty")
	}
}

func TestInsert(t *testing.T) {
	m := New[string, Animal]()

	elephant := Animal{"elephant"}
	m.Set("elephant", elephant)

	monkey := Animal{"monkey"}
	m.Set("monkey", monkey)

	if m.Count() != 2 {
		t.Error("map must contain exactly two elements")
	}
}

func TestInsertAbsent(t *testing.T) {
	m := New[string, Animal]()

	elephant := Animal{"elephant"}
	m.SetIfAbsent("elephant", elephant)

	monkey := Animal{"monkey"}
	if ok := m.SetIfAbsent("elephant", monkey); ok {
		t.Error("new value has been set, but the element is already present")
	}
	if val, _ := m.Get("elephant"); val != elephant {
		t.Error("element has been modified after second SetIfAbsent")
	}
}

func TestGet(t *testing.T) {
	m := New[string, Animal]()

	// Get a missing element.
	val, ok := m.Get("Money")
	if ok {
		t.Error("ok must be false when an element is missing from the map")
	}
	if val != (Animal{}) {
		t.Error("missing values must return as default")
	}

	elephant := Animal{"elephant"}
	m.Set("elephant", elephant)

	// Retrieve inserted element.
	elephant, ok = m.Get("elephant")
	if !ok {
		t.Error("ok must be true for an element stored within the map")
	}
	if elephant.name != "elephant" {
		t.Error("element has been modified after Get")
	}
}

func TestHas(t *testing.T) {
	m := New[string, Animal]()

	// Get a missing element.
	if m.Has("Money") {
		t.Error("element must not exist")
	}

	elephant := Animal{"elephant"}
	m.Set("elephant", elephant)

	if !m.Has("elephant") {
		t.Error("element doesn't exist, expected Has to return true")
	}
}

func TestRemove(t *testing.T) {
	m := New[string, Animal]()

	monkey := Animal{"monkey"}
	m.Set("monkey", monkey)
	m.Remove("monkey")

	if m.Count() != 0 {
		t.Error("expected count to be zero once an element has been removed")
	}

	temp, ok := m.Get("monkey")
	if ok {
		t.Error("expected ok to be false for missing elements")
	}
	if temp != (Animal{}) {
		t.Error("expected an element to be default after its removal")
	}

	// Remove a non-existent element.
	m.Remove("none")
	if _, ok = m.Get("none"); ok {
		t.Error("element has been created after Remove")
	}
}

func TestRemoveCb(t *testing.T) {
	m := New[string, Animal]()

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

	// Monkey must be removed
	result := m.RemoveCb("monkey", cb)
	if !result {
		t.Error("element must be removed")
	}
	if mapKey != "monkey" {
		t.Errorf("wrong key has been provided to the callback: expected=monkey got=%s", mapKey)
	}
	if mapVal != monkey {
		t.Errorf("wrong value has been provided to the callback: expected=%v got=%v", monkey, mapVal)
	}
	if !wasFound {
		t.Error("key must be found")
	}
	if m.Has("monkey") {
		t.Error("key has not been removed")
	}

	// Elephant must not be removed
	result = m.RemoveCb("elephant", cb)
	if result {
		t.Error("element must not be removed")
	}
	if mapKey != "elephant" {
		t.Errorf("wrong key has been provided to the callback: expected=elephant got=%s", mapKey)
	}
	if mapVal != elephant {
		t.Errorf("wrong value has been provided to the callback: expected=%v got=%v", elephant, mapVal)
	}
	if !wasFound {
		t.Error("key must be found")
	}
	if !m.Has("elephant") {
		t.Error("key has been removed")
	}

	// Unset key must remain unset
	result = m.RemoveCb("horse", cb)
	if result {
		t.Error("element must not be removed")
	}
	if mapKey != "horse" {
		t.Errorf("wrong key has been provided to the callback: expected=horse got=%s", mapKey)
	}
	if mapVal != (Animal{}) {
		t.Errorf("wrong value has been provided to the callback: expected=%v got=%v", Animal{}, mapVal)
	}
	if wasFound {
		t.Error("key must not be found")
	}
	if m.Has("horse") {
		t.Error("element has been created")
	}
}

func TestPop(t *testing.T) {
	m := New[string, Animal]()

	monkey := Animal{"monkey"}
	m.Set("monkey", monkey)

	v, exists := m.Pop("monkey")
	if !exists || v != monkey {
		t.Error("element has not been removed")
	}

	v2, exists2 := m.Pop("monkey")
	if exists2 || v2 == monkey {
		t.Error("element has been removed twice")
	}

	if m.Count() != 0 {
		t.Error("expected Count to return zero")
	}

	temp, ok := m.Get("monkey")
	if ok {
		t.Error("element has been found even though it has been removed")
	}
	if temp != (Animal{}) {
		t.Error("removed elements must return as default")
	}
}

func TestCount(t *testing.T) {
	m := New[string, Animal]()
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	if m.Count() != 100 {
		t.Errorf("expected the map to contain 100 elements, instead got %d", m.Count())
	}
}

func TestIsEmpty(t *testing.T) {
	m := New[string, Animal]()
	if !m.IsEmpty() {
		t.Error("new map must be empty")
	}
	m.Set("elephant", Animal{"elephant"})
	if m.IsEmpty() {
		t.Error("map must not be empty")
	}
}

func TestClear(t *testing.T) {
	m := New[string, Animal]()
	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	m.Clear()
	if m.Count() != 0 {
		t.Errorf("expected the map to be empty, instead got %d elements", m.Count())
	}
}

func TestIter(t *testing.T) {
	m := New[string, Animal]()
	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	counter := 0
	// Iterate over elements.
	m.Iter(func(_ string, val Animal) bool {
		if val == (Animal{}) {
			t.Error("expecting a valid object")
		}
		counter++
		return true
	})
	if counter != 100 {
		t.Errorf("expected to iterate over 100 elements, instead got %d", counter)
	}
}

func TestSeq(t *testing.T) {
	m := New[string, Animal]()
	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	counter := 0
	// Iterate over elements.
	for _, val := range m.Seq() {
		if val == (Animal{}) {
			t.Error("expecting a valid object")
		}
		counter++
	}
	if counter != 100 {
		t.Errorf("expected to iterate over 100 elements, instead got %d", counter)
	}
}

func TestItems(t *testing.T) {
	m := New[string, Animal]()
	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	if n := len(m.Items()); n != 100 {
		t.Errorf("expected 100 elements, got %d", n)
	}
}

func TestKeys(t *testing.T) {
	m := New[string, Animal]()
	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	if n := len(m.Keys()); n != 100 {
		t.Errorf("expected 100 keys, got %d", n)
	}
}

func TestRemoveFunc(t *testing.T) {
	testcases := []string{"elephant", "dolphin", "whale", "tiger", "lion"}

	type test struct {
		removed []string
		remains []string
	}

	tests := make([]test, 0)
	for n := 1; n <= len(testcases); n++ {
		for perm := range permute(testcases, n) {
			remains := slices.DeleteFunc(slices.Clone(testcases), func(s string) bool {
				return slices.Contains(perm, s)
			})
			slices.Sort(remains)
			tests = append(tests, test{
				removed: perm,
				remains: remains,
			})
		}
	}

	for _, tt := range tests {
		m := New[string, Animal]()
		for _, tc := range testcases {
			m.Set(tc, Animal{tc})
		}

		m.RemoveFunc(func(key string, v Animal) bool {
			return slices.Contains(tt.removed, key)
		})

		got := m.Keys()
		slices.Sort(got)

		if !slices.Equal(tt.remains, got) {
			t.Errorf("wrong keys after RemoveFunc: input=%v expected=%v got=%v",
				tt.removed, tt.remains, got)
		}
	}
}

func FuzzUpsert(f *testing.F) {
	testcases := []string{"elephant", "dolphin", "whale", "tiger", "lion"}

	for i := range testcases {
		for j := range testcases {
			f.Add(testcases[i], testcases[i])
			f.Add(testcases[i], testcases[j])
			f.Add(testcases[j], testcases[i])
			f.Add(testcases[j], testcases[j])
		}
	}

	cb := func(exists bool, valueInMap, newValue Animal) Animal {
		if !exists {
			return newValue
		}
		valueInMap.name += newValue.name
		return valueInMap
	}

	f.Fuzz(func(t *testing.T, xName, yName string) {
		x := Animal{xName}
		y := Animal{yName}

		m := New[string, Animal]()

		m.Upsert("upsert-x", x, cb)
		if animal, ok := m.Get("upsert-x"); ok {
			if animal.name != x.name {
				t.Errorf("upsert-x failed: got=%s expected=%s", animal.name, x.name)
			}
		} else {
			t.Error("upsert-x doesn't exist")
		}

		m.Upsert("upsert-y", y, cb)
		if animal, ok := m.Get("upsert-y"); ok {
			if animal.name != y.name {
				t.Errorf("upsert-y failed: got=%s expected=%s", animal.name, y.name)
			}
		} else {
			t.Error("upsert-y doesn't exist")
		}

		m.Set("set-upsert", x)
		m.Upsert("set-upsert", y, cb)
		if animal, ok := m.Get("set-upsert"); ok {
			expected := x.name + y.name
			if animal.name != expected {
				t.Errorf("set-upsert failed: got=%s expected=%s", animal.name, expected)
			}
		} else {
			t.Error("set-upsert doesn't exist")
		}

		m.Upsert("upsert-set", y, cb)
		m.Set("upsert-set", x)
		if animal, ok := m.Get("upsert-set"); ok {
			if animal.name != x.name {
				t.Errorf("upsert-set failed: got=%s expected=%s", animal.name, x.name)
			}
		} else {
			t.Error("upsert-set doesn't exist")
		}

		m.Upsert("upsert-upsert", x, cb)
		m.Upsert("upsert-upsert", y, cb)
		if animal, ok := m.Get("upsert-upsert"); ok {
			expected := x.name + y.name
			if animal.name != expected {
				t.Errorf("upsert-upsert failed: got=%s expected=%s", animal.name, expected)
			}
		} else {
			t.Error("upsert-upsert doesn't exist")
		}
	})
}

func FuzzUpdate(f *testing.F) {
	testcases := []string{"elephant", "dolphin", "whale", "tiger", "lion"}

	for i := range testcases {
		for j := range testcases {
			f.Add(testcases[i], testcases[i])
			f.Add(testcases[i], testcases[j])
			f.Add(testcases[j], testcases[i])
			f.Add(testcases[j], testcases[j])
		}
	}

	cb := func(valueInMap, newValue Animal) Animal {
		valueInMap.name += newValue.name
		return valueInMap
	}

	f.Fuzz(func(t *testing.T, xName, yName string) {
		x := Animal{xName}
		y := Animal{yName}

		m := New[string, Animal]()

		m.Update("update-x", x, cb)
		if _, ok := m.Get("update-x"); ok {
			t.Error("update-x exists")
		}

		m.Update("update-y", y, cb)
		if _, ok := m.Get("update-y"); ok {
			t.Error("update-y exists")
		}

		m.Set("set-update", x)
		m.Update("set-update", y, cb)
		if animal, ok := m.Get("set-update"); ok {
			expected := x.name + y.name
			if animal.name != expected {
				t.Errorf("set-update failed: got=%s expected=%s", animal.name, expected)
			}
		} else {
			t.Error("set-update doesn't exist")
		}

		m.Update("update-set", y, cb)
		m.Set("update-set", x)
		if animal, ok := m.Get("update-set"); ok {
			if animal.name != x.name {
				t.Errorf("update-set failed: got=%s expected=%s", animal.name, x.name)
			}
		} else {
			t.Error("update-set doesn't exist")
		}

		m.Update("update-update", x, cb)
		m.Update("update-update", y, cb)
		if _, ok := m.Get("update-update"); ok {
			t.Error("update-update")
		}
	})
}

func TestKeysWhenRemoving(t *testing.T) {
	m := New[string, Animal]()
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
			t.Error("got empty key")
		}
	}
}

func TestConcurrent(t *testing.T) {
	m := New[string, int]()
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
		t.Errorf("expected 1000 elements, got %d", m.Count())
	}
	// Make sure all inserted values we're fetched from map.
	for i := 0; i < iterations; i++ {
		if i != a[i] {
			t.Errorf("missing value %d", i)
		}
	}
}

func permute(cases []string, n int) iter.Seq[[]string] {
	return func(yield func([]string) bool) {
		idxs := make([]int, n)
		mins := make([]int, n)
		maxs := make([]int, n)

		for i := range n {
			idxs[i] = i
			mins[i] = i
			maxs[i] = i + 1 + (len(cases) - n)
		}

		for {
			result := make([]string, n)
			for i, idx := range idxs {
				result[i] = cases[idx]
			}
			if !yield(result) {
				return
			}
			for k := len(idxs) - 1; k >= 0; k-- {
				if idxs[k]+1 < maxs[k] {
					idxs[k]++
					break
				} else if k == 0 {
					return
				}
				if mins[k]+1 < maxs[k] {
					mins[k]++
					idxs[k] = mins[k]
				}
			}
		}
	}
}
