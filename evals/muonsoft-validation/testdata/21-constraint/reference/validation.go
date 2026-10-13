package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Signup) Validate(ctx context.Context, v *validation.Validator, c CodeConstraint) error {
	return v.Validate(ctx, validation.NilStringProperty("code", x.Code, it.IsNotBlank(), c), validation.StringProperty("label", x.Label, it.IsNotBlank()))
}
func (x Transfer) Validate(ctx context.Context, v *validation.Validator, c CodeConstraint) error {
	return v.Validate(ctx, validation.StringProperty("destinationCode", x.Destination, it.IsNotBlank(), c), validation.StringProperty("label", x.Label, it.IsNotBlank()))
}
