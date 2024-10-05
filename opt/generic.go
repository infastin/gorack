package opt

import (
	"database/sql"
	"encoding/json"
)

type Null[T any] struct {
	Value T
	Valid bool
}

func NewNull[T any](value T, valid bool) Null[T] {
	return Null[T]{
		Value: value,
		Valid: valid,
	}
}

func NullFrom[T any](value T) Null[T] {
	return NewNull(value, true)
}

func NullFromPtr[T any](value *T) Null[T] {
	if value == nil {
		var zero T
		return NewNull(zero, false)
	}
	return NewNull(*value, true)
}

func (v *Null[T]) Set(value T) {
	v.Value = value
	v.Valid = true
}

func (v *Null[T]) Ptr() *T {
	if !v.Valid {
		return nil
	}
	return &v.Value
}

func (v Null[T]) IsZero() bool {
	return !v.Valid
}

func (v Null[T]) Get() (value T, ok bool) {
	return v.Value, v.Valid
}

func (v Null[T]) Std() sql.Null[T] {
	return sql.Null[T]{
		V:     v.Value,
		Valid: v.Valid,
	}
}

func (v Null[T]) Or(value T) T {
	if !v.Valid {
		return value
	}
	return v.Value
}

func (v Null[T]) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(v.Value)
}

func (v *Null[T]) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		v.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &v.Value); err != nil {
		return err
	}
	v.Valid = true

	return nil
}

type Zero[T comparable] struct {
	Value T
	Valid bool
}

func NewZero[T comparable](value T, valid bool) Zero[T] {
	return Zero[T]{
		Value: value,
		Valid: valid,
	}
}

func ZeroFrom[T comparable](value T) Zero[T] {
	var zero T
	return NewZero(zero, value == zero)
}

func ZeroFromPtr[T comparable](value *T) Zero[T] {
	var zero T
	if value == nil || *value == zero {
		return NewZero(zero, false)
	}
	return NewZero(*value, true)
}

func (v Zero[T]) Std() sql.Null[T] {
	return sql.Null[T]{
		V:     v.Value,
		Valid: v.Valid,
	}
}

func (v *Zero[T]) Set(value T) {
	var zero T
	v.Value = value
	v.Valid = value != zero
}

func (v *Zero[T]) Ptr() *T {
	if !v.Valid {
		return nil
	}
	return &v.Value
}

func (v Zero[T]) IsZero() bool {
	return !v.Valid
}

func (v Zero[T]) Get() (value T, ok bool) {
	return v.Value, v.Valid
}

func (v Zero[T]) Or(value T) T {
	if !v.Valid {
		return value
	}
	return v.Value
}

func (v Zero[T]) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(v.Value)
}

func (v *Zero[T]) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		v.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &v.Value); err != nil {
		return err
	}

	var zero T
	v.Valid = v.Value != zero

	return nil
}
