package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

func (x Signup) Validate(ctx context.Context, v *validation.Validator, policy EmailPolicy) error {
	return nil
}
