package scenario

import (
	"context"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Input) Validate(ctx context.Context, v *validation.Validator) error {
	codes := make([]string, len(x.Entries))
	for i, entry := range x.Entries {
		codes[i] = entry.Code
	}
	return v.Validate(ctx, validation.StringProperty("title", x.Title, it.IsNotBlank()), validation.AtProperty("entries",
		validation.Countable(len(x.Entries), it.HasMinCount(1), it.HasMaxCount(3)),
		validation.Comparables(codes, it.HasUniqueValues[string]()), validation.ValidSlice(x.Entries)))
}
func (x Entry) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(ctx, validation.StringProperty("code", x.Code, it.IsNotBlank()), validation.NumberProperty("quantity", x.Quantity, it.IsPositive[int]()))
}
