package textopt

import (
	"database/sql"
	"database/sql/driver"

	"github.com/infastin/gorack/opt/v2"
	"github.com/infastin/gorack/opt/v2/internal"
)

type Undefined[T any] opt.Undefined[T]

func NewUndefined[T any](value T, valid bool) Undefined[T] {
	return Undefined[T](opt.NewUndefined(value, valid))
}

func UndefinedFrom[T any](value T) Undefined[T] {
	return Undefined[T](opt.UndefinedFrom(value))
}

func UndefinedFromPtr[T any](value *T) Undefined[T] {
	return Undefined[T](opt.UndefinedFromPtr(value))
}

func UndefinedFromFunc[T, U any](value *U, f func(U) T) Undefined[T] {
	return Undefined[T](opt.UndefinedFromFunc(value, f))
}

func UndefinedFromFuncPtr[T, U any](value *U, f func(*U) T) Undefined[T] {
	return Undefined[T](opt.UndefinedFromFuncPtr(value, f))
}

func (v *Undefined[T]) Set(value T) {
	(*opt.Undefined[T])(v).Set(value)
}

func (v *Undefined[T]) Reset() {
	(*opt.Undefined[T])(v).Reset()
}

func (v *Undefined[T]) Ptr() *T {
	return (*opt.Undefined[T])(v).Ptr()
}

func (v Undefined[T]) IsZero() bool {
	return opt.Undefined[T](v).IsZero()
}

func (v Undefined[T]) Get() (value T, ok bool) {
	return opt.Undefined[T](v).Get()
}

func (v Undefined[T]) ToSQL() sql.Null[T] {
	return opt.Undefined[T](v).ToSQL()
}

func (v Undefined[T]) Or(value T) T {
	return opt.Undefined[T](v).Or(value)
}

func (v Undefined[T]) MarshalJSON() ([]byte, error) {
	return opt.Undefined[T](v).MarshalJSON()
}

func (v *Undefined[T]) UnmarshalJSON(data []byte) error {
	return (*opt.Undefined[T])(v).UnmarshalJSON(data)
}

func (v Undefined[T]) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return internal.MarshalText(v.V)
}

func (v *Undefined[T]) UnmarshalText(data []byte) error {
	if err := internal.UnmarshalText(&v.V, data); err != nil {
		return err
	}
	v.Valid = true
	return nil
}

func (v Undefined[T]) Value() (driver.Value, error) {
	return opt.Undefined[T](v).Value()
}

func (v *Undefined[T]) Scan(value any) error {
	return (*opt.Undefined[T])(v).Scan(value)
}
