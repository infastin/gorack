package opt

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/infastin/gorack/constraints"
	"github.com/infastin/gorack/fastconv"
)

type NullFloat[T constraints.Float] struct {
	Value T
	Valid bool
}

func NewNullFloat[T constraints.Float](value T, valid bool) NullFloat[T] {
	return NullFloat[T]{
		Value: value,
		Valid: valid,
	}
}

func NullFloatFrom[T constraints.Float](value T) NullFloat[T] {
	return NewNullFloat(value, true)
}

func NullFloatFromPtr[T constraints.Float](value *T) NullFloat[T] {
	if value == nil {
		return NewNullFloat[T](0, false)
	}
	return NullFloatFrom(*value)
}

func NullFloatFromFunc[T constraints.Float, U any](value U, f func(U) T) NullFloat[T] {
	return NullFloatFrom(f(value))
}

func NullFloatFromPtrFunc[T constraints.Float, U any](value *U, f func(U) T) NullFloat[T] {
	if value == nil {
		return NewNullFloat[T](0, false)
	}
	return NullFloatFrom(f(*value))
}

func (v NullFloat[T]) Std() sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: float64(v.Value),
		Valid:   v.Valid,
	}
}

func (v *NullFloat[T]) Set(value T) {
	v.Value = value
	v.Valid = true
}

func (v *NullFloat[T]) Ptr() *T {
	if !v.Valid {
		return nil
	}
	return &v.Value
}

func (v NullFloat[T]) IsZero() bool {
	return !v.Valid
}

func (v NullFloat[T]) Get() (value T, ok bool) {
	return v.Value, v.Valid
}

func (v NullFloat[T]) Or(value T) T {
	if !v.Valid {
		return value
	}
	return v.Value
}

func (v NullFloat[T]) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`null`), nil
	}
	return strconv.AppendFloat(nil, float64(v.Value), 'f', -1, 64), nil
}

func (v *NullFloat[T]) UnmarshalJSON(data []byte) error {
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

func (v NullFloat[T]) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return strconv.AppendFloat(nil, float64(v.Value), 'f', -1, 64), nil
}

func (v *NullFloat[T]) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		v.Valid = false
		return nil
	}

	value, err := strconv.ParseFloat(fastconv.String(data), 64)
	if err != nil {
		return err
	}

	v.Value = T(value)
	v.Valid = true

	return nil
}

type ZeroFloat[T constraints.Float] struct {
	Value T
	Valid bool
}

func NewZeroFloat[T constraints.Float](value T, valid bool) ZeroFloat[T] {
	return ZeroFloat[T]{
		Value: value,
		Valid: valid,
	}
}

func ZeroFloatFrom[T constraints.Float](value T) ZeroFloat[T] {
	return NewZeroFloat(value, value != 0)
}

func ZeroFloatFromPtr[T constraints.Float](value *T) ZeroFloat[T] {
	if value == nil {
		return NewZeroFloat[T](0, false)
	}
	return ZeroFloatFrom(*value)
}

func ZeroFloatFromFunc[T constraints.Float, U any](value U, f func(U) T) ZeroFloat[T] {
	return ZeroFloatFrom(f(value))
}

func ZeroFloatFromPtrFunc[T constraints.Float, U any](value *U, f func(U) T) ZeroFloat[T] {
	if value == nil {
		return NewZeroFloat[T](0, false)
	}
	return ZeroFloatFrom(f(*value))
}

func (v ZeroFloat[T]) Std() sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: float64(v.Value),
		Valid:   v.Valid,
	}
}

func (v *ZeroFloat[T]) Set(value T) {
	v.Value = value
	v.Valid = value != 0
}

func (v *ZeroFloat[T]) Ptr() *T {
	if !v.Valid {
		return nil
	}
	return &v.Value
}

func (v ZeroFloat[T]) IsZero() bool {
	return !v.Valid
}

func (v ZeroFloat[T]) Get() (value T, ok bool) {
	return v.Value, v.Valid
}

func (v ZeroFloat[T]) Or(value T) T {
	if !v.Valid {
		return value
	}
	return v.Value
}

func (v ZeroFloat[T]) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`null`), nil
	}
	return strconv.AppendFloat(nil, float64(v.Value), 'f', -1, 64), nil
}

func (v *ZeroFloat[T]) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		v.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &v.Value); err != nil {
		return err
	}
	v.Valid = v.Value != 0

	return nil
}

func (v ZeroFloat[T]) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return strconv.AppendFloat(nil, float64(v.Value), 'f', -1, 64), nil
}

func (v *ZeroFloat[T]) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		v.Valid = false
		return nil
	}

	value, err := strconv.ParseFloat(fastconv.String(data), 64)
	if err != nil {
		return err
	}

	v.Value = T(value)
	v.Valid = value != 0

	return nil
}
