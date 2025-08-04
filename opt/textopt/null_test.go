package textopt_test

import (
	"testing"

	"github.com/infastin/gorack/opt/v2/textopt"
)

func TestNull_MarshalText(t *testing.T) {
	tests := []struct {
		input    textopt.Null[string]
		expected string
	}{
		{textopt.NewNull("foo", true), `foo`},
		{textopt.NewNull("", true), ``},
		{textopt.NewNull("foo", false), ``},
		{textopt.NewNull("", false), ``},
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
		expected textopt.Null[string]
	}{
		{`foo`, textopt.NewNull("foo", true)},
		{``, textopt.NewNull("", true)},
		{`null`, textopt.NewNull("null", true)},
	}

	for _, tt := range tests {
		var got textopt.Null[string]
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
