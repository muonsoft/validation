package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

func (x Field) Validate(ctx context.Context, v *validation.Validator, depth, maxDepth int) error {
	return nil
}
