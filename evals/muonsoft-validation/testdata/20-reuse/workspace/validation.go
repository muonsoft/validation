package scenario

import (
	"context"

	"github.com/muonsoft/validation"
)

func (x Profile) Validate(ctx context.Context, v *validation.Validator, region validation.StringConstraint) error {
	if x.DisplayName == "" {
		return v.AtProperty("displayName").CreateViolation(ctx, validation.ErrIsBlank, validation.ErrIsBlank.Message())
	}
	return nil
}
