package it

import (
	"context"
	"errors"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validate"
)

// NoSuspiciousCharactersConstraint rejects strings that look like common Unicode spoofing
// (invisible/format controls, mixed-script decimal digits, risky combining sequences, optional script/locale rules).
// Empty and nil strings are ignored; combine with [IsNotBlank] for required fields.
//
// Behavior is inspired by Symfony’s [NoSuspiciousCharacters] (ICU Spoofchecker) but implemented without CGO;
// edge cases may differ from ICU.
//
// [NoSuspiciousCharacters]: https://symfony.com/doc/current/reference/constraints/NoSuspiciousCharacters.html
type NoSuspiciousCharactersConstraint struct {
	isIgnored bool
	groups    []string

	checks      uint
	restriction validate.SuspiciousRestriction
	locales     []string

	invisibleErr          error
	invisibleTemplate     string
	invisibleParams       validation.TemplateParameterList
	mixedNumbersErr       error
	mixedNumbersTemplate  string
	mixedNumbersParams    validation.TemplateParameterList
	hiddenOverlayErr      error
	hiddenOverlayTemplate string
	hiddenOverlayParams   validation.TemplateParameterList
	restrictionErr        error
	restrictionTemplate   string
	restrictionParams     validation.TemplateParameterList
}

// NoSuspiciousCharacters creates a constraint with default checks (invisible, mixed numbers, hidden overlay)
// and no locale/script restriction. Use [NoSuspiciousCharactersConstraint.WithChecks],
// [NoSuspiciousCharactersConstraint.WithSuspiciousRestriction], and [NoSuspiciousCharactersConstraint.WithSuspiciousLocales] to customize.
func NoSuspiciousCharacters() NoSuspiciousCharactersConstraint {
	return NoSuspiciousCharactersConstraint{
		checks:                validate.DefaultSuspiciousChecks,
		restriction:           validate.SuspiciousRestrictionNone,
		invisibleErr:          validation.ErrSuspiciousInvisible,
		invisibleTemplate:     validation.ErrSuspiciousInvisible.Message(),
		mixedNumbersErr:       validation.ErrSuspiciousMixedNumbers,
		mixedNumbersTemplate:  validation.ErrSuspiciousMixedNumbers.Message(),
		hiddenOverlayErr:      validation.ErrSuspiciousHiddenOverlay,
		hiddenOverlayTemplate: validation.ErrSuspiciousHiddenOverlay.Message(),
		restrictionErr:        validation.ErrSuspiciousCharactersRestriction,
		restrictionTemplate:   validation.ErrSuspiciousCharactersRestriction.Message(),
	}
}

// WithChecks sets the check bitmask ([validate.CheckSuspiciousInvisible], [validate.CheckSuspiciousMixedNumbers],
// [validate.CheckSuspiciousHiddenOverlay]). Use 0 to disable all checks (only restriction rules still apply if set).
func (c NoSuspiciousCharactersConstraint) WithChecks(checks uint) NoSuspiciousCharactersConstraint {
	c.checks = checks
	return c
}

// WithSuspiciousRestriction sets script/locale restriction mode ([validate.SuspiciousRestrictionLocales],
// [validate.SuspiciousRestrictionSingleScript], or [validate.SuspiciousRestrictionNone]).
func (c NoSuspiciousCharactersConstraint) WithSuspiciousRestriction(r validate.SuspiciousRestriction) NoSuspiciousCharactersConstraint {
	c.restriction = r
	return c
}

// WithSuspiciousLocales sets BCP 47 tags for [validate.SuspiciousRestrictionLocales] (e.g. "en", "ru-RU").
func (c NoSuspiciousCharactersConstraint) WithSuspiciousLocales(locales ...string) NoSuspiciousCharactersConstraint {
	c.locales = append([]string(nil), locales...)
	return c
}

// WithInvisibleError overrides the error for invisible/format-control detection.
func (c NoSuspiciousCharactersConstraint) WithInvisibleError(err error) NoSuspiciousCharactersConstraint {
	c.invisibleErr = err
	return c
}

