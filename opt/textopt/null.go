package textopt

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"unsafe"

	"github.com/infastin/gorack/opt/v2"
	"github.com/infastin/gorack/opt/v2/internal"
)

type Null[T any] opt.Null[T]

func NewNull[T any](value T, valid bool) Null[T] {
	return Null[T](opt.NewNull(value, valid))
}

func NullFrom[T any](value T) Null[T] {
	return Null[T](opt.NullFrom(value))
}

func NullFromPtr[T any](value *T) Null[T] {
	return Null[T](opt.NullFromPtr(value))
}

func NullFromFunc[T, U any](value *U, f func(U) T) Null[T] {
	return Null[T](opt.NullFromFunc(value, f))
}

func NullFromFuncPtr[T, U any](value *U, f func(*U) T) Null[T] {
	return Null[T](opt.NullFromFuncPtr(value, f))
}

func (v *Null[T]) Set(value T) {
	(*opt.Null[T])(v).Set(value)
}

func (v *Null[T]) Reset() {
	(*opt.Null[T])(v).Reset()
}

func (v *Null[T]) Ptr() *T {
	return (*opt.Null[T])(v).Ptr()
}

func (v Null[T]) IsZero() bool {
	return opt.Null[T](v).IsZero()
}

func (v Null[T]) Get() (value T, ok bool) {
	return opt.Null[T](v).Get()
}

func (v Null[T]) ToSQL() sql.Null[T] {
	return opt.Null[T](v).ToSQL()
}

func (v Null[T]) Or(value T) T {
	return opt.Null[T](v).Or(value)
}

func (v Null[T]) MarshalJSON() ([]byte, error) {
	return opt.Null[T](v).MarshalJSON()
}

func (v *Null[T]) UnmarshalJSON(data []byte) error {
	return (*opt.Null[T])(v).UnmarshalJSON(data)
}

func (v Null[T]) MarshalText() ([]byte, error) {
	if !v.Valid {
		return []byte{}, nil
	}
	return internal.MarshalText(v.V)
}

func (v *Null[T]) UnmarshalText(data []byte) error {
	if reflect.TypeFor[T]().Kind() == reflect.String {
		*(*string)(unsafe.Pointer(&v.V)) = string(data)
		v.Valid = true
		return nil
	}

	if internal.IsNullText(data) {
		var zero T
		v.V, v.Valid = zero, false
		return nil
	}

	if err := internal.UnmarshalText(&v.V, data); err != nil {
		return err
	}
	v.Valid = true

	return nil
}

func (v Null[T]) Value() (driver.Value, error) {
	return opt.Null[T](v).Value()
}

func (v *Null[T]) Scan(value any) error {
	return (*opt.Null[T])(v).Scan(value)
}
