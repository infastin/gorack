package opt_test

import (
	"database/sql"
	"testing"

	"github.com/infastin/gorack/opt/v2"
)

func TestNewUndefined(t *testing.T) {
	tests := []struct {
		value    string
		valid    bool
		expected opt.Undefined[string]
	}{
		{"foo", true, opt.Undefined[string]{"foo", true}},
		{"foo", false, opt.Undefined[string]{"foo", false}},
		{"", true, opt.Undefined[string]{"", true}},
		{"", false, opt.Undefined[string]{"", false}},
	}

	for _, tt := range tests {
		got := opt.NewUndefined(tt.value, tt.valid)
		if got != tt.expected {
			t.Errorf("must be equal: input=(%v, %v) expected=%+v got=%+v",
				tt.value, tt.valid, got, tt.expected)
		}
	}
}

func TestUndefinedFrom(t *testing.T) {
	tests := []struct {
		value    string
		expected opt.Undefined[string]
	}{
		{"foo", opt.NewUndefined("foo", true)},
		{"", opt.NewUndefined("", true)},
	}

	for _, tt := range tests {
		got := opt.UndefinedFrom(tt.value)
		if got != tt.expected {
			t.Errorf("must be equal: input=%v expected=%+v got=%+v",
				tt.value, got, tt.expected)
		}
	}
}

func TestUndefinedFromPtr(t *testing.T) {
	tests := []struct {
		value    *string
		expected opt.Undefined[string]
	}{
		{opt.Ptr("foo"), opt.NewUndefined("foo", true)},
		{opt.Ptr(""), opt.NewUndefined("", true)},
		{nil, opt.NewUndefined("", true)},
	}

	for _, tt := range tests {
		got := opt.UndefinedFromPtr(tt.value)
		if got != tt.expected {
			t.Errorf("must be equal: input=%s expected=%+v got=%+v",
				PtrString(tt.value), got, tt.expected)
		}
	}
}

func TestUndefinedFromFunc(t *testing.T) {
	fn := func(s string) string {
		return s + "bar"
	}

	tests := []struct {
		value    *string
		expected opt.Undefined[string]
	}{
		{opt.Ptr("foo"), opt.NewUndefined("foobar", true)},
		{opt.Ptr(""), opt.NewUndefined("bar", true)},
		{nil, opt.NewUndefined("bar", true)},
	}

	for _, tt := range tests {
		got := opt.UndefinedFromFunc(tt.value, fn)
		if got != tt.expected {
			t.Errorf("must be equal: input=%s expected=%+v got=%+v",
				PtrString(tt.value), got, tt.expected)
		}
	}
}

func TestUndefinedFromFuncPtr(t *testing.T) {
	fn := func(s *string) string {
		return *s + "bar"
	}

	tests := []struct {
		value    *string
		expected opt.Undefined[string]
	}{
		{opt.Ptr("foo"), opt.NewUndefined("foobar", true)},
		{opt.Ptr(""), opt.NewUndefined("bar", true)},
		{nil, opt.NewUndefined("bar", true)},
	}

	for _, tt := range tests {
		got := opt.UndefinedFromFuncPtr(tt.value, fn)
		if got != tt.expected {
			t.Errorf("must be equal: input=%s expected=%+v got=%+v",
				PtrString(tt.value), got, tt.expected)
		}
	}
}

func TestUndefined_Set(t *testing.T) {
	tests := []struct {
		input    opt.Undefined[string]
		value    string
		expected opt.Undefined[string]
	}{
		{opt.NewUndefined("", false), "foo", opt.NewUndefined("foo", true)},
		{opt.NewUndefined("", false), "", opt.NewUndefined("", true)},
		{opt.NewUndefined("bar", true), "foo", opt.NewUndefined("foo", true)},
		{opt.NewUndefined("bar", true), "", opt.NewUndefined("", true)},
		{opt.NewUndefined("bar", false), "foo", opt.NewUndefined("foo", true)},
		{opt.NewUndefined("bar", false), "", opt.NewUndefined("", true)},
	}

	for _, tt := range tests {
		got := tt.input
		got.Set(tt.value)
		if got != tt.expected {
			t.Errorf("must be equal: input=(%+v, %v) expected=%+v got=%+v",
				tt.input, tt.value, got, tt.expected)
		}
	}
}

func TestUndefined_Reset(t *testing.T) {
	tests := []struct {
		input    opt.Undefined[string]
		expected opt.Undefined[string]
	}{
		{opt.NewUndefined("", false), opt.NewUndefined("", false)},
		{opt.NewUndefined("bar", true), opt.NewUndefined("", false)},
		{opt.NewUndefined("bar", false), opt.NewUndefined("", false)},
	}

	for _, tt := range tests {
		got := tt.input
		got.Reset()
		if got != tt.expected {
			t.Errorf("must be equal: input=%+v expected=%+v got=%+v",
				tt.input, got, tt.expected)
		}
	}
}

func TestUndefined_Ptr(t *testing.T) {
	tests := []struct {
		input    opt.Undefined[string]
		expected *string
	}{
		{opt.NewUndefined("foo", true), opt.Ptr("foo")},
		{opt.NewUndefined("", true), opt.Ptr("")},
		{opt.NewUndefined("foo", false), nil},
		{opt.NewUndefined("", false), nil},
	}

	for _, tt := range tests {
		got := tt.input.Ptr()
		if !PtrEqual(tt.expected, got) {
			t.Errorf("must be equal: input=%+v expected=%s got=%s",
				tt.input, PtrString(tt.expected), PtrString(got))
		}
	}
}

