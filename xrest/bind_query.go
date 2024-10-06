package xrest

import (
	"encoding"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/infastin/gorack/fastconv"
)

type InvalidBindQueryError struct {
	Type reflect.Type
}

func (e *InvalidBindQueryError) Error() string {
	if e.Type == nil {
		return "query: invalid out argument: nil"
	}
	if e.Type.Kind() != reflect.Pointer {
		return "query: invalid out argument: non-pointer " + e.Type.String()
	}
	if e.Type.Elem().Kind() != reflect.Struct {
		return "query: invalid out argument: pointer to non-struct " + e.Type.String()
	}
	return "query: invalid out argument: nil " + e.Type.String()
}

type BindQueryTypeError struct {
	Type   reflect.Type
	Struct string
	Field  string
}

func (e *BindQueryTypeError) Error() string {
	return "query: cannot decode into Go struct field " + e.Struct + "." + e.Field + " of type " + e.Type.String()
}

type BindQueryValueError struct {
	Type  reflect.Type
	Value string
	Err   error
}

func (e *BindQueryValueError) Error() string {
	return "query: cannot decode " + strconv.Quote(e.Value) + " into " + e.Type.String() + ": " + e.Err.Error()
}

func (e *BindQueryValueError) Unwrap() error {
	return e.Err
}

func BindQuery(r *http.Request, out any) error {
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidBindQueryError{Type: rv.Type()}
	}

	v := rv.Elem()
	if v.Kind() != reflect.Struct {
		return &InvalidBindQueryError{Type: rv.Type()}
	}

	d := &queryDecoder{
		query:   r.URL.Query(),
		out:     v,
		outType: v.Type(),
	}

	return d.decode()
}

type queryDecoder struct {
	query     url.Values
	out       reflect.Value
	outType   reflect.Type
	fieldType reflect.StructField
	field     reflect.Value
}

func (d *queryDecoder) decode() error {
	for i := range d.outType.NumField() {
		d.fieldType = d.outType.Field(i)

		tag := d.fieldType.Tag.Get("query")
		if tag == "" || !d.query.Has(tag) {
			continue
		}

		d.field = d.out.Field(i)

		var unmarshaler encoding.TextUnmarshaler
		if i, ok := d.field.Interface().(encoding.TextUnmarshaler); ok {
			unmarshaler = i
		}
		if d.field.CanAddr() {
			if i, ok := d.field.Addr().Interface().(encoding.TextUnmarshaler); ok {
				unmarshaler = i
			}
		}
		if unmarshaler != nil {
			if err := unmarshaler.UnmarshalText(fastconv.Bytes(d.query.Get(tag))); err != nil {
				return &BindQueryValueError{Type: d.fieldType.Type, Value: d.query.Get(tag), Err: err}
			}
			continue
		}

		switch d.field.Kind() {
		case reflect.Slice:
			for _, val := range d.query[tag] {
				sv := reflect.New(d.fieldType.Type.Elem())
				if err := d.parseBasicTypes(sv.Elem(), val); err != nil {
					return err
				}
				d.field.Set(reflect.Append(d.field, sv.Elem()))
			}
		case reflect.Array:
			for i, val := range d.query[tag] {
				if i == d.field.Len() {
					break
				}
				if err := d.parseBasicTypes(d.field.Index(i), val); err != nil {
					return err
				}
			}
		default:
			if err := d.parseBasicTypes(d.field, d.query.Get(tag)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *queryDecoder) parseBasicTypes(rv reflect.Value, val string) error {
	switch rv.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return &BindQueryValueError{
				Type:  rv.Type(),
				Value: val,
				Err:   err,
			}
		}
		rv.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return &BindQueryValueError{Type: rv.Type(), Value: val, Err: err}
		}
		rv.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return &BindQueryValueError{Type: rv.Type(), Value: val, Err: err}
		}
		rv.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return &BindQueryValueError{Type: rv.Type(), Value: val, Err: err}
		}
		rv.SetFloat(f)
	case reflect.String:
		rv.SetString(val)
	case reflect.Complex64, reflect.Complex128:
		c, err := strconv.ParseComplex(val, 128)
		if err != nil {
			return &BindQueryValueError{Type: rv.Type(), Value: val, Err: err}
		}
		rv.SetComplex(c)
	case reflect.Pointer:
		if rv.IsNil() {
			elem := reflect.New(rv.Type().Elem())
			if err := d.parseBasicTypes(elem, val); err != nil {
				return err
			}
			rv.Set(elem)
		} else {
			return d.parseBasicTypes(rv.Elem(), val)
		}
	default:
		return &BindQueryTypeError{Type: rv.Type(), Struct: d.outType.Name(), Field: d.fieldType.Name}
	}
	return nil
}
