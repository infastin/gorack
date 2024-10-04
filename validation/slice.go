package validation

import "slices"

type sliceValidatorData[T any] struct {
	value []T
	name  string
}

type SliceValidator[T any] struct {
	data  *sliceValidatorData[T]
	rules []SliceRule[T]
	scope validatorScope
}

func Slice[T any](s []T, name string) SliceValidator[T] {
	return SliceValidator[T]{
		data: &sliceValidatorData[T]{
			value: s,
			name:  name,
		},
		rules: make([]SliceRule[T], 0),
		scope: nil,
	}
}

func SliceI[T any](s []T) SliceValidator[T] {
	return SliceValidator[T]{
		data: &sliceValidatorData[T]{
			value: s,
			name:  "",
		},
		rules: make([]SliceRule[T], 0),
		scope: nil,
	}
}

func SliceV[T any]() SliceValidator[T] {
	return SliceValidator[T]{
		data:  nil,
		rules: make([]SliceRule[T], 0),
		scope: nil,
	}
}

func (sv SliceValidator[T]) If(condition bool) SliceValidator[T] {
	if sv.scope.Ok() {
		sv.scope = sv.scope.Push(condition)
	}
	return sv
}

func (sv SliceValidator[T]) ElseIf(condition bool) SliceValidator[T] {
	if !sv.scope.Ok() {
		sv.scope.Set(condition)
	}
	return sv
}

func (sv SliceValidator[T]) Else() SliceValidator[T] {
	if !sv.scope.Ok() {
		sv.scope.Set(true)
	}
	return sv
}

func (sv SliceValidator[T]) Break(condition bool) SliceValidator[T] {
	if !sv.scope.Empty() && condition {
		sv.scope.Set(false)
	}
	return sv
}

func (sv SliceValidator[T]) EndIf() SliceValidator[T] {
	if !sv.scope.Empty() {
		sv.scope = sv.scope.Pop()
	}
	return sv
}

func (sv SliceValidator[T]) Required(condition bool) SliceValidator[T] {
	if sv.scope.Ok() {
		sv.rules = append(sv.rules, RequiredSlice[T](condition))
	}
	return sv
}

func (sv SliceValidator[T]) NilOrNotEmpty(condition bool) SliceValidator[T] {
	if sv.scope.Ok() {
		sv.rules = append(sv.rules, NilOrNotEmptySlice[T](condition))
	}
	return sv
}

func (sv SliceValidator[T]) Empty(condition bool) SliceValidator[T] {
	if sv.scope.Ok() {
		sv.rules = append(sv.rules, EmptySlice[T](condition))
	}
	return sv
}

func (sv SliceValidator[T]) NotNil(condition bool) SliceValidator[T] {
	if sv.scope.Ok() {
		sv.rules = append(sv.rules, NotNilSlice[T](condition))
	}
	return sv
}

func (sv SliceValidator[T]) Nil(condition bool) SliceValidator[T] {
	if sv.scope.Ok() {
		sv.rules = append(sv.rules, NilSlice[T](condition))
	}
	return sv
}

func (sv SliceValidator[T]) Length(min, max int) SliceValidator[T] {
	if sv.scope.Ok() {
		sv.rules = append(sv.rules, LengthSlice[T](min, max))
	}
	return sv
}

func (sv SliceValidator[T]) With(fns ...func(s []T) error) SliceValidator[T] {
	if sv.scope.Ok() {
		sv.rules = slices.Grow(sv.rules, len(fns))
		for _, fn := range fns {
			sv.rules = append(sv.rules, SliceRuleFunc[T](fn))
		}
	}
	return sv
}

func (sv SliceValidator[T]) By(rules ...SliceRule[T]) SliceValidator[T] {
	if sv.scope.Ok() {
		sv.rules = append(sv.rules, rules...)
	}
	return sv
}

func (sv SliceValidator[T]) ValuesWith(fns ...func(v T) error) SliceValidator[T] {
	if sv.scope.Ok() {
		sv.rules = append(sv.rules, SliceRuleFunc[T](func(s []T) error {
			for i := range s {
				for _, fn := range fns {
					if err := fn(s[i]); err != nil {
						return NewIndexError(i, err)
					}
				}
			}
			return nil
		}))
	}
	return sv
}

func (sv SliceValidator[T]) ValuesBy(rules ...AnyRule[T]) SliceValidator[T] {
	if sv.scope.Ok() {
		sv.rules = append(sv.rules, SliceRuleFunc[T](func(s []T) error {
			for i := range s {
				for _, rule := range rules {
					if err := rule.Validate(s[i]); err != nil {
						return NewIndexError(i, err)
					}
				}
			}
			return nil
		}))
	}
	return sv
}

func (sv SliceValidator[T]) ValuesPtrBy(rules ...AnyRule[*T]) SliceValidator[T] {
	if sv.scope.Ok() {
		sv.rules = append(sv.rules, SliceRuleFunc[T](func(s []T) error {
			for i := range s {
				for _, rule := range rules {
					if err := rule.Validate(&s[i]); err != nil {
						return NewIndexError(i, err)
					}
				}
			}
			return nil
		}))
	}
	return sv
}

func (sv SliceValidator[T]) ValuesPtrWith(fns ...func(v *T) error) SliceValidator[T] {
	if sv.scope.Ok() {
		sv.rules = append(sv.rules, SliceRuleFunc[T](func(s []T) error {
			for i := range s {
				for _, fn := range fns {
					if err := fn(&s[i]); err != nil {
						return NewIndexError(i, err)
					}
				}
			}
			return nil
		}))
	}
	return sv
}

func (sv SliceValidator[T]) Valid() error {
	for _, rule := range sv.rules {
		if err := rule.Validate(sv.data.value); err != nil {
			if sv.data.name != "" {
				err = NewValueError(sv.data.name, err)
			}
			return err
		}
	}
	return nil
}

func (sv SliceValidator[T]) Validate(v []T) error {
	for _, rule := range sv.rules {
		if err := rule.Validate(v); err != nil {
			return err
		}
	}
	return nil
}
