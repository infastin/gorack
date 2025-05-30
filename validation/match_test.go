package validation_test

import (
	"testing"

	"github.com/infastin/gorack/validation"
)

func Test_Match_Validate(t *testing.T) {
	tests := []struct {
		value string
		expr  string
		valid bool
	}{
		{"123", "^[0-9]+$", true},
		{"123", "^[a-z]+$", false},
		{"hello", "^[a-z]+$", true},
		{"hello", "^[0-9]+$", false},
		{"hello123", "^[a-z0-9]+$", true},
	}
	for _, tt := range tests {
		r := validation.Match[string](tt.expr)
		err := r.Validate(tt.value)
		if err != nil && tt.valid {
			t.Errorf("unexpected error: value=%s expr=%s error=%s", tt.value, tt.expr, err.Error())
		} else if err == nil && !tt.valid {
			t.Errorf("expected an error: value=%s expr=%s", tt.value, tt.expr)
		}
	}
}

func Test_NotMatch_Validate(t *testing.T) {
	tests := []struct {
		value string
		expr  string
		valid bool
	}{
		{"123", "^[0-9]+$", false},
		{"123", "^[a-z]+$", true},
		{"hello", "^[a-z]+$", false},
		{"hello", "^[0-9]+$", true},
		{"hello123", "^[a-z0-9]+$", false},
	}
	for _, tt := range tests {
		r := validation.NotMatch[string](tt.expr)
		err := r.Validate(tt.value)
		if err != nil && tt.valid {
			t.Errorf("unexpected error: value=%s expr=%s error=%s", tt.value, tt.expr, err.Error())
		} else if err == nil && !tt.valid {
			t.Errorf("expected an error: value=%s expr=%s", tt.value, tt.expr)
		}
	}
}
