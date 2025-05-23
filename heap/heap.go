package heap

// Heap is a generic min-heap backed by a slice.
type Heap[E any] struct {
	data []E
	cmp  func(a, b E) int
}

// Creates a new Heap with elements of type E
// with the comparison function cmp.
// cmp(a, b) should return a negative number when a < b, a positive number when
// a > b and zero when a == b or a and b are incomparable in the sense of
// a strict weak ordering.
func New[E any](cmp func(a, b E) int) *Heap[E] {
	return &Heap[E]{
		data: make([]E, 0),
		cmp:  cmp,
	}
}

// Returns the element at the index i from the heap.
func (h *Heap[E]) At(i int) E {
	return h.data[i]
}

// Changes the element at the index i in the heap.
// Fix should be called afterwards if this change breaks heap ordering.
func (h *Heap[E]) Set(i int, x E) {
	h.data[i] = x
}

// Returns the number of elements in the heap.
func (h *Heap[E]) Len() int {
	return len(h.data)
}

// Pushes the element x onto the heap.
func (h *Heap[E]) Push(x E) {
	h.data = append(h.data, x)
	h.up(len(h.data) - 1)
}

// Pop removes and returns the minimum element (according to cmp function) from the heap.
// Pop is equivalent to Remove(0).
func (h *Heap[E]) Pop() E {
	n := len(h.data) - 1
	if n != 0 {
		h.swap(0, n)
		h.down(0, n)
	}
	return h.pop()
}

// Removes and returns the element at index i from the heap.
func (h *Heap[E]) Remove(i int) E {
	n := len(h.data) - 1
	if i != n {
		h.swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	return h.pop()
}

// Re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(i) followed by a Push of the new value.
func (h *Heap[E]) Fix(i int) {
	if !h.down(i, len(h.data)) {
		h.up(i)
	}
}

func (h *Heap[E]) pop() E {
	var res E
	n := len(h.data) - 1
	res, h.data[n] = h.data[n], res
	h.data = h.data[:n]
	return res
}

func (h *Heap[E]) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

// Fixes heap ordering starting from index i and ending at index 0.
func (h *Heap[E]) up(i int) {
	for {
		p := (i - 1) / 2 // parent
		if p == i || h.cmp(h.data[i], h.data[p]) >= 0 {
			break
		}
		h.swap(i, p)
		i = p
	}
}

// Fixes heap ordering starting from index i0 and ending at index n.
// Returns false if the element at index i0 is not greater that its children
// or has no children at all.
func (h *Heap[E]) down(i0, n int) bool {
	i := i0
	for {
		c := (2 * i) + 1     // left child
		if c >= n || c < 0 { // c < 0 after int overflow
			break
		}
		if c2 := c + 1; c2 < n && h.cmp(h.data[c2], h.data[c]) < 0 {
			c = c2 // right child
		}
		if h.cmp(h.data[c], h.data[i]) >= 0 {
			break
		}
		h.swap(i, c)
		i = c
	}
	return i > i0
}
