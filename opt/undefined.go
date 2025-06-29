package opt

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type Undefined[T any] struct {
	V     T
	Valid bool
}

func NewUndefined[T any](value T, valid bool) Undefined[T] {
	return Undefined[T]{
		V:     value,
		Valid: valid,
	}
}

func UndefinedFrom[T any](value T) Undefined[T] {
	return NewUndefined(value, true)
}

func UndefinedFromPtr[T any](value *T) Undefined[T] {
	if value == nil {
		var zero T
		return UndefinedFrom(zero)
	}
	return UndefinedFrom(*value)
}

func UndefinedFromFunc[T, U any](value *U, f func(U) T) Undefined[T] {
	if value == nil {
		var zero U
		return UndefinedFrom(f(zero))
	}
	return UndefinedFrom(f(*value))
}

func UndefinedFromFuncPtr[T, U any](value *U, f func(*U) T) Undefined[T] {
	if value == nil {
		var zero U
		return UndefinedFrom(f(&zero))
	}
	return UndefinedFrom(f(value))
}

func (v *Undefined[T]) Set(value T) {
	v.V, v.Valid = value, true
}

func (v *Undefined[T]) Reset() {
	var zero T
	v.V, v.Valid = zero, false
}

func (v *Undefined[T]) Ptr() *T {
	if !v.Valid {
		return nil
	}
	return &v.V
}

func (v Undefined[T]) IsZero() bool {
	return !v.Valid
}

func (v Undefined[T]) Get() (value T, ok bool) {
	return v.V, v.Valid
}

func (v Undefined[T]) ToSQL() sql.Null[T] {
	return sql.Null[T](v)
}

func (v Undefined[T]) Or(value T) T {
	if !v.Valid {
		return value
	}
	return v.V
}

func (v Undefined[T]) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(v.V)
}

func (v *Undefined[T]) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		var zero T
		v.V, v.Valid = zero, true
		return nil
	}

	if err := json.Unmarshal(data, &v.V); err != nil {
		return err
	}
	v.Valid = true

	return nil
}

func (v Undefined[T]) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return marshalText(v.V)
}

func (v *Undefined[T]) UnmarshalText(data []byte) error {
	if err := unmarshalText(&v.V, data); err != nil {
		return err
	}
	v.Valid = true
	return nil
}

func (v Undefined[T]) Value() (driver.Value, error) {
	return sql.Null[T](v).Value()
}

func (v *Undefined[T]) Scan(value any) error {
	if value == nil {
		var zero T
		v.V, v.Valid = zero, true
		return nil
	}
	return (*sql.Null[T])(v).Scan(value)
}
