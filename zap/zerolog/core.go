package zapzerolog

import (
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/zap/zapcore"
)

var logLevels = map[zapcore.Level]zerolog.Level{
	zapcore.DebugLevel:  zerolog.DebugLevel,
	zapcore.InfoLevel:   zerolog.InfoLevel,
	zapcore.WarnLevel:   zerolog.WarnLevel,
	zapcore.ErrorLevel:  zerolog.ErrorLevel,
	zapcore.DPanicLevel: zerolog.PanicLevel,
	zapcore.PanicLevel:  zerolog.PanicLevel,
	zapcore.FatalLevel:  zerolog.FatalLevel,
}

type core struct {
	lg  zerolog.Logger
	cfg *config
}

func New(lg zerolog.Logger, opts ...Option) zapcore.Core {
	c := &core{
		lg:  lg,
		cfg: defaultConfig(),
	}
	for _, opt := range opts {
		opt(c.cfg)
	}
	return c
}

func (c *core) Enabled(level zapcore.Level) bool {
	lvl := logLevels[level]
	return lvl >= c.lg.GetLevel() && lvl >= zerolog.GlobalLevel()
}

func (c *core) With(fields []zapcore.Field) zapcore.Core {
	return &core{
		lg:  ctxFields(c.lg.With(), fields).Logger(),
		cfg: c.cfg,
	}
}

func (c *core) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return nil
}

func (c *core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	ev := c.lg.WithLevel(logLevels[entry.Level])
	if entry.LoggerName != "" && c.cfg.nameKey != "" {
		ev.Str(c.cfg.nameKey, entry.LoggerName)
	}
	if c.cfg.copyTimestamp {
		ev.Time(zerolog.TimestampFieldName, entry.Time)
	}
	if entry.Stack != "" && c.cfg.copyStack {
		ev.Str(zerolog.ErrorStackFieldName, entry.Stack)
	}
	if entry.Caller.Defined && c.cfg.copyCaller {
		ev.Str(zerolog.CallerFieldName,
			zerolog.CallerMarshalFunc(entry.Caller.PC, entry.Caller.File, entry.Caller.Line))
	}
	if err := eventFields(ev, fields); err != nil {
		return err
	}
	ev.Msg(entry.Message)
	return nil
}

func (*core) Sync() error {
	return nil
}

func ctxFields(ctx zerolog.Context, fields []zapcore.Field) zerolog.Context {
	for i := range fields {
		field := &fields[i]
		if field.Type == zapcore.NamespaceType {
			dict := zerolog.Dict()
			_ = eventFields(dict, fields[i+1:])
			ctx = ctx.Dict(field.Key, dict)
			break
		}
		ctx = ctxField(ctx, field)
	}
	return ctx
}

func ctxField(ctx zerolog.Context, field *zapcore.Field) zerolog.Context {
	switch field.Type {
	case zapcore.ArrayMarshalerType:
		ctx = ctx.Array(field.Key, arrayMarshalerFunc(func(a *zerolog.Array) {
			arr := arrayFromZerolog(a)
			_ = field.Interface.(zapcore.ArrayMarshaler).MarshalLogArray(arr)
		}))
	case zapcore.ObjectMarshalerType:
		ctx = ctx.Object(field.Key, objectMarshalerFunc(func(ev *zerolog.Event) {
			obj := objectFromZerolog(ev)
			_ = field.Interface.(zapcore.ObjectMarshaler).MarshalLogObject(obj)
			obj.close()
		}))
	case zapcore.InlineMarshalerType:
		ctx = ctx.EmbedObject(objectMarshalerFunc(func(ev *zerolog.Event) {
			obj := objectFromZerolog(ev)
			_ = field.Interface.(zapcore.ObjectMarshaler).MarshalLogObject(obj)
			obj.close()
		}))
	case zapcore.BinaryType:
		ctx = ctx.Str(field.Key, base64.StdEncoding.EncodeToString(field.Interface.([]byte)))
	case zapcore.BoolType:
		ctx = ctx.Bool(field.Key, field.Integer == 1)
	case zapcore.ByteStringType:
		ctx = ctx.Bytes(field.Key, field.Interface.([]byte))
	case zapcore.Complex128Type:
		ctx = ctx.Str(field.Key, strconv.FormatComplex(field.Interface.(complex128), 'f', -1, 128))
	case zapcore.Complex64Type:
		ctx = ctx.Str(field.Key, strconv.FormatComplex(complex128(field.Interface.(complex64)), 'f', -1, 64))
	case zapcore.DurationType:
		ctx = ctx.Dur(field.Key, time.Duration(field.Integer))
	case zapcore.Float64Type:
		ctx = ctx.Float64(field.Key, math.Float64frombits(uint64(field.Integer)))
	case zapcore.Float32Type:
		ctx = ctx.Float32(field.Key, math.Float32frombits(uint32(field.Integer)))
	case zapcore.Int64Type:
		ctx = ctx.Int64(field.Key, field.Integer)
	case zapcore.Int32Type:
		ctx = ctx.Int32(field.Key, int32(field.Integer))
	case zapcore.Int16Type:
		ctx = ctx.Int16(field.Key, int16(field.Integer))
	case zapcore.Int8Type:
		ctx = ctx.Int8(field.Key, int8(field.Integer))
	case zapcore.StringType:
		ctx = ctx.Str(field.Key, field.String)
	case zapcore.TimeType:
		if field.Interface != nil {
			ctx = ctx.Time(field.Key, time.Unix(0, field.Integer).In(field.Interface.(*time.Location)))
		} else {
			ctx = ctx.Time(field.Key, time.Unix(0, field.Integer))
		}
	case zapcore.TimeFullType:
		ctx = ctx.Time(field.Key, field.Interface.(time.Time))
	case zapcore.Uint64Type:
		ctx = ctx.Uint64(field.Key, uint64(field.Integer))
	case zapcore.Uint32Type:
		ctx = ctx.Uint32(field.Key, uint32(field.Integer))
	case zapcore.Uint16Type:
		ctx = ctx.Uint16(field.Key, uint16(field.Integer))
	case zapcore.Uint8Type:
		ctx = ctx.Uint8(field.Key, uint8(field.Integer))
	case zapcore.UintptrType:
		ctx = ctx.Uint(field.Key, uint(field.Integer))
	case zapcore.ReflectType:
		ctx = ctx.Type(field.Key, field.Interface)
	case zapcore.StringerType:
		ctx = ctx.Stringer(field.Key, field.Interface.(fmt.Stringer))
	case zapcore.ErrorType:
		ctx = ctx.AnErr(field.Key, field.Interface.(error))
	case zapcore.NamespaceType:
		// handled in ctxFields
	case zapcore.UnknownType, zapcore.SkipType:
		// noop
	}
	return ctx
}

