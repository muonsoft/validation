package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Product) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx,
		validation.StringProperty("name", x.Name, it.IsNotBlank()),
		validation.NumberProperty("quantity", x.Quantity, it.IsPositive[int]().WithError(ErrQuantity).WithMessage(ErrQuantity.Message())),
		validation.NumberProperty("price", x.Price, it.IsGreaterThanOrEqual(float64(0)).WithError(ErrPrice).WithMessage(ErrPrice.Message())))
}
