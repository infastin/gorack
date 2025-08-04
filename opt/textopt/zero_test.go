package textopt_test

import (
	"testing"
	"time"

	"github.com/infastin/gorack/opt/v2/textopt"
)

func TestZero_MarshalText(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			input    textopt.Zero[string]
			expected string
		}{
			{textopt.NewZero("foo", true), `foo`},
			{textopt.NewZero("", true), ``},
			{textopt.NewZero("foo", false), ``},
			{textopt.NewZero("", false), ``},
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
			input    textopt.Zero[time.Time]
			expected string
		}{
			{
				input:    textopt.NewZero(time.Date(2006, time.January, 2, 15, 04, 05, 0, mst), true),
				expected: `2006-01-02T15:04:05-07:00`,
			},
			{
				input:    textopt.NewZero(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), true),
				expected: `0000-01-01T00:00:00Z`,
			},
			{
				input:    textopt.NewZero(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), false),
				expected: `0001-01-01T00:00:00Z`,
			},
			{
				input:    textopt.NewZero(time.Time{}, false),
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
			expected textopt.Zero[string]
		}{
			{`foo`, textopt.NewZero("foo", true)},
			{``, textopt.NewZero("", false)},
			{`null`, textopt.NewZero("null", true)},
		}

		for _, tt := range tests {
			var got textopt.Zero[string]
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
			expected textopt.Zero[time.Time]
		}{
			{
				input:    `2006-01-02T15:04:05-07:00`,
				expected: textopt.NewZero(time.Date(2006, time.January, 2, 15, 04, 05, 0, unnamedMST), true),
			},
			{
				input:    `0000-01-01T00:00:00Z`,
				expected: textopt.NewZero(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), true),
			},
			{
				input:    `0001-01-01T00:00:00Z`,
				expected: textopt.NewZero(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), false),
			},
			{
				input:    `0001-01-01T00:00:00Z`,
				expected: textopt.NewZero(time.Time{}, false),
			},
		}

		for _, tt := range tests {
			var got textopt.Zero[time.Time]
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
