package it_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/muonsoft/language"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message/translations/english"
	"github.com/muonsoft/validation/message/translations/russian"
	"github.com/stretchr/testify/require"
)

func TestFileConstraints(t *testing.T) {
	v, err := validation.NewValidator(validation.Translations(english.Messages), validation.Translations(russian.Messages))
	require.NoError(t, err)
	ctx := context.Background()
	cases := []struct {
		name, value, ru string
		constraint      validation.Constraint[string]
		sentinel        error
	}{
		{"name", "../a", "Значение не является допустимым именем файла.", it.IsFileName(), validation.ErrInvalidFileName},
		{"extension", "a.exe", "Расширение файла недопустимо.", it.HasFileExtension("png"), validation.ErrInvalidFileExtension},
		{"mime", "text/plain", "MIME-тип некорректен или недопустим.", it.IsMIMEType("image/png"), validation.ErrInvalidMIMEType},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.ErrorIs(t, v.Validate(ctx, validation.This(tc.value, tc.constraint)), tc.sentinel)
			err := v.WithLanguage(language.Russian).Validate(ctx, validation.EachProperty("files", []string{tc.value}, tc.constraint))
			require.ErrorIs(t, err, tc.sentinel)
			require.Contains(t, err.Error(), tc.ru)
			require.Contains(t, err.Error(), "files[0]")
		})
	}
	content := it.HasContentType("image/png")
	err = v.WithLanguage(language.Russian).Validate(ctx, validation.EachProperty("files", [][]byte{[]byte("PRIVATE_CONTENT")}, content))
	require.ErrorIs(t, err, validation.ErrInvalidContentType)
	require.Contains(t, err.Error(), "Тип содержимого недопустим.")
	var violation validation.Violation
	// Inspect the single violation directly, without the aggregate list.
	err = content.Validate(ctx, v, []byte("PRIVATE_CONTENT"))
	require.ErrorAs(t, err, &violation)
	require.NotContains(t, fmt.Sprint(violation.Parameters()), "PRIVATE_CONTENT")
	require.NotContains(t, err.Error(), "PRIVATE_CONTENT")
}

func TestFileConstraintOptions(t *testing.T) {
	v, err := validation.NewValidator()
	require.NoError(t, err)
	ctx := context.Background()
	custom := errors.New("custom file error")
	param := validation.TemplateParameter{Key: "{{ example }}", Value: "custom"}
	name := it.IsFileName()
	require.NoError(t, name.ValidateString(ctx, v, nil))
	require.NoError(t, name.Validate(ctx, v, "CON"))
	require.ErrorIs(t, name.WithWindowsRestrictions().Validate(ctx, v, "CON"), validation.ErrInvalidFileName)
	require.NoError(t, name.When(false).Validate(ctx, v, ".."))
	require.NoError(t, name.WhenGroups("upload").Validate(ctx, v, ".."))
	require.ErrorIs(t, name.WhenGroups("upload").Validate(ctx, v.WithGroups("upload"), ".."), validation.ErrInvalidFileName)
	require.EqualError(t, name.WithError(custom).WithMessage("{{ example }}", param).Validate(ctx, v, ".."), `violation: "custom"`)

	ext := it.HasFileExtension("jpg")
	require.NoError(t, ext.ValidateString(ctx, v, nil))
	require.NoError(t, ext.When(false).Validate(ctx, v, "a.exe"))
	require.NoError(t, ext.WhenGroups("upload").Validate(ctx, v, "a.exe"))
	require.ErrorIs(t, ext.WhenGroups("upload").Validate(ctx, v.WithGroups("upload"), "a.exe"), validation.ErrInvalidFileExtension)
	require.ErrorIs(t, ext.WithError(custom).WithMessage("{{ example }}", param).Validate(ctx, v, "a.exe"), custom)

	mime := it.IsMIMEType("image/png")
	require.NoError(t, mime.ValidateString(ctx, v, nil))
	require.NoError(t, mime.When(false).Validate(ctx, v, "text/plain"))
	require.NoError(t, mime.WhenGroups("upload").Validate(ctx, v, "text/plain"))
	require.ErrorIs(t, mime.WhenGroups("upload").Validate(ctx, v.WithGroups("upload"), "text/plain"), validation.ErrInvalidMIMEType)
	require.ErrorIs(t, mime.WithError(custom).WithMessage("{{ example }}", param).Validate(ctx, v, "text/plain"), custom)

	calls := 0
	content := it.HasContentType("image/png").WithDetector(func([]byte) string { calls++; return "image/png" })
	require.NoError(t, content.Validate(ctx, v, nil))
	require.NoError(t, content.When(false).Validate(ctx, v, []byte("a")))
	require.NoError(t, content.WhenGroups("upload").Validate(ctx, v, []byte("a")))
	require.Zero(t, calls)
	require.NoError(t, content.WhenGroups("upload").Validate(ctx, v.WithGroups("upload"), []byte("a")))
	require.Equal(t, 1, calls)
	require.NoError(t, content.Validate(ctx, v, []byte("a")))
	require.ErrorIs(t, content.WithDetector(nil).Validate(ctx, v, []byte("a")), validation.ErrInvalidContentType)
	err = content.WithDetector(nil).WithError(custom).WithMessage("{{ example }}", param).Validate(ctx, v, []byte("a"))
	require.ErrorIs(t, err, custom)
	require.EqualError(t, err, `violation: "custom"`)
}

