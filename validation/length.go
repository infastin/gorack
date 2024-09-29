package validation

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

type lengthStringRule[T ~string] struct {
	min, max   int
	rune       bool
	buildError func() error
}

func LengthString[T ~string](min, max int) lengthStringRule[T] {
	return lengthStringRule[T]{
		min:  min,
		max:  max,
		rune: false,
		buildError: func() error {
			return buildLengthError(min, max)
		},
	}
}

func LengthStringRune[T ~string](min, max int) lengthStringRule[T] {
	return lengthStringRule[T]{
		min:  min,
		max:  max,
		rune: true,
		buildError: func() error {
			return buildLengthError(min, max)
		},
	}
}

func (r lengthStringRule[T]) Validate(v T) error {
	var length int
	if r.rune {
		length = utf8.RuneCountInString(string(v))
	} else {
		length = len(v)
	}

	if !isValidLength(length, r.min, r.max) {
		return r.buildError()
	}

	return nil
}

type lengthSliceRule[T any] struct {
	min, max   int
	buildError func() error
}

func LengthSlice[T any](min, max int) lengthSliceRule[T] {
	return lengthSliceRule[T]{
		min: min,
		max: max,
		buildError: func() error {
			return buildLengthError(min, max)
		},
	}
}

func (r lengthSliceRule[T]) Validate(slice []T) error {
	if !isValidLength(len(slice), r.min, r.max) {
		return r.buildError()
	}
	return nil
}

type lengthMapRule[T any] struct {
	min, max   int
	buildError func() error
}

func LengthMap[T any](min, max int) lengthMapRule[T] {
	return lengthMapRule[T]{
		min: min,
		max: max,
		buildError: func() error {
			return buildLengthError(min, max)
		},
	}
}

func (r lengthMapRule[T]) Validate(m map[string]T) error {
	if !isValidLength(len(m), r.min, r.max) {
		return r.buildError()
	}
	return nil
}

func isValidLength(l, min, max int) bool {
	return (min == 0 || l >= min) && (max == 0 || l <= max) && (min != 0 || max != 0 || l == 0)
}

func buildLengthError(min, max int) error {
	var (
		code    string
		message strings.Builder
	)

	switch {
	case min == 0 && max > 0:
		code = "length_too_long"
		message.WriteString("the length must be no more than ")
		message.WriteString(strconv.Itoa(max))
	case min > 0 && max == 0:
		code = "length_too_short"
		message.WriteString("the length must be no less than ")
		message.WriteString(strconv.Itoa(min))
	case min == 0 && max == 0:
		code = "length_empty_required"
		message.WriteString("the value must be empty")
	case min == max:
		code = "length_invalid"
		message.WriteString("the length must be exactly ")
		message.WriteString(strconv.Itoa(min))
	default:
		code = "length_out_of_range"
		message.WriteString("the length must be between ")
		message.WriteString(strconv.Itoa(min))
		message.WriteString(" and ")
		message.WriteString(strconv.Itoa(max))
	}

	return NewRuleError(code, message.String())
}
