package opt

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/infastin/gorack/fastconv"
)

type NullBool[T ~bool] struct {
	Value T
	Valid bool
}

func NewNullBool[T ~bool](value T, valid bool) NullBool[T] {
	return NullBool[T]{
		Value: value,
		Valid: valid,
	}
}

func NullBoolFrom[T ~bool](value T) NullBool[T] {
	return NewNullBool(value, true)
}

func NullBoolFromPtr[T ~bool](value *T) NullBool[T] {
	if value == nil {
		return NewNullBool[T](false, false)
	}
	return NewNullBool(*value, true)
}

func (v NullBool[T]) Std() sql.NullBool {
	return sql.NullBool{
		Bool:  bool(v.Value),
		Valid: v.Valid,
	}
}

func (v *NullBool[T]) Set(value T) {
	v.Value = value
	v.Valid = true
}

func (v *NullBool[T]) Ptr() *T {
	if !v.Valid {
		return nil
	}
	return &v.Value
}

func (v NullBool[T]) IsZero() bool {
	return !v.Valid
}

func (v NullBool[T]) Get() (value T, ok bool) {
	return v.Value, v.Valid
}

func (v NullBool[T]) Or(value T) T {
	if !v.Valid {
		return value
	}
	return v.Value
}

func (v NullBool[T]) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`null`), nil
	}
	return strconv.AppendBool(nil, bool(v.Value)), nil
}

func (v *NullBool[T]) UnmarshalJSON(data []byte) error {
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

func (v NullBool[T]) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return strconv.AppendBool(nil, bool(v.Value)), nil
}

func (v *NullBool[T]) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		v.Valid = false
		return nil
	}

	value, err := strconv.ParseBool(fastconv.String(data))
	if err != nil {
		return err
	}

	v.Value = T(value)
	v.Valid = true

	return nil
}

type ZeroBool[T ~bool] struct {
	Value T
	Valid bool
}

func NewZeroBool[T ~bool](value T, valid bool) ZeroBool[T] {
	return ZeroBool[T]{
		Value: value,
		Valid: valid,
	}
}

func ZeroBoolFrom[T ~bool](value T) ZeroBool[T] {
	return NewZeroBool(value, true)
}

func ZeroBoolFromPtr[T ~bool](value *T) ZeroBool[T] {
	if value == nil {
		return NewZeroBool[T](false, false)
	}
	return NewZeroBool(*value, true)
}

func (v ZeroBool[T]) Std() sql.NullBool {
	return sql.NullBool{
		Bool:  bool(v.Value),
		Valid: v.Valid,
	}
}

func (v *ZeroBool[T]) Set(value T) {
	v.Value = value
	v.Valid = !bool(value)
}

func (v *ZeroBool[T]) Ptr() *T {
	if !v.Valid {
		return nil
	}
	return &v.Value
}

func (v ZeroBool[T]) IsZeroBool() bool {
	return !v.Valid
}

func (v ZeroBool[T]) Get() (value T, ok bool) {
	return v.Value, v.Valid
}

func (v ZeroBool[T]) Or(value T) T {
	if !v.Valid {
		return value
	}
	return v.Value
}

func (v ZeroBool[T]) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`null`), nil
	}
	return strconv.AppendBool(nil, bool(v.Value)), nil
}

func (v *ZeroBool[T]) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		v.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &v.Value); err != nil {
		return err
	}
	v.Valid = !bool(v.Value)

	return nil
}

func (v ZeroBool[T]) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return strconv.AppendBool(nil, bool(v.Value)), nil
}

func (v *ZeroBool[T]) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		v.Valid = false
		return nil
	}

	value, err := strconv.ParseBool(fastconv.String(data))
	if err != nil {
		return err
	}

	v.Value = T(value)
	v.Valid = !bool(value)

	return nil
}
