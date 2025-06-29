package xslices

import (
	"cmp"
	"slices"

	"github.com/infastin/gorack/xalg"
)

// SetUnion computes union of two sorted slices.
func SetUnion[E cmp.Ordered, Slice ~[]E](s1, s2 Slice) Slice {
	return slices.Collect(xalg.SetUnion(slices.Values(s1), slices.Values(s2)))
}

// SetUnionFunc computes union of two sorted slices using a comparison function.
func SetUnionFunc[E any, Slice ~[]E](s1, s2 Slice, comp func(E, E) int) Slice {
	return slices.Collect(xalg.SetUnionFunc(slices.Values(s1), slices.Values(s2), comp))
}

// SetIntersection computes intersection of two sorted slices.
func SetIntersection[E cmp.Ordered, Slice ~[]E](s1, s2 Slice) Slice {
	return slices.Collect(xalg.SetIntersection(slices.Values(s1), slices.Values(s2)))
}

// SetIntersectionFunc computes intersection of two sorted slices using a comparison function.
func SetIntersectionFunc[E any, Slice ~[]E](s1, s2 Slice, comp func(E, E) int) Slice {
	return slices.Collect(xalg.SetIntersectionFunc(slices.Values(s1), slices.Values(s2), comp))
}

// SetDifference returns a slice that contains elements from the first sorted slice
// that are not found in the second sorted slice.
func SetDifference[E cmp.Ordered, Slice ~[]E](s1, s2 Slice) Slice {
	return slices.Collect(xalg.SetDifference(slices.Values(s1), slices.Values(s2)))
}

// SetDifferenceFunc returns a slice that contains elements from the first sorted slice
// that are not found in the second sorted slice using a comparison function.
func SetDifferenceFunc[E any, Slice ~[]E](s1, s2 Slice, comp func(E, E) int) Slice {
	return slices.Collect(xalg.SetDifferenceFunc(slices.Values(s1), slices.Values(s2), comp))
}

// SetSymmetricDifference computes symmetric difference of two sorted slices.
func SetSymmetricDifference[E cmp.Ordered, Slice ~[]E](s1, s2 Slice) (first, second Slice) {
	for elem, typ := range xalg.SetSymmetricDifference(slices.Values(s1), slices.Values(s2)) {
		switch typ {
		case xalg.DiffElemFirst:
			first = append(first, elem)
		case xalg.DiffElemSecond:
			second = append(second, elem)
		}
	}
	return first, second
}

// Computes symmetric difference of two sorted slices using a comparison function.
func SetSymmetricDifferenceFunc[E any, Slice ~[]E](s1, s2 Slice, comp func(E, E) int) (first, second Slice) {
	for elem, typ := range xalg.SetSymmetricDifferenceFunc(slices.Values(s1), slices.Values(s2), comp) {
		switch typ {
		case xalg.DiffElemFirst:
			first = append(first, elem)
		case xalg.DiffElemSecond:
			second = append(second, elem)
		}
	}
	return first, second
}
