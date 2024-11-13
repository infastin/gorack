package xalg

import (
	"cmp"
	"iter"
)

func Union[E cmp.Ordered](seqs ...iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		seen := make(map[E]struct{})
		for i := range seqs {
			for v := range seqs[i] {
				if _, ok := seen[v]; ok {
					continue
				}
				seen[v] = struct{}{}
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Intersection[E cmp.Ordered](s1, s2 iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		seen := make(map[E]struct{})

		for v := range s1 {
			seen[v] = struct{}{}
		}

		for v := range s2 {
			if _, ok := seen[v]; ok && !yield(v) {
				return
			}
		}
	}
}

func Difference[E cmp.Ordered](s1, s2 iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		seen := make(map[E]struct{})

		for v := range s2 {
			seen[v] = struct{}{}
		}

		for v := range s1 {
			if _, ok := seen[v]; !ok && !yield(v) {
				return
			}
		}
	}
}

func SymmetricDifference[E cmp.Ordered](s1, s2 iter.Seq[E]) iter.Seq2[E, DiffElemType] {
	return func(yield func(E, DiffElemType) bool) {
		firstSeen := make(map[E]struct{})
		secondSeen := make(map[E]struct{})

		first := make([]E, 0)
		second := make([]E, 0)

		for v := range s1 {
			firstSeen[v] = struct{}{}
			first = append(first, v)
		}

		for v := range s2 {
			secondSeen[v] = struct{}{}
			second = append(second, v)
		}

		for i := range first {
			if _, ok := secondSeen[first[i]]; !ok && !yield(first[i], DiffElemFirst) {
				return
			}
		}

		for i := range second {
			if _, ok := firstSeen[second[i]]; !ok && !yield(second[i], DiffElemSecond) {
				return
			}
		}
	}
}
