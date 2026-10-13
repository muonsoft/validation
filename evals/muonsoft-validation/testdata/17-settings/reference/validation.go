package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Document) Validate(ctx context.Context, v *validation.Validator, maximum int) error {
	return v.Validate(ctx,
		validation.StringProperty("title", x.Title, it.IsNotBlank()),
		validation.ValidProperty("schedule", validation.ValidatableFunc(func(ctx context.Context, child *validation.Validator) error {
			return x.Schedule.Validate(ctx, child, maximum)
		})))
}