func eventFields(ev *zerolog.Event, fields []zapcore.Field) error {
	for i := range fields {
		field := &fields[i]
		if field.Type == zapcore.NamespaceType {
			dict := zerolog.Dict()
			if err := eventFields(dict, fields[i+1:]); err != nil {
				return err
			}
			ev.Dict(field.Key, dict)
			break
		}
		if err := eventField(ev, field); err != nil {
			return err
		}
	}
	return nil
}

func eventField(ev *zerolog.Event, field *zapcore.Field) (err error) {
	switch field.Type {
	case zapcore.ArrayMarshalerType:
		ev.Array(field.Key, arrayMarshalerFunc(func(a *zerolog.Array) {
			arr := arrayFromZerolog(a)
			if err = field.Interface.(zapcore.ArrayMarshaler).MarshalLogArray(arr); err != nil {
				return
			}
		}))
	case zapcore.ObjectMarshalerType:
		ev.Object(field.Key, objectMarshalerFunc(func(ev *zerolog.Event) {
			obj := objectFromZerolog(ev)
			if err = field.Interface.(zapcore.ObjectMarshaler).MarshalLogObject(obj); err != nil {
				return
			}
			obj.close()
		}))
	case zapcore.InlineMarshalerType:
		ev.EmbedObject(objectMarshalerFunc(func(ev *zerolog.Event) {
			obj := objectFromZerolog(ev)
			if err = field.Interface.(zapcore.ObjectMarshaler).MarshalLogObject(obj); err != nil {
				return
			}
			obj.close()
		}))
	case zapcore.BinaryType:
		ev.Str(field.Key, base64.StdEncoding.EncodeToString(field.Interface.([]byte)))
	case zapcore.BoolType:
		ev.Bool(field.Key, field.Integer == 1)
	case zapcore.ByteStringType:
		ev.Bytes(field.Key, field.Interface.([]byte))
	case zapcore.Complex128Type:
		ev.Str(field.Key, strconv.FormatComplex(field.Interface.(complex128), 'f', -1, 128))
	case zapcore.Complex64Type:
		ev.Str(field.Key, strconv.FormatComplex(complex128(field.Interface.(complex64)), 'f', -1, 64))
	case zapcore.DurationType:
		ev.Dur(field.Key, time.Duration(field.Integer))
	case zapcore.Float64Type:
		ev.Float64(field.Key, math.Float64frombits(uint64(field.Integer)))
	case zapcore.Float32Type:
		ev.Float32(field.Key, math.Float32frombits(uint32(field.Integer)))
	case zapcore.Int64Type:
		ev.Int64(field.Key, field.Integer)
	case zapcore.Int32Type:
		ev.Int32(field.Key, int32(field.Integer))
	case zapcore.Int16Type:
		ev.Int16(field.Key, int16(field.Integer))
	case zapcore.Int8Type:
		ev.Int8(field.Key, int8(field.Integer))
	case zapcore.StringType:
		ev.Str(field.Key, field.String)
	case zapcore.TimeType:
		if field.Interface != nil {
			ev.Time(field.Key, time.Unix(0, field.Integer).In(field.Interface.(*time.Location)))
		} else {
			ev.Time(field.Key, time.Unix(0, field.Integer))
		}
	case zapcore.TimeFullType:
		ev.Time(field.Key, field.Interface.(time.Time))
	case zapcore.Uint64Type:
		ev.Uint64(field.Key, uint64(field.Integer))
	case zapcore.Uint32Type:
		ev.Uint32(field.Key, uint32(field.Integer))
	case zapcore.Uint16Type:
		ev.Uint16(field.Key, uint16(field.Integer))
	case zapcore.Uint8Type:
		ev.Uint8(field.Key, uint8(field.Integer))
	case zapcore.UintptrType:
		ev.Uint(field.Key, uint(field.Integer))
	case zapcore.ReflectType:
		ev.Type(field.Key, field.Interface)
	case zapcore.StringerType:
		ev.Stringer(field.Key, field.Interface.(fmt.Stringer))
	case zapcore.ErrorType:
		ev.AnErr(field.Key, field.Interface.(error))
	case zapcore.NamespaceType:
		// handled in eventFields
	case zapcore.UnknownType, zapcore.SkipType:
		// noop
	}
	return err
}
