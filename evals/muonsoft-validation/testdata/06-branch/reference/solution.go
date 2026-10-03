package scenario

import (
	"context"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func Validate(ctx context.Context, v *validation.Validator, x Input, c validation.Constraint[string]) error {
	return v.Validate(ctx, validation.StringProperty("name", x.Name, it.IsNotBlank()), validation.Sequentially(validation.StringProperty("reference", x.Reference, it.IsNotBlank()), validation.This(x.Reference, c).At(validation.PropertyName("reference"))))
}
