package opt

import (
	"image"
	"testing"
	"time"
)

type Rect image.Rectangle

func (r Rect) IsZero() bool {
	return image.Rectangle(r).Empty()
}

type RectPtr image.Rectangle

func (r *RectPtr) IsZero() bool {
	return (*image.Rectangle)(r).Empty()
}

func Test_isZero_string(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"foo", false},
		{"bar", false},
		{"", true},
	}

	for _, tt := range tests {
		got := isZero(tt.input)
		if tt.expected != got {
			t.Errorf("must be equal: input=%s expected=%v got=%v",
				tt.input, tt.expected, got)
		}
	}
}

func Test_isZero_Time(t *testing.T) {
	mst, err := time.LoadLocation("MST")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		input    time.Time
		expected bool
	}{
		{time.Date(2006, time.January, 2, 15, 04, 05, 0, mst), false},
		{time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), false},
		{time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), true},
		{time.Time{}, true},
	}

	for _, tt := range tests {
		got := isZero(tt.input)
		if tt.expected != got {
			t.Errorf("must be equal: input=%v expected=%v got=%v",
				tt.input, tt.expected, got)
		}
	}
}

func Test_isZero_Rect(t *testing.T) {
	tests := []struct {
		input    Rect
		expected bool
	}{
		{Rect{image.Pt(0, 0), image.Pt(1, 1)}, false},
		{Rect{image.Pt(0, 1), image.Pt(1, 2)}, false},
		{Rect{image.Pt(0, 1), image.Pt(1, 1)}, true},
		{Rect{image.Pt(1, 1), image.Pt(1, 1)}, true},
		{Rect{}, true},
	}

	for _, tt := range tests {
		got := isZero(tt.input)
		if tt.expected != got {
			t.Errorf("must be equal: input=%v expected=%v got=%v",
				tt.input, tt.expected, got)
		}
	}
}

func Test_isZero_PtrRect(t *testing.T) {
	tests := []struct {
		input    *Rect
		expected bool
	}{
		{&Rect{image.Pt(0, 0), image.Pt(1, 1)}, false},
		{&Rect{image.Pt(0, 1), image.Pt(1, 2)}, false},
		{&Rect{image.Pt(0, 1), image.Pt(1, 1)}, true},
		{&Rect{image.Pt(1, 1), image.Pt(1, 1)}, true},
		{&Rect{}, true},
		{nil, true},
	}

	for _, tt := range tests {
		got := isZero(tt.input)
		if tt.expected != got {
			t.Errorf("must be equal: input=%v expected=%v got=%v",
				tt.input, tt.expected, got)
		}
	}
}

func Test_isZero_RectPtr(t *testing.T) {
	tests := []struct {
		input    RectPtr
		expected bool
	}{
		{RectPtr{image.Pt(0, 0), image.Pt(1, 1)}, false},
		{RectPtr{image.Pt(0, 1), image.Pt(1, 2)}, false},
		{RectPtr{image.Pt(0, 1), image.Pt(1, 1)}, true},
		{RectPtr{image.Pt(1, 1), image.Pt(1, 1)}, true},
		{RectPtr{}, true},
	}

	for _, tt := range tests {
		got := isZero(tt.input)
		if tt.expected != got {
			t.Errorf("must be equal: input=%v expected=%v got=%v",
				tt.input, tt.expected, got)
		}
	}
}

func Test_isZero_iface(t *testing.T) {
	tests := []struct {
		input    isZeroer
		expected bool
	}{
		{time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), false},
		{time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), true},
		{time.Time{}, true},

		{Ptr(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)), false},
		{Ptr(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)), true},
		{&time.Time{}, true},
		{(*time.Time)(nil), true},

		{Rect{image.Pt(0, 0), image.Pt(1, 1)}, false},
		{Rect{image.Pt(0, 1), image.Pt(1, 2)}, false},
		{Rect{image.Pt(0, 1), image.Pt(1, 1)}, true},
		{Rect{image.Pt(1, 1), image.Pt(1, 1)}, true},
		{Rect{}, true},

		{&Rect{image.Pt(0, 0), image.Pt(1, 1)}, false},
		{&Rect{image.Pt(0, 1), image.Pt(1, 2)}, false},
		{&Rect{image.Pt(0, 1), image.Pt(1, 1)}, true},
		{&Rect{image.Pt(1, 1), image.Pt(1, 1)}, true},
		{&Rect{}, true},
		{(*Rect)(nil), true},
	}

	for _, tt := range tests {
		got := isZero(tt.input)
		if tt.expected != got {
			t.Errorf("must be equal: input=%v expected=%v got=%v",
				tt.input, tt.expected, got)
		}
	}
}
