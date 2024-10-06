package xrest

import (
	"encoding"
	"errors"
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
	Name  string
	Value string
	Err   error
}

func (e *BindQueryValueError) Error() string {
	if numError := (*strconv.NumError)(nil); errors.As(e.Err, &numError) {
		return "query: " + strconv.Quote(e.Name) + ": cannot decode " + strconv.Quote(e.Value) + ": " + numError.Err.Error()
	}
	return "query: " + strconv.Quote(e.Name) + ": cannot decode " + strconv.Quote(e.Value) + ": " + e.Err.Error()
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
	name      string
}

func (d *queryDecoder) decode() error {
	for i := range d.outType.NumField() {
		d.fieldType = d.outType.Field(i)

		d.name = d.fieldType.Tag.Get("query")
		if d.name == "" || !d.query.Has(d.name) {
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
			if err := unmarshaler.UnmarshalText(fastconv.Bytes(d.query.Get(d.name))); err != nil {
				return d.valueError(d.query.Get(d.name), err)
			}
			continue
		}

		switch d.field.Kind() {
		case reflect.Slice:
			for _, val := range d.query[d.name] {
				sv := reflect.New(d.fieldType.Type.Elem())
				if err := d.parseBasicTypes(sv.Elem(), val); err != nil {
					return err
				}
				d.field.Set(reflect.Append(d.field, sv.Elem()))
			}
		case reflect.Array:
			for i, val := range d.query[d.name] {
				if i == d.field.Len() {
					break
				}
				if err := d.parseBasicTypes(d.field.Index(i), val); err != nil {
					return err
				}
			}
		default:
			if err := d.parseBasicTypes(d.field, d.query.Get(d.name)); err != nil {
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
			return d.valueError(val, err)
		}
		rv.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return d.valueError(val, err)
		}
		rv.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return d.valueError(val, err)
		}
		rv.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return d.valueError(val, err)
		}
		rv.SetFloat(f)
	case reflect.String:
		rv.SetString(val)
	case reflect.Complex64, reflect.Complex128:
		c, err := strconv.ParseComplex(val, 128)
		if err != nil {
			return d.valueError(val, err)
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
		return d.typeError(rv.Type())
	}
	return nil
}

func (d *queryDecoder) typeError(typ reflect.Type) error {
	return &BindQueryTypeError{
		Type:   typ,
		Struct: d.outType.Name(),
		Field:  d.fieldType.Name,
	}
}

func (d *queryDecoder) valueError(val string, err error) error {
	return &BindQueryValueError{
		Name:  d.name,
		Value: val,
		Err:   err,
	}
}
