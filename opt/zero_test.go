package opt_test

import (
	"database/sql"
	"image"
	"testing"
	"time"

	"github.com/infastin/gorack/opt/v2"
)

func TestNewZero(t *testing.T) {
	tests := []struct {
		value    string
		valid    bool
		expected opt.Zero[string]
	}{
		{"foo", true, opt.NewZero("foo", true)},
		{"foo", false, opt.NewZero("foo", false)},
		{"", true, opt.NewZero("", true)},
		{"", false, opt.NewZero("", false)},
	}

	for _, tt := range tests {
		got := opt.NewZero(tt.value, tt.valid)
		if got != tt.expected {
			t.Errorf("must be equal: input=(%v, %v) expected=%+v got=%+v",
				tt.value, tt.valid, got, tt.expected)
		}
	}
}

func TestZeroFrom(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			input  string
			isZero bool
		}{
			{"foo", false},
			{"bar", false},
			{"", true},
		}

		for _, tt := range tests {
			expected := opt.NewZero(tt.input, !tt.isZero)
			got := opt.ZeroFrom(tt.input)
			if expected != got {
				t.Errorf("must be equal: input=%s expected=%+v got=%+v",
					tt.input, expected, got)
			}
		}
	})

	t.Run("Time", func(t *testing.T) {
		mst, err := time.LoadLocation("MST")
		if err != nil {
			t.Fatal(err)
		}

		tests := []struct {
			input  time.Time
			isZero bool
		}{
			{time.Date(2006, time.January, 2, 15, 04, 05, 0, mst), false},
			{time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), false},
			{time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), true},
			{time.Time{}, true},
		}

		for _, tt := range tests {
			expected := opt.NewZero(tt.input, !tt.isZero)
			got := opt.ZeroFrom(tt.input)
			if expected != got {
				t.Errorf("must be equal: input=%v expected=%+v got=%+v",
					tt.input, expected, got)
			}
		}
	})

	t.Run("Rect", func(t *testing.T) {
		tests := []struct {
			input  opt.Rect
			isZero bool
		}{
			{opt.Rect{image.Pt(0, 0), image.Pt(1, 1)}, false},
			{opt.Rect{image.Pt(0, 1), image.Pt(1, 2)}, false},
			{opt.Rect{image.Pt(0, 1), image.Pt(1, 1)}, true},
			{opt.Rect{image.Pt(1, 1), image.Pt(1, 1)}, true},
			{opt.Rect{}, true},
		}

		for _, tt := range tests {
			expected := opt.NewZero(tt.input, !tt.isZero)
			got := opt.ZeroFrom(tt.input)
			if expected != got {
				t.Errorf("must be equal: input=%v expected=%+v got=%+v",
					tt.input, expected, got)
			}
		}
	})

	t.Run("RectPtr", func(t *testing.T) {
		tests := []struct {
			input  opt.RectPtr
			isZero bool
		}{
			{opt.RectPtr{image.Pt(0, 0), image.Pt(1, 1)}, false},
			{opt.RectPtr{image.Pt(0, 1), image.Pt(1, 2)}, false},
			{opt.RectPtr{image.Pt(0, 1), image.Pt(1, 1)}, true},
			{opt.RectPtr{image.Pt(1, 1), image.Pt(1, 1)}, true},
			{opt.RectPtr{}, true},
		}

		for _, tt := range tests {
			expected := opt.NewZero(tt.input, !tt.isZero)
			got := opt.ZeroFrom(tt.input)
			if expected != got {
				t.Errorf("must be equal: input=%v expected=%+v got=%+v",
					tt.input, expected, got)
			}
		}
	})
}

