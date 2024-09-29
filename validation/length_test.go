package validation_test

import (
	"testing"

	"github.com/infastin/go-rack/validation"
)

func Test_LengthString_Validate(t *testing.T) {
	type params struct {
		min, max int
		str      string
	}
	tests := []struct {
		name    string
		params  params
		want    string
		wantErr bool
	}{
		{"between", params{2, 4, "abc"}, "", false},
		{"between empty", params{2, 4, ""}, "the length must be between 2 and 4", true},
		{"between out of range", params{2, 4, "abcdef"}, "the length must be between 2 and 4", true},
		{"below max", params{0, 4, "foo"}, "", false},
		{"above max", params{0, 4, "hello"}, "the length must be no more than 4", true},
		{"below min", params{4, 0, "foo"}, "the length must be no less than 4", true},
		{"above min", params{4, 0, "hello"}, "", false},
		{"exact", params{3, 3, "foo"}, "", false},
		{"not exact", params{3, 3, "hello"}, "the length must be exactly 3", true},
		{"empty", params{0, 0, ""}, "", false},
		{"not empty", params{0, 0, "abcd"}, "the value must be empty", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validation.LengthString[string](tt.params.min, tt.params.max)
			got := r.Validate(tt.params.str)
			if (got != nil) != tt.wantErr {
				t.Errorf("LengthString.Validate() error = %v, wantErr %v", got, tt.wantErr)
				return
			}
			if got != nil && got.Error() != tt.want {
				t.Errorf("LengthString.Validate() = %v, want %v", got.Error(), tt.want)
			}
		})
	}
}
