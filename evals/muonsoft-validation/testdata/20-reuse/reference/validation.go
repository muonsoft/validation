package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Profile) Validate(ctx context.Context, v *validation.Validator, region validation.StringConstraint) error {
	return v.Validate(ctx,
		validation.StringProperty("displayName", x.DisplayName, it.IsNotBlank(), it.HasMaxLength(20)),
		validation.StringProperty("email", x.Email, it.IsNotBlank(), it.IsEmail()),
		validation.NumberProperty("age", x.Age, it.IsBetween(18, 120)),
		validation.ValidProperty("address", validation.ValidatableFunc(func(ctx context.Context, child *validation.Validator) error {
			return x.Address.Validate(ctx, child, region)
		})),
		validation.ValidProperty("preferences", x.Preferences))
}