func TestZeroFromPtr(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			input  *string
			isZero bool
		}{
			{opt.Ptr("foo"), false},
			{opt.Ptr("bar"), false},
			{opt.Ptr(""), true},
			{nil, true},
		}

		for _, tt := range tests {
			expected := opt.NewZero(opt.Deref(tt.input, ""), !tt.isZero)
			got := opt.ZeroFromPtr(tt.input)
			if expected != got {
				t.Errorf("must be equal: input=%s expected=%+v got=%+v",
					PtrString(tt.input), expected, got)
			}
		}
	})

	t.Run("Time", func(t *testing.T) {
		mst, err := time.LoadLocation("MST")
		if err != nil {
			t.Fatal(err)
		}

		tests := []struct {
			input  *time.Time
			isZero bool
		}{
			{opt.Ptr(time.Date(2006, time.January, 2, 15, 04, 05, 0, mst)), false},
			{opt.Ptr(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)), false},
			{opt.Ptr(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)), true},
			{&time.Time{}, true},
			{nil, true},
		}

		for _, tt := range tests {
			expected := opt.NewZero(opt.Deref(tt.input, time.Time{}), !tt.isZero)
			got := opt.ZeroFromPtr(tt.input)
			if expected != got {
				t.Errorf("must be equal: input=%v expected=%+v got=%+v",
					PtrString(tt.input), expected, got)
			}
		}
	})

	t.Run("Rect", func(t *testing.T) {
		tests := []struct {
			input  *opt.Rect
			isZero bool
		}{
			{&opt.Rect{image.Pt(0, 0), image.Pt(1, 1)}, false},
			{&opt.Rect{image.Pt(0, 1), image.Pt(1, 2)}, false},
			{&opt.Rect{image.Pt(0, 1), image.Pt(1, 1)}, true},
			{&opt.Rect{image.Pt(1, 1), image.Pt(1, 1)}, true},
			{&opt.Rect{}, true},
			{nil, true},
		}

		for _, tt := range tests {
			expected := opt.NewZero(opt.Deref(tt.input, opt.Rect{}), !tt.isZero)
			got := opt.ZeroFromPtr(tt.input)
			if expected != got {
				t.Errorf("must be equal: input=%v expected=%+v got=%+v",
					PtrString(tt.input), expected, got)
			}
		}
	})

	t.Run("RectPtr", func(t *testing.T) {
		tests := []struct {
			input  *opt.RectPtr
			isZero bool
		}{
			{&opt.RectPtr{image.Pt(0, 0), image.Pt(1, 1)}, false},
			{&opt.RectPtr{image.Pt(0, 1), image.Pt(1, 2)}, false},
			{&opt.RectPtr{image.Pt(0, 1), image.Pt(1, 1)}, true},
			{&opt.RectPtr{image.Pt(1, 1), image.Pt(1, 1)}, true},
			{&opt.RectPtr{}, true},
			{nil, true},
		}

		for _, tt := range tests {
			expected := opt.NewZero(opt.Deref(tt.input, opt.RectPtr{}), !tt.isZero)
			got := opt.ZeroFromPtr(tt.input)
			if expected != got {
				t.Errorf("must be equal: input=%v expected=%+v got=%+v",
					PtrString(tt.input), expected, got)
			}
		}
	})
}

func TestZeroFromFunc(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		fn := func(s string) string {
			return s + "bar"
		}

		tests := []struct {
			value    *string
			expected opt.Zero[string]
		}{
			{opt.Ptr("foo"), opt.NewZero("foobar", true)},
			{opt.Ptr(""), opt.NewZero("bar", true)},
			{nil, opt.NewZero("", false)},
		}

		for _, tt := range tests {
			got := opt.ZeroFromFunc(tt.value, fn)
			if got != tt.expected {
				t.Errorf("must be equal: input=%s expected=%+v got=%+v",
					PtrString(tt.value), got, tt.expected)
			}
		}
	})

	t.Run("Time", func(t *testing.T) {
		mst, err := time.LoadLocation("MST")
		if err != nil {
			t.Fatal(err)
		}

		fn := func(t time.Time) time.Time {
			return t.AddDate(0, 0, 1)
		}

		tests := []struct {
			value    *time.Time
			expected opt.Zero[time.Time]
		}{
			{
				value:    opt.Ptr(time.Date(2006, time.January, 2, 15, 04, 05, 0, mst)),
				expected: opt.NewZero(time.Date(2006, time.January, 3, 15, 04, 05, 0, mst), true),
			},
			{
				value:    opt.Ptr(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)),
				expected: opt.NewZero(time.Date(0, time.January, 2, 0, 0, 0, 0, time.UTC), true),
			},
			{
				value:    opt.Ptr(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)),
				expected: opt.NewZero(time.Date(1, time.January, 2, 0, 0, 0, 0, time.UTC), true),
			},
			{
				value:    nil,
				expected: opt.NewZero(time.Time{}, false),
			},
		}

		for _, tt := range tests {
			got := opt.ZeroFromFunc(tt.value, fn)
			if got != tt.expected {
				t.Errorf("must be equal: input=%s expected=%+v got=%+v",
					PtrString(tt.value), got, tt.expected)
			}
		}
	})
}

