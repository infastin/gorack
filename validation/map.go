package validation

import "slices"

type mapValidatorData[T any] struct {
	value map[string]T
	name  string
}

type MapValidator[T any] struct {
	data  *mapValidatorData[T]
	rules []MapRule[T]
	scope validatorScope
}

func Map[T any](m map[string]T, name string) MapValidator[T] {
	return MapValidator[T]{
		data: &mapValidatorData[T]{
			value: m,
			name:  name,
		},
		rules: make([]MapRule[T], 0),
		scope: nil,
	}
}

func MapI[T any](m map[string]T) MapValidator[T] {
	return MapValidator[T]{
		data: &mapValidatorData[T]{
			value: m,
			name:  "",
		},
		rules: make([]MapRule[T], 0),
		scope: nil,
	}
}

func MapV[T any]() MapValidator[T] {
	return MapValidator[T]{
		data:  nil,
		rules: make([]MapRule[T], 0),
		scope: nil,
	}
}

func (mv MapValidator[T]) If(condition bool) MapValidator[T] {
	if mv.scope.Ok() {
		mv.scope = mv.scope.Push(condition)
	}
	return mv
}

func (mv MapValidator[T]) ElseIf(condition bool) MapValidator[T] {
	if !mv.scope.Ok() {
		mv.scope.Set(condition)
	}
	return mv
}

func (mv MapValidator[T]) Else() MapValidator[T] {
	if !mv.scope.Ok() {
		mv.scope.Set(true)
	}
	return mv
}

func (mv MapValidator[T]) Break(condition bool) MapValidator[T] {
	if !mv.scope.Empty() && condition {
		mv.scope.Set(false)
	}
	return mv
}

func (mv MapValidator[T]) EndIf() MapValidator[T] {
	if !mv.scope.Empty() {
		mv.scope = mv.scope.Pop()
	}
	return mv
}

func (mv MapValidator[T]) Required(condition bool) MapValidator[T] {
	if mv.scope.Ok() {
		mv.rules = append(mv.rules, RequiredMap[T](condition))
	}
	return mv
}

func (mv MapValidator[T]) NilOrNotEmpty(condition bool) MapValidator[T] {
	if mv.scope.Ok() {
		mv.rules = append(mv.rules, NilOrNotEmptyMap[T](condition))
	}
	return mv
}

func (mv MapValidator[T]) Empty(condition bool) MapValidator[T] {
	if mv.scope.Ok() {
		mv.rules = append(mv.rules, EmptyMap[T](condition))
	}
	return mv
}

func (mv MapValidator[T]) NotNil(condition bool) MapValidator[T] {
	if mv.scope.Ok() {
		mv.rules = append(mv.rules, NotNilMap[T](condition))
	}
	return mv
}

func (mv MapValidator[T]) Nil(condition bool) MapValidator[T] {
	if mv.scope.Ok() {
		mv.rules = append(mv.rules, NilMap[T](condition))
	}
	return mv
}

func (mv MapValidator[T]) Length(min, max int) MapValidator[T] {
	if mv.scope.Ok() {
		mv.rules = append(mv.rules, LengthMap[T](min, max))
	}
	return mv
}

func (mv MapValidator[T]) With(fns ...func(s map[string]T) error) MapValidator[T] {
	if mv.scope.Ok() {
		mv.rules = slices.Grow(mv.rules, len(fns))
		for _, fn := range fns {
			mv.rules = append(mv.rules, MapRuleFunc[T](fn))
		}
	}
	return mv
}

func (mv MapValidator[T]) By(rules ...MapRule[T]) MapValidator[T] {
	if mv.scope.Ok() {
		mv.rules = append(mv.rules, rules...)
	}
	return mv
}

func (mv MapValidator[T]) Valid() error {
	for _, rule := range mv.rules {
		if err := rule.Validate(mv.data.value); err != nil {
			if mv.data.name != "" {
				err = NewValueError(mv.data.name, err)
			}
			return err
		}
	}
	return nil
}

func (mv MapValidator[T]) Validate(m map[string]T) error {
	for _, rule := range mv.rules {
		if err := rule.Validate(m); err != nil {
			return err
		}
	}
	return nil
}
