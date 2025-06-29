package errdefer

import "context"

type CloseFunc interface {
	func() | func() error | func(context.Context) | func(context.Context) error
}

// Close calls fn when *err is not nil.
// Returns an error returned from fn.
func Close[F CloseFunc](err *error, fn F) error {
	if *err == nil {
		return nil
	}
	switch fn := any(fn).(type) {
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
	}
	return nil
}

type CloseContextFunc interface {
	func(context.Context) | func(context.Context) error
}

// CloseContext calls fn with provided Context when *err is not nil.
// Returns an error returned from fn.
func CloseContext[F CloseContextFunc](ctx context.Context, err *error, fn F) error {
	if *err == nil {
		return nil
	}
	switch fn := any(fn).(type) {
	case func(context.Context):
		fn(ctx)
	case func(context.Context) error:
		if err := fn(ctx); err != nil {
			return nil
		}
	}
	return nil
}