func TestZeroFromFuncPtr(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		fn := func(s *string) string {
			return *s + "bar"
		}

		tests := []struct {
			value    *string
			expected opt.Zero[string]
		}{
			{opt.Ptr("foo"), opt.NewZero("foobar", true)},
			{opt.Ptr(""), opt.NewZero("bar", true)},
			{nil, opt.NewZero("", false)},
		}

		for _, tt := range tests {
			got := opt.ZeroFromFuncPtr(tt.value, fn)
			if got != tt.expected {
				t.Errorf("must be equal: input=%s expected=%+v got=%+v",
					PtrString(tt.value), got, tt.expected)
			}
		}
	})

	t.Run("Time", func(t *testing.T) {
		mst, err := time.LoadLocation("MST")
		if err != nil {
			t.Fatal(err)
		}

		fn := func(t *time.Time) time.Time {
			return t.AddDate(0, 0, 1)
		}

		tests := []struct {
			value    *time.Time
			expected opt.Zero[time.Time]
		}{
			{
				value:    opt.Ptr(time.Date(2006, time.January, 2, 15, 04, 05, 0, mst)),
				expected: opt.NewZero(time.Date(2006, time.January, 3, 15, 04, 05, 0, mst), true),
			},
			{
				value:    opt.Ptr(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)),
				expected: opt.NewZero(time.Date(0, time.January, 2, 0, 0, 0, 0, time.UTC), true),
			},
			{
				value:    opt.Ptr(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)),
				expected: opt.NewZero(time.Date(1, time.January, 2, 0, 0, 0, 0, time.UTC), true),
			},
			{
				value:    nil,
				expected: opt.NewZero(time.Time{}, false),
			},
		}

		for _, tt := range tests {
			got := opt.ZeroFromFuncPtr(tt.value, fn)
			if got != tt.expected {
				t.Errorf("must be equal: input=%s expected=%+v got=%+v",
					PtrString(tt.value), got, tt.expected)
			}
		}
	})
}

