package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validate"
)

func (x PasswordUpdate) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx, validation.NilStringProperty("password", x.Password,
		it.HasPasswordStrength().WithMinScore(validate.PasswordStrengthStrong).WithEstimator(x.Estimate)))
}
