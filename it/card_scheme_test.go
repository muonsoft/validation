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

func TestCardSchemeCopiesAndGroups(t *testing.T) {
	v, err := validation.NewValidator()
	require.NoError(t, err)
	schemes := []validate.CardSchemeName{validate.CardSchemeVisa}
	c := it.IsCardScheme(schemes...).WhenGroups("payment")
	schemes[0] = validate.CardSchemeMastercard
	require.NoError(t, c.Validate(context.Background(), v.WithGroups("payment"), "4111111111111111"))
	require.NoError(t, c.Validate(context.Background(), v, "bad"))
	require.ErrorIs(t, c.Validate(context.Background(), v.WithGroups("payment"), "bad"), validation.ErrInvalidCardScheme)
	require.NoError(t, v.Validate(context.Background(), validation.This("4111111111111111", it.IsCardScheme(validate.CardSchemeVisa))))
	require.ErrorIs(t, v.Validate(context.Background(), validation.Each([]string{"4111111111111111", "bad"}, it.IsCardScheme(validate.CardSchemeVisa))), validation.ErrInvalidCardScheme)
	require.Error(t, v.Validate(context.Background(), validation.String("4111111111111112", it.IsCardScheme(validate.CardSchemeVisa), it.IsLUHN())))
}

func TestCardSchemeRussian(t *testing.T) {
	v, err := validation.NewValidator(validation.Translations(russian.Messages))
	require.NoError(t, err)
	err = it.IsCardScheme(validate.CardSchemeVisa).Validate(context.Background(), v.WithLanguage(language.Russian), "bad")
	require.EqualError(t, err, `violation: "Неподдерживаемый тип карты или неверный номер карты."`)
}
