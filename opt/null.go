package opt

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type Null[T any] struct {
	V     T
	Valid bool
}

func NewNull[T any](value T, valid bool) Null[T] {
	return Null[T]{
		V:     value,
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
	return NullFrom(*value)
}

func NullFromFunc[T, U any](value *U, f func(U) T) Null[T] {
	if value == nil {
		var zero T
		return NewNull(zero, false)
	}
	return NullFrom(f(*value))
}

func NullFromFuncPtr[T, U any](value *U, f func(*U) T) Null[T] {
	if value == nil {
		var zero T
		return NewNull(zero, false)
	}
	return NullFrom(f(value))
}

func (v *Null[T]) Set(value T) {
	v.V, v.Valid = value, true
}

func (v *Null[T]) Reset() {
	var zero T
	v.V, v.Valid = zero, false
}

func (v *Null[T]) Ptr() *T {
	if !v.Valid {
		return nil
	}
	return &v.V
}

func (v Null[T]) IsZero() bool {
	return !v.Valid
}

func (v Null[T]) Get() (value T, ok bool) {
	return v.V, v.Valid
}

func (v Null[T]) ToSQL() sql.Null[T] {
	return sql.Null[T](v)
}

func (v Null[T]) Or(value T) T {
	if !v.Valid {
		return value
	}
	return v.V
}

func (v Null[T]) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(v.V)
}

func (v *Null[T]) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		var zero T
		v.V, v.Valid = zero, false
		return nil
	}

	if err := json.Unmarshal(data, &v.V); err != nil {
		return err
	}
	v.Valid = true

	return nil
}

func (v Null[T]) Value() (driver.Value, error) {
	return sql.Null[T](v).Value()
}

func (v *Null[T]) Scan(value any) error {
	return (*sql.Null[T])(v).Scan(value)
}
