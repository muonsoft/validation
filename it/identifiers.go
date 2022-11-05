package it

import (
	"context"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/validate"
)

// IsULID validates whether the value is a valid ULID (Universally Unique Lexicographically Sortable Identifier).
// See https://github.com/ulid/spec for ULID specifications.
func IsULID() validation.StringFuncConstraint {
	return validation.OfStringBy(is.ULID).
		WithError(validation.ErrInvalidULID).
		WithMessage(validation.ErrInvalidULID.Message())
}

// UUIDConstraint validates whether a string value is a valid UUID (also known as GUID).
//
// By default, it uses strict mode and checks the UUID as specified in RFC 4122.
// To parse additional formats, use the [UUIDConstraint.NonCanonical] method.
//
// In addition, it checks if the UUID version matches one of
// the registered versions: 1, 2, 3, 4, 5, 6 or 7.
// Use [UUIDConstraint.WithVersions] to validate for a specific set of versions.
//
// Nil UUID ("00000000-0000-0000-0000-000000000000") values are considered as valid.
// Use [UUIDConstraint.NotNil] to disallow nil value.
//
// See http://tools.ietf.org/html/rfc4122 for specifications.
type UUIDConstraint struct {
	isIgnored         bool
	groups            []string
	options           []func(o *validate.UUIDOptions)
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsUUID validates whether a string value is a valid UUID (also known as GUID).
// See [UUIDConstraint] for more info.
func IsUUID() UUIDConstraint {
	return UUIDConstraint{
		err:             validation.ErrInvalidUUID,
		messageTemplate: validation.ErrInvalidUUID.Message(),
	}
}

// NotNil used to treat nil UUID ("00000000-0000-0000-0000-000000000000") value as invalid.
func (c UUIDConstraint) NotNil() UUIDConstraint {
	c.options = append(c.options, validate.DenyNilUUID())
	return c
}

// NonCanonical used to enable parsing UUID value from non-canonical formats.
//
// Following formats are supported:
//   - "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
//   - "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
//   - "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
//   - "6ba7b8109dad11d180b400c04fd430c8"
//   - "{6ba7b8109dad11d180b400c04fd430c8}",
//   - "urn:uuid:6ba7b8109dad11d180b400c04fd430c8".
func (c UUIDConstraint) NonCanonical() UUIDConstraint {
	c.options = append(c.options, validate.AllowNonCanonicalUUIDFormats())
	return c
}

// WithVersions used to set valid versions of the UUID value.
// By default, the UUID will be checked for compliance with the default
// registered versions: 1, 2, 3, 4, 5, 6 or 7.
func (c UUIDConstraint) WithVersions(versions ...byte) UUIDConstraint {
	c.options = append(c.options, validate.AllowUUIDVersions(versions...))
	return c
}

// WithError overrides default error for produced violation.
func (c UUIDConstraint) WithError(err error) UUIDConstraint {
	c.err = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ value }} - the current (invalid) value.
func (c UUIDConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) UUIDConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c UUIDConstraint) When(condition bool) UUIDConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c UUIDConstraint) WhenGroups(groups ...string) UUIDConstraint {
	c.groups = groups
	return c
}

func (c UUIDConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}
	if is.UUID(*value, c.options...) {
		return nil
	}

	return validator.BuildViolation(ctx, c.err, c.messageTemplate).
		WithParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		Create()
}