func TestUndefined_IsZero(t *testing.T) {
	tests := []struct {
		input    opt.Undefined[string]
		expected bool
	}{
		{opt.NewUndefined("foo", true), false},
		{opt.NewUndefined("", true), false},
		{opt.NewUndefined("foo", false), true},
		{opt.NewUndefined("", false), true},
	}

	for _, tt := range tests {
		got := tt.input.IsZero()
		if tt.expected != got {
			t.Errorf("must be equal: input=%+v expected=%v got=%v",
				tt.input, tt.expected, got)
		}
	}
}

func TestUndefined_Get(t *testing.T) {
	tests := []struct {
		input         opt.Undefined[string]
		expectedValue string
		expectedValid bool
	}{
		{opt.NewUndefined("foo", true), "foo", true},
		{opt.NewUndefined("", true), "", true},
		{opt.NewUndefined("foo", false), "foo", false},
		{opt.NewUndefined("", false), "", false},
	}

	for _, tt := range tests {
		gotValue, gotValid := tt.input.Get()
		if tt.expectedValue != gotValue {
			t.Errorf("values must be equal: input=%+v expected=%v got=%v",
				tt.input, tt.expectedValid, gotValue)
		}
		if tt.expectedValid != gotValid {
			t.Errorf("valids must be equal: input=%+v expected=%v got=%v",
				tt.input, tt.expectedValid, gotValid)
		}
	}
}

func TestUndefined_ToSQL(t *testing.T) {
	tests := []struct {
		input    opt.Undefined[string]
		expected sql.Null[string]
	}{
		{opt.NewUndefined("foo", true), sql.Null[string]{"foo", true}},
		{opt.NewUndefined("", true), sql.Null[string]{"", true}},
		{opt.NewUndefined("foo", false), sql.Null[string]{"foo", false}},
		{opt.NewUndefined("", false), sql.Null[string]{"", false}},
	}

	for _, tt := range tests {
		got := tt.input.ToSQL()
		if tt.expected != got {
			t.Errorf("must be equal: input=%+v expected=%+v got=%+v",
				tt.input, tt.expected, got)
		}
	}
}

func TestUndefined_Or(t *testing.T) {
	tests := []struct {
		input    opt.Undefined[string]
		value    string
		expected string
	}{
		{opt.NewUndefined("foo", true), "bar", "foo"},
		{opt.NewUndefined("", true), "bar", ""},
		{opt.NewUndefined("foo", false), "bar", "bar"},
		{opt.NewUndefined("", false), "bar", "bar"},
	}

	for _, tt := range tests {
		got := tt.input.Or(tt.value)
		if tt.expected != got {
			t.Errorf("must be equal: input=%+v expected=%v got=%v",
				tt.input, tt.expected, got)
		}
	}
}

func TestUndefined_MarshalJSON(t *testing.T) {
	tests := []struct {
		input    opt.Undefined[string]
		expected string
	}{
		{opt.NewUndefined("foo", true), `"foo"`},
		{opt.NewUndefined("", true), `""`},
		{opt.NewUndefined("foo", false), `null`},
		{opt.NewUndefined("", false), `null`},
	}

	for _, tt := range tests {
		got, err := tt.input.MarshalJSON()
		if err != nil {
			t.Errorf("unexpected error: input=%+v error=%s",
				tt.input, err.Error())
			continue
		}
		if tt.expected != string(got) {
			t.Errorf("must be equal: input=%+v expected=%s got=%s",
				tt.input, tt.expected, got)
		}
	}
}

func TestUndefined_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected opt.Undefined[string]
	}{
		{`"foo"`, opt.NewUndefined("foo", true)},
		{`""`, opt.NewUndefined("", true)},
		{`null`, opt.NewUndefined("", true)},
	}

	for _, tt := range tests {
		var got opt.Undefined[string]
		if err := got.UnmarshalJSON([]byte(tt.input)); err != nil {
			t.Errorf("unexpected error: input=%s error=%s",
				tt.input, err.Error())
			continue
		}
		if tt.expected != got {
			t.Errorf("must be equal: input=%s expected=%+v got=%+v",
				tt.input, tt.expected, got)
		}
	}
}

func TestUndefined_Scan(t *testing.T) {
	tests := []struct {
		input    opt.Undefined[string]
		value    any
		expected opt.Undefined[string]
	}{
		{opt.NewUndefined("foo", true), "bar", opt.NewUndefined("bar", true)},
		{opt.NewUndefined("", false), "bar", opt.NewUndefined("bar", true)},
		{opt.NewUndefined("foo", true), "", opt.NewUndefined("", true)},
		{opt.NewUndefined("", false), "", opt.NewUndefined("", true)},
		{opt.NewUndefined("foo", true), nil, opt.NewUndefined("", true)},
		{opt.NewUndefined("", false), nil, opt.NewUndefined("", true)},
	}

	for _, tt := range tests {
		got := tt.input
		got.Scan(tt.value)
		if tt.expected != got {
			t.Errorf("must be equal: input=(%+v, %v) expected=%+v got=%+v",
				tt.input, tt.value, tt.expected, got)
		}
	}
}
