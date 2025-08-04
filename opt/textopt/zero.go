package textopt

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"unsafe"

	"github.com/infastin/gorack/opt/v2"
	"github.com/infastin/gorack/opt/v2/internal"
)

type Zero[T any] opt.Zero[T]

func NewZero[T any](value T, valid bool) Zero[T] {
	return Zero[T](opt.NewZero(value, valid))
}

func ZeroFrom[T any](value T) Zero[T] {
	return Zero[T](opt.ZeroFrom(value))
}

func ZeroFromPtr[T any](value *T) Zero[T] {
	return Zero[T](opt.ZeroFromPtr(value))
}

func ZeroFromFunc[T, U any](value *U, f func(U) T) Zero[T] {
	return Zero[T](opt.ZeroFromFunc(value, f))
}

func ZeroFromFuncPtr[T, U any](value *U, f func(*U) T) Zero[T] {
	return Zero[T](opt.ZeroFromFuncPtr(value, f))
}

func (v *Zero[T]) Set(value T) {
	(*opt.Zero[T])(v).Set(value)
}

func (v *Zero[T]) Reset() {
	(*opt.Zero[T])(v).Reset()
}

func (v *Zero[T]) Ptr() *T {
	return (*opt.Zero[T])(v).Ptr()
}

func (v Zero[T]) IsZero() bool {
	return opt.Zero[T](v).IsZero()
}

func (v Zero[T]) Get() (value T, ok bool) {
	return opt.Zero[T](v).Get()
}

func (v Zero[T]) ToSQL() sql.Null[T] {
	return opt.Zero[T](v).ToSQL()
}

func (v Zero[T]) Or(value T) T {
	return opt.Zero[T](v).Or(value)
}

func (v Zero[T]) MarshalJSON() ([]byte, error) {
	return opt.Zero[T](v).MarshalJSON()
}

func (v *Zero[T]) UnmarshalJSON(data []byte) error {
	return (*opt.Zero[T])(v).UnmarshalJSON(data)
}

func (v Zero[T]) MarshalText() ([]byte, error) {
	var value T
	if v.Valid {
		value = v.V
	}
	return internal.MarshalText(value)
}

func (v *Zero[T]) UnmarshalText(data []byte) error {
	if reflect.TypeFor[T]().Kind() == reflect.String {
		*(*string)(unsafe.Pointer(&v.V)) = string(data)
		v.Valid = !internal.IsZero(v.V)
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
	v.Valid = !internal.IsZero(v.V)

	return nil
}

func (v Zero[T]) Value() (driver.Value, error) {
	return opt.Zero[T](v).Value()
}

func (v *Zero[T]) Scan(value any) error {
	return (*opt.Zero[T])(v).Scan(value)
}
