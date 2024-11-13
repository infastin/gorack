package xalg

import (
	"cmp"
	"iter"
)

// Returns an iterator that produces elements from the union of two sorted ranges.
// Reference: https://en.cppreference.com/w/cpp/algorithm/set_union.
func SetUnion[E cmp.Ordered](s1, s2 iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		next1, stop1 := iter.Pull(s1)
		defer stop1()

		next2, stop2 := iter.Pull(s2)
		defer stop2()

		v1, ok1 := next1()
		v2, ok2 := next2()

		for ok1 {
			if !ok2 {
				for ; ok1; v1, ok1 = next1() {
					if !yield(v1) {
						return
					}
				}
				return
			}

			if v2 < v1 {
				if !yield(v2) {
					return
				}
				v2, ok2 = next2()
			} else {
				if !yield(v1) {
					return
				}

				if v1 == v2 {
					v2, ok2 = next2()
				}
				v1, ok1 = next1()
			}
		}

		for ; ok2; v2, ok2 = next2() {
			if !yield(v2) {
				return
			}
		}
	}
}

// Returns an iterator that produces elements from the union of two sorted ranges.
// Reference: https://en.cppreference.com/w/cpp/algorithm/set_union.
func SetUnionFunc[E any](s1, s2 iter.Seq[E], comp func(E, E) int) iter.Seq[E] {
	return func(yield func(E) bool) {
		next1, stop1 := iter.Pull(s1)
		defer stop1()

		next2, stop2 := iter.Pull(s2)
		defer stop2()

		v1, ok1 := next1()
		v2, ok2 := next2()

		for ok1 {
			if !ok2 {
				for ; ok1; v1, ok1 = next1() {
					if !yield(v1) {
						return
					}
				}
				return
			}

			if comp(v2, v1) < 0 {
				if !yield(v2) {
					return
				}
				v2, ok2 = next2()
			} else {
				if !yield(v1) {
					return
				}

				if comp(v1, v2) == 0 {
					v2, ok2 = next2()
				}
				v1, ok1 = next1()
			}
		}

		for ; ok2; v2, ok2 = next2() {
			if !yield(v2) {
				return
			}
		}
	}
}

// Returns an iterator that produces elements from the intersection of two sorted ranges.
// Reference: https://en.cppreference.com/w/cpp/algorithm/set_intersection.
func SetIntersection[E cmp.Ordered](s1, s2 iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		next1, stop1 := iter.Pull(s1)
		defer stop1()

		next2, stop2 := iter.Pull(s2)
		defer stop2()

		v1, ok1 := next1()
		v2, ok2 := next2()

		for ok1 && ok2 {
			if v1 < v2 {
				v1, ok1 = next1()
			} else {
				if v1 == v2 {
					if !yield(v1) {
						return
					}
					v1, ok1 = next1()
				}
				v2, ok2 = next2()
			}
		}
	}
}

// Returns an iterator that produces elements from the intersection of two sorted ranges.
// Reference: https://en.cppreference.com/w/cpp/algorithm/set_intersection.
func SetIntersectionFunc[E any](s1, s2 iter.Seq[E], comp func(E, E) int) iter.Seq[E] {
	return func(yield func(E) bool) {
		next1, stop1 := iter.Pull(s1)
		defer stop1()

		next2, stop2 := iter.Pull(s2)
		defer stop2()

		v1, ok1 := next1()
		v2, ok2 := next2()

		for ok1 && ok2 {
			if comp(v1, v2) < 0 {
				v1, ok1 = next1()
			} else {
				if comp(v1, v2) == 0 {
					if !yield(v1) {
						return
					}
					v1, ok1 = next1()
				}
				v2, ok2 = next2()
			}
		}
	}
}

// Returns an iterator that produces elements from the first sorted range
// that are not found in the second sorted range.
// Reference: https://en.cppreference.com/w/cpp/algorithm/set_difference.
func SetDifference[E cmp.Ordered](s1, s2 iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		next1, stop1 := iter.Pull(s1)
		defer stop1()

		next2, stop2 := iter.Pull(s2)
		defer stop2()

		v1, ok1 := next1()
		v2, ok2 := next2()

		for ok1 {
			if !ok2 {
				for ; ok1; v1, ok1 = next1() {
					if !yield(v1) {
						return
					}
				}
				return
			}

			if v1 < v2 {
				if !yield(v1) {
					return
				}
				v1, ok1 = next1()
			} else {
				if v1 == v2 {
					v1, ok1 = next1()
				}
				v2, ok2 = next2()
			}
		}
	}
}

