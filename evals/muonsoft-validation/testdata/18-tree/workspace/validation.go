package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

func (x Form) Validate(ctx context.Context, v *validation.Validator, maxDepth int) error { return nil }
