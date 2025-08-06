package errdefer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/infastin/gorack/errdefer"
)

func TestClose(t *testing.T) {
	const (
		msg    = "42"
		errMsg = "cat"
	)

	err := errors.New(errMsg)
	var value string

	errdefer.Close(&err, func() {
		value = msg
	})
	if value != msg {
		t.Errorf("fn(): values %q != %q", value, msg)
	}

	gotErr := errdefer.Close(&err, func() error {
		value = msg
		return errors.New(errMsg)
	})
	if value != msg {
		t.Errorf("fn() error: values %q != %q", value, msg)
	}
	if gotErr == nil {
		t.Error("fn() error: must return non-nil error")
	} else if gotErr.Error() != err.Error() {
		t.Errorf("fn() error: errors %q != %q", gotErr.Error(), err.Error())
	}

	errdefer.Close(&err, func(ctx context.Context) {
		value = msg
	})
	if value != msg {
		t.Errorf("fn(context.Context): values %q != %q", value, msg)
	}

	gotErr = errdefer.Close(&err, func(ctx context.Context) error {
		value = msg
		return errors.New(errMsg)
	})
	if value != msg {
		t.Errorf("fn(context.Context) error: values %q != %q", value, msg)
	}
	if gotErr == nil {
		t.Error("fn(context.Context) error: must return non-nil error")
	} else if gotErr.Error() != err.Error() {
		t.Errorf("fn(context.Context) error: errors %q != %q", gotErr.Error(), err.Error())
	}
}

func TestClose_reflect(t *testing.T) {
	type (
		fn   func()
		fne  func() error
		fnc  func(context.Context)
		fnce func(context.Context) error
	)

	const (
		msg    = "42"
		errMsg = "cat"
	)

	err := errors.New(errMsg)
	var value string

	errdefer.Close(&err, fn(func() {
		value = msg
	}))
	if value != msg {
		t.Errorf("fn(): values %q != %q", value, msg)
	}

	gotErr := errdefer.Close(&err, fne(func() error {
		value = msg
		return errors.New(errMsg)
	}))
	if value != msg {
		t.Errorf("fn() error: values %q != %q", value, msg)
	}
	if gotErr.Error() != err.Error() {
		t.Errorf("fn() error: errors %q != %q", gotErr.Error(), err.Error())
	}

	errdefer.Close(&err, fnc(func(ctx context.Context) {
		value = msg
	}))
	if value != msg {
		t.Errorf("fn(context.Context): values %q != %q", value, msg)
	}

	gotErr = errdefer.Close(&err, fnce(func(ctx context.Context) error {
		value = msg
		return errors.New(errMsg)
	}))
	if value != msg {
		t.Errorf("fn(context.Context) error: values %q != %q", value, msg)
	}
	if gotErr == nil {
		t.Error("fn(context.Context) error: must return non-nil error")
	} else if gotErr.Error() != err.Error() {
		t.Errorf("fn(context.Context) error: errors %q != %q", gotErr.Error(), err.Error())
	}
}

func TestCloseContext(t *testing.T) {
	type ctxKey struct{}

	const (
		errMsg = "cat"
		ctxMsg = "owl"
	)

	ctx := context.WithValue(context.Background(), ctxKey{}, ctxMsg)

	err := errors.New(errMsg)
	var value string

	errdefer.CloseContext(ctx, &err, func(ctx context.Context) {
		value = ctx.Value(ctxKey{}).(string)
	})
	if value != ctxMsg {
		t.Errorf("fn(context.Context): values %q != %q", value, ctxMsg)
	}

	gotErr := errdefer.CloseContext(ctx, &err, func(ctx context.Context) error {
		value = ctx.Value(ctxKey{}).(string)
		return errors.New(errMsg)
	})
	if value != ctxMsg {
		t.Errorf("fn(context.Context) error: values %q != %q", value, ctxMsg)
	}
	if gotErr == nil {
		t.Error("fn(context.Context) error: must return non-nil error")
	} else if gotErr.Error() != err.Error() {
		t.Errorf("fn(context.Context) error: errors %q != %q", gotErr.Error(), err.Error())
	}
}

func TestCloseContext_reflect(t *testing.T) {
	type (
		ctxKey struct{}
		fnc    func(context.Context)
		fnce   func(context.Context) error
	)

	const (
		errMsg = "cat"
		ctxMsg = "owl"
	)

	ctx := context.WithValue(context.Background(), ctxKey{}, ctxMsg)

	err := errors.New(errMsg)
	var value string

	errdefer.CloseContext(ctx, &err, fnc(func(ctx context.Context) {
		value = ctx.Value(ctxKey{}).(string)
	}))
	if value != ctxMsg {
		t.Errorf("fn(context.Context): values %q != %q", value, ctxMsg)
	}

	gotErr := errdefer.CloseContext(ctx, &err, fnce(func(ctx context.Context) error {
		value = ctx.Value(ctxKey{}).(string)
		return errors.New(errMsg)
	}))
	if value != ctxMsg {
		t.Errorf("fn(context.Context) error: values %q != %q", value, ctxMsg)
	}
	if gotErr == nil {
		t.Error("fn(context.Context) error: must return non-nil error")
	} else if gotErr.Error() != err.Error() {
		t.Errorf("fn(context.Context) error: errors %q != %q", gotErr.Error(), err.Error())
	}
}
