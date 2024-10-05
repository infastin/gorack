package validation

import (
	"fmt"
	"time"

	"github.com/infastin/gorack/constraints"
)

type compareRule[T any] struct {
	comp       func(x T) bool
	buildError func() error
}

func Equal[T comparable](v T) compareRule[T] {
	return compareRule[T]{
		comp: func(x T) bool {
			return x == v
		},
		buildError: func() error {
			return buildEqualError(v)
		},
	}
}

func EqualAny[T any](eq func(a, b T) bool, v T) compareRule[T] {
	return compareRule[T]{
		comp: func(x T) bool {
			return eq(x, v)
		},
		buildError: func() error {
			return buildEqualError(v)
		},
	}
}

func Less[T constraints.Ordered](v T) compareRule[T] {
	return compareRule[T]{
		comp: func(x T) bool {
			return x < v
		},
		buildError: func() error {
			return buildLessError(v)
		},
	}
}

func LessAny[T any](cmp func(a, b T) int, v T) compareRule[T] {
	return compareRule[T]{
		comp: func(x T) bool {
			return cmp(x, v) < 0
		},
		buildError: func() error {
			return buildLessError(v)
		},
	}
}

func LessEqual[T constraints.Ordered](v T) compareRule[T] {
	return compareRule[T]{
		comp: func(x T) bool {
			return x <= v
		},
		buildError: func() error {
			return buildLessEqualError(v)
		},
	}
}

func LessEqualAny[T any](cmp func(a, b T) int, v T) compareRule[T] {
	return compareRule[T]{
		comp: func(x T) bool {
			return cmp(x, v) <= 0
		},
		buildError: func() error {
			return buildLessEqualError(v)
		},
	}
}

func Greater[T constraints.Ordered](v T) compareRule[T] {
	return compareRule[T]{
		comp: func(x T) bool {
			return x > v
		},
		buildError: func() error {
			return buildGreaterError(v)
		},
	}
}

func GreaterAny[T any](cmp func(a, b T) int, v T) compareRule[T] {
	return compareRule[T]{
		comp: func(x T) bool {
			return cmp(x, v) > 0
		},
		buildError: func() error {
			return buildGreaterError(v)
		},
	}
}

func GreaterEqual[T constraints.Ordered](v T) compareRule[T] {
	return compareRule[T]{
		comp: func(x T) bool {
			return x >= v
		},
		buildError: func() error {
			return buildGreaterEqualError(v)
		},
	}
}

func GreaterEqualAny[T any](cmp func(a, b T) int, v T) compareRule[T] {
	return compareRule[T]{
		comp: func(x T) bool {
			return cmp(x, v) >= 0
		},
		buildError: func() error {
			return buildGreaterEqualError(v)
		},
	}
}

func Between[T constraints.Ordered](a, b T) compareRule[T] {
	return compareRule[T]{
		comp: func(x T) bool {
			return x > a && x < b
		},
		buildError: func() error {
			return buildBetweenError(a, b)
		},
	}
}

func BetweenAny[T any](cmp func(a, b T) int, a, b T) compareRule[T] {
	return compareRule[T]{
		comp: func(x T) bool {
			return cmp(x, a) > 0 && cmp(x, b) < 0
		},
		buildError: func() error {
			return buildBetweenError(a, b)
		},
	}
}

func BetweenEqual[T constraints.Ordered](a, b T) compareRule[T] {
	return compareRule[T]{
		comp: func(x T) bool {
			return x >= a && x <= b
		},
		buildError: func() error {
			return buildBetweenEqualError(a, b)
		},
	}
}

func BetweenEqualAny[T any](cmp func(a, b T) int, a, b T) compareRule[T] {
	return compareRule[T]{
		comp: func(x T) bool {
			return cmp(x, a) >= 0 && cmp(x, b) <= 0
		},
		buildError: func() error {
			return buildBetweenEqualError(a, b)
		},
	}
}

func (r compareRule[T]) Validate(v T) error {
	if !r.comp(v) {
		return r.buildError()
	}
	return nil
}

type compareTimeRule struct {
	comp       func(x time.Time) bool
	buildError func() error
}

func EqualTime(v time.Time) compareTimeRule {
	return compareTimeRule{
		comp: func(x time.Time) bool {
			return x.Equal(v)
		},
		buildError: func() error {
			return buildEqualError(v)
		},
	}
}

func LessTime(v time.Time) compareTimeRule {
	return compareTimeRule{
		comp: func(x time.Time) bool {
			return x.Compare(v) < 0
		},
		buildError: func() error {
			return buildLessError(v)
		},
	}
}

func LessEqualTime(v time.Time) compareTimeRule {
	return compareTimeRule{
		comp: func(x time.Time) bool {
			return x.Compare(v) <= 0
		},
		buildError: func() error {
			return buildLessEqualError(v)
		},
	}
}

func GreaterTime(v time.Time) compareTimeRule {
	return compareTimeRule{
		comp: func(x time.Time) bool {
			return x.Compare(v) > 0
		},
		buildError: func() error {
			return buildGreaterError(v)
		},
	}
}

func GreaterEqualTime(v time.Time) compareTimeRule {
	return compareTimeRule{
		comp: func(x time.Time) bool {
			return x.Compare(v) >= 0
		},
		buildError: func() error {
			return buildGreaterEqualError(v)
		},
	}
}

func BetweenTime(a, b time.Time) compareTimeRule {
	return compareTimeRule{
		comp: func(x time.Time) bool {
			return x.After(a) && x.Before(b)
		},
		buildError: func() error {
			return buildBetweenError(a, b)
		},
	}
}

func BetweenEqualTime(a, b time.Time) compareTimeRule {
	return compareTimeRule{
		comp: func(x time.Time) bool {
			return x.Compare(a) >= 0 && x.Compare(b) <= 0
		},
		buildError: func() error {
			return buildBetweenEqualError(a, b)
		},
	}
}

func (r compareTimeRule) Validate(v time.Time) error {
	if !r.comp(v) {
		return r.buildError()
	}
	return nil
}

func buildEqualError(v any) error {
	return NewRuleError("equal", fmt.Sprintf("must be equal to %v", v))
}

func buildLessError(v any) error {
	return NewRuleError("less", fmt.Sprintf("must be less than %v", v))
}

func buildLessEqualError(v any) error {
	return NewRuleError("less_equal", fmt.Sprintf("must be no greater than %v", v))
}

func buildGreaterError(v any) error {
	return NewRuleError("greater", fmt.Sprintf("must be greater than %v", v))
}

func buildGreaterEqualError(v any) error {
	return NewRuleError("greater_equal", fmt.Sprintf("must be no less than %v", v))
}

func buildBetweenError(a, b any) error {
	return NewRuleError("between", fmt.Sprintf("must exclusively be between %v and %v", a, b))
}

func buildBetweenEqualError(a, b any) error {
	return NewRuleError("between_equal", fmt.Sprintf("must inclusively be between %v and %v", a, b))
}
