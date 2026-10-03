package scenario

import (
	"context"
	"github.com/muonsoft/validation"
)

type Lookup interface {
	Find(context.Context, []string) (map[string]bool, error)
}

var ErrMissing = validation.NewError("missing_reference", "Reference is missing.")
