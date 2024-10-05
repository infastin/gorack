package validation

import (
	"time"

	"github.com/infastin/gorack/validation/constraints"
)

type Validatable interface {
	Validate() error
}

type StringRule[T ~string] interface {
	Validate(s T) error
}

type StringRuleFunc[T ~string] func(s T) error

func (fn StringRuleFunc[T]) Validate(s T) error {
	return fn(s)
}

type NumberRule[T constraints.Number] interface {
	Validate(n T) error
}

type NumberRuleFunc[T constraints.Number] func(n T) error

func (fn NumberRuleFunc[T]) Validate(n T) error {
	return fn(n)
}

type TimeRule interface {
	Validate(t time.Time) error
}

type TimeRuleFunc func(t time.Time) error

func (fn TimeRuleFunc) Validate(t time.Time) error {
	return fn(t)
}

type PtrRule[T any] interface {
	Validate(p *T) error
}

type PtrRuleFunc[T any] func(p *T) error

func (fn PtrRuleFunc[T]) Validate(p *T) error {
	return fn(p)
}

type SliceRule[T any] interface {
	Validate(s []T) error
}

type SliceRuleFunc[T any] func(s []T) error

func (fn SliceRuleFunc[T]) Validate(s []T) error {
	return fn(s)
}

type MapRule[T any] interface {
	Validate(m map[string]T) error
}

type MapRuleFunc[T any] func(m map[string]T) error

func (fn MapRuleFunc[T]) Validate(m map[string]T) error {
	return fn(m)
}

type AnyRule[T any] interface {
	Validate(v T) error
}

type AnyRuleFunc[T any] func(v T) error

func (fn AnyRuleFunc[T]) Validate(v T) error {
	return fn(v)
}

type ComparableRule[T comparable] interface {
	Validate(v T) error
}

type ComparableRuleFunc[T comparable] func(v T) error

func (fn ComparableRuleFunc[T]) Validate(v T) error {
	return fn(v)
}

type Validator interface {
	Valid() error
}

func All(validators ...Validator) error {
	var errs Errors
	for _, v := range validators {
		if err := v.Valid(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return errs
	}
	return nil
}
