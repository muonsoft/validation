package scenario

import (
	"context"
	"io"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Upload) Validate(ctx context.Context, v *validation.Validator, body io.Reader) error {
	if err := v.Validate(ctx,
		validation.StringProperty("name", x.Name, it.IsNotBlank(), it.IsFileName().WithWindowsRestrictions(), it.HasFileExtension("txt")),
		validation.StringProperty("mime", x.MIME, it.IsNotBlank(), it.IsMIMEType("text/plain")),
	); err != nil {
		return err
	}
	prefix, err := io.ReadAll(io.LimitReader(body, 512))
	if err != nil {
		return err
	}
	return v.Validate(ctx,
		validation.CheckProperty("content", len(prefix) > 0).WithError(validation.ErrIsBlank).WithMessage(validation.ErrIsBlank.Message()),
		validation.This(prefix, it.HasContentType("text/plain")).At(validation.PropertyName("content")),
	)
}
