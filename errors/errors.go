package errors

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

type Error struct {
	kind       error
	cause      error
	message    string
	stackTrace StackTrace
}

func New(s string) error {
	return &Error{
		message:    s,
		stackTrace: callers(3),
	}
}

func Errorf(format string, args ...any) error {
	return &Error{
		message:    fmt.Sprintf(format, args...),
		stackTrace: callers(3),
	}
}

func Wrap(err error, msg string) error {
	e, ok := Into[*Error](err)
	if !ok {
		e = &Error{cause: err}
	}

	if e.message == "" {
		e.message = msg
	} else {
		e.message = msg + ": " + e.message
	}

	if e.stackTrace == nil {
		e.stackTrace = callers(3)
	}

	return e
}

func Wrapf(err error, format string, args ...any) error {
	e, ok := Into[*Error](err)
	if !ok {
		e = &Error{cause: err}
	}

	msg := fmt.Sprintf(format, args...)
	if e.message == "" {
		e.message = msg
	} else {
		e.message = msg + ": " + e.message
	}

	if e.stackTrace == nil {
		e.stackTrace = callers(3)
	}

	return e
}

// Creates a new error with a given kind.
// List of valid function signatures:
//
//	Pack(kind error)
//	Pack(kind error, cause error)
//	Pack(kind error, cause error, s string)
//	Pack(kind error, cause error, format string, args ...any)
//	Pack(kind error, s string)
//	Pack(kind error, format string string, args ...any)
func Pack(kind error, args ...any) error {
	e := &Error{kind: kind, stackTrace: callers(3)}

	if len(args) >= 1 {
		if arg, ok := args[0].(error); ok {
			e.cause = arg
			args = args[1:]
		}
	}

	if len(args) >= 1 {
		if arg, ok := args[0].(string); ok {
			if len(args) >= 2 {
				e.message = fmt.Sprintf(arg, args[1:]...)
			} else {
				e.message = arg
			}
		}
	}

	return e
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

// Formats the stack of Frames according to the fmt.Formatter interface.
//
//	%s    prints Error as is
//	%v    equivalent to %s
//	%q    prints Error in quotes
//
// Accepts flags that alter the printing of some verbs, as follows:
//
//	%+s   prints Error and it's stack trace (if it has one)
//	%+v   equivalent to %+s
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		io.WriteString(s, e.Error())
		if s.Flag('+') {
			io.WriteString(s, "\n")
			if e.stackTrace != nil {
				e.stackTrace.Format(s, verb)
			}
			return
		}
	case 'q':
		io.WriteString(s, strconv.Quote(e.Error()))
	}
}

func (e *Error) Unwrap() []error {
	var errs []error
	if e.kind != nil {
		errs = append(errs, e.kind)
	}
	if e.cause != nil {
		errs = append(errs, e.cause)
	}
	return errs
}

func (e *Error) Kind() error {
	return e.kind
}

func (e *Error) Cause() error {
	return e.cause
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) StackTrace() StackTrace {
	return slices.Clone(e.stackTrace)
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

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}

func Cause(err error) error {
	for err != nil {
		cause, ok := err.(interface{ Cause() error })
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
