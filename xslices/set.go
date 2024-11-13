package xslices

import (
	"cmp"
	"iter"
	"slices"

	"github.com/infastin/gorack/xalg"
)

func Union[E cmp.Ordered, Slice ~[]E](s ...Slice) Slice {
	seqs := make([]iter.Seq[E], len(s))
	for i := range s {
		seqs = append(seqs, slices.Values(s[i]))
	}
	return slices.Collect(xalg.Union(seqs...))
}

func Intersection[E cmp.Ordered, Slice ~[]E](s1, s2 Slice) Slice {
	return slices.Collect(xalg.Intersection(slices.Values(s1), slices.Values(s2)))
}

func Difference[E cmp.Ordered, Slice ~[]E](s1, s2 Slice) Slice {
	return slices.Collect(xalg.Difference(slices.Values(s1), slices.Values(s2)))
}

func SymmetricDifference[E cmp.Ordered, Slice ~[]E](s1, s2 Slice) (first, second Slice) {
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
