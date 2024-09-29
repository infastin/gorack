package validation

type notNilPtrRule[T any] struct {
	condition bool
}

var ErrNotNil = NewRuleError("not_nil", "is required")

func NotNilPtr[T any](condition bool) notNilPtrRule[T] {
	return notNilPtrRule[T]{
		condition: condition,
	}
}

func (r notNilPtrRule[T]) Validate(p *T) error {
	if r.condition && p == nil {
		return ErrNotNil
	}
	return nil
}

type notNilSliceRule[T any] struct {
	condition bool
}

func NotNilSlice[T any](condition bool) notNilSliceRule[T] {
	return notNilSliceRule[T]{
		condition: condition,
	}
}

func (r notNilSliceRule[T]) Validate(s []T) error {
	if r.condition && s == nil {
		return ErrNotNil
	}
	return nil
}

type notNilMapRule[T any] struct {
	condition bool
}

func NotNilMap[T any](condition bool) notNilMapRule[T] {
	return notNilMapRule[T]{
		condition: condition,
	}
}

func (r notNilMapRule[T]) Validate(m map[string]T) error {
	if r.condition && m == nil {
		return ErrNotNil
	}
	return nil
}