func TestZero_Set(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			input    opt.Zero[string]
			value    string
			expected opt.Zero[string]
		}{
			{opt.NewZero("", false), "foo", opt.NewZero("foo", true)},
			{opt.NewZero("", false), "", opt.NewZero("", false)},
			{opt.NewZero("bar", true), "foo", opt.NewZero("foo", true)},
			{opt.NewZero("bar", true), "", opt.NewZero("", false)},
			{opt.NewZero("bar", false), "foo", opt.NewZero("foo", true)},
			{opt.NewZero("bar", false), "", opt.NewZero("", false)},
		}

		for _, tt := range tests {
			got := tt.input
			got.Set(tt.value)
			if got != tt.expected {
				t.Errorf("must be equal: input=(%+v, %v) expected=%+v got=%+v",
					tt.input, tt.value, got, tt.expected)
			}
		}
	})

	t.Run("Time", func(t *testing.T) {
		mst, err := time.LoadLocation("MST")
		if err != nil {
			t.Fatal(err)
		}

		tests := []struct {
			input    opt.Zero[time.Time]
			value    time.Time
			expected opt.Zero[time.Time]
		}{
			{
				input:    opt.NewZero(time.Time{}, false),
				value:    time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC),
				expected: opt.NewZero(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), true),
			},
			{
				input:    opt.NewZero(time.Date(2006, time.January, 2, 15, 04, 05, 0, mst), true),
				value:    time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC),
				expected: opt.NewZero(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), true),
			},
			{
				input:    opt.NewZero(time.Date(2006, time.January, 2, 15, 04, 05, 0, mst), true),
				value:    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				expected: opt.NewZero(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), false),
			},
			{
				input:    opt.NewZero(time.Date(2006, time.January, 2, 15, 04, 05, 0, mst), true),
				value:    time.Time{},
				expected: opt.NewZero(time.Time{}, false),
			},
			{
				input:    opt.NewZero(time.Time{}, false),
				value:    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				expected: opt.NewZero(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), false),
			},
		}

		for _, tt := range tests {
			got := tt.input
			got.Set(tt.value)
			if got != tt.expected {
				t.Errorf("must be equal: input=(%+v, %v) expected=%+v got=%+v",
					tt.input, tt.value, got, tt.expected)
			}
		}
	})

	t.Run("Rect", func(t *testing.T) {
		tests := []struct {
			input    opt.Zero[opt.Rect]
			value    opt.Rect
			expected opt.Zero[opt.Rect]
		}{
			{
				input:    opt.NewZero(opt.Rect{image.Pt(1, 1), image.Pt(1, 1)}, false),
				value:    opt.Rect{image.Pt(0, 1), image.Pt(1, 2)},
				expected: opt.NewZero(opt.Rect{image.Pt(0, 1), image.Pt(1, 2)}, true),
			},
			{
				input:    opt.NewZero(opt.Rect{image.Pt(0, 0), image.Pt(1, 1)}, true),
				value:    opt.Rect{image.Pt(0, 1), image.Pt(1, 2)},
				expected: opt.NewZero(opt.Rect{image.Pt(0, 1), image.Pt(1, 2)}, true),
			},
			{
				input:    opt.NewZero(opt.Rect{image.Pt(0, 0), image.Pt(1, 1)}, true),
				value:    opt.Rect{image.Pt(0, 1), image.Pt(1, 1)},
				expected: opt.NewZero(opt.Rect{image.Pt(0, 1), image.Pt(1, 1)}, false),
			},
			{
				input:    opt.NewZero(opt.Rect{image.Pt(0, 0), image.Pt(1, 1)}, true),
				value:    opt.Rect{image.Pt(0, 1), image.Pt(1, 1)},
				expected: opt.NewZero(opt.Rect{image.Pt(0, 1), image.Pt(1, 1)}, false),
			},
			{
				input:    opt.NewZero(opt.Rect{image.Pt(1, 1), image.Pt(1, 1)}, false),
				value:    opt.Rect{image.Pt(0, 0), image.Pt(0, 0)},
				expected: opt.NewZero(opt.Rect{image.Pt(0, 0), image.Pt(0, 0)}, false),
			},
		}

		for _, tt := range tests {
			got := tt.input
			got.Set(tt.value)
			if got != tt.expected {
				t.Errorf("must be equal: input=(%+v, %v) expected=%+v got=%+v",
					tt.input, tt.value, got, tt.expected)
			}
		}
	})
}

