package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Patch) Validate(ctx context.Context, v *validation.Validator) error {
	args := []validation.Argument{validation.NilStringProperty("note", x.Note, it.IsNotBlank().WithAllowedNil(), it.HasMaxLength(5))}
	if x.Mode.Present {
		if x.Mode.Value != nil {
			args = append(args, validation.NumberProperty("mode", *x.Mode.Value, it.IsOneOf(1, 2).WithoutBlank().WithError(ErrMode).WithMessage(ErrMode.Message())))
		} else if !v.IsIgnoredForGroups("activate") {
			args = append(args, validation.CheckProperty("mode", false).WithError(ErrClear).WithMessage(ErrClear.Message()))
		}
	}
	return v.WithGroups(validation.DefaultGroup).Validate(ctx, args...)
}
