package heap_test

import (
	"math/rand"
	"testing"

	"github.com/infastin/gorack/heap"
)

func cmpInt(a, b int) int {
	return a - b
}

func verify(t *testing.T, h *heap.Heap[int], i int) {
	t.Helper()
	n := h.Len()
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if h.At(j1) < h.At(i) {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h.At(i), j1, h.At(j1))
			return
		}
		verify(t, h, j1)
	}
	if j2 < n {
		if h.At(j2) < h.At(i) {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h.At(i), j1, h.At(j2))
			return
		}
		verify(t, h, j2)
	}
}

func Test(t *testing.T) {
	h := heap.New(cmpInt)
	verify(t, h, 0)

	for i := 20; i > 10; i-- {
		h.Push(i)
	}
	verify(t, h, 0)

	for i := 10; i > 0; i-- {
		h.Push(i)
		verify(t, h, 0)
	}

	for i := 1; h.Len() > 0; i++ {
		x := h.Pop()
		if i < 20 {
			h.Push(20 + i)
		}
		verify(t, h, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestRemove0(t *testing.T) {
	h := heap.New(cmpInt)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	verify(t, h, 0)

	for h.Len() > 0 {
		i := h.Len() - 1
		x := h.Remove(i)
		if x != i {
			t.Errorf("Remove(%d) got %d; want %d", i, x, i)
		}
		verify(t, h, 0)
	}
}

func TestRemove1(t *testing.T) {
	h := heap.New(cmpInt)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	verify(t, h, 0)

	for i := 0; h.Len() > 0; i++ {
		x := h.Remove(0)
		if x != i {
			t.Errorf("Remove(0) got %d; want %d", x, i)
		}
		verify(t, h, 0)
	}
}

func TestRemove2(t *testing.T) {
	N := 10

	h := heap.New(cmpInt)
	for i := 0; i < N; i++ {
		h.Push(i)
	}
	verify(t, h, 0)

	m := make(map[int]bool)
	for h.Len() > 0 {
		m[h.Remove((h.Len()-1)/2)] = true
		verify(t, h, 0)
	}

	if len(m) != N {
		t.Errorf("len(m) = %d; want %d", len(m), N)
	}
	for i := 0; i < len(m); i++ {
		if !m[i] {
			t.Errorf("m[%d] doesn't exist", i)
		}
	}
}

func TestFix(t *testing.T) {
	h := heap.New(cmpInt)
	verify(t, h, 0)

	for i := 200; i > 0; i -= 10 {
		h.Push(i)
	}
	verify(t, h, 0)

	if h.At(0) != 10 {
		t.Fatalf("Expected head to be 10, was %d", h.At(0))
	}
	h.Set(0, 120)
	h.Fix(0)
	verify(t, h, 0)

	for i := 100; i > 0; i-- {
		elem := rand.Intn(h.Len())
		if i&1 == 0 {
			h.Set(elem, h.At(elem)*2)
		} else {
			h.Set(elem, h.At(elem)/2)
		}
		h.Fix(elem)
		verify(t, h, 0)
	}
}
