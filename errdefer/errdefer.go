package errdefer

import (
	"context"
	"reflect"
	"unsafe"
)

var (
	fnType       = reflect.TypeFor[func()]()
	fnErrType    = reflect.TypeFor[func() error]()
	fnCtxType    = reflect.TypeFor[func(context.Context)]()
	fnCtxErrType = reflect.TypeFor[func(context.Context) error]()
)

type CloseFunc interface {
	~func() | ~func() error | ~func(context.Context) | ~func(context.Context) error
}

// Close calls fn when *err is not nil.
// Returns an error returned from fn.
func Close[F CloseFunc](err *error, closeFn F) error {
	if *err == nil {
		return nil
	}
	switch fn := any(closeFn).(type) {
	case func():
		fn()
	case func() error:
		if err := fn(); err != nil {
			return err
		}
	case func(context.Context):
		fn(context.Background())
	case func(context.Context) error:
		if err := fn(context.Background()); err != nil {
			return err
		}
	default:
		switch typ := reflect.TypeFor[F](); {
		case typ.ConvertibleTo(fnType):
			fn := *(*func())(unsafe.Pointer(&closeFn))
			fn()
		case typ.ConvertibleTo(fnErrType):
			fn := *(*func() error)(unsafe.Pointer(&closeFn))
			if err := fn(); err != nil {
				return err
			}
		case typ.ConvertibleTo(fnCtxType):
			fn := *(*func(context.Context))(unsafe.Pointer(&closeFn))
			fn(context.Background())
		case typ.ConvertibleTo(fnCtxErrType):
			fn := *(*func(context.Context) error)(unsafe.Pointer(&closeFn))
			if err := fn(context.Background()); err != nil {
				return err
			}
		}
	}
	return nil
}

type CloseContextFunc interface {
	~func(context.Context) | ~func(context.Context) error
}

// CloseContext calls fn with provided Context when *err is not nil.
// Returns an error returned from fn.
func CloseContext[F CloseContextFunc](ctx context.Context, err *error, closeFn F) error {
	if *err == nil {
		return nil
	}
	switch fn := any(closeFn).(type) {
	case func(context.Context):
		fn(ctx)
	case func(context.Context) error:
		if err := fn(ctx); err != nil {
			return err
		}
	default:
		switch typ := reflect.TypeFor[F](); {
		case typ.ConvertibleTo(fnCtxType):
			fn := *(*func(context.Context))(unsafe.Pointer(&closeFn))
			fn(ctx)
		case typ.ConvertibleTo(fnCtxErrType):
			fn := *(*func(context.Context) error)(unsafe.Pointer(&closeFn))
			if err := fn(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}
