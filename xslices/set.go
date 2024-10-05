package xslices

import (
	"cmp"
	"slices"

	"github.com/infastin/gorack/xalg"
)

func Union[E cmp.Ordered](s1, s2 []E) []E {
	return slices.Collect(xalg.Union(slices.Values(s1), slices.Values(s2)))
}

func UnionFunc[E any](s1, s2 []E, comp func(E, E) int) []E {
	return slices.Collect(xalg.UnionFunc(slices.Values(s1), slices.Values(s2), comp))
}

func Intersection[E cmp.Ordered](s1, s2 []E) []E {
	return slices.Collect(xalg.Intersection(slices.Values(s1), slices.Values(s2)))
}

func IntersectionFunc[E any](s1, s2 []E, comp func(E, E) int) []E {
	return slices.Collect(xalg.IntersectionFunc(slices.Values(s1), slices.Values(s2), comp))
}

func SymmetricDifference[E cmp.Ordered](s1, s2 []E) (first, second []E) {
	for elem, typ := range xalg.SymmetricDifference(slices.Values(s1), slices.Values(s2)) {
		switch typ {
		case xalg.DiffElemFirst:
			first = append(first, elem)
		case xalg.DiffElemSecond:
			second = append(second, elem)
		}
	}
	return first, second
}

func SymmetricDifferenceFunc[E any](s1, s2 []E, comp func(E, E) int) (first, second []E) {
	for elem, typ := range xalg.SymmetricDifferenceFunc(slices.Values(s1), slices.Values(s2), comp) {
		switch typ {
		case xalg.DiffElemFirst:
			first = append(first, elem)
		case xalg.DiffElemSecond:
			second = append(second, elem)
		}
	}
	return first, second
}

func Difference[E cmp.Ordered](s1, s2 []E) []E {
	return slices.Collect(xalg.Difference(slices.Values(s1), slices.Values(s2)))
}

func DifferenceFunc[E any](s1, s2 []E, comp func(E, E) int) []E {
	return slices.Collect(xalg.DifferenceFunc(slices.Values(s1), slices.Values(s2), comp))
}
