package validation

import (
	"slices"

	"github.com/infastin/go-rack/validation/constraints"
)

type numberValidatorData[T constraints.Number] struct {
	value T
	name  string
}

type NumberValidator[T constraints.Number] struct {
	data  *numberValidatorData[T]
	rules []NumberRule[T]
	scope validatorScope
}

func Number[T constraints.Number](n T, name string) NumberValidator[T] {
	return NumberValidator[T]{
		data: &numberValidatorData[T]{
			value: n,
			name:  name,
		},
		rules: make([]NumberRule[T], 0),
		scope: nil,
	}
}

func NumberI[T constraints.Number](n T) NumberValidator[T] {
	return NumberValidator[T]{
		data: &numberValidatorData[T]{
			value: n,
			name:  "",
		},
		rules: make([]NumberRule[T], 0),
		scope: nil,
	}
}

func NumberV[T constraints.Number]() NumberValidator[T] {
	return NumberValidator[T]{
		data:  nil,
		rules: make([]NumberRule[T], 0),
		scope: nil,
	}
}

func (nv NumberValidator[T]) If(condition bool) NumberValidator[T] {
	if nv.scope.Ok() {
		nv.scope = nv.scope.Push(condition)
	}
	return nv
}

func (nv NumberValidator[T]) ElseIf(condition bool) NumberValidator[T] {
	if !nv.scope.Ok() {
		nv.scope.Set(condition)
	}
	return nv
}

func (nv NumberValidator[T]) Else() NumberValidator[T] {
	if !nv.scope.Ok() {
		nv.scope.Set(true)
	}
	return nv
}

func (nv NumberValidator[T]) Break(condition bool) NumberValidator[T] {
	if !nv.scope.Empty() && condition {
		nv.scope.Set(false)
	}
	return nv
}

func (nv NumberValidator[T]) EndIf() NumberValidator[T] {
	if !nv.scope.Empty() {
		nv.scope = nv.scope.Pop()
	}
	return nv
}

func (nv NumberValidator[T]) Required(condition bool) NumberValidator[T] {
	if nv.scope.Ok() {
		nv.rules = append(nv.rules, Required[T](condition))
	}
	return nv
}

func (nv NumberValidator[T]) In(elements ...T) NumberValidator[T] {
	if nv.scope.Ok() {
		nv.rules = append(nv.rules, In(elements...))
	}
	return nv
}

func (nv NumberValidator[T]) NotIn(elements ...T) NumberValidator[T] {
	if nv.scope.Ok() {
		nv.rules = append(nv.rules, NotIn(elements...))
	}
	return nv
}

func (nv NumberValidator[T]) Equal(v T) NumberValidator[T] {
	if nv.scope.Ok() {
		nv.rules = append(nv.rules, Equal(v))
	}
	return nv
}

func (nv NumberValidator[T]) Less(v T) NumberValidator[T] {
	if nv.scope.Ok() {
		nv.rules = append(nv.rules, Less(v))
	}
	return nv
}

func (nv NumberValidator[T]) LessEqual(v T) NumberValidator[T] {
	if nv.scope.Ok() {
		nv.rules = append(nv.rules, LessEqual(v))
	}
	return nv
}

func (nv NumberValidator[T]) Greater(v T) NumberValidator[T] {
	if nv.scope.Ok() {
		nv.rules = append(nv.rules, Greater(v))
	}
	return nv
}

func (nv NumberValidator[T]) GreaterEqual(v T) NumberValidator[T] {
	if nv.scope.Ok() {
		nv.rules = append(nv.rules, GreaterEqual(v))
	}
	return nv
}

func (nv NumberValidator[T]) Between(a, b T) NumberValidator[T] {
	if nv.scope.Ok() {
		nv.rules = append(nv.rules, Between(a, b))
	}
	return nv
}

func (nv NumberValidator[T]) BetweenEqual(a, b T) NumberValidator[T] {
	if nv.scope.Ok() {
		nv.rules = append(nv.rules, BetweenEqual(a, b))
	}
	return nv
}

func (nv NumberValidator[T]) With(fns ...func(n T) error) NumberValidator[T] {
	if nv.scope.Ok() {
		nv.rules = slices.Grow(nv.rules, len(fns))
		for _, fn := range fns {
			nv.rules = append(nv.rules, NumberRuleFunc[T](fn))
		}
	}
	return nv
}

func (nv NumberValidator[T]) By(rules ...NumberRule[T]) NumberValidator[T] {
	if nv.scope.Ok() {
		nv.rules = append(nv.rules, rules...)
	}
	return nv
}

func (nv NumberValidator[T]) Valid() error {
	for _, rule := range nv.rules {
		if err := rule.Validate(nv.data.value); err != nil {
			if nv.data.name != "" {
				err = NewValueError(nv.data.name, err)
			}
			return err
		}
	}
	return nil
}

func (nv NumberValidator[T]) Validate(v T) error {
	for _, rule := range nv.rules {
		if err := rule.Validate(v); err != nil {
			return err
		}
	}
	return nil
}
