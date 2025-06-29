package opt_test

import (
	"database/sql"
	"testing"

	"github.com/infastin/gorack/opt/v2"
)

func TestNewNull(t *testing.T) {
	tests := []struct {
		value    string
		valid    bool
		expected opt.Null[string]
	}{
		{"foo", true, opt.Null[string]{"foo", true}},
		{"foo", false, opt.Null[string]{"foo", false}},
		{"", true, opt.Null[string]{"", true}},
		{"", false, opt.Null[string]{"", false}},
	}

	for _, tt := range tests {
		got := opt.NewNull(tt.value, tt.valid)
		if got != tt.expected {
			t.Errorf("must be equal: input=(%v, %v) expected=%+v got=%+v",
				tt.value, tt.valid, got, tt.expected)
		}
	}
}

func TestNullFrom(t *testing.T) {
	tests := []struct {
		value    string
		expected opt.Null[string]
	}{
		{"foo", opt.NewNull("foo", true)},
		{"", opt.NewNull("", true)},
	}

	for _, tt := range tests {
		got := opt.NullFrom(tt.value)
		if got != tt.expected {
			t.Errorf("must be equal: input=%v expected=%+v got=%+v",
				tt.value, got, tt.expected)
		}
	}
}

func TestNullFromPtr(t *testing.T) {
	tests := []struct {
		value    *string
		expected opt.Null[string]
	}{
		{opt.Ptr("foo"), opt.NewNull("foo", true)},
		{opt.Ptr(""), opt.NewNull("", true)},
		{nil, opt.NewNull("", false)},
	}

	for _, tt := range tests {
		got := opt.NullFromPtr(tt.value)
		if got != tt.expected {
			t.Errorf("must be equal: input=%s expected=%+v got=%+v",
				PtrString(tt.value), got, tt.expected)
		}
	}
}

func TestNullFromFunc(t *testing.T) {
	fn := func(s string) string {
		return s + "bar"
	}

	tests := []struct {
		value    *string
		expected opt.Null[string]
	}{
		{opt.Ptr("foo"), opt.NewNull("foobar", true)},
		{opt.Ptr(""), opt.NewNull("bar", true)},
		{nil, opt.NewNull("", false)},
	}

	for _, tt := range tests {
		got := opt.NullFromFunc(tt.value, fn)
		if got != tt.expected {
			t.Errorf("must be equal: input=%s expected=%+v got=%+v",
				PtrString(tt.value), got, tt.expected)
		}
	}
}

func TestNullFromFuncPtr(t *testing.T) {
	fn := func(s *string) string {
		return *s + "bar"
	}

	tests := []struct {
		value    *string
		expected opt.Null[string]
	}{
		{opt.Ptr("foo"), opt.NewNull("foobar", true)},
		{opt.Ptr(""), opt.NewNull("bar", true)},
		{nil, opt.NewNull("", false)},
	}

	for _, tt := range tests {
		got := opt.NullFromFuncPtr(tt.value, fn)
		if got != tt.expected {
			t.Errorf("must be equal: input=%s expected=%+v got=%+v",
				PtrString(tt.value), got, tt.expected)
		}
	}
}

