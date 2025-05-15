package xhttp

import (
	"encoding"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/infastin/gorack/fastconv"
)

type ParamLocation string

const (
	ParamLocationQuery  ParamLocation = "query"
	ParamLocationHeader ParamLocation = "header"
	ParamLocationPath   ParamLocation = "path"
	ParamLocationInline ParamLocation = "inline"
)

type InvalidBindParamsError struct {
	Location ParamLocation
	Type     reflect.Type
}

func (e *InvalidBindParamsError) Error() string {
	if e.Type == nil {
		return fmt.Sprintf("%s: invalid argument: nil", e.Location)
	}
	if e.Type.Kind() != reflect.Pointer {
		return fmt.Sprintf("%s: invalid output argument: non-pointer %s", e.Location, e.Type.String())
	}
	if e.Type.Elem().Kind() != reflect.Struct {
		return fmt.Sprintf("%s: invalid output argument: pointer to non-struct %s", e.Location, e.Type.String())
	}
	return fmt.Sprintf("%s: invalid output argument: nil %s", e.Location, e.Type.String())
}

type BindParamsTypeError struct {
	Location ParamLocation
	Type     reflect.Type
	Struct   string
	Field    string
}

func (e *BindParamsTypeError) Error() string {
	if e.Location == ParamLocationInline {
		return fmt.Sprintf("cannot inline Go struct field %s.%s of type %s", e.Struct, e.Field, e.Type.String())
	}
	return fmt.Sprintf("%s: cannot decode into Go struct field %s.%s of type %s", e.Location, e.Struct, e.Field, e.Type.String())
}

type BindParamsValueError struct {
	Location ParamLocation
	Name     string
	Value    string
	Err      error
}

func (e *BindParamsValueError) Error() string {
	if numError := (*strconv.NumError)(nil); errors.As(e.Err, &numError) {
		return fmt.Sprintf("%s: %q: cannot decode %q: %s", e.Location, e.Name, e.Value, numError.Err.Error())
	}
	return fmt.Sprintf("%s: %q: cannot decode %q: %s", e.Location, e.Name, e.Value, e.Err.Error())
}

func (e *BindParamsValueError) Unwrap() error {
	return e.Err
}

// Decodes request parameters (query, header, path) into a structure
// using structure tags.
func BindParams(r *http.Request, output any) error {
	rv := reflect.ValueOf(output)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidBindParamsError{Type: rv.Type()}
	}

	v := rv.Elem()
	if v.Kind() != reflect.Struct {
		return &InvalidBindParamsError{Type: rv.Type()}
	}

	d := &paramsDecoder{
		req:        r,
		output:     v,
		outputType: v.Type(),
	}

	return d.decode()
}

type paramsDecoder struct {
	req        *http.Request
	output     reflect.Value
	outputType reflect.Type
	location   ParamLocation
	fieldType  reflect.StructField
	field      reflect.Value
	name       string
}

func (d *paramsDecoder) decode() error {
	query := d.req.URL.Query()

	for i := range d.outputType.NumField() {
		d.fieldType = d.outputType.Field(i)
		tag := parseStructTag(d.fieldType.Tag)

		var values []string
		if name, ok := tag["query"]; ok {
			d.location = ParamLocationQuery
			d.name = name
			values = query[d.name]
		} else if name, ok := tag["path"]; ok {
			d.location = ParamLocationPath
			d.name = name
			values = []string{d.req.PathValue(d.name)}
		} else if name, ok := tag["header"]; ok {
			d.location = ParamLocationHeader
			d.name = http.CanonicalHeaderKey(name)
			values = d.req.Header[d.name]
		} else if _, ok := tag["inline"]; ok || d.fieldType.Anonymous {
			d.location = ParamLocationInline

			rv := d.output.Field(i)
			rt := rv.Type()

			for rv.Kind() == reflect.Pointer {
				if rv.IsNil() {
					elem := reflect.New(rt.Elem())
					rv.Set(elem)
				}
				rv = rv.Elem()
				rt = rv.Type()
			}

			if rv.Kind() != reflect.Struct {
				return d.typeError(rt)
			}

			inline := &paramsDecoder{
				req:        d.req,
				output:     rv,
				outputType: rt,
			}

			if err := inline.decode(); err != nil {
				return err
			}
		}

		if len(values) != 0 {
			d.field = d.output.Field(i)
			if err := d.decodeValues(values); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *paramsDecoder) decodeValues(vs []string) error {
	rv := d.field
	rt := rv.Type()

	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			elem := reflect.New(rt.Elem())
			rv.Set(elem)
		}
		rv = rv.Elem()
		rt = rv.Type()
	}

	var unmarshaler encoding.TextUnmarshaler
	if i, ok := rv.Interface().(encoding.TextUnmarshaler); ok {
		unmarshaler = i
	} else if rv.CanAddr() {
		if i, ok := rv.Addr().Interface().(encoding.TextUnmarshaler); ok {
			unmarshaler = i
		}
	}
	if unmarshaler != nil {
		return unmarshaler.UnmarshalText(fastconv.Bytes(vs[0]))
	}

	switch rv.Kind() {
	case reflect.Slice:
		rte := rt.Elem()
		for _, val := range vs {
			sv := reflect.New(rte)
			if err := d.decodeBasicTypes(sv.Elem(), val); err != nil {
				return err
			}
			rv.Set(reflect.Append(rv, sv.Elem()))
		}
	case reflect.Array:
		for i, val := range vs {
			if i == rv.Len() {
				break
			}
			if err := d.decodeBasicTypes(rv.Index(i), val); err != nil {
				return err
			}
		}
	default:
		if err := d.decodeBasicTypes(rv, vs[0]); err != nil {
			return err
		}
	}

	return nil
}

func (d *paramsDecoder) decodeBasicTypes(rv reflect.Value, val string) error {
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
			if err := d.decodeBasicTypes(elem, val); err != nil {
				return err
			}
			rv.Set(elem)
		} else {
			return d.decodeBasicTypes(rv.Elem(), val)
		}
	default:
		return d.typeError(rv.Type())
	}
	return nil
}

func (d *paramsDecoder) typeError(typ reflect.Type) error {
	return &BindParamsTypeError{
		Location: d.location,
		Type:     typ,
		Struct:   d.outputType.Name(),
		Field:    d.fieldType.Name,
	}
}

func (d *paramsDecoder) valueError(val string, err error) error {
	return &BindParamsValueError{
		Location: d.location,
		Name:     d.name,
		Value:    val,
		Err:      err,
	}
}

func parseStructTag(tag reflect.StructTag) map[string]string {
	result := make(map[string]string)

	for tag != "" {
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}

		value, err := strconv.Unquote(string(tag[:i+1]))
		if err != nil {
			break
		}
		tag = tag[i+1:]

		result[name] = value
	}

	return result
}
