package validation

import (
	"time"
)

type requiredRule[T comparable] struct {
	condition bool
}

var (
	ErrRequired      = NewRuleError("required", "cannot be blank")
	ErrNilOrNotEmpty = NewRuleError("nil_or_not_empty_required", "cannot be blank")
)

func Required[T comparable](condition bool) requiredRule[T] {
	return requiredRule[T]{
		condition: condition,
	}
}

func (r requiredRule[T]) Validate(v T) error {
	if r.condition && v == *new(T) {
		return ErrRequired
	}
	return nil
}

type requiredTimeRule struct {
	condition bool
}

func RequiredTime(condition bool) requiredTimeRule {
	return requiredTimeRule{
		condition: condition,
	}
}

func (r requiredTimeRule) Validate(v time.Time) error {
	if r.condition && v.IsZero() {
		return ErrRequired
	}
	return nil
}

type requiredSliceRule[T any] struct {
	condition bool
	skipNil   bool
}

func RequiredSlice[T any](condition bool) requiredSliceRule[T] {
	return requiredSliceRule[T]{
		condition: condition,
		skipNil:   false,
	}
}

func NilOrNotEmptySlice[T any](condition bool) requiredSliceRule[T] {
	return requiredSliceRule[T]{
		condition: condition,
		skipNil:   true,
	}
}

func (r requiredSliceRule[T]) Validate(s []T) error {
	if r.condition && len(s) == 0 {
		if !r.skipNil {
			return ErrRequired
		} else if s != nil {
			return ErrNilOrNotEmpty
		}
	}
	return nil
}

type requiredMapRule[T any] struct {
	condition bool
	checkNil  bool
}

func RequiredMap[T any](condition bool) requiredMapRule[T] {
	return requiredMapRule[T]{
		condition: condition,
		checkNil:  false,
	}
}

func NilOrNotEmptyMap[T any](condition bool) requiredMapRule[T] {
	return requiredMapRule[T]{
		condition: condition,
		checkNil:  true,
	}
}

func (r requiredMapRule[T]) Validate(s map[string]T) error {
	if r.condition && len(s) == 0 {
		if !r.checkNil {
			return ErrRequired
		} else if s != nil {
			return ErrNilOrNotEmpty
		}
	}
	return nil
}

type requiredAnyRule[T any] struct {
	condition bool
	isDefault func(a T) bool
}

func RequiredAny[T any](condition bool, isDefault func(a T) bool) requiredAnyRule[T] {
	return requiredAnyRule[T]{
		condition: condition,
		isDefault: isDefault,
	}
}

func (r requiredAnyRule[T]) Validate(v T) error {
	if r.condition && r.isDefault(v) {
		return ErrRequired
	}
	return nil
}
