package validation_test

import (
	"testing"

	"github.com/infastin/go-rack/validation"
)

func Test_RequiredSlice_Validate(t *testing.T) {
	type params struct {
		slice []string
	}
	tests := []struct {
		name   string
		params params
		want   error
	}{
		{"nil", params{nil}, validation.ErrRequired},
		{"empty", params{[]string{}}, validation.ErrRequired},
		{"not empty", params{make([]string, 10)}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validation.RequiredSlice[string](true)
			if got := r.Validate(tt.params.slice); got != tt.want {
				t.Errorf("RequiredSlice.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_NilOrNotEmptySlice_Validate(t *testing.T) {
	type params struct {
		slice []string
	}
	tests := []struct {
		name   string
		params params
		want   error
	}{
		{"nil", params{nil}, nil},
		{"empty", params{[]string{}}, validation.ErrNilOrNotEmpty},
		{"not empty", params{make([]string, 10)}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validation.NilOrNotEmptySlice[string](true)
			if got := r.Validate(tt.params.slice); got != tt.want {
				t.Errorf("NilOrNotEmptySlice.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RequiredMap_Validate(t *testing.T) {
	type params struct {
		m map[string]string
	}
	tests := []struct {
		name   string
		params params
		want   error
	}{
		{"nil", params{nil}, validation.ErrRequired},
		{"empty", params{make(map[string]string)}, validation.ErrRequired},
		{"not empty", params{map[string]string{"foo": "bar"}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validation.RequiredMap[string](true)
			if got := r.Validate(tt.params.m); got != tt.want {
				t.Errorf("RequiredMap.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_NilOrNotEmptyMap_Validate(t *testing.T) {
	type params struct {
		m map[string]string
	}
	tests := []struct {
		name   string
		params params
		want   error
	}{
		{"nil", params{nil}, nil},
		{"empty", params{make(map[string]string)}, validation.ErrNilOrNotEmpty},
		{"not empty", params{map[string]string{"foo": "bar"}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validation.NilOrNotEmptyMap[string](true)
			if got := r.Validate(tt.params.m); got != tt.want {
				t.Errorf("NilOrNotEmptyMap.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
