package opt

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/infastin/gorack/constraints"
	"github.com/infastin/gorack/fastconv"
)

type NullInt[T constraints.Int] struct {
	Value T
	Valid bool
}

func NewNullInt[T constraints.Int](value T, valid bool) NullInt[T] {
	return NullInt[T]{
		Value: value,
		Valid: valid,
	}
}

func NullIntFrom[T constraints.Int](value T) NullInt[T] {
	return NewNullInt(value, true)
}

func NullIntFromPtr[T constraints.Int](value *T) NullInt[T] {
	if value == nil {
		return NewNullInt[T](0, false)
	}
	return NullIntFrom(*value)
}

func NullIntFromFunc[T constraints.Int, U any](value *U, f func(U) T) NullInt[T] {
	if value == nil {
		return NewNullInt[T](0, false)
	}
	return NullIntFrom(f(*value))
}

func NullIntFromFuncPtr[T constraints.Int, U any](value *U, f func(*U) T) NullInt[T] {
	if value == nil {
		return NewNullInt[T](0, false)
	}
	return NullIntFrom(f(value))
}

func (v NullInt[T]) Std16() sql.NullInt16 {
	return sql.NullInt16{
		Int16: int16(v.Value),
		Valid: v.Valid,
	}
}

func (v NullInt[T]) Std32() sql.NullInt32 {
	return sql.NullInt32{
		Int32: int32(v.Value),
		Valid: v.Valid,
	}
}

func (v NullInt[T]) Std64() sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(v.Value),
		Valid: v.Valid,
	}
}

func (v NullInt[T]) Std() sql.NullInt64 {
	return v.Std64()
}

func (v *NullInt[T]) Set(value T) {
	v.Value = value
	v.Valid = true
}

func (v *NullInt[T]) Ptr() *T {
	if !v.Valid {
		return nil
	}
	return &v.Value
}

func (v NullInt[T]) IsZero() bool {
	return !v.Valid
}

func (v NullInt[T]) Get() (value T, ok bool) {
	return v.Value, v.Valid
}

func (v NullInt[T]) Or(value T) T {
	if !v.Valid {
		return value
	}
	return v.Value
}

func (v NullInt[T]) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`null`), nil
	}
	return strconv.AppendInt(nil, int64(v.Value), 10), nil
}

func (v *NullInt[T]) UnmarshalJSON(data []byte) error {
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

func (v NullInt[T]) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return strconv.AppendInt(nil, int64(v.Value), 10), nil
}

func (v *NullInt[T]) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		v.Valid = false
		return nil
	}

	value, err := strconv.ParseInt(fastconv.String(data), 10, 64)
	if err != nil {
		return err
	}

	v.Value = T(value)
	v.Valid = true

	return nil
}

type ZeroInt[T constraints.Int] struct {
	Value T
	Valid bool
}

func NewZeroInt[T constraints.Int](value T, valid bool) ZeroInt[T] {
	return ZeroInt[T]{
		Value: value,
		Valid: valid,
	}
}

func ZeroIntFrom[T constraints.Int](value T) ZeroInt[T] {
	return NewZeroInt(value, true)
}

func ZeroIntFromPtr[T constraints.Int](value *T) ZeroInt[T] {
	if value == nil {
		return NewZeroInt[T](0, false)
	}
	return ZeroIntFrom(*value)
}

func ZeroIntFromFunc[T constraints.Int, U any](value *U, f func(U) T) ZeroInt[T] {
	if value == nil {
		return NewZeroInt[T](0, false)
	}
	return ZeroIntFrom(f(*value))
}

func ZeroIntFromFuncPtr[T constraints.Int, U any](value *U, f func(*U) T) ZeroInt[T] {
	if value == nil {
		return NewZeroInt[T](0, false)
	}
	return ZeroIntFrom(f(value))
}

func (v ZeroInt[T]) Std16() sql.NullInt16 {
	return sql.NullInt16{
		Int16: int16(v.Value),
		Valid: v.Valid,
	}
}

func (v ZeroInt[T]) Std32() sql.NullInt32 {
	return sql.NullInt32{
		Int32: int32(v.Value),
		Valid: v.Valid,
	}
}

func (v ZeroInt[T]) Std64() sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(v.Value),
		Valid: v.Valid,
	}
}

func (v ZeroInt[T]) Std() sql.NullInt64 {
	return v.Std64()
}

func (v *ZeroInt[T]) Set(value T) {
	v.Value = value
	v.Valid = value != 0
}

func (v *ZeroInt[T]) Ptr() *T {
	if !v.Valid {
		return nil
	}
	return &v.Value
}

func (v ZeroInt[T]) IsZero() bool {
	return !v.Valid
}

func (v ZeroInt[T]) Get() (value T, ok bool) {
	return v.Value, v.Valid
}

func (v ZeroInt[T]) Or(value T) T {
	if !v.Valid {
		return value
	}
	return v.Value
}

func (v ZeroInt[T]) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`null`), nil
	}
	return strconv.AppendInt(nil, int64(v.Value), 10), nil
}

func (v *ZeroInt[T]) UnmarshalJSON(data []byte) error {
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

func (v ZeroInt[T]) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return strconv.AppendInt(nil, int64(v.Value), 10), nil
}

func (v *ZeroInt[T]) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		v.Valid = false
		return nil
	}

	value, err := strconv.ParseInt(fastconv.String(data), 10, 64)
	if err != nil {
		return err
	}

	v.Value = T(value)
	v.Valid = value != 0

	return nil
}
