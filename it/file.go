package it

import (
	"context"
	"errors"
	"slices"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validate"
)

// FileNameConstraint validates file names. See [validate.FileName] for exact
// semantics. Nil and empty values are skipped; requiredness is separate.
type FileNameConstraint struct {
	options           []func(*validate.FileNameOptions)
	isIgnored         bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsFileName creates a [FileNameConstraint]. See [validate.FileName].
func IsFileName() FileNameConstraint {
	return FileNameConstraint{
		err:             validation.ErrInvalidFileName,
		messageTemplate: validation.ErrInvalidFileName.Message(),
	}
}

// WithError overrides the error attached to a validation violation.
func (c FileNameConstraint) WithError(err error) FileNameConstraint {
	c.err = err
	return c
}

// WithMessage overrides the message and its parameters. Input values are not
// automatically included in messages or parameters.
func (c FileNameConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) FileNameConstraint {
	c.messageTemplate = template
	c.messageParameters = slices.Clone(parameters)
	return c
}

// When enables or disables this constraint.
func (c FileNameConstraint) When(condition bool) FileNameConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups restricts validation to the supplied groups.
func (c FileNameConstraint) WhenGroups(groups ...string) FileNameConstraint {
	c.groups = slices.Clone(groups)
	return c
}

// WithWindowsRestrictions adds conservative Windows name restrictions on every OS.
func (c FileNameConstraint) WithWindowsRestrictions() FileNameConstraint {
	c.options = append(slices.Clip(c.options), validate.WithWindowsFileNameRestrictions())
	return c
}

// ValidateString implements [validation.StringConstraint].
func (c FileNameConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if value == nil {
		return c.Validate(ctx, validator, "")
	}
	return c.Validate(ctx, validator, *value)
}

// Validate implements [validation.Constraint] for use with This and Each.
func (c FileNameConstraint) Validate(ctx context.Context, validator *validation.Validator, value string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil
	}
	err := validate.FileName(value, c.options...)
	if err == nil {
		return nil
	}
	return validator.BuildViolation(ctx, c.err, c.messageTemplate).WithParameters(c.messageParameters...).Create()
}

// FileExtensionConstraint validates file extensions. See [validate.FileExtension] for exact
// semantics. Nil and empty values are skipped; requiredness is separate.
type FileExtensionConstraint struct {
	extensions        []string
	isIgnored         bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// HasFileExtension creates a [FileExtensionConstraint]. See [validate.FileExtension].
func HasFileExtension(extensions ...string) FileExtensionConstraint {
	return FileExtensionConstraint{
		extensions:      slices.Clone(extensions),
		err:             validation.ErrInvalidFileExtension,
		messageTemplate: validation.ErrInvalidFileExtension.Message(),
	}
}

// WithError overrides the error attached to a validation violation.
func (c FileExtensionConstraint) WithError(err error) FileExtensionConstraint {
	c.err = err
	return c
}

// WithMessage overrides the message and its parameters. Input values are not
// automatically included in messages or parameters.
func (c FileExtensionConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) FileExtensionConstraint {
	c.messageTemplate = template
	c.messageParameters = slices.Clone(parameters)
	return c
}

// When enables or disables this constraint.
func (c FileExtensionConstraint) When(condition bool) FileExtensionConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups restricts validation to the supplied groups.
func (c FileExtensionConstraint) WhenGroups(groups ...string) FileExtensionConstraint {
	c.groups = slices.Clone(groups)
	return c
}

// ValidateString implements [validation.StringConstraint].
func (c FileExtensionConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if value == nil {
		return c.Validate(ctx, validator, "")
	}
	return c.Validate(ctx, validator, *value)
}

// Validate implements [validation.Constraint] for use with This and Each.
func (c FileExtensionConstraint) Validate(ctx context.Context, validator *validation.Validator, value string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil
	}
	err := validate.FileExtension(value, c.extensions...)
	if err == nil {
		return nil
	}
	if errors.Is(err, validate.ErrInvalidFileExtensions) {
		return validator.CreateConstraintError("FileExtensionConstraint", err.Error())
	}
	return validator.BuildViolation(ctx, c.err, c.messageTemplate).WithParameters(c.messageParameters...).Create()
}

