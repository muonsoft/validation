package scenario

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Article) Validate(ctx context.Context, v *validation.Validator) error {
	groups := []string{validation.DefaultGroup, "publish"}
	return v.Validate(ctx,
		validation.StringProperty("title", x.Title,
			it.IsNotBlank().WithError(ErrTitleRequired).WithMessage(ErrTitleRequired.Message()).WhenGroups(groups...),
			it.HasLengthBetween(3, 12).WithMinError(ErrTitleLength).WithMaxError(ErrTitleLength).WithMinMessage(ErrTitleLength.Message()).WithMaxMessage(ErrTitleLength.Message()).WhenGroups(groups...)),
		validation.StringProperty("contact", x.Contact,
			it.IsNotBlank().WithError(ErrContactRequired).WithMessage(ErrContactRequired.Message()).WhenGroups("publish"),
			it.IsEmail().WithError(ErrContactFormat).WithMessage(ErrContactFormat.Message()).WhenGroups(groups...)),
	)
}
