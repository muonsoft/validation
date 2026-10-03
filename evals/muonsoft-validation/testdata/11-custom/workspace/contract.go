package scenario

import (
	"github.com/muonsoft/validation"
)

var ErrPrefix = validation.NewError("invalid_prefix", "Value must start with {{ prefix }}.")

type PrefixConstraint struct{ Prefix string }
