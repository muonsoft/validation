package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

func (x Visit) Validate(ctx context.Context, v *validation.Validator) error { return nil }
func (x Task) Validate(ctx context.Context, v *validation.Validator) error  { return nil }