// WithInvisibleMessage sets the violation template for invisible/format-control detection.
func (c NoSuspiciousCharactersConstraint) WithInvisibleMessage(template string, parameters ...validation.TemplateParameter) NoSuspiciousCharactersConstraint {
	c.invisibleTemplate = template
	c.invisibleParams = parameters
	return c
}

// WithMixedNumbersError overrides the error for mixed decimal digit scripts.
func (c NoSuspiciousCharactersConstraint) WithMixedNumbersError(err error) NoSuspiciousCharactersConstraint {
	c.mixedNumbersErr = err
	return c
}

// WithMixedNumbersMessage sets the violation template for mixed decimal digit scripts.
func (c NoSuspiciousCharactersConstraint) WithMixedNumbersMessage(template string, parameters ...validation.TemplateParameter) NoSuspiciousCharactersConstraint {
	c.mixedNumbersTemplate = template
	c.mixedNumbersParams = parameters
	return c
}

// WithHiddenOverlayError overrides the error for hidden combining overlay sequences.
func (c NoSuspiciousCharactersConstraint) WithHiddenOverlayError(err error) NoSuspiciousCharactersConstraint {
	c.hiddenOverlayErr = err
	return c
}

// WithHiddenOverlayMessage sets the violation template for hidden overlay sequences.
func (c NoSuspiciousCharactersConstraint) WithHiddenOverlayMessage(template string, parameters ...validation.TemplateParameter) NoSuspiciousCharactersConstraint {
	c.hiddenOverlayTemplate = template
	c.hiddenOverlayParams = parameters
	return c
}

// WithRestrictionError overrides the error for locale/script restriction failures.
func (c NoSuspiciousCharactersConstraint) WithRestrictionError(err error) NoSuspiciousCharactersConstraint {
	c.restrictionErr = err
	return c
}

// WithRestrictionMessage sets the violation template for locale/script restriction failures.
func (c NoSuspiciousCharactersConstraint) WithRestrictionMessage(template string, parameters ...validation.TemplateParameter) NoSuspiciousCharactersConstraint {
	c.restrictionTemplate = template
	c.restrictionParams = parameters
	return c
}

// When enables conditional validation of this constraint.
func (c NoSuspiciousCharactersConstraint) When(condition bool) NoSuspiciousCharactersConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation by validation groups.
func (c NoSuspiciousCharactersConstraint) WhenGroups(groups ...string) NoSuspiciousCharactersConstraint {
	c.groups = groups
	return c
}

// ValidateString implements [validation.StringConstraint].
func (c NoSuspiciousCharactersConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}
	opts := []validate.NoSuspiciousCharactersOption{
		validate.WithSuspiciousChecks(c.checks),
		validate.WithSuspiciousRestriction(c.restriction),
	}
	if len(c.locales) > 0 {
		opts = append(opts, validate.WithSuspiciousLocales(c.locales...))
	}
	err := validate.NoSuspiciousCharacters(*value, opts...)
	if err == nil {
		return nil
	}
	tmpl, params, verr := c.templateForSuspiciousValidateError(err)
	return validator.BuildViolation(ctx, verr, tmpl).
		WithParameters(
			params.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		Create()
}

func (c NoSuspiciousCharactersConstraint) templateForSuspiciousValidateError(err error) (string, validation.TemplateParameterList, error) {
	if errors.Is(err, validate.ErrSuspiciousInvisible) {
		return c.invisibleTemplate, c.invisibleParams, c.invisibleErr
	}
	if errors.Is(err, validate.ErrSuspiciousMixedNumbers) {
		return c.mixedNumbersTemplate, c.mixedNumbersParams, c.mixedNumbersErr
	}
	if errors.Is(err, validate.ErrSuspiciousHiddenOverlay) {
		return c.hiddenOverlayTemplate, c.hiddenOverlayParams, c.hiddenOverlayErr
	}
	if errors.Is(err, validate.ErrSuspiciousRestriction) {
		return c.restrictionTemplate, c.restrictionParams, c.restrictionErr
	}
	return validation.ErrNotValid.Message(), nil, validation.ErrNotValid
}
