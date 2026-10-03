package scenario

import (
	"context"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Input) Validate(ctx context.Context, v *validation.Validator) error {
	args := []validation.Argument{validation.NilStringProperty("note", x.Note, it.IsNotBlank().WithAllowedNil()), validation.ComparableProperty("status", x.Status, it.IsOneOf("draft", "ready").WithoutBlank())}
	if x.Lower != nil && x.Upper != nil {
		args = append(args, validation.CheckProperty("upper", *x.Lower < *x.Upper).WithError(ErrRange).WithMessage(ErrRange.Message()))
	}
	return v.Validate(ctx, args...)
}
