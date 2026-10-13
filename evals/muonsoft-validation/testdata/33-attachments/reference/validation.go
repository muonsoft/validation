package scenario

import (
	"context"
	"io"
	"path"
	"strings"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

func (x Attachment) Validate(ctx context.Context, v *validation.Validator, body io.Reader) error {
	if err := v.Validate(ctx,
		validation.StringProperty("name", x.Name, it.IsNotBlank(), it.IsFileName().WithWindowsRestrictions(), it.HasFileExtension("png", "jpg", "jpeg")),
		validation.StringProperty("mime", x.MIME, it.IsNotBlank(), it.IsMIMEType("image/png", "image/jpeg")),
		validation.NumberProperty("size", x.Size, it.IsBetween(int64(1), int64(1048576)).WithError(ErrSize).WithMessage(ErrSize.Message())),
	); err != nil {
		return err
	}
	contentType := "image/jpeg"
	if strings.ToLower(path.Ext(x.Name)) == ".png" {
		contentType = "image/png"
	}
	if err := v.Validate(ctx, validation.StringProperty("mime", x.MIME, it.IsMIMEType(contentType))); err != nil {
		return err
	}
	prefix, err := io.ReadAll(io.LimitReader(body, 512))
	if err != nil {
		return err
	}
	return v.Validate(ctx,
		validation.CountableProperty("content", len(prefix), it.IsNotBlank()),
		validation.This(prefix, it.HasContentType(contentType)).At(validation.PropertyName("content")))
}
