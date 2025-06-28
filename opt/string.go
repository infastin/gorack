package opt

import (
	"database/sql"
	"encoding/json"
)

type NullString[T ~string] struct {
	Value T
	Valid bool
}

func NewNullString[T ~string](value T, valid bool) NullString[T] {
	return NullString[T]{
		Value: value,
		Valid: valid,
	}
}

func NullStringFrom[T ~string](value T) NullString[T] {
	return NewNullString(value, true)
}

func NullStringFromPtr[T ~string](value *T) NullString[T] {
	if value == nil {
		return NewNullString[T]("", false)
	}
	return NullStringFrom(*value)
}

func NullStringFromFunc[T ~string, U any](value *U, f func(U) T) NullString[T] {
	if value == nil {
		return NewNullString[T]("", false)
	}
	return NullStringFrom(f(*value))
}

func NullStringFromFuncPtr[T ~string, U any](value *U, f func(*U) T) NullString[T] {
	if value == nil {
		return NewNullString[T]("", false)
	}
	return NullStringFrom(f(value))
}

func (v NullString[T]) Std() sql.NullString {
	return sql.NullString{
		String: string(v.Value),
		Valid:  v.Valid,
	}
}

func (v *NullString[T]) Set(value T) {
	v.Value = value
	v.Valid = true
}

func (v *NullString[T]) Ptr() *T {
	if !v.Valid {
		return nil
	}
	return &v.Value
}

func (v NullString[T]) IsZero() bool {
	return !v.Valid
}

func (v NullString[T]) Get() (value T, ok bool) {
	return v.Value, v.Valid
}

func (v NullString[T]) Or(value T) T {
	if !v.Valid {
		return value
	}
	return v.Value
}

func (v NullString[T]) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`null`), nil
	}
	return []byte(v.Value), nil
}

func (v *NullString[T]) UnmarshalJSON(data []byte) error {
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

func (v NullString[T]) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return []byte(v.Value), nil
}

func (v *NullString[T]) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		v.Valid = false
		return nil
	}

	v.Value = T(data)
	v.Valid = true

	return nil
}

type ZeroString[T ~string] struct {
	Value T
	Valid bool
}

func NewZeroString[T ~string](value T, valid bool) ZeroString[T] {
	return ZeroString[T]{
		Value: value,
		Valid: valid,
	}
}

func ZeroStringFrom[T ~string](value T) ZeroString[T] {
	return NewZeroString(value, value != "")
}

func ZeroStringFromPtr[T ~string](value *T) ZeroString[T] {
	if value == nil {
		return NewZeroString[T]("", false)
	}
	return ZeroStringFrom(*value)
}

func ZeroStringFromFunc[T ~string, U any](value *U, f func(U) T) ZeroString[T] {
	if value == nil {
		return NewZeroString[T]("", false)
	}
	return ZeroStringFrom(f(*value))
}

func ZeroStringFromFuncPtr[T ~string, U any](value *U, f func(*U) T) ZeroString[T] {
	if value == nil {
		return NewZeroString[T]("", false)
	}
	return ZeroStringFrom(f(value))
}

func (v ZeroString[T]) Std() sql.NullString {
	return sql.NullString{
		String: string(v.Value),
		Valid:  v.Valid,
	}
}

func (v *ZeroString[T]) Set(value T) {
	v.Value = value
	v.Valid = value != ""
}

func (v *ZeroString[T]) Ptr() *T {
	if !v.Valid {
		return nil
	}
	return &v.Value
}

func (v ZeroString[T]) IsZero() bool {
	return !v.Valid
}

func (v ZeroString[T]) Get() (value T, ok bool) {
	return v.Value, v.Valid
}

func (v ZeroString[T]) Or(value T) T {
	if !v.Valid {
		return value
	}
	return v.Value
}

func (v ZeroString[T]) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`null`), nil
	}
	return []byte(v.Value), nil
}

func (v *ZeroString[T]) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		v.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &v.Value); err != nil {
		return err
	}
	v.Valid = v.Value != ""

	return nil
}

func (v ZeroString[T]) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return []byte(v.Value), nil
}

func (v *ZeroString[T]) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		v.Valid = false
		return nil
	}

	v.Value = T(data)
	v.Valid = len(data) != 0

	return nil
}
