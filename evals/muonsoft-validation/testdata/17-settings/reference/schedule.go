package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x *Schedule) Validate(ctx context.Context, v *validation.Validator, maximum int) error {
	if x == nil {
		return nil
	}
	ordered := x.Lower == nil || x.Upper == nil || *x.Lower < *x.Upper
	return v.Validate(ctx,
		validation.NilStringProperty("note", x.Note, it.IsNotBlank().WithAllowedNil(), it.HasMaxLength(12)),
		validation.ComparableProperty("mode", x.Mode, it.IsOneOf("quiet", "active"), it.IsOneOf("quiet", "active").WithoutBlank().WhenGroups("publish").When(x.Mode == "")),
		validation.NilNumberProperty("lower", x.Lower, it.IsPositiveOrZero[int]()),
		validation.NilNumberProperty("upper", x.Upper, it.IsPositiveOrZero[int]()),
		validation.CheckProperty("upper", ordered).WithError(ErrRange).WithMessage(ErrRange.Message()),
		validation.ValidProperty("limits", validation.ValidatableFunc(func(ctx context.Context, child *validation.Validator) error {
			return x.Limits.Validate(ctx, child, maximum)
		})))
}
