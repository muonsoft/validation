package it

import (
	"context"
	"errors"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validate"
)

// HasNoSuspiciousCharactersConstraint requires that the string has no common Unicode spoofing patterns
// (invisible/format controls, mixed-script decimal digits, risky combining sequences, optional script/locale rules).
// Empty and nil strings are ignored; combine with [IsNotBlank] for required fields.
//
// Behavior is inspired by Symfony’s [NoSuspiciousCharacters] (ICU Spoofchecker) but implemented without CGO;
// edge cases may differ from ICU.
//
// Configure which checks run with [HasNoSuspiciousCharactersConstraint.CheckInvisible],
// [HasNoSuspiciousCharactersConstraint.CheckMixedNumbers], [HasNoSuspiciousCharactersConstraint.CheckHiddenOverlay]
// and the matching Without* methods (bitmask-based internally). If none of these are called, all standard checks run
// (same as Symfony default).
//
// [NoSuspiciousCharacters]: https://symfony.com/doc/current/reference/constraints/NoSuspiciousCharacters.html
type HasNoSuspiciousCharactersConstraint struct {
	isIgnored bool
	groups    []string

	checks    uint
	checksSet bool

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

// HasNoSuspiciousCharacters creates a constraint with default checks (invisible, mixed numbers, hidden overlay)
// and no locale/script restriction. Use Check* / Without* methods to customize the check bitmask, and
// [HasNoSuspiciousCharactersConstraint.WithSuspiciousRestriction] / [HasNoSuspiciousCharactersConstraint.WithSuspiciousLocales] for script rules.
func HasNoSuspiciousCharacters() HasNoSuspiciousCharactersConstraint {
	return HasNoSuspiciousCharactersConstraint{
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

func (c HasNoSuspiciousCharactersConstraint) startExplicitChecks() HasNoSuspiciousCharactersConstraint {
	if !c.checksSet {
		c.checks = 0
		c.checksSet = true
	}
	return c
}

func (c HasNoSuspiciousCharactersConstraint) startFromDefaultChecks() HasNoSuspiciousCharactersConstraint {
	if !c.checksSet {
		c.checks = validate.DefaultSuspiciousChecks
		c.checksSet = true
	}
	return c
}

// CheckInvisible enables detection of invisible and format-control characters (bit [validate.CheckSuspiciousInvisible]).
func (c HasNoSuspiciousCharactersConstraint) CheckInvisible() HasNoSuspiciousCharactersConstraint {
	c = c.startExplicitChecks()
	c.checks |= validate.CheckSuspiciousInvisible
	return c
}

// CheckMixedNumbers enables detection of mixed decimal digit scripts (bit [validate.CheckSuspiciousMixedNumbers]).
func (c HasNoSuspiciousCharactersConstraint) CheckMixedNumbers() HasNoSuspiciousCharactersConstraint {
	c = c.startExplicitChecks()
	c.checks |= validate.CheckSuspiciousMixedNumbers
	return c
}

// CheckHiddenOverlay enables detection of risky combining overlay sequences (bit [validate.CheckSuspiciousHiddenOverlay]).
func (c HasNoSuspiciousCharactersConstraint) CheckHiddenOverlay() HasNoSuspiciousCharactersConstraint {
	c = c.startExplicitChecks()
	c.checks |= validate.CheckSuspiciousHiddenOverlay
	return c
}

// WithoutInvisible turns off invisible/format-control detection (clears [validate.CheckSuspiciousInvisible] in the mask).
func (c HasNoSuspiciousCharactersConstraint) WithoutInvisible() HasNoSuspiciousCharactersConstraint {
	c = c.startFromDefaultChecks()
	c.checks &^= validate.CheckSuspiciousInvisible
	return c
}

// WithoutMixedNumbers turns off mixed decimal digit script detection (clears [validate.CheckSuspiciousMixedNumbers] in the mask).
func (c HasNoSuspiciousCharactersConstraint) WithoutMixedNumbers() HasNoSuspiciousCharactersConstraint {
	c = c.startFromDefaultChecks()
	c.checks &^= validate.CheckSuspiciousMixedNumbers
	return c
}

// WithoutHiddenOverlay turns off hidden overlay detection (clears [validate.CheckSuspiciousHiddenOverlay] in the mask).
func (c HasNoSuspiciousCharactersConstraint) WithoutHiddenOverlay() HasNoSuspiciousCharactersConstraint {
	c = c.startFromDefaultChecks()
	c.checks &^= validate.CheckSuspiciousHiddenOverlay
	return c
}

// WithSuspiciousRestriction sets script/locale restriction mode ([validate.SuspiciousRestrictionLocales],
// [validate.SuspiciousRestrictionSingleScript], or [validate.SuspiciousRestrictionNone]).
func (c HasNoSuspiciousCharactersConstraint) WithSuspiciousRestriction(r validate.SuspiciousRestriction) HasNoSuspiciousCharactersConstraint {
	c.restriction = r
	return c
}

// WithSuspiciousLocales sets BCP 47 tags for [validate.SuspiciousRestrictionLocales] (e.g. "en", "ru-RU").
func (c HasNoSuspiciousCharactersConstraint) WithSuspiciousLocales(locales ...string) HasNoSuspiciousCharactersConstraint {
	c.locales = append([]string(nil), locales...)
	return c
}

// WithInvisibleError overrides the error for invisible/format-control detection.
func (c HasNoSuspiciousCharactersConstraint) WithInvisibleError(err error) HasNoSuspiciousCharactersConstraint {
	c.invisibleErr = err
	return c
}

// WithInvisibleMessage sets the violation template for invisible/format-control detection.
func (c HasNoSuspiciousCharactersConstraint) WithInvisibleMessage(template string, parameters ...validation.TemplateParameter) HasNoSuspiciousCharactersConstraint {
	c.invisibleTemplate = template
	c.invisibleParams = parameters
	return c
}

// WithMixedNumbersError overrides the error for mixed decimal digit scripts.
func (c HasNoSuspiciousCharactersConstraint) WithMixedNumbersError(err error) HasNoSuspiciousCharactersConstraint {
	c.mixedNumbersErr = err
	return c
}

// WithMixedNumbersMessage sets the violation template for mixed decimal digit scripts.
func (c HasNoSuspiciousCharactersConstraint) WithMixedNumbersMessage(template string, parameters ...validation.TemplateParameter) HasNoSuspiciousCharactersConstraint {
	c.mixedNumbersTemplate = template
	c.mixedNumbersParams = parameters
	return c
}

// WithHiddenOverlayError overrides the error for hidden combining overlay sequences.
func (c HasNoSuspiciousCharactersConstraint) WithHiddenOverlayError(err error) HasNoSuspiciousCharactersConstraint {
	c.hiddenOverlayErr = err
	return c
}

// WithHiddenOverlayMessage sets the violation template for hidden overlay sequences.
func (c HasNoSuspiciousCharactersConstraint) WithHiddenOverlayMessage(template string, parameters ...validation.TemplateParameter) HasNoSuspiciousCharactersConstraint {
	c.hiddenOverlayTemplate = template
	c.hiddenOverlayParams = parameters
	return c
}

// WithRestrictionError overrides the error for locale/script restriction failures.
func (c HasNoSuspiciousCharactersConstraint) WithRestrictionError(err error) HasNoSuspiciousCharactersConstraint {
	c.restrictionErr = err
	return c
}

// WithRestrictionMessage sets the violation template for locale/script restriction failures.
func (c HasNoSuspiciousCharactersConstraint) WithRestrictionMessage(template string, parameters ...validation.TemplateParameter) HasNoSuspiciousCharactersConstraint {
	c.restrictionTemplate = template
	c.restrictionParams = parameters
	return c
}

// When enables conditional validation of this constraint.
func (c HasNoSuspiciousCharactersConstraint) When(condition bool) HasNoSuspiciousCharactersConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation by validation groups.
func (c HasNoSuspiciousCharactersConstraint) WhenGroups(groups ...string) HasNoSuspiciousCharactersConstraint {
	c.groups = groups
	return c
}

// ValidateString implements [validation.StringConstraint].
func (c HasNoSuspiciousCharactersConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}
	var opts []validate.NoSuspiciousCharactersOption
	if c.checksSet {
		opts = append(opts, validate.WithSuspiciousChecks(c.checks))
	}
	opts = append(opts, validate.WithSuspiciousRestriction(c.restriction))
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

func (c HasNoSuspiciousCharactersConstraint) templateForSuspiciousValidateError(err error) (string, validation.TemplateParameterList, error) {
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
