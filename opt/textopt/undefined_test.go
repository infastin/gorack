package textopt_test

import (
	"testing"

	"github.com/infastin/gorack/opt/v2/textopt"
)

func TestUndefined_MarshalText(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			input    textopt.Undefined[string]
			expected string
		}{
			{textopt.NewUndefined("foo", true), `foo`},
			{textopt.NewUndefined("", true), ``},
			{textopt.NewUndefined("foo", false), ``},
			{textopt.NewUndefined("", false), ``},
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

	t.Run("int", func(t *testing.T) {
		tests := []struct {
			input    textopt.Undefined[int]
			expected string
		}{
			{textopt.NewUndefined(0, false), ``},
			{textopt.NewUndefined(0, true), `0`},
			{textopt.NewUndefined(42, false), ``},
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

func TestUndefined_UnmarshalText(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			input    string
			expected textopt.Undefined[string]
		}{
			{`foo`, textopt.NewUndefined("foo", true)},
			{``, textopt.NewUndefined("", true)},
			{`null`, textopt.NewUndefined("null", true)},
		}

		for _, tt := range tests {
			var got textopt.Undefined[string]
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

	t.Run("int", func(t *testing.T) {
		tests := []struct {
			input    string
			expected textopt.Undefined[int]
		}{
			{``, textopt.NewUndefined(0, true)},
			{`null`, textopt.NewUndefined(0, true)},
			{`0`, textopt.NewUndefined(0, true)},
			{`42`, textopt.NewUndefined(42, true)},
		}

		for _, tt := range tests {
			var got textopt.Undefined[int]
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
