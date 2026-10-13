package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Form) Validate(ctx context.Context, v *validation.Validator, maxDepth int) error {
	return v.Validate(ctx, validation.ValidProperty("fields", fieldList(x.Fields, 1, maxDepth)))
}
func fieldList(fields []Field, depth, maxDepth int) validation.ValidatableFunc {
	return func(ctx context.Context, v *validation.Validator) error {
		return v.Validate(ctx,
			validation.Countable(len(fields), it.HasCountBetween(1, 3)),
			validation.Slice(fields, it.HasUniqueValuesBy(func(f Field) string { return f.Name }).At(validation.PropertyName("name"))),
			validation.Each(fields, validation.Func[Field](func(ctx context.Context, child *validation.Validator, f Field) error {
				return f.Validate(ctx, child, depth, maxDepth)
			})))
	}
}
