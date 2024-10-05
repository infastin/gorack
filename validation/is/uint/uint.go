package isuint

import (
	"github.com/infastin/gorack/constraints"
	"github.com/infastin/gorack/validation"
)

var (
	ErrPort = validation.NewRuleError("is_port", "must be a valid port number")
)

func Port[T constraints.Uint](i T) error {
	if i == 0 || uint32(i) >= 65536 {
		return ErrPort
	}
	return nil
}
