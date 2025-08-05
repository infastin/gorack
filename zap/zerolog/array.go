package zapzerolog

import (
	"errors"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/zap/zapcore"
)

type arrayMarshalerFunc func(arr *zerolog.Array)

func (f arrayMarshalerFunc) MarshalZerologArray(arr *zerolog.Array) {
	f(arr)
}

type array struct {
	arr *zerolog.Array
}

func newArray() *array {
	return &array{arr: zerolog.Arr()}
}

func arrayFromZerolog(arr *zerolog.Array) *array {
	return &array{arr: arr}
}

func (a *array) unwrap() *zerolog.Array {
	return a.arr
}

func (a *array) AppendBool(value bool) {
	a.arr.Bool(value)
}

func (a *array) AppendByteString(value []byte) {
	a.arr.Bytes(value)
}

func (a *array) AppendComplex128(value complex128) {
	a.arr.Str(strconv.FormatComplex(value, 'f', -1, 128))
}

func (a *array) AppendComplex64(value complex64) {
	a.arr.Str(strconv.FormatComplex(complex128(value), 'f', -1, 128))
}

func (a *array) AppendFloat64(value float64) {
	a.arr.Float64(value)
}

func (a *array) AppendFloat32(value float32) {
	a.arr.Float32(value)
}

func (a *array) AppendInt(value int) {
	a.arr.Int(value)
}

func (a *array) AppendInt64(value int64) {
	a.arr.Int64(value)
}

func (a *array) AppendInt32(value int32) {
	a.arr.Int32(value)
}

func (a *array) AppendInt16(value int16) {
	a.arr.Int16(value)
}

func (a *array) AppendInt8(value int8) {
	a.arr.Int8(value)
}

func (a *array) AppendString(value string) {
	a.arr.Str(value)
}

func (a *array) AppendUint(value uint) {
	a.arr.Uint(value)
}

func (a *array) AppendUint64(value uint64) {
	a.arr.Uint64(value)
}

func (a *array) AppendUint32(value uint32) {
	a.arr.Uint32(value)
}

func (a *array) AppendUint16(value uint16) {
	a.arr.Uint16(value)
}

func (a *array) AppendUint8(value uint8) {
	a.arr.Uint8(value)
}

func (a *array) AppendUintptr(value uintptr) {
	a.arr.Uint(uint(value))
}

func (a *array) AppendDuration(value time.Duration) {
	a.arr.Dur(value)
}

func (a *array) AppendTime(value time.Time) {
	a.arr.Time(value)
}

func (*array) AppendArray(marshaler zapcore.ArrayMarshaler) error {
	// ???
	return errors.New("nested arrays are not supported by zerolog")
}

func (a *array) AppendObject(marshaler zapcore.ObjectMarshaler) error {
	obj := newObject()
	if err := marshaler.MarshalLogObject(obj); err != nil {
		return err
	}
	a.arr.Dict(obj.close())
	return nil
}

func (a *array) AppendReflected(value any) error {
	a.arr.Interface(value)
	return nil
}
