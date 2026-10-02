package it

import (
	"context"
	"errors"
	"strconv"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validate"
)

// PasswordStrengthConstraint checks a minimum heuristic password score.
// Nil values are skipped. Empty strings are evaluated, unlike most string constraints.
// Passwords are never automatically included in violation messages or parameters.
type PasswordStrengthConstraint struct {
	isIgnored         bool
	groups            []string
	minScore          validate.PasswordStrengthScore
	estimator         validate.PasswordStrengthEstimator
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// HasPasswordStrength requires at least Medium strength using Symfony's byte-based
// heuristic. See [validate.EstimatePasswordStrength] for its limitations.
func HasPasswordStrength() PasswordStrengthConstraint {
	return PasswordStrengthConstraint{
		minScore:        validate.PasswordStrengthMedium,
		err:             validation.ErrPasswordTooWeak,
		messageTemplate: validation.ErrPasswordTooWeak.Message(),
	}
}

// WithMinScore sets the inclusive minimum, from Weak (1) to VeryStrong (4).
func (c PasswordStrengthConstraint) WithMinScore(score validate.PasswordStrengthScore) PasswordStrengthConstraint {
	c.minScore = score
	return c
}

// WithEstimator sets a custom estimator returning 0–4. Nil restores the default.
// A custom estimator must be safe for concurrent calls if the constraint is shared.
func (c PasswordStrengthConstraint) WithEstimator(estimator validate.PasswordStrengthEstimator) PasswordStrengthConstraint {
	c.estimator = estimator
	return c
}

// WithError overrides the error for a password below the required strength.
func (c PasswordStrengthConstraint) WithError(err error) PasswordStrengthConstraint {
	c.err = err
	return c
}

// WithMessage sets a custom message with {{ strength }} and {{ minScore }} parameters.
// There is no automatic {{ value }} parameter containing the password.
func (c PasswordStrengthConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) PasswordStrengthConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When disables validation when condition is false.
func (c PasswordStrengthConstraint) When(condition bool) PasswordStrengthConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups limits validation to the specified groups.
func (c PasswordStrengthConstraint) WhenGroups(groups ...string) PasswordStrengthConstraint {
	c.groups = groups
	return c
}

// ValidateString validates an optional password without storing it in a violation.
func (c PasswordStrengthConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil
	}
	if c.minScore < validate.PasswordStrengthWeak || c.minScore > validate.PasswordStrengthVeryStrong {
		return validator.CreateConstraintError("PasswordStrengthConstraint", "minimum score must be between 1 and 4")
	}
	if value == nil {
		return nil
	}
	score, err := c.check(*value)
	if err == nil {
		return nil
	}
	if errors.Is(err, validate.ErrInvalidPasswordStrengthScore) {
		return validator.CreateConstraintError("PasswordStrengthConstraint", "estimator score must be between 0 and 4")
	}
	return validator.BuildViolation(ctx, c.err, c.messageTemplate).WithParameters(c.messageParameters.Prepend(
		validation.TemplateParameter{Key: "{{ strength }}", Value: strconv.Itoa(int(score))},
		validation.TemplateParameter{Key: "{{ minScore }}", Value: strconv.Itoa(int(c.minScore))},
	)...).Create()
}

func (c PasswordStrengthConstraint) check(password string) (validate.PasswordStrengthScore, error) {
	estimator := c.estimator
	if estimator == nil {
		estimator = validate.EstimatePasswordStrength
	}
	var score validate.PasswordStrengthScore
	err := validate.PasswordStrength(password, validate.WithMinPasswordStrength(c.minScore),
		validate.WithPasswordStrengthEstimator(func(value string) validate.PasswordStrengthScore {
			score = estimator(value)
			return score
		}))
	return score, err
}

// Validate implements [validation.Constraint][string] for [validation.This] and [validation.Each].
func (c PasswordStrengthConstraint) Validate(ctx context.Context, validator *validation.Validator, value string) error {
	return c.ValidateString(ctx, validator, &value)
}
