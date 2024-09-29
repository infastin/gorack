package isint

import (
	"github.com/infastin/go-rack/validation"
	"github.com/infastin/go-rack/validation/constraints"
)

var (
	ErrPort = validation.NewRuleError("is_port", "must be a valid port number")
)

func Port[T constraints.Int](i T) error {
	if i <= 0 || int32(i) >= 65536 {
		return ErrPort
	}
	return nil
}
