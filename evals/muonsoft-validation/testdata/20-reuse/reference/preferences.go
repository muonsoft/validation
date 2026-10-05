package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x *Preferences) Validate(ctx context.Context, v *validation.Validator) error {
	if x == nil {
		return nil
	}
	return v.Validate(ctx,
		validation.ComparableProperty("theme", x.Theme, it.IsOneOf("light", "dark").WithoutBlank()),
		validation.NilNumberProperty("quota", x.Quota, it.IsBetween(0, 100)),
		validation.AtProperty("tags", validation.Countable(len(x.Tags), it.HasMaxCount(3)), validation.Comparables(x.Tags, it.HasUniqueValues[string]()), validation.EachString(x.Tags, it.IsNotBlank(), it.HasMaxLength(8))))
}
