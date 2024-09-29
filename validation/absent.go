package validation

import "time"

type absentRule[T comparable] struct {
	condition bool
}

var (
	ErrEmpty = NewRuleError("empty", "must be blank")
	ErrNil   = NewRuleError("nil", "must be blank")
)

func Empty[T comparable](condition bool) absentRule[T] {
	return absentRule[T]{
		condition: condition,
	}
}

func (r absentRule[T]) Validate(v T) error {
	if r.condition && v != *new(T) {
		return ErrEmpty
	}
	return nil
}

type absentTimeRule struct {
	condition bool
}

func EmptyTime(condition bool) absentTimeRule {
	return absentTimeRule{
		condition: condition,
	}
}

func (r absentTimeRule) Validate(v time.Time) error {
	if r.condition && !v.IsZero() {
		return ErrEmpty
	}
	return nil
}

type absentPtrRule[T any] struct {
	condition bool
}

func NilPtr[T any](condition bool) absentPtrRule[T] {
	return absentPtrRule[T]{
		condition: condition,
	}
}

func (r absentPtrRule[T]) Validate(p *T) error {
	if r.condition && p != nil {
		return ErrNil
	}
	return nil
}

type absentSliceRule[T any] struct {
	condition bool
	checkNil  bool
}

func EmptySlice[T any](condition bool) absentSliceRule[T] {
	return absentSliceRule[T]{
		condition: condition,
		checkNil:  false,
	}
}

func NilSlice[T any](condition bool) absentSliceRule[T] {
	return absentSliceRule[T]{
		condition: condition,
		checkNil:  true,
	}
}

func (r absentSliceRule[T]) Validate(s []T) error {
	if r.condition {
		if !r.checkNil && len(s) != 0 {
			return ErrEmpty
		} else if r.checkNil && s != nil {
			return ErrNil
		}
	}
	return nil
}

type absentMapRule[T any] struct {
	condition bool
	checkNil  bool
}

func EmptyMap[T any](condition bool) absentMapRule[T] {
	return absentMapRule[T]{
		condition: condition,
		checkNil:  false,
	}
}

func NilMap[T any](condition bool) absentMapRule[T] {
	return absentMapRule[T]{
		condition: condition,
		checkNil:  true,
	}
}

func (r absentMapRule[T]) Validate(s map[string]T) error {
	if r.condition {
		if !r.checkNil && len(s) != 0 {
			return ErrEmpty
		} else if r.checkNil && s != nil {
			return ErrNil
		}
	}
	return nil
}
