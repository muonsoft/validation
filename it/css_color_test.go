package it_test

import (
	"context"
	"testing"

	"github.com/muonsoft/language"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message/translations/russian"
	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/require"
)

func TestCSSColorCopiesAndGroups(t *testing.T) {
	v, err := validation.NewValidator()
	require.NoError(t, err)
	formats := []validate.CSSColorFormat{validate.CSSColorHexShort}
	c := it.IsCSSColor(formats...).WhenGroups("theme")
	formats[0] = validate.CSSColorRGB
	require.NoError(t, c.Validate(context.Background(), v.WithGroups("theme"), "#abc"))
	require.NoError(t, c.Validate(context.Background(), v, "bad"))
	require.ErrorIs(t, c.Validate(context.Background(), v.WithGroups("theme"), "bad"), validation.ErrInvalidCSSColor)
	require.NoError(t, v.Validate(context.Background(), validation.This("red", it.IsCSSColor())))
	require.ErrorIs(t, v.Validate(context.Background(), validation.Each([]string{"red", "bad"}, it.IsCSSColor())), validation.ErrInvalidCSSColor)
}

func TestCSSColorRussian(t *testing.T) {
	v, err := validation.NewValidator(validation.Translations(russian.Messages))
	require.NoError(t, err)
	require.EqualError(t, it.IsCSSColor().Validate(context.Background(), v.WithLanguage(language.Russian), "bad"), `violation: "Значение не является допустимым цветом CSS."`)
}