func TestZero_Reset(t *testing.T) {
	tests := []struct {
		input    opt.Zero[string]
		expected opt.Zero[string]
	}{
		{opt.NewZero("", false), opt.NewZero("", false)},
		{opt.NewZero("bar", true), opt.NewZero("", false)},
		{opt.NewZero("bar", false), opt.NewZero("", false)},
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

func TestZero_Ptr(t *testing.T) {
	tests := []struct {
		input    opt.Zero[string]
		expected *string
	}{
		{opt.NewZero("foo", true), opt.Ptr("foo")},
		{opt.NewZero("", true), opt.Ptr("")},
		{opt.NewZero("foo", false), nil},
		{opt.NewZero("", false), nil},
	}

	for _, tt := range tests {
		got := tt.input.Ptr()
		if !PtrEqual(tt.expected, got) {
			t.Errorf("must be equal: input=%+v expected=%s got=%s",
				tt.input, PtrString(tt.expected), PtrString(got))
		}
	}
}

func TestZero_IsZero(t *testing.T) {
	tests := []struct {
		input    opt.Zero[string]
		expected bool
	}{
		{opt.NewZero("foo", true), false},
		{opt.NewZero("", true), false},
		{opt.NewZero("foo", false), true},
		{opt.NewZero("", false), true},
	}

	for _, tt := range tests {
		got := tt.input.IsZero()
		if tt.expected != got {
			t.Errorf("must be equal: input=%+v expected=%v got=%v",
				tt.input, tt.expected, got)
		}
	}
}

func TestZero_Get(t *testing.T) {
	tests := []struct {
		input         opt.Zero[string]
		expectedValue string
		expectedValid bool
	}{
		{opt.NewZero("foo", true), "foo", true},
		{opt.NewZero("", true), "", true},
		{opt.NewZero("foo", false), "foo", false},
		{opt.NewZero("", false), "", false},
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

func TestZero_ToSQL(t *testing.T) {
	tests := []struct {
		input    opt.Zero[string]
		expected sql.Null[string]
	}{
		{opt.NewZero("foo", true), sql.Null[string]{"foo", true}},
		{opt.NewZero("", true), sql.Null[string]{"", true}},
		{opt.NewZero("foo", false), sql.Null[string]{"foo", false}},
		{opt.NewZero("", false), sql.Null[string]{"", false}},
	}

	for _, tt := range tests {
		got := tt.input.ToSQL()
		if tt.expected != got {
			t.Errorf("must be equal: input=%+v expected=%+v got=%+v",
				tt.input, tt.expected, got)
		}
	}
}

func TestZero_Or(t *testing.T) {
	tests := []struct {
		input    opt.Zero[string]
		value    string
		expected string
	}{
		{opt.NewZero("foo", true), "bar", "foo"},
		{opt.NewZero("", true), "bar", ""},
		{opt.NewZero("foo", false), "bar", "bar"},
		{opt.NewZero("", false), "bar", "bar"},
	}

	for _, tt := range tests {
		got := tt.input.Or(tt.value)
		if tt.expected != got {
			t.Errorf("must be equal: input=%+v expected=%v got=%v",
				tt.input, tt.expected, got)
		}
	}
}

func TestZero_MarshalJSON(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			input    opt.Zero[string]
			expected string
		}{
			{opt.NewZero("foo", true), `"foo"`},
			{opt.NewZero("", true), `""`},
			{opt.NewZero("foo", false), `""`},
			{opt.NewZero("", false), `""`},
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
	})

	t.Run("Time", func(t *testing.T) {
		mst, err := time.LoadLocation("MST")
		if err != nil {
			t.Fatal(err)
		}

		tests := []struct {
			input    opt.Zero[time.Time]
			expected string
		}{
			{
				input:    opt.NewZero(time.Date(2006, time.January, 2, 15, 04, 05, 0, mst), true),
				expected: `"2006-01-02T15:04:05-07:00"`,
			},
			{
				input:    opt.NewZero(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), true),
				expected: `"0000-01-01T00:00:00Z"`,
			},
			{
				input:    opt.NewZero(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), false),
				expected: `"0001-01-01T00:00:00Z"`,
			},
			{
				input:    opt.NewZero(time.Time{}, false),
				expected: `"0001-01-01T00:00:00Z"`,
			},
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
	})
}

func TestZero_UnmarshalJSON(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			input    string
			expected opt.Zero[string]
		}{
			{`"foo"`, opt.NewZero("foo", true)},
			{`""`, opt.NewZero("", false)},
			{`null`, opt.NewZero("", false)},
		}

		for _, tt := range tests {
			var got opt.Zero[string]
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
	})

	t.Run("Time", func(t *testing.T) {
		unnamedMST := time.FixedZone("", -7*60*60)

		tests := []struct {
			input    string
			expected opt.Zero[time.Time]
		}{
			{
				input:    `"2006-01-02T15:04:05-07:00"`,
				expected: opt.NewZero(time.Date(2006, time.January, 2, 15, 04, 05, 0, unnamedMST), true),
			},
			{
				input:    `"0000-01-01T00:00:00Z"`,
				expected: opt.NewZero(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), true),
			},
			{
				input:    `"0001-01-01T00:00:00Z"`,
				expected: opt.NewZero(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), false),
			},
			{
				input:    `"0001-01-01T00:00:00Z"`,
				expected: opt.NewZero(time.Time{}, false),
			},
		}

		for _, tt := range tests {
			var got opt.Zero[time.Time]
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
	})
}

func TestZero_MarshalText(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			input    opt.Zero[string]
			expected string
		}{
			{opt.NewZero("foo", true), `foo`},
			{opt.NewZero("", true), ``},
			{opt.NewZero("foo", false), ``},
			{opt.NewZero("", false), ``},
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
	})

	t.Run("Time", func(t *testing.T) {
		mst, err := time.LoadLocation("MST")
		if err != nil {
			t.Fatal(err)
		}

		tests := []struct {
			input    opt.Zero[time.Time]
			expected string
		}{
			{
				input:    opt.NewZero(time.Date(2006, time.January, 2, 15, 04, 05, 0, mst), true),
				expected: `2006-01-02T15:04:05-07:00`,
			},
			{
				input:    opt.NewZero(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), true),
				expected: `0000-01-01T00:00:00Z`,
			},
			{
				input:    opt.NewZero(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), false),
				expected: `0001-01-01T00:00:00Z`,
			},
			{
				input:    opt.NewZero(time.Time{}, false),
				expected: `0001-01-01T00:00:00Z`,
			},
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
	})
}