// Returns an iterator that produces elements from the first sorted range
// that are not found in the second sorted range.
// Reference: https://en.cppreference.com/w/cpp/algorithm/set_difference.
func SetDifferenceFunc[E any](s1, s2 iter.Seq[E], comp func(E, E) int) iter.Seq[E] {
	return func(yield func(E) bool) {
		next1, stop1 := iter.Pull(s1)
		defer stop1()

		next2, stop2 := iter.Pull(s2)
		defer stop2()

		v1, ok1 := next1()
		v2, ok2 := next2()

		for ok1 {
			if !ok2 {
				for ; ok1; v1, ok1 = next1() {
					if !yield(v1) {
						return
					}
				}
				return
			}

			if comp(v1, v2) < 0 {
				if !yield(v1) {
					return
				}
				v1, ok1 = next1()
			} else {
				if comp(v1, v2) == 0 {
					v1, ok1 = next1()
				}
				v2, ok2 = next2()
			}
		}
	}
}

// Indicates to which set an element of the symmetric difference belongs.
type DiffElemType int8

const (
	// The element belongs to the first set.
	DiffElemFirst DiffElemType = iota
	// The element belongs to the second set.
	DiffElemSecond
)

// Returns an iterator that produces elements from the symmetric difference of two sorted ranges.
// Reference: https://en.cppreference.com/w/cpp/algorithm/set_symmetric_difference.
func SetSymmetricDifference[E cmp.Ordered](s1, s2 iter.Seq[E]) iter.Seq2[E, DiffElemType] {
	return func(yield func(E, DiffElemType) bool) {
		next1, stop1 := iter.Pull(s1)
		defer stop1()

		next2, stop2 := iter.Pull(s2)
		defer stop2()

		v1, ok1 := next1()
		v2, ok2 := next2()

		for ok1 {
			if !ok2 {
				for ; ok1; v1, ok1 = next1() {
					if !yield(v1, DiffElemFirst) {
						return
					}
				}
				break
			}

			if v1 < v2 {
				if !yield(v1, DiffElemFirst) {
					return
				}
				v1, ok1 = next1()
			} else {
				if v2 < v1 {
					if !yield(v2, DiffElemSecond) {
						return
					}
				} else {
					v1, ok1 = next1()
				}
				v2, ok2 = next2()
			}
		}

		for ; ok2; v2, ok2 = next2() {
			if !yield(v2, DiffElemSecond) {
				return
			}
		}
	}
}

// Returns an iterator that produces elements from the symmetric difference of two sorted ranges.
// Reference: https://en.cppreference.com/w/cpp/algorithm/set_symmetric_difference.
func SetSymmetricDifferenceFunc[E any](s1, s2 iter.Seq[E], comp func(E, E) int) iter.Seq2[E, DiffElemType] {
	return func(yield func(E, DiffElemType) bool) {
		next1, stop1 := iter.Pull(s1)
		defer stop1()

		next2, stop2 := iter.Pull(s2)
		defer stop2()

		v1, ok1 := next1()
		v2, ok2 := next2()

		for ok1 {
			if !ok2 {
				for ; ok1; v1, ok1 = next1() {
					if !yield(v1, DiffElemFirst) {
						return
					}
				}
				break
			}

			if comp(v1, v2) < 0 {
				if !yield(v1, DiffElemFirst) {
					return
				}
				v1, ok1 = next1()
			} else {
				if comp(v2, v1) < 0 {
					if !yield(v2, DiffElemSecond) {
						return
					}
				} else {
					v1, ok1 = next1()
				}
				v2, ok2 = next2()
			}
		}

		for ; ok2; v2, ok2 = next2() {
			if !yield(v2, DiffElemSecond) {
				return
			}
		}
	}
}
