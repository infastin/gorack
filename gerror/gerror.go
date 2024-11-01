package gerror

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

var Trace = false

type Error struct {
	kind    error
	cause   error
	message string
	stack   []byte
}

type errorOption func(e *Error)

func WithKind(kind error) errorOption {
	return func(e *Error) {
		e.kind = kind
	}
}

func WithCause(cause error) errorOption {
	return func(e *Error) {
		e.cause = cause
	}
}

func WithMessage(msg string) errorOption {
	return func(e *Error) {
		e.message = msg
	}
}

func WithMessagef(format string, args ...any) errorOption {
	return func(e *Error) {
		e.message = fmt.Sprintf(format, args...)
	}
}

func WithStack() errorOption {
	return func(e *Error) {
		stack := make([]byte, 1<<16)
		stack = stack[:runtime.Stack(stack, false)]
		e.stack = stack
	}
}

func New(opts ...errorOption) error {
	e := &Error{}
	for _, opt := range opts {
		opt(e)
	}
	if Trace && e.stack == nil {
		stack := make([]byte, 1<<16)
		stack = stack[:runtime.Stack(stack, false)]
		e.stack = stack
	}
	return e
}

func Text(msg string) error {
	return New(WithMessage(msg))
}

func Textf(format string, args ...any) error {
	return Text(fmt.Sprintf(format, args...))
}

func (e *Error) Error() string {
	parts := []string{}
	if e.kind != nil {
		parts = append(parts, e.kind.Error())
	}
	if e.message != "" {
		parts = append(parts, e.message)
	}
	if e.cause != nil {
		parts = append(parts, e.cause.Error())
	}
	return strings.Join(parts, ": ")
}

func (e *Error) Kind() error {
	return e.kind
}

func (e *Error) Stack() []byte {
	return e.stack
}

func (e *Error) Unwrap() error {
	return e.cause
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As[T error](err error, target *T) bool {
	return errors.As(err, target)
}

func Into[T error](err error) (target T, ok bool) {
	ok = As(err, &target)
	return target, ok
}

func Of(err, kind error) bool {
	if err == nil || kind == nil {
		return err == kind
	}

	e, ok := Into[*Error](err)
	if !ok {
		return false
	}

	isComparable := reflect.TypeOf(kind).Comparable()
	if isComparable && e.kind == kind {
		return true
	}

	if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(kind) {
		return true
	}

	return false
}

func Wrap(err error, msg string) error {
	e, ok := Into[*Error](err)
	if !ok {
		return New(WithCause(err), WithMessage(msg))
	}

	if e.message == "" {
		e.message = msg
	} else {
		e.message = msg + ": " + e.message
	}

	if Trace && e.stack == nil {
		stack := make([]byte, 1<<16)
		stack = stack[:runtime.Stack(stack, false)]
		e.stack = stack
	}

	return e
}

func Wrapf(err error, format string, args ...any) error {
	return Wrap(err, fmt.Sprintf(format, args...))
}

func Update(err error, opts ...errorOption) error {
	e, ok := Into[*Error](err)
	if !ok {
		e = &Error{cause: err}
	}

	for _, opt := range opts {
		opt(e)
	}

	if Trace && e.stack == nil {
		stack := make([]byte, 1<<16)
		stack = stack[:runtime.Stack(stack, false)]
		e.stack = stack
	}

	return e
}
