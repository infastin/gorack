package slogzerolog

import (
	"context"
	"log/slog"

	"github.com/rs/zerolog"
)

var logLevels = map[slog.Level]zerolog.Level{
	slog.LevelDebug: zerolog.DebugLevel,
	slog.LevelInfo:  zerolog.InfoLevel,
	slog.LevelWarn:  zerolog.WarnLevel,
	slog.LevelError: zerolog.ErrorLevel,
}

type handler struct {
	lg     zerolog.Logger
	groups []string
}

func New(lg zerolog.Logger) slog.Handler {
	return &handler{lg: lg, groups: nil}
}

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
	return logLevels[level] >= h.lg.GetLevel()
}

func (h *handler) Handle(ctx context.Context, record slog.Record) error {
	var root, group *zerolog.Event
	for i, name := range h.groups {
		if i == 0 {
			root = zerolog.Dict()
			group = root
		} else {
			dict := zerolog.Dict()
			group.Dict(name, dict)
			group = dict
		}
	}
	lg := h.lg
	ev := lg.WithLevel(logLevels[record.Level]).Ctx(ctx)
	if group == nil {
		group = ev
	}
	record.Attrs(func(attr slog.Attr) bool {
		eventAttr(group, attr)
		return true
	})
	if root != nil {
		ev.Dict(h.groups[0], root)
	}
	ev.Msg(record.Message)
	return ev.GetCtx().Err()
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	lgCtx := h.lg.With()
	if len(h.groups) != 0 {
		var root, group *zerolog.Event
		for i, name := range h.groups {
			if i == 0 {
				root = zerolog.Dict()
				group = root
			} else {
				dict := zerolog.Dict()
				group.Dict(name, dict)
				group = dict
			}
		}
		for _, attr := range attrs {
			eventAttr(group, attr)
		}
		lgCtx = lgCtx.Dict(h.groups[0], root)
	} else {
		for _, attr := range attrs {
			lgCtx = ctxAttr(lgCtx, attr)
		}
	}
	return &handler{lg: lgCtx.Logger()}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		lg:     h.lg,
		groups: append(h.groups, name),
	}
}

func ctxAttr(ctx zerolog.Context, attr slog.Attr) zerolog.Context {
	switch attr.Value.Kind() {
	case slog.KindAny:
		ctx = ctx.Any(attr.Key, attr.Value.Any())
	case slog.KindBool:
		ctx = ctx.Bool(attr.Key, attr.Value.Bool())
	case slog.KindDuration:
		ctx = ctx.Dur(attr.Key, attr.Value.Duration())
	case slog.KindFloat64:
		ctx = ctx.Float64(attr.Key, attr.Value.Float64())
	case slog.KindInt64:
		ctx = ctx.Int64(attr.Key, attr.Value.Int64())
	case slog.KindString:
		ctx = ctx.Str(attr.Key, attr.Value.String())
	case slog.KindTime:
		ctx = ctx.Time(attr.Key, attr.Value.Time())
	case slog.KindUint64:
		ctx = ctx.Uint64(attr.Key, attr.Value.Uint64())
	case slog.KindGroup:
		dict := zerolog.Dict()
		for _, attr := range attr.Value.Group() {
			eventAttr(dict, attr)
		}
		ctx = ctx.Dict(attr.Key, dict)
	case slog.KindLogValuer:
		ctx = ctxAttr(ctx, slog.Attr{
			Key:   attr.Key,
			Value: attr.Value.Resolve(),
		})
	}
	return ctx
}

func eventAttr(ev *zerolog.Event, attr slog.Attr) {
	switch attr.Value.Kind() {
	case slog.KindAny:
		ev.Any(attr.Key, attr.Value.Any())
	case slog.KindBool:
		ev.Bool(attr.Key, attr.Value.Bool())
	case slog.KindDuration:
		ev.Dur(attr.Key, attr.Value.Duration())
	case slog.KindFloat64:
		ev.Float64(attr.Key, attr.Value.Float64())
	case slog.KindInt64:
		ev.Int64(attr.Key, attr.Value.Int64())
	case slog.KindString:
		ev.Str(attr.Key, attr.Value.String())
	case slog.KindTime:
		ev.Time(attr.Key, attr.Value.Time())
	case slog.KindUint64:
		ev.Uint64(attr.Key, attr.Value.Uint64())
	case slog.KindGroup:
		dict := zerolog.Dict()
		for _, attr := range attr.Value.Group() {
			eventAttr(dict, attr)
		}
		ev.Dict(attr.Key, dict)
	case slog.KindLogValuer:
		eventAttr(ev, slog.Attr{
			Key:   attr.Key,
			Value: attr.Value.Resolve(),
		})
	}
}
