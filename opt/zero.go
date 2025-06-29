package opt

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type Zero[T any] struct {
	V     T
	Valid bool
}

func NewZero[T any](value T, valid bool) Zero[T] {
	return Zero[T]{
		V:     value,
		Valid: valid,
	}
}

func ZeroFrom[T any](value T) Zero[T] {
	return NewZero(value, !isZero(value))
}

func ZeroFromPtr[T any](value *T) Zero[T] {
	if value == nil {
		var zero T
		return NewZero(zero, false)
	}
	return ZeroFrom(*value)
}

func ZeroFromFunc[T, U any](value *U, f func(U) T) Zero[T] {
	if value == nil {
		var zero T
		return NewZero(zero, false)
	}
	return ZeroFrom(f(*value))
}

func ZeroFromFuncPtr[T, U any](value *U, f func(*U) T) Zero[T] {
	if value == nil {
		var zero T
		return NewZero(zero, false)
	}
	return ZeroFrom(f(value))
}

func (v *Zero[T]) Set(value T) {
	v.V, v.Valid = value, !isZero(value)
}

func (v *Zero[T]) Reset() {
	var zero T
	v.V, v.Valid = zero, false
}

func (v *Zero[T]) Ptr() *T {
	if !v.Valid {
		return nil
	}
	return &v.V
}

func (v Zero[T]) IsZero() bool {
	return !v.Valid
}

func (v Zero[T]) Get() (value T, ok bool) {
	return v.V, v.Valid
}

func (v Zero[T]) ToSQL() sql.Null[T] {
	return sql.Null[T](v)
}

func (v Zero[T]) Or(value T) T {
	if !v.Valid {
		return value
	}
	return v.V
}

func (v Zero[T]) MarshalJSON() ([]byte, error) {
	var value T
	if v.Valid {
		value = v.V
	}
	return json.Marshal(value)
}

func (v *Zero[T]) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		var zero T
		v.V, v.Valid = zero, false
		return nil
	}

	if err := json.Unmarshal(data, &v.V); err != nil {
		return err
	}
	v.Valid = !isZero(v.V)

	return nil
}

func (v Zero[T]) MarshalText() ([]byte, error) {
	var value T
	if v.Valid {
		value = v.V
	}
	return marshalText(value)
}

func (v *Zero[T]) UnmarshalText(data []byte) error {
	if err := unmarshalText(&v.V, data); err != nil {
		return err
	}
	v.Valid = !isZero(v.V)
	return nil
}

func (v Zero[T]) Value() (driver.Value, error) {
	return sql.Null[T](v).Value()
}

func (v *Zero[T]) Scan(value any) error {
	if err := (*sql.Null[T])(v).Scan(value); err != nil {
		return err
	}
	v.Valid = !isZero(v.V)
	return nil
}
