package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

func (x Signup) Validate(ctx context.Context, v *validation.Validator, c CodeConstraint) error {
	return nil
}
func (x Transfer) Validate(ctx context.Context, v *validation.Validator, c CodeConstraint) error {
	return nil
}
