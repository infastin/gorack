package lifecycle

import (
	"context"
	"log/slog"
)

// discardHandler discards all log output.
// TODO: replace with slog.DiscardHandler in Go 1.24.
type discardHandler struct{}

func (dh discardHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (dh discardHandler) Handle(context.Context, slog.Record) error { return nil }
func (dh discardHandler) WithAttrs(attrs []slog.Attr) slog.Handler  { return dh }
func (dh discardHandler) WithGroup(name string) slog.Handler        { return dh }