func TestZero_UnmarshalText(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			input    string
			expected opt.Zero[string]
		}{
			{`foo`, opt.NewZero("foo", true)},
			{``, opt.NewZero("", false)},
			{`null`, opt.NewZero("null", true)},
		}

		for _, tt := range tests {
			var got opt.Zero[string]
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
	})

	t.Run("Time", func(t *testing.T) {
		unnamedMST := time.FixedZone("", -7*60*60)

		tests := []struct {
			input    string
			expected opt.Zero[time.Time]
		}{
			{
				input:    `2006-01-02T15:04:05-07:00`,
				expected: opt.NewZero(time.Date(2006, time.January, 2, 15, 04, 05, 0, unnamedMST), true),
			},
			{
				input:    `0000-01-01T00:00:00Z`,
				expected: opt.NewZero(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), true),
			},
			{
				input:    `0001-01-01T00:00:00Z`,
				expected: opt.NewZero(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), false),
			},
			{
				input:    `0001-01-01T00:00:00Z`,
				expected: opt.NewZero(time.Time{}, false),
			},
		}

		for _, tt := range tests {
			var got opt.Zero[time.Time]
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
	})
}

func TestZero_Scan(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			input    opt.Zero[string]
			value    any
			expected opt.Zero[string]
		}{
			{opt.NewZero("foo", true), "bar", opt.NewZero("bar", true)},
			{opt.NewZero("", false), "bar", opt.NewZero("bar", true)},
			{opt.NewZero("foo", true), "", opt.NewZero("", false)},
			{opt.NewZero("", false), "", opt.NewZero("", false)},
			{opt.NewZero("foo", true), nil, opt.NewZero("", false)},
			{opt.NewZero("", false), nil, opt.NewZero("", false)},
		}

		for _, tt := range tests {
			got := tt.input
			got.Scan(tt.value)
			if tt.expected != got {
				t.Errorf("must be equal: input=(%+v, %v) expected=%+v got=%+v",
					tt.input, tt.value, tt.expected, got)
			}
		}
	})

	t.Run("Time", func(t *testing.T) {
		mst, err := time.LoadLocation("MST")
		if err != nil {
			t.Fatal(err)
		}

		tests := []struct {
			input    opt.Zero[time.Time]
			value    any
			expected opt.Zero[time.Time]
		}{
			{
				input:    opt.NewZero(time.Time{}, false),
				value:    time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC),
				expected: opt.NewZero(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), true),
			},
			{
				input:    opt.NewZero(time.Date(2006, time.January, 2, 15, 04, 05, 0, mst), true),
				value:    time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC),
				expected: opt.NewZero(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), true),
			},
			{
				input:    opt.NewZero(time.Date(2006, time.January, 2, 15, 04, 05, 0, mst), true),
				value:    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				expected: opt.NewZero(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), false),
			},
			{
				input:    opt.NewZero(time.Date(2006, time.January, 2, 15, 04, 05, 0, mst), true),
				value:    time.Time{},
				expected: opt.NewZero(time.Time{}, false),
			},
			{
				input:    opt.NewZero(time.Time{}, false),
				value:    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				expected: opt.NewZero(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), false),
			},
			{
				input:    opt.NewZero(time.Date(2006, time.January, 2, 15, 04, 05, 0, mst), true),
				value:    nil,
				expected: opt.NewZero(time.Time{}, false),
			},
			{
				input:    opt.NewZero(time.Time{}, false),
				value:    nil,
				expected: opt.NewZero(time.Time{}, false),
			},
		}

		for _, tt := range tests {
			got := tt.input
			got.Scan(tt.value)
			if tt.expected != got {
				t.Errorf("must be equal: input=(%+v, %v) expected=%+v got=%+v",
					tt.input, tt.value, tt.expected, got)
			}
		}
	})
}
