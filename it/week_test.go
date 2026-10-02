package it_test

import (
	"context"
	"testing"

	"github.com/muonsoft/language"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message/translations/russian"
	"github.com/stretchr/testify/require"
)

func TestWeekTranslations(t *testing.T) {
	v, err := validation.NewValidator(validation.Translations(russian.Messages))
	require.NoError(t, err)
	tests := []struct {
		value      string
		constraint it.WeekConstraint
		want       string
	}{
		{"bad", it.IsWeek(), "Значение не является допустимой ISO-неделей в формате YYYY-Www."},
		{"2021-W53", it.IsWeek(), "Неделя 2021-W53 не существует в указанном ISO-году."},
		{"2020-W53", it.IsWeek().WithMin("2021-W01"), "Значение должно быть не раньше недели 2021-W01."},
		{"2021-W01", it.IsWeek().WithMax("2020-W53"), "Значение должно быть не позже недели 2020-W53."},
	}
	for _, tt := range tests {
		t.Run(tt.value+tt.want, func(t *testing.T) {
			require.EqualError(t, tt.constraint.Validate(context.Background(), v.WithLanguage(language.Russian), tt.value), `violation: "`+tt.want+`"`)
		})
	}
}

func TestWeekGroupsAndCopies(t *testing.T) {
	v, err := validation.NewValidator()
	require.NoError(t, err)
	base := it.IsWeek().WithMin("2020-W53").WhenGroups("schedule")
	derived := base.WithMin("2021-W01")
	require.NoError(t, derived.Validate(context.Background(), v, "2020-W53"))
	require.NoError(t, base.Validate(context.Background(), v.WithGroups("schedule"), "2020-W53"))
	require.ErrorIs(t, derived.Validate(context.Background(), v.WithGroups("schedule"), "2020-W53"), validation.ErrWeekTooEarly)
	require.NoError(t, v.Validate(context.Background(), validation.This("2020-W53", it.IsWeek())))
	require.ErrorIs(t, v.Validate(context.Background(), validation.Each([]string{"2020-W53", "2021-W53"}, it.IsWeek())), validation.ErrInvalidWeekNumber)
	require.NoError(t, it.IsWeek().WithMin("2021-W01").WithMin("").Validate(context.Background(), v, "2020-W53"))
	require.NoError(t, it.IsWeek().WithMax("2020-W53").WithMax("").Validate(context.Background(), v, "2021-W01"))
}

func TestWeekMessageParameters(t *testing.T) {
	v, err := validation.NewValidator()
	require.NoError(t, err)
	template := "{{ value }}: {{ min }}..{{ max }} {{ custom }}"
	parameter := validation.TemplateParameter{Key: "{{ custom }}", Value: "custom"}
	c := it.IsWeek().WithMin("2020-W01").WithMax("2020-W53").
		WithFormatMessage(template, parameter).WithWeekNumberMessage(template, parameter).
		WithMinMessage(template, parameter).WithMaxMessage(template, parameter)
	for _, value := range []string{"bad", "2021-W53", "2019-W52", "2021-W01"} {
		require.EqualError(t, c.Validate(context.Background(), v, value), `violation: "`+value+`: 2020-W01..2020-W53 custom"`)
	}
}
