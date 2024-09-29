package validation_test

import (
	"testing"

	"github.com/infastin/go-rack/validation"
)

func Test_NilSlice_Validation(t *testing.T) {
	type params struct {
		slice []string
	}
	tests := []struct {
		name   string
		params params
		want   error
	}{
		{"nil", params{nil}, nil},
		{"empty", params{[]string{}}, validation.ErrNil},
		{"not empty", params{make([]string, 10)}, validation.ErrNil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validation.NilSlice[string](true)
			if got := r.Validate(tt.params.slice); got != tt.want {
				t.Errorf("NilSlice.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_EmptySlice_Validation(t *testing.T) {
	type params struct {
		slice []string
	}
	tests := []struct {
		name   string
		params params
		want   error
	}{
		{"nil", params{nil}, nil},
		{"empty", params{[]string{}}, nil},
		{"not empty", params{make([]string, 10)}, validation.ErrEmpty},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validation.EmptySlice[string](true)
			if got := r.Validate(tt.params.slice); got != tt.want {
				t.Errorf("EmptySlice.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_NilMap_Validation(t *testing.T) {
	type params struct {
		m map[string]string
	}
	tests := []struct {
		name   string
		params params
		want   error
	}{
		{"nil", params{nil}, nil},
		{"empty", params{make(map[string]string)}, validation.ErrNil},
		{"not empty", params{map[string]string{"foo": "bar"}}, validation.ErrNil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validation.NilMap[string](true)
			if got := r.Validate(tt.params.m); got != tt.want {
				t.Errorf("NilMap.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_EmptyMap_Validation(t *testing.T) {
	type params struct {
		m map[string]string
	}
	tests := []struct {
		name   string
		params params
		want   error
	}{
		{"nil", params{nil}, nil},
		{"empty", params{make(map[string]string)}, nil},
		{"not empty", params{map[string]string{"foo": "bar"}}, validation.ErrEmpty},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validation.EmptyMap[string](true)
			if got := r.Validate(tt.params.m); got != tt.want {
				t.Errorf("EmptyMap.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
