package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x ProductBatch) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx, validation.AtProperty("items",
		validation.Countable(len(x.Items), it.HasCountBetween(1, 3)),
		validation.Slice(x.Items, it.HasUniqueValuesBy(func(p Product) string { return p.SKU }).SkipEmptyKeys().At(validation.PropertyName("sku"))),
		validation.ValidSlice(x.Items)))
}