func TestNull_Set(t *testing.T) {
	tests := []struct {
		input    opt.Null[string]
		value    string
		expected opt.Null[string]
	}{
		{opt.NewNull("", false), "foo", opt.NewNull("foo", true)},
		{opt.NewNull("", false), "", opt.NewNull("", true)},
		{opt.NewNull("bar", true), "foo", opt.NewNull("foo", true)},
		{opt.NewNull("bar", true), "", opt.NewNull("", true)},
		{opt.NewNull("bar", false), "foo", opt.NewNull("foo", true)},
		{opt.NewNull("bar", false), "", opt.NewNull("", true)},
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

func TestNull_Reset(t *testing.T) {
	tests := []struct {
		input    opt.Null[string]
		expected opt.Null[string]
	}{
		{opt.NewNull("", false), opt.NewNull("", false)},
		{opt.NewNull("bar", true), opt.NewNull("", false)},
		{opt.NewNull("bar", false), opt.NewNull("", false)},
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

func TestNull_Ptr(t *testing.T) {
	tests := []struct {
		input    opt.Null[string]
		expected *string
	}{
		{opt.NewNull("foo", true), opt.Ptr("foo")},
		{opt.NewNull("", true), opt.Ptr("")},
		{opt.NewNull("foo", false), nil},
		{opt.NewNull("", false), nil},
	}

	for _, tt := range tests {
		got := tt.input.Ptr()
		if !PtrEqual(tt.expected, got) {
			t.Errorf("must be equal: input=%+v expected=%s got=%s",
				tt.input, PtrString(tt.expected), PtrString(got))
		}
	}
}

func TestNull_IsZero(t *testing.T) {
	tests := []struct {
		input    opt.Null[string]
		expected bool
	}{
		{opt.NewNull("foo", true), false},
		{opt.NewNull("", true), false},
		{opt.NewNull("foo", false), true},
		{opt.NewNull("", false), true},
	}

	for _, tt := range tests {
		got := tt.input.IsZero()
		if tt.expected != got {
			t.Errorf("must be equal: input=%+v expected=%v got=%v",
				tt.input, tt.expected, got)
		}
	}
}

func TestNull_Get(t *testing.T) {
	tests := []struct {
		input         opt.Null[string]
		expectedValue string
		expectedValid bool
	}{
		{opt.NewNull("foo", true), "foo", true},
		{opt.NewNull("", true), "", true},
		{opt.NewNull("foo", false), "foo", false},
		{opt.NewNull("", false), "", false},
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

func TestNull_ToSQL(t *testing.T) {
	tests := []struct {
		input    opt.Null[string]
		expected sql.Null[string]
	}{
		{opt.NewNull("foo", true), sql.Null[string]{"foo", true}},
		{opt.NewNull("", true), sql.Null[string]{"", true}},
		{opt.NewNull("foo", false), sql.Null[string]{"foo", false}},
		{opt.NewNull("", false), sql.Null[string]{"", false}},
	}

	for _, tt := range tests {
		got := tt.input.ToSQL()
		if tt.expected != got {
			t.Errorf("must be equal: input=%+v expected=%+v got=%+v",
				tt.input, tt.expected, got)
		}
	}
}

func TestNull_Or(t *testing.T) {
	tests := []struct {
		input    opt.Null[string]
		value    string
		expected string
	}{
		{opt.NewNull("foo", true), "bar", "foo"},
		{opt.NewNull("", true), "bar", ""},
		{opt.NewNull("foo", false), "bar", "bar"},
		{opt.NewNull("", false), "bar", "bar"},
	}

	for _, tt := range tests {
		got := tt.input.Or(tt.value)
		if tt.expected != got {
			t.Errorf("must be equal: input=%+v expected=%v got=%v",
				tt.input, tt.expected, got)
		}
	}
}

func TestNull_MarshalJSON(t *testing.T) {
	tests := []struct {
		input    opt.Null[string]
		expected string
	}{
		{opt.NewNull("foo", true), `"foo"`},
		{opt.NewNull("", true), `""`},
		{opt.NewNull("foo", false), `null`},
		{opt.NewNull("", false), `null`},
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

func TestNull_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected opt.Null[string]
	}{
		{`"foo"`, opt.NewNull("foo", true)},
		{`""`, opt.NewNull("", true)},
		{`null`, opt.NewNull("", false)},
	}

	for _, tt := range tests {
		var got opt.Null[string]
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

func TestNull_MarshalText(t *testing.T) {
	tests := []struct {
		input    opt.Null[string]
		expected string
	}{
		{opt.NewNull("foo", true), `foo`},
		{opt.NewNull("", true), ``},
		{opt.NewNull("foo", false), ``},
		{opt.NewNull("", false), ``},
	}

	for _, tt := range tests {
		got, err := tt.input.MarshalText()
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

func TestNull_UnmarshalText(t *testing.T) {
	tests := []struct {
		input    string
		expected opt.Null[string]
	}{
		{`foo`, opt.NewNull("foo", true)},
		{``, opt.NewNull("", true)},
		{`null`, opt.NewNull("null", true)},
	}

	for _, tt := range tests {
		var got opt.Null[string]
		if err := got.UnmarshalText([]byte(tt.input)); err != nil {
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

func TestNull_Scan(t *testing.T) {
	tests := []struct {
		input    opt.Null[string]
		value    any
		expected opt.Null[string]
	}{
		{opt.NewNull("foo", true), "bar", opt.NewNull("bar", true)},
		{opt.NewNull("", false), "bar", opt.NewNull("bar", true)},
		{opt.NewNull("foo", true), "", opt.NewNull("", true)},
		{opt.NewNull("", false), "", opt.NewNull("", true)},
		{opt.NewNull("foo", true), nil, opt.NewNull("", false)},
		{opt.NewNull("", false), nil, opt.NewNull("", false)},
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