func TestFileConstraintConfigurationErrors(t *testing.T) {
	v, err := validation.NewValidator()
	require.NoError(t, err)
	ctx := context.Background()
	for _, constraint := range []validation.StringConstraint{it.HasFileExtension(), it.HasFileExtension("*"), it.IsMIMEType("image/*")} {
		err := v.Validate(ctx, validation.NilString(nil, constraint))
		require.Error(t, err)
		_, violations := validation.UnwrapViolations(err)
		require.False(t, violations)
	}
	for _, constraint := range []it.ContentTypeConstraint{it.HasContentType(), it.HasContentType("image/*"), it.HasContentType("image/png").WithDetector(func([]byte) string { return "" })} {
		err := constraint.Validate(ctx, v, []byte("a"))
		require.Error(t, err)
		_, violations := validation.UnwrapViolations(err)
		require.False(t, violations)
	}
	require.NoError(t, it.HasFileExtension().When(false).Validate(ctx, v, "a"))
	require.NoError(t, it.IsMIMEType("*").When(false).Validate(ctx, v, "a"))
	require.NoError(t, it.HasContentType().When(false).Validate(ctx, v, []byte("a")))
}

func TestFileConstraintOwnership(t *testing.T) {
	t.Parallel()
	v, err := validation.NewValidator()
	require.NoError(t, err)
	ctx := context.Background()
	extensions := []string{"jpg"}
	types := []string{"image/png"}
	groups := []string{"upload"}
	params := []validation.TemplateParameter{{Key: "{{ example }}", Value: "original"}}
	ext := it.HasFileExtension(extensions...)
	mime := it.IsMIMEType(types...)
	content := it.HasContentType(types...).WhenGroups(groups...).WithMessage("{{ example }}", params...)
	extensions[0], types[0], groups[0], params[0].Value = "exe", "text/plain", "other", "changed"
	require.NoError(t, ext.Validate(ctx, v, "a.jpg"))
	require.NoError(t, mime.Validate(ctx, v, "image/png"))
	require.EqualError(t, content.Validate(ctx, v.WithGroups("upload"), []byte("a")), `violation: "original"`)
	for i := range 8 {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()
			require.NoError(t, ext.Validate(ctx, v, "a.jpg"))
			require.NoError(t, mime.Validate(ctx, v, "image/png"))
			require.ErrorIs(t, content.Validate(ctx, v.WithGroups("upload"), []byte("a")), validation.ErrInvalidContentType)
		})
	}
}
