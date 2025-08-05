package zapzerolog

import (
	"encoding/base64"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/zap/zapcore"
)

type objectMarshalerFunc func(ev *zerolog.Event)

func (f objectMarshalerFunc) MarshalZerologObject(ev *zerolog.Event) {
	f(ev)
}

type object struct {
	parent *object
	key    string
	dict   *zerolog.Event
}

func newObject() *object {
	return &object{
		parent: nil,
		key:    "",
		dict:   zerolog.Dict(),
	}
}

func objectFromZerolog(ev *zerolog.Event) *object {
	return &object{
		parent: nil,
		key:    "",
		dict:   ev,
	}
}

func (o *object) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	arr := newArray()
	if err := marshaler.MarshalLogArray(arr); err != nil {
		return err
	}
	o.dict.Array(key, arr.unwrap())
	return nil
}

func (o *object) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	obj := newObject()
	if err := marshaler.MarshalLogObject(obj); err != nil {
		return err
	}
	o.dict.Dict(key, obj.unwrap())
	return nil
}

func (o *object) AddBinary(key string, value []byte) {
	o.dict.Str(key, base64.StdEncoding.EncodeToString(value))
}

func (o *object) AddByteString(key string, value []byte) {
	o.dict.Bytes(key, value)
}

func (o *object) AddBool(key string, value bool) {
	o.dict.Bool(key, value)
}

func (o *object) AddComplex128(key string, value complex128) {
	o.dict.Str(key, strconv.FormatComplex(value, 'f', -1, 128))
}

func (o *object) AddComplex64(key string, value complex64) {
	o.dict.Str(key, strconv.FormatComplex(complex128(value), 'f', -1, 128))
}

func (o *object) AddDuration(key string, value time.Duration) {
	o.dict.Dur(key, value)
}

func (o *object) AddFloat64(key string, value float64) {
	o.dict.Float64(key, value)
}

func (o *object) AddFloat32(key string, value float32) {
	o.dict.Float32(key, value)
}

func (o *object) AddInt(key string, value int) {
	o.dict.Int(key, value)
}

func (o *object) AddInt64(key string, value int64) {
	o.dict.Int64(key, value)
}

func (o *object) AddInt32(key string, value int32) {
	o.dict.Int32(key, value)
}

func (o *object) AddInt16(key string, value int16) {
	o.dict.Int16(key, value)
}

func (o *object) AddInt8(key string, value int8) {
	o.dict.Int8(key, value)
}

func (o *object) AddString(key, value string) {
	o.dict.Str(key, value)
}

func (o *object) AddTime(key string, value time.Time) {
	o.dict.Time(key, value)
}

func (o *object) AddUint(key string, value uint) {
	o.dict.Uint(key, value)
}

func (o *object) AddUint64(key string, value uint64) {
	o.dict.Uint64(key, value)
}

func (o *object) AddUint32(key string, value uint32) {
	o.dict.Uint32(key, value)
}

func (o *object) AddUint16(key string, value uint16) {
	o.dict.Uint16(key, value)
}

func (o *object) AddUint8(key string, value uint8) {
	o.dict.Uint8(key, value)
}

func (o *object) AddUintptr(key string, value uintptr) {
	o.dict.Uint(key, uint(value))
}

func (o *object) AddReflected(key string, value any) error {
	o.dict.Interface(key, value)
	return nil
}

func (o *object) OpenNamespace(key string) {
	*o = object{
		parent: &object{
			parent: o.parent,
			key:    o.key,
			dict:   o.dict,
		},
		key:  key,
		dict: zerolog.Dict(),
	}
}

func (o *object) unwrap() *zerolog.Event {
	current := o
	for current.parent != nil {
		parent := current.parent
		parent.dict.Dict(current.key, current.dict)
		current.parent, current = nil, parent
	}
	return current.dict
}
