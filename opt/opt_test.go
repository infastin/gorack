package opt_test

import (
	"fmt"
	"image"
	"testing"
	"time"

	"github.com/infastin/gorack/opt/v2"
)

func PtrEqual[T comparable](p1, p2 *T) bool {
	if (p1 == nil) != (p2 == nil) {
		return false
	}
	if p1 == nil || p2 == nil {
		return true
	}
	return *p1 == *p2
}

func PtrString[T any](p *T) string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprint(*p)
}

func TestConvert(t *testing.T) {
	fn := func(s string) string {
		return s + "bar"
	}

	tests := []struct {
		input    opt.Null[string]
		expected string
	}{
		{opt.NewNull("foo", true), "foobar"},
		{opt.NewNull("", true), "bar"},
		{opt.NewNull("foo", false), ""},
		{opt.NewNull("", false), ""},
	}

	for _, tt := range tests {
		got := opt.Convert(tt.input, fn)
		if tt.expected != got {
			t.Errorf("must be equal: input=%+v expected=%s got=%s",
				tt.input, tt.expected, got)
		}
	}
}

func TestPtr(t *testing.T) {
	tests := []string{"foo", "bar", ""}
	for _, tt := range tests {
		got := opt.Ptr(tt)
		if got == nil {
			t.Errorf("must not be nil: input=%s", tt)
			continue
		}
		if *got != tt {
			t.Errorf("must be equal: expected=%s got=%s", tt, *got)
		}
	}
}

func TestZeroPtr(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			input string
			isNil bool
		}{
			{"foo", false},
			{"bar", false},
			{"", true},
		}

		for _, tt := range tests {
			var expected *string
			if !tt.isNil {
				expected = opt.Ptr(tt.input)
			}
			got := opt.ZeroPtr(tt.input)
			if !PtrEqual(expected, got) {
				t.Errorf("must be equal: input=%s expected=%s got=%s",
					tt.input, PtrString(expected), PtrString(got))
			}
		}
	})

	t.Run("Time", func(t *testing.T) {
		mst, err := time.LoadLocation("MST")
		if err != nil {
			t.Fatal(err)
		}

		tests := []struct {
			input time.Time
			isNil bool
		}{
			{time.Date(2006, time.January, 2, 15, 04, 05, 0, mst), false},
			{time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), false},
			{time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), true},
			{time.Time{}, true},
		}

		for _, tt := range tests {
			var expected *time.Time
			if !tt.isNil {
				expected = opt.Ptr(tt.input)
			}
			got := opt.ZeroPtr(tt.input)
			if !PtrEqual(expected, got) {
				t.Errorf("must be equal: input=%v expected=%s got=%s",
					tt.input, PtrString(expected), PtrString(got))
			}
		}
	})

	t.Run("Rect", func(t *testing.T) {
		tests := []struct {
			input opt.Rect
			isNil bool
		}{
			{opt.Rect{image.Pt(0, 0), image.Pt(1, 1)}, false},
			{opt.Rect{image.Pt(0, 1), image.Pt(1, 2)}, false},
			{opt.Rect{image.Pt(0, 1), image.Pt(1, 1)}, true},
			{opt.Rect{image.Pt(1, 1), image.Pt(1, 1)}, true},
			{opt.Rect{}, true},
		}

		for _, tt := range tests {
			var expected *opt.Rect
			if !tt.isNil {
				expected = opt.Ptr(tt.input)
			}
			got := opt.ZeroPtr(tt.input)
			if !PtrEqual(expected, got) {
				t.Errorf("must be equal: input=%v expected=%s got=%s",
					tt.input, PtrString(expected), PtrString(got))
			}
		}
	})

	t.Run("RectPtr", func(t *testing.T) {
		tests := []struct {
			input opt.RectPtr
			isNil bool
		}{
			{opt.RectPtr{image.Pt(0, 0), image.Pt(1, 1)}, false},
			{opt.RectPtr{image.Pt(0, 1), image.Pt(1, 2)}, false},
			{opt.RectPtr{image.Pt(0, 1), image.Pt(1, 1)}, true},
			{opt.RectPtr{image.Pt(1, 1), image.Pt(1, 1)}, true},
			{opt.RectPtr{}, true},
		}

		for _, tt := range tests {
			var expected *opt.RectPtr
			if !tt.isNil {
				expected = opt.Ptr(tt.input)
			}
			got := opt.ZeroPtr(tt.input)
			if !PtrEqual(expected, got) {
				t.Errorf("must be equal: input=%v expected=%s got=%s",
					tt.input, PtrString(expected), PtrString(got))
			}
		}
	})
}

func TestConvertPtr(t *testing.T) {
	fn := func(s string) string {
		return s + "bar"
	}

	tests := []struct {
		input    *string
		expected *string
	}{
		{opt.Ptr("foo"), opt.Ptr("foobar")},
		{opt.Ptr(""), opt.Ptr("bar")},
		{nil, nil},
	}

	for _, tt := range tests {
		got := opt.ConvertPtr(tt.input, fn)
		if !PtrEqual(tt.expected, got) {
			t.Errorf("must be equal: input=%s expected=%s got=%s",
				PtrString(tt.input), PtrString(tt.expected), PtrString(got))
		}
	}
}

func TestDeref(t *testing.T) {
	tests := []struct {
		input    *string
		def      string
		expected string
	}{
		{opt.Ptr("foo"), "bar", "foo"},
		{opt.Ptr(""), "bar", ""},
		{nil, "bar", "bar"},
	}

	for _, tt := range tests {
		got := opt.Deref(tt.input, tt.def)
		if got != tt.expected {
			t.Errorf("must be equal: input=(%s, %s) expected=%s got=%s",
				PtrString(tt.input), tt.def, tt.expected, got)
		}
	}
}
