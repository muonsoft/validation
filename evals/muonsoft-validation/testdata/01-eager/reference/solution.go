package scenario

import (
	"context"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Input) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx, validation.StringProperty("name", x.Name, it.IsNotBlank()), validation.NumberProperty("count", x.Count, it.IsPositive[int]()))
}
