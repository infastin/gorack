package validation

import "slices"

type comparableValidatorData[T comparable] struct {
	value T
	name  string
}

type ComparableValidator[T comparable] struct {
	data  *comparableValidatorData[T]
	rules []ComparableRule[T]
	scope validatorScope
}

func Comparable[T comparable](v T, name string) ComparableValidator[T] {
	return ComparableValidator[T]{
		data: &comparableValidatorData[T]{
			value: v,
			name:  name,
		},
		rules: make([]ComparableRule[T], 0),
		scope: nil,
	}
}

func ComparableI[T comparable](v T) ComparableValidator[T] {
	return ComparableValidator[T]{
		data: &comparableValidatorData[T]{
			value: v,
			name:  "",
		},
		rules: make([]ComparableRule[T], 0),
		scope: nil,
	}
}

func ComparableV[T comparable]() ComparableValidator[T] {
	return ComparableValidator[T]{
		data:  nil,
		rules: make([]ComparableRule[T], 0),
		scope: nil,
	}
}

func (cv ComparableValidator[T]) If(condition bool) ComparableValidator[T] {
	if cv.scope.Ok() {
		cv.scope = cv.scope.Push(condition)
	}
	return cv
}

func (cv ComparableValidator[T]) ElseIf(condition bool) ComparableValidator[T] {
	if !cv.scope.Ok() {
		cv.scope.Set(condition)
	}
	return cv
}

func (cv ComparableValidator[T]) Else() ComparableValidator[T] {
	if !cv.scope.Ok() {
		cv.scope.Set(true)
	}
	return cv
}

func (cv ComparableValidator[T]) Break(condition bool) ComparableValidator[T] {
	if !cv.scope.Empty() && condition {
		cv.scope.Set(false)
	}
	return cv
}

func (cv ComparableValidator[T]) EndIf() ComparableValidator[T] {
	if !cv.scope.Empty() {
		cv.scope = cv.scope.Pop()
	}
	return cv
}

func (cv ComparableValidator[T]) Required(condition bool) ComparableValidator[T] {
	if cv.scope.Ok() {
		cv.rules = append(cv.rules, Required[T](condition))
	}
	return cv
}

func (cv ComparableValidator[T]) In(elements ...T) ComparableValidator[T] {
	if cv.scope.Ok() {
		cv.rules = append(cv.rules, In(elements...))
	}
	return cv
}

func (cv ComparableValidator[T]) NotIn(elements ...T) ComparableValidator[T] {
	if cv.scope.Ok() {
		cv.rules = append(cv.rules, NotIn(elements...))
	}
	return cv
}

func (cv ComparableValidator[T]) Equal(v T) ComparableValidator[T] {
	if cv.scope.Ok() {
		cv.rules = append(cv.rules, Equal(v))
	}
	return cv
}

func (cv ComparableValidator[T]) With(fns ...func(v T) error) ComparableValidator[T] {
	if cv.scope.Ok() {
		cv.rules = slices.Grow(cv.rules, len(fns))
		for _, fn := range fns {
			cv.rules = append(cv.rules, ComparableRuleFunc[T](fn))
		}
	}
	return cv
}

func (cv ComparableValidator[T]) By(rules ...ComparableRule[T]) ComparableValidator[T] {
	if cv.scope.Ok() {
		cv.rules = append(cv.rules, rules...)
	}
	return cv
}

func (cv ComparableValidator[T]) Valid() error {
	for _, rule := range cv.rules {
		if err := rule.Validate(cv.data.value); err != nil {
			if cv.data.name != "" {
				err = NewValueError(cv.data.name, err)
			}
			return err
		}
	}
	return nil
}

func (cv ComparableValidator[T]) Validate(v T) error {
	for _, rule := range cv.rules {
		if err := rule.Validate(v); err != nil {
			return err
		}
	}
	return nil
}
