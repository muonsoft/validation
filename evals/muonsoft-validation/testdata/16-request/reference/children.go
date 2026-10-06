package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Visit) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx,
		validation.StringProperty("address", x.Address, it.IsNotBlank()),
		validation.AtProperty("tasks",
			validation.Countable(len(x.Tasks), it.HasCountBetween(1, 2)),
			validation.Slice(x.Tasks, it.HasUniqueValuesBy(func(t Task) [2]string { return [2]string{t.ServiceCode, t.Variant} })),
			validation.ValidSlice(x.Tasks)))
}
func (x Task) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx,
		validation.StringProperty("serviceCode", x.ServiceCode, it.IsNotBlank()),
		validation.ComparableProperty("variant", x.Variant, it.IsOneOf("basic", "extended").WithoutBlank()),
		validation.NumberProperty("duration", x.Duration, it.IsPositive[int]()))
}
