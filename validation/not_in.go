package validation

import "time"

var ErrNotInInvalid = NewRuleError("not_in_invalid", "must be a valid value")

type notInRule[T comparable] struct {
	elements []T
}

func NotIn[T comparable](elements ...T) notInRule[T] {
	return notInRule[T]{
		elements: elements,
	}
}

func (r notInRule[T]) Validate(v T) error {
	for i := range r.elements {
		if r.elements[i] == v {
			return ErrNotInInvalid
		}
	}
	return nil
}

type notInAnyRule[T any] struct {
	elements []T
	eq       func(a, b T) bool
}

func NotInAny[T any](eq func(a, b T) bool, elements ...T) notInAnyRule[T] {
	return notInAnyRule[T]{
		elements: elements,
		eq:       eq,
	}
}

func (r notInAnyRule[T]) Validate(v T) error {
	for i := range r.elements {
		if r.eq(r.elements[i], v) {
			return ErrNotInInvalid
		}
	}
	return nil
}

type notInTimeRule struct {
	elements []time.Time
}

func NotInTime(elements ...time.Time) notInTimeRule {
	return notInTimeRule{
		elements: elements,
	}
}

func (r notInTimeRule) Validate(t time.Time) error {
	for i := range r.elements {
		if t.Equal(r.elements[i]) {
			return ErrNotInInvalid
		}
	}
	return nil
}
