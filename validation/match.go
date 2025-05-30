package validation

import "regexp"

var (
	ErrNoMatch = NewRuleError("no_match", "must be a valid value")
	ErrMatch   = NewRuleError("match", "must be a valid value")
)

type matchRule[T ~string] struct {
	re          *regexp.Regexp
	failOnMatch bool
}

func Match[T ~string](expr string) matchRule[T] {
	return matchRule[T]{
		re:          regexp.MustCompile(expr),
		failOnMatch: false,
	}
}

func NotMatch[T ~string](expr string) matchRule[T] {
	return matchRule[T]{
		re:          regexp.MustCompile(expr),
		failOnMatch: true,
	}
}

func (r matchRule[T]) Validate(str T) error {
	match := r.re.MatchString(string(str))
	if match == r.failOnMatch {
		if match {
			return ErrMatch
		}
		return ErrNoMatch
	}
	return nil
}
