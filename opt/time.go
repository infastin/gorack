package opt

import (
	"database/sql"
	"time"
)

type NullTime struct {
	Value time.Time
	Valid bool
}

func NewNullTime(value time.Time, valid bool) NullTime {
	return NullTime{
		Value: value,
		Valid: valid,
	}
}

func NullTimeFrom(value time.Time) NullTime {
	return NewNullTime(value, true)
}

func NullTimeFromPtr(value *time.Time) NullTime {
	if value == nil {
		return NewNullTime(time.Time{}, false)
	}
	return NullTimeFrom(*value)
}

func NullTimeFromFunc[U any](value U, f func(U) time.Time) NullTime {
	return NullTimeFrom(f(value))
}

func NullTimeFromPtrFunc[U any](value *U, f func(U) time.Time) NullTime {
	if value == nil {
		return NewNullTime(time.Time{}, false)
	}
	return NullTimeFrom(f(*value))
}

func (v NullTime) Std() sql.NullTime {
	return sql.NullTime{
		Time:  v.Value,
		Valid: v.Valid,
	}
}

func (v *NullTime) Set(value time.Time) {
	v.Value = value
	v.Valid = true
}

func (v *NullTime) Ptr() *time.Time {
	if !v.Valid {
		return nil
	}
	return &v.Value
}

func (v NullTime) IsZero() bool {
	return !v.Valid
}

func (v NullTime) Get() (value time.Time, ok bool) {
	return v.Value, v.Valid
}

func (v NullTime) Or(value time.Time) time.Time {
	if !v.Valid {
		return value
	}
	return v.Value
}

func (v NullTime) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`null`), nil
	}
	return v.Value.MarshalJSON()
}

func (v *NullTime) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		v.Valid = false
		return nil
	}

	if err := v.Value.UnmarshalJSON(data); err != nil {
		return err
	}
	v.Valid = true

	return nil
}

func (v NullTime) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return v.Value.MarshalText()
}

func (v *NullTime) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		v.Valid = false
		return nil
	}

	if err := v.Value.UnmarshalText(data); err != nil {
		return err
	}
	v.Valid = true

	return nil
}

type ZeroTime struct {
	Value time.Time
	Valid bool
}

func NewZeroTime(value time.Time, valid bool) ZeroTime {
	return ZeroTime{
		Value: value,
		Valid: valid,
	}
}

func ZeroTimeFrom(value time.Time) ZeroTime {
	return NewZeroTime(value, true)
}

func ZeroTimeFromPtr(value *time.Time) ZeroTime {
	if value == nil {
		return NewZeroTime(time.Time{}, false)
	}
	return NewZeroTime(*value, true)
}

func ZeroTimeFromFunc[U any](value U, f func(U) time.Time) ZeroTime {
	return ZeroTimeFrom(f(value))
}

func ZeroTimeFromPtrFunc[U any](value *U, f func(U) time.Time) ZeroTime {
	if value == nil {
		return NewZeroTime(time.Time{}, false)
	}
	return ZeroTimeFrom(f(*value))
}

func (v ZeroTime) Std() sql.NullTime {
	return sql.NullTime{
		Time:  v.Value,
		Valid: v.Valid,
	}
}

func (v *ZeroTime) Set(value time.Time) {
	v.Value = value
	v.Valid = !value.IsZero()
}

func (v *ZeroTime) Ptr() *time.Time {
	if !v.Valid {
		return nil
	}
	return &v.Value
}

func (v ZeroTime) IsZero() bool {
	return !v.Valid
}

func (v ZeroTime) Get() (value time.Time, ok bool) {
	return v.Value, v.Valid
}

func (v ZeroTime) Or(value time.Time) time.Time {
	if !v.Valid {
		return value
	}
	return v.Value
}

func (v ZeroTime) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`null`), nil
	}
	return v.Value.MarshalJSON()
}

func (v *ZeroTime) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		v.Valid = false
		return nil
	}

	if err := v.Value.UnmarshalJSON(data); err != nil {
		return err
	}
	v.Valid = !v.Value.IsZero()

	return nil
}

func (v ZeroTime) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return v.Value.MarshalText()
}

func (v *ZeroTime) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		v.Valid = false
		return nil
	}

	if err := v.Value.UnmarshalText(data); err != nil {
		return err
	}
	v.Valid = !v.Value.IsZero()

	return nil
}
