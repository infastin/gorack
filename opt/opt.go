package opt

import "github.com/infastin/gorack/opt/v2/internal"

type getter[T any] interface {
	Get() (value T, ok bool)
}

func Convert[T, U any, Opt getter[T]](opt Opt, f func(T) U) U {
	value, ok := opt.Get()
	if !ok {
		var zero U
		return zero
	}
	return f(value)
}

func Ptr[T any](value T) *T {
	ptr := new(T)
	*ptr = value
	return ptr
}

func ZeroPtr[T any](value T) *T {
	if internal.IsZero(value) {
		return nil
	}
	return Ptr(value)
}

func ConvertPtr[T, U any](ptr *T, fn func(T) U) *U {
	if ptr == nil {
		return nil
	}
	return Ptr(fn(*ptr))
}

func Deref[T any](ptr *T, def T) T {
	if ptr == nil {
		return def
	}
	return *ptr
}
