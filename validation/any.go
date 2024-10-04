package validation

import "slices"

type anyValidatorData[T any] struct {
	value T
	name  string
}

type AnyValidator[T any] struct {
	data  *anyValidatorData[T]
	rules []AnyRule[T]
	scope validatorScope
}

func Any[T any](v T, name string) AnyValidator[T] {
	return AnyValidator[T]{
		data: &anyValidatorData[T]{
			value: v,
			name:  name,
		},
		rules: make([]AnyRule[T], 0),
		scope: nil,
	}
}

func AnyI[T any](v T) AnyValidator[T] {
	return AnyValidator[T]{
		data: &anyValidatorData[T]{
			value: v,
			name:  "",
		},
		rules: make([]AnyRule[T], 0),
		scope: nil,
	}
}

func AnyV[T any]() AnyValidator[T] {
	return AnyValidator[T]{
		data:  nil,
		rules: make([]AnyRule[T], 0),
		scope: nil,
	}
}

func (av AnyValidator[T]) If(condition bool) AnyValidator[T] {
	if av.scope.Ok() {
		av.scope = av.scope.Push(condition)
	}
	return av
}

func (av AnyValidator[T]) ElseIf(condition bool) AnyValidator[T] {
	if !av.scope.Ok() {
		av.scope.Set(condition)
	}
	return av
}

func (av AnyValidator[T]) Else() AnyValidator[T] {
	if !av.scope.Ok() {
		av.scope.Set(true)
	}
	return av
}

func (av AnyValidator[T]) Break(condition bool) AnyValidator[T] {
	if !av.scope.Empty() && condition {
		av.scope.Set(false)
	}
	return av
}

func (av AnyValidator[T]) EndIf() AnyValidator[T] {
	if !av.scope.Empty() {
		av.scope = av.scope.Pop()
	}
	return av
}

func (av AnyValidator[T]) Required(condition bool, isDefault func(v T) bool) AnyValidator[T] {
	if av.scope.Ok() {
		av.rules = append(av.rules, RequiredAny(condition, isDefault))
	}
	return av
}

func (av AnyValidator[T]) In(eq func(a, b T) bool, elements ...T) AnyValidator[T] {
	if av.scope.Ok() {
		av.rules = append(av.rules, InAny(eq, elements...))
	}
	return av
}

func (av AnyValidator[T]) NotIn(eq func(a, b T) bool, elements ...T) AnyValidator[T] {
	if av.scope.Ok() {
		av.rules = append(av.rules, NotInAny(eq, elements...))
	}
	return av
}

func (av AnyValidator[T]) Equal(eq func(a, b T) bool, v T) AnyValidator[T] {
	if av.scope.Ok() {
		av.rules = append(av.rules, EqualAny(eq, v))
	}
	return av
}

func (av AnyValidator[T]) Less(cmp func(a, b T) int, v T) AnyValidator[T] {
	if av.scope.Ok() {
		av.rules = append(av.rules, LessAny(cmp, v))
	}
	return av
}

func (av AnyValidator[T]) LessEqual(cmp func(a, b T) int, v T) AnyValidator[T] {
	if av.scope.Ok() {
		av.rules = append(av.rules, LessEqualAny(cmp, v))
	}
	return av
}

func (av AnyValidator[T]) Greater(cmp func(a, b T) int, v T) AnyValidator[T] {
	if av.scope.Ok() {
		av.rules = append(av.rules, GreaterAny(cmp, v))
	}
	return av
}

func (av AnyValidator[T]) GreaterEqual(cmp func(a, b T) int, v T) AnyValidator[T] {
	if av.scope.Ok() {
		av.rules = append(av.rules, GreaterEqualAny(cmp, v))
	}
	return av
}

func (av AnyValidator[T]) Between(cmp func(a, b T) int, a, b T) AnyValidator[T] {
	if av.scope.Ok() {
		av.rules = append(av.rules, BetweenAny(cmp, a, b))
	}
	return av
}

func (av AnyValidator[T]) BetweenEqual(cmp func(a, b T) int, a, b T) AnyValidator[T] {
	if av.scope.Ok() {
		av.rules = append(av.rules, BetweenEqualAny(cmp, a, b))
	}
	return av
}

func (av AnyValidator[T]) With(fns ...func(v T) error) AnyValidator[T] {
	if av.scope.Ok() {
		av.rules = slices.Grow(av.rules, len(fns))
		for _, fn := range fns {
			av.rules = append(av.rules, AnyRuleFunc[T](fn))
		}
	}
	return av
}

func (av AnyValidator[T]) By(rules ...AnyRule[T]) AnyValidator[T] {
	if av.scope.Ok() {
		av.rules = append(av.rules, rules...)
	}
	return av
}

func (av AnyValidator[T]) Valid() error {
	for _, rule := range av.rules {
		if err := rule.Validate(av.data.value); err != nil {
			if av.data.name != "" {
				err = NewValueError(av.data.name, err)
			}
			return err
		}
	}
	return nil
}

func (av AnyValidator[T]) Validate(v T) error {
	for _, rule := range av.rules {
		if err := rule.Validate(v); err != nil {
			return err
		}
	}
	return nil
}
