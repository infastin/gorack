package opt

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

type isZeroer interface {
	IsZero() bool
}

var (
	isZeroerType        = reflect.TypeFor[isZeroer]()
	textMarshalerType   = reflect.TypeFor[encoding.TextMarshaler]()
	textUnmarshalerType = reflect.TypeFor[encoding.TextUnmarshaler]()
)

func isZero(v any) bool {
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	if rt.Implements(isZeroerType) {
		switch rt.Kind() {
		case reflect.Interface:
			return rv.IsNil() ||
				(rv.Elem().Kind() == reflect.Pointer && rv.Elem().IsNil()) ||
				rv.Interface().(isZeroer).IsZero()
		case reflect.Pointer:
			return rv.IsNil() || rv.Interface().(isZeroer).IsZero()
		default:
			return rv.Interface().(isZeroer).IsZero()
		}
	} else if reflect.PointerTo(rt).Implements(isZeroerType) {
		if !rv.CanAddr() {
			tmp := reflect.New(rt).Elem()
			tmp.Set(rv)
			rv = tmp
		}
		return rv.Addr().Interface().(isZeroer).IsZero()
	}

	return reflect.ValueOf(v).IsZero()
}

func marshalText(v any) ([]byte, error) {
	return valueMarshalText(reflect.ValueOf(v))
}

func valueMarshalText(rv reflect.Value) ([]byte, error) {
	rt := rv.Type()

	if rt.Implements(textMarshalerType) {
		return rv.Interface().(encoding.TextMarshaler).MarshalText()
	}
	if rt.Kind() != reflect.Pointer && rv.CanAddr() && reflect.PointerTo(rt).Implements(textMarshalerType) {
		return rv.Addr().Interface().(encoding.TextMarshaler).MarshalText()
	}

	switch rt.Kind() {
	case reflect.Bool:
		return strconv.AppendBool(nil, rv.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(nil, rv.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.AppendUint(nil, rv.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.AppendFloat(nil, rv.Float(), 'g', -1, rt.Bits()), nil
	case reflect.String:
		return []byte(rv.String()), nil
	case reflect.Pointer:
		if rv.IsNil() {
			return []byte{}, nil
		}
		return valueMarshalText(rv.Elem())
	default:
		return nil, fmt.Errorf("opt: unsupported type: %s", rt.String())
	}
}

func unmarshalText(v any, data []byte) error {
	return valueUnmarshalText(reflect.ValueOf(v), data)
}

func valueUnmarshalText(rv reflect.Value, data []byte) error {
	rt := rv.Type()

	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("opt: must be a pointer: %s", rt.String())
	}
	if rv.IsNil() {
		return fmt.Errorf("opt: must be non-nil: %s", rt.String())
	}

	if rt.Implements(textUnmarshalerType) {
		return rv.Interface().(encoding.TextUnmarshaler).UnmarshalText(data)
	}

	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			elem := reflect.New(rv.Type().Elem())
			rv.Set(elem)
		}
		rv = rv.Elem()
	}
	rt = rv.Type()

	text := unsafe.String(unsafe.SliceData(data), len(data))

	switch rt.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(text)
		if err != nil {
			return fmt.Errorf("opt: %w", err)
		}
		rv.SetBool(b)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(text, 10, rt.Bits())
		if err != nil {
			return fmt.Errorf("opt: %w", err)
		}
		rv.SetInt(i)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u, err := strconv.ParseUint(text, 10, rt.Bits())
		if err != nil {
			return fmt.Errorf("opt: %w", err)
		}
		rv.SetUint(u)
		return nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(text, rt.Bits())
		if err != nil {
			return fmt.Errorf("opt: %w", err)
		}
		rv.SetFloat(f)
		return nil
	case reflect.String:
		rv.SetString(text)
		return nil
	default:
		return fmt.Errorf("opt: unsupported type: %s", rt.String())
	}
}