// MIMETypeConstraint validates MIME metadata. See [validate.MIMEType] for exact
// semantics. Nil and empty values are skipped; requiredness is separate.
type MIMETypeConstraint struct {
	types             []string
	isIgnored         bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// IsMIMEType creates a [MIMETypeConstraint]. See [validate.MIMEType].
func IsMIMEType(types ...string) MIMETypeConstraint {
	return MIMETypeConstraint{
		types:           slices.Clone(types),
		err:             validation.ErrInvalidMIMEType,
		messageTemplate: validation.ErrInvalidMIMEType.Message(),
	}
}

// WithError overrides the error attached to a validation violation.
func (c MIMETypeConstraint) WithError(err error) MIMETypeConstraint {
	c.err = err
	return c
}

// WithMessage overrides the message and its parameters. Input values are not
// automatically included in messages or parameters.
func (c MIMETypeConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) MIMETypeConstraint {
	c.messageTemplate = template
	c.messageParameters = slices.Clone(parameters)
	return c
}

// When enables or disables this constraint.
func (c MIMETypeConstraint) When(condition bool) MIMETypeConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups restricts validation to the supplied groups.
func (c MIMETypeConstraint) WhenGroups(groups ...string) MIMETypeConstraint {
	c.groups = slices.Clone(groups)
	return c
}

// ValidateString implements [validation.StringConstraint].
func (c MIMETypeConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if value == nil {
		return c.Validate(ctx, validator, "")
	}
	return c.Validate(ctx, validator, *value)
}

// Validate implements [validation.Constraint] for use with This and Each.
func (c MIMETypeConstraint) Validate(ctx context.Context, validator *validation.Validator, value string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil
	}
	err := validate.MIMEType(value, c.types...)
	if err == nil {
		return nil
	}
	if errors.Is(err, validate.ErrInvalidMIMETypes) {
		return validator.CreateConstraintError("MIMETypeConstraint", err.Error())
	}
	return validator.BuildViolation(ctx, c.err, c.messageTemplate).WithParameters(c.messageParameters...).Create()
}

// ContentTypeConstraint validates detected content types. See [validate.ContentType] for exact
// semantics. Nil and empty values are skipped; requiredness is separate.
type ContentTypeConstraint struct {
	types             []string
	detector          validate.ContentTypeDetector
	isIgnored         bool
	groups            []string
	err               error
	messageTemplate   string
	messageParameters validation.TemplateParameterList
}

// HasContentType creates a [ContentTypeConstraint]. See [validate.ContentType].
func HasContentType(types ...string) ContentTypeConstraint {
	return ContentTypeConstraint{
		types:           slices.Clone(types),
		err:             validation.ErrInvalidContentType,
		messageTemplate: validation.ErrInvalidContentType.Message(),
	}
}

// WithError overrides the error attached to a validation violation.
func (c ContentTypeConstraint) WithError(err error) ContentTypeConstraint {
	c.err = err
	return c
}

// WithMessage overrides the message and its parameters. Input values are not
// automatically included in messages or parameters.
func (c ContentTypeConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) ContentTypeConstraint {
	c.messageTemplate = template
	c.messageParameters = slices.Clone(parameters)
	return c
}

// When enables or disables this constraint.
func (c ContentTypeConstraint) When(condition bool) ContentTypeConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups restricts validation to the supplied groups.
func (c ContentTypeConstraint) WhenGroups(groups ...string) ContentTypeConstraint {
	c.groups = slices.Clone(groups)
	return c
}

// WithDetector replaces content detection. Nil restores http.DetectContentType.
// The detector receives all supplied bytes, must not modify them, and must be
// safe for concurrent calls when the constraint is shared.
func (c ContentTypeConstraint) WithDetector(detector validate.ContentTypeDetector) ContentTypeConstraint {
	c.detector = detector
	return c
}

// Validate implements [validation.Constraint] for use with This and Each.
func (c ContentTypeConstraint) Validate(ctx context.Context, validator *validation.Validator, value []byte) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) {
		return nil
	}
	err := validate.ContentType(value, c.types, validate.WithContentTypeDetector(c.detector))
	if err == nil {
		return nil
	}
	if errors.Is(err, validate.ErrInvalidMIMETypes) || errors.Is(err, validate.ErrInvalidDetectedContentType) {
		return validator.CreateConstraintError("ContentTypeConstraint", err.Error())
	}
	return validator.BuildViolation(ctx, c.err, c.messageTemplate).WithParameters(c.messageParameters...).Create()
}
