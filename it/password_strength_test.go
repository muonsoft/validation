package it_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/muonsoft/language"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message/translations/russian"
	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/require"
)

func TestPasswordStrengthConstraint(t *testing.T) {
	v, err := validation.NewValidator(validation.Translations(russian.Messages))
	require.NoError(t, err)
	ctx := context.Background()
	calls := 0
	base := it.HasPasswordStrength().WithEstimator(func(string) validate.PasswordStrengthScore { calls++; return 1 })
	require.NoError(t, base.ValidateString(ctx, v, nil))
	require.Equal(t, 0, calls)
	custom := errors.New("custom")
	err = base.WithError(custom).WithMessage("{{ strength }}/{{ minScore }} {{ custom }}", validation.TemplateParameter{Key: "{{ custom }}", Value: "test"}).Validate(ctx, v, "PRIVATE_PASSWORD")
	require.ErrorIs(t, err, custom)
	require.EqualError(t, err, `violation: "1/2 test"`)
	require.Equal(t, 1, calls)
	var violation validation.Violation
	require.ErrorAs(t, err, &violation)
	require.NotContains(t, fmt.Sprint(violation.Parameters()), "PRIVATE_PASSWORD")
	require.NotContains(t, fmt.Sprint(violation.Parameters()), "{{ value }}")
	require.ErrorIs(t, base.Validate(ctx, v, ""), validation.ErrPasswordTooWeak)
	require.NoError(t, base.WithMinScore(validate.PasswordStrengthWeak).Validate(ctx, v, "test"))
	require.ErrorIs(t, base.Validate(ctx, v, "test"), validation.ErrPasswordTooWeak)
	require.NoError(t, base.When(false).Validate(ctx, v, "test"))
	require.NoError(t, base.WhenGroups("password").Validate(ctx, v, "test"))
	require.ErrorIs(t, base.WhenGroups("password").Validate(ctx, v.WithGroups("password"), "test"), validation.ErrPasswordTooWeak)
	require.EqualError(t, it.HasPasswordStrength().Validate(ctx, v.WithLanguage(language.Russian), "short"), `violation: "Пароль слишком слабый. Используйте более надёжный пароль."`)
	require.NoError(t, base.WithEstimator(nil).Validate(ctx, v, "Reasonable-pwd"))
	require.NoError(t, v.Validate(ctx, validation.This("Reasonable-pwd", it.HasPasswordStrength())))
	require.ErrorIs(t, v.Validate(ctx, validation.Each([]string{"Reasonable-pwd", ""}, it.HasPasswordStrength())), validation.ErrPasswordTooWeak)
}

func TestPasswordStrengthConfigurationErrors(t *testing.T) {
	v, err := validation.NewValidator()
	require.NoError(t, err)
	for _, score := range []validate.PasswordStrengthScore{-1, 0, 5} {
		require.EqualError(t, it.HasPasswordStrength().WithMinScore(score).ValidateString(context.Background(), v, nil), "validate by PasswordStrengthConstraint: minimum score must be between 1 and 4")
	}
	for _, score := range []validate.PasswordStrengthScore{-1, 5} {
		require.EqualError(t, it.HasPasswordStrength().WithEstimator(func(string) validate.PasswordStrengthScore { return score }).Validate(context.Background(), v, "secret"), "validate by PasswordStrengthConstraint: estimator score must be between 0 and 4")
	}
}
