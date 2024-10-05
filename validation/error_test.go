package validation_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/infastin/gorack/validation"
)

func Test_ruleError_Error(t *testing.T) {
	type fields struct {
		code    string
		message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"ab", fields{"a", "b"}, "b"},
		{"hello", fields{"hello", "world"}, "world"},
		{"foobar", fields{"foo", "bar"}, "bar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := validation.NewRuleError(tt.fields.code, tt.fields.message)
			if got := re.Error(); got != tt.want {
				t.Errorf("ruleError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ruleError_Code(t *testing.T) {
	type fields struct {
		code    string
		message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"ab", fields{"a", "b"}, "a"},
		{"hello", fields{"hello", "world"}, "hello"},
		{"foobar", fields{"foo", "bar"}, "foo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := validation.NewRuleError(tt.fields.code, tt.fields.message)
			if got := re.Code(); got != tt.want {
				t.Errorf("ruleError.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ruleError_Message(t *testing.T) {
	type fields struct {
		code    string
		message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"ab", fields{"a", "b"}, "b"},
		{"hello", fields{"hello", "world"}, "world"},
		{"foobar", fields{"foo", "bar"}, "bar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := validation.NewRuleError(tt.fields.code, tt.fields.message)
			if got := re.Message(); got != tt.want {
				t.Errorf("ruleError.Message() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_valueError_Error(t *testing.T) {
	type fields struct {
		name   string
		nested error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"ab", fields{"a", errors.New("b")}, "a: b"},
		{"hello", fields{"hello", errors.New("world")}, "hello: world"},
		{"foobar", fields{"foo", errors.New("bar")}, "foo: bar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ve := validation.NewValueError(tt.fields.name, tt.fields.nested)
			if got := ve.Error(); got != tt.want {
				t.Errorf("valueError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_valueError_Name(t *testing.T) {
	type fields struct {
		name   string
		nested error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"ab", fields{"a", errors.New("b")}, "a"},
		{"hello", fields{"hello", errors.New("world")}, "hello"},
		{"foobar", fields{"foo", errors.New("bar")}, "foo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ve := validation.NewValueError(tt.fields.name, tt.fields.nested)
			if got := ve.Name(); got != tt.want {
				t.Errorf("valueError.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexError_Error(t *testing.T) {
	type fields struct {
		index  int
		nested error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"0b", fields{0, errors.New("b")}, "[0]: b"},
		{"hello", fields{1337, errors.New("hello world")}, "[1337]: hello world"},
		{"foobar", fields{0xdeadbeef, errors.New("foobar")}, "[3735928559]: foobar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ie := validation.NewIndexError(tt.fields.index, tt.fields.nested)
			if got := ie.Error(); got != tt.want {
				t.Errorf("indexError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexError_Index(t *testing.T) {
	type fields struct {
		index  int
		nested error
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"0b", fields{0, errors.New("b")}, 0},
		{"hello", fields{1337, errors.New("hello world")}, 1337},
		{"foobar", fields{0xdeadbeef, errors.New("foobar")}, 0xdeadbeef},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ie := validation.NewIndexError(tt.fields.index, tt.fields.nested)
			if got := ie.Index(); got != tt.want {
				t.Errorf("indexError.Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrors_Error(t *testing.T) {
	tests := []struct {
		name string
		es   validation.Errors
		want string
	}{
		{"foobarbaz", []error{
			errors.New("foo"),
			errors.New("bar"),
			errors.New("baz"),
		}, ""},
		{"validation", []error{
			validation.NewValueError("foo", errors.New("bar")),
			validation.NewIndexError(13, errors.New("out of bounds")),
		}, "foo: bar"},
		{"rules", []error{
			validation.NewValueError("type", validation.NewRuleError("foo", "bar")),
			validation.NewValueError("data", validation.NewRuleError("baz", "quux")),
		}, "type: bar; data: quux"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.es.Error(); got != tt.want {
				t.Errorf("Errors.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrors_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		es      validation.Errors
		want    []byte
		wantErr bool
	}{
		{"nested", []error{
			validation.NewValueError("foo", validation.NewIndexError(0, errors.New("bar"))),
			errors.New("B"),
			validation.NewValueError("baz", errors.New("quux")),
			errors.New("A"),
		}, []byte(`{"foo":{"0":"bar"},"baz":"quux"}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.es.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Errors.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Errors.MarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}
