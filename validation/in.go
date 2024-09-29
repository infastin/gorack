package validation

import "time"

var ErrInInvalid = NewRuleError("in_invalid", "must be a valid value")

type inRule[T comparable] struct {
	elements []T
}

func In[T comparable](elements ...T) inRule[T] {
	return inRule[T]{
		elements: elements,
	}
}

func (r inRule[T]) Validate(v T) error {
	for i := range r.elements {
		if r.elements[i] == v {
			return nil
		}
	}
	return ErrInInvalid
}

type inAnyRule[T any] struct {
	elements []T
	eq       func(a, b T) bool
}

func InAny[T any](eq func(a, b T) bool, elements ...T) inAnyRule[T] {
	return inAnyRule[T]{
		elements: elements,
		eq:       eq,
	}
}

func (r inAnyRule[T]) Validate(v T) error {
	for i := range r.elements {
		if r.eq(r.elements[i], v) {
			return nil
		}
	}
	return ErrInInvalid
}

type inTimeRule struct {
	elements []time.Time
}

func InTime(elements ...time.Time) inTimeRule {
	return inTimeRule{
		elements: elements,
	}
}

func (r inTimeRule) Validate(t time.Time) error {
	for i := range r.elements {
		if t.Equal(r.elements[i]) {
			return nil
		}
	}
	return ErrInInvalid
}
