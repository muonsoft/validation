package it

import (
	"context"
	"errors"
	"net"
	"net/url"
	"regexp"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/validate"
)

// IsEmail is used for simplified validation of an email address. It allows all values
// with an "@" symbol in, and a "." in the second host part of the email address.
func IsEmail() validation.StringFuncConstraint {
	return validation.OfStringBy(is.Email).
		WithError(validation.ErrInvalidEmail).
		WithMessage(validation.ErrInvalidEmail.Message())
}

// IsHTML5Email is used for validation of an email address based on pattern for HTML5
// (see https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address).
func IsHTML5Email() validation.StringFuncConstraint {
	return validation.OfStringBy(is.HTML5Email).
		WithError(validation.ErrInvalidEmail).
		WithMessage(validation.ErrInvalidEmail.Message())
}

// IsHostname validates that a value is a valid hostname. It checks that:
//   - each label within a valid hostname may be no more than 63 octets long;
//   - the total length of the hostname must not exceed 255 characters;
//   - hostname is fully qualified and include its top-level domain name
//     (for instance, example.com is valid but example is not);
//   - checks for reserved top-level domains according to RFC 2606
//     (hostnames containing them are not considered valid:
//     .example, .invalid, .localhost, and .test).
//
// If you do not want to check for top-level domains use IsLooseHostname version of constraint.
func IsHostname() validation.StringFuncConstraint {
	return validation.OfStringBy(is.StrictHostname).
		WithError(validation.ErrInvalidHostname).
		WithMessage(validation.ErrInvalidHostname.Message())
}

// IsLooseHostname validates that a value is a valid hostname. It checks that:
//   - each label within a valid hostname may be no more than 63 octets long;
//   - the total length of the hostname must not exceed 255 characters.
func IsLooseHostname() validation.StringFuncConstraint {
	return validation.OfStringBy(is.Hostname).
		WithError(validation.ErrInvalidHostname).
		WithMessage(validation.ErrInvalidHostname.Message())
}

// URLConstraint is used to validate URL string. This constraint doesn’t check that the host of the
// given URL really exists, because the information of the DNS records is not reliable.
//
// This constraint doesn't check the length of the URL. Use LengthConstraint to check the length of the given value.
type URLConstraint struct {
	isIgnored                   bool
	supportsRelativeSchema      bool
	schemas                     []string
	hosts                       []string
	hostPattern                 *regexp.Regexp
	restrictions                []validate.URLRestriction
	groups                      []string
	invalidErr                  error
	prohibitedErr               error
	invalidMessageTemplate      string
	invalidMessageParameters    validation.TemplateParameterList
	prohibitedMessageTemplate   string
	prohibitedMessageParameters validation.TemplateParameterList
}

// IsURL creates a URLConstraint to validate an URL. By default, constraint checks
// only for the http:// and https:// schemas. Use the WithSchemas method to configure
// the list of expected schemas. Also, you can use WithRelativeSchema to enable support
// of the relative schema (without schema, e.g. "//example.com").
func IsURL() URLConstraint {
	return URLConstraint{
		schemas:                   []string{"http", "https"},
		invalidErr:                validation.ErrInvalidURL,
		invalidMessageTemplate:    validation.ErrInvalidURL.Message(),
		prohibitedErr:             validation.ErrProhibitedURL,
		prohibitedMessageTemplate: validation.ErrProhibitedURL.Message(),
	}
}

// WithRelativeSchema enables support of relative URL schema, which means that URL value
// may be treated as relative (without schema, e.g. "//example.com").
func (c URLConstraint) WithRelativeSchema() URLConstraint {
	c.supportsRelativeSchema = true
	return c
}

// WithSchemas is used to set up a list of accepted schemas. For example, if you also consider the ftp:// type URLs
// to be valid, redefine the schemas list, listing http, https, and also ftp.
// If the list is empty, then an error will be returned.
func (c URLConstraint) WithSchemas(schemas ...string) URLConstraint {
	c.schemas = schemas
	return c
}

// WithHosts is used to set up a list of accepted hosts. If the list is empty, then any host will be treated as valid.
func (c URLConstraint) WithHosts(hosts ...string) URLConstraint {
	c.hosts = hosts
	return c
}

// WithHostMatches is used to set up restricted host validation by regexp pattern.
// If pattern is nil, then validation will be ignored.
func (c URLConstraint) WithHostMatches(pattern *regexp.Regexp) URLConstraint {
	c.hostPattern = pattern
	return c
}

// WithRestriction is used to additionally check parsed URL by callback function.
func (c URLConstraint) WithRestriction(isAllowed func(u *url.URL) bool) URLConstraint {
	c.restrictions = append(c.restrictions, func(u *url.URL) error {
		if isAllowed(u) {
			return nil
		}

		return validate.ErrProhibited
	})

	return c
}

// WithError overrides default error for produced violation.
func (c URLConstraint) WithError(err error) URLConstraint {
	c.invalidErr = err
	return c
}

// WithMessage sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ value }} - the current (invalid) value.
func (c URLConstraint) WithMessage(template string, parameters ...validation.TemplateParameter) URLConstraint {
	c.invalidMessageTemplate = template
	c.invalidMessageParameters = parameters
	return c
}

// WithProhibitedError overrides default error for produced violation on a prohibited value.
func (c URLConstraint) WithProhibitedError(err error) URLConstraint {
	c.prohibitedErr = err
	return c
}

// WithProhibitedMessage sets the violation message template for violation on a prohibited value.
// You can set custom template parameters for injecting its values into the final message.
// Also, you can use default parameters:
//
//	{{ value }} - the current (invalid) value.
func (c URLConstraint) WithProhibitedMessage(template string, parameters ...validation.TemplateParameter) URLConstraint {
	c.prohibitedMessageTemplate = template
	c.prohibitedMessageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c URLConstraint) When(condition bool) URLConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c URLConstraint) WhenGroups(groups ...string) URLConstraint {
	c.groups = groups
	return c
}

func (c URLConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if len(c.schemas) == 0 {
		return validator.CreateConstraintError("URLConstraint", "empty list of schemas")
	}
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}

	if err := validate.URL(*value, c.getRestrictions()...); err != nil {
		if errors.Is(err, validate.ErrRestrictedHost) || errors.Is(err, validate.ErrProhibited) {
			return c.newProhibitedViolation(ctx, validator, *value)
		}

		return c.newInvalidViolation(ctx, validator, *value)
	}

	return nil
}

func (c URLConstraint) getRestrictions() []validate.URLRestriction {
	schemas := c.schemas
	if c.supportsRelativeSchema {
		schemas = append(schemas, "")
	}

	restrictions := []validate.URLRestriction{validate.RestrictURLSchemas(schemas...)}
	if len(c.hosts) > 0 {
		restrictions = append(restrictions, validate.RestrictURLHosts(c.hosts...))
	}
	if c.hostPattern != nil {
		restrictions = append(restrictions, validate.RestrictURLHostByPattern(c.hostPattern))
	}
	restrictions = append(restrictions, c.restrictions...)

	return restrictions
}

func (c URLConstraint) newInvalidViolation(ctx context.Context, validator *validation.Validator, value string) error {
	return validator.BuildViolation(ctx, c.invalidErr, c.invalidMessageTemplate).
		WithParameters(
			c.invalidMessageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: value},
			)...,
		).
		Create()
}

func (c URLConstraint) newProhibitedViolation(ctx context.Context, validator *validation.Validator, value string) error {
	return validator.BuildViolation(ctx, c.prohibitedErr, c.prohibitedMessageTemplate).
		WithParameters(
			c.prohibitedMessageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: value},
			)...,
		).
		Create()
}

// IPConstraint is used to validate IP address. You can check for different versions
// and restrict some ranges by additional options.
type IPConstraint struct {
	isIgnored    bool
	validate     func(value string, restrictions ...validate.IPRestriction) error
	restrictions []validate.IPRestriction

	groups []string

	invalidErr    error
	prohibitedErr error

	invalidMessageTemplate      string
	invalidMessageParameters    validation.TemplateParameterList
	prohibitedMessageTemplate   string
	prohibitedMessageParameters validation.TemplateParameterList
}

// IsIP creates an IPConstraint to validate an IP address (IPv4 or IPv6).
func IsIP() IPConstraint {
	return newIPConstraint(validate.IP)
}

// IsIPv4 creates an IPConstraint to validate an IPv4 address.
func IsIPv4() IPConstraint {
	return newIPConstraint(validate.IPv4)
}

// IsIPv6 creates an IPConstraint to validate an IPv4 address.
func IsIPv6() IPConstraint {
	return newIPConstraint(validate.IPv6)
}

func newIPConstraint(validate func(value string, restrictions ...validate.IPRestriction) error) IPConstraint {
	return IPConstraint{
		validate:                  validate,
		invalidErr:                validation.ErrInvalidIP,
		prohibitedErr:             validation.ErrProhibitedIP,
		invalidMessageTemplate:    validation.ErrInvalidIP.Message(),
		prohibitedMessageTemplate: validation.ErrProhibitedIP.Message(),
	}
}

// DenyPrivateIP denies using of private IPs according to RFC 1918 (IPv4 addresses)
// and RFC 4193 (IPv6 addresses).
func (c IPConstraint) DenyPrivateIP() IPConstraint {
	c.restrictions = append(c.restrictions, validate.DenyPrivateIP())
	return c
}

// DenyIP can be used to deny custom range of IP addresses.
func (c IPConstraint) DenyIP(restrict func(ip net.IP) bool) IPConstraint {
	c.restrictions = append(c.restrictions, restrict)
	return c
}

// WithInvalidError overrides default underlying error for violation produced on invalid IP case.
func (c IPConstraint) WithInvalidError(err error) IPConstraint {
	c.invalidErr = err
	return c
}

// WithProhibitedError overrides default underlying error for violation produced on prohibited IP case.
func (c IPConstraint) WithProhibitedError(err error) IPConstraint {
	c.prohibitedErr = err
	return c
}

// WithInvalidMessage sets the violation message template for invalid IP case.
// You can set custom template parameters for injecting its values into the final message.
// Also, you can use default parameters:
//
//	{{ value }} - the current (invalid) value.
func (c IPConstraint) WithInvalidMessage(template string, parameters ...validation.TemplateParameter) IPConstraint {
	c.invalidMessageTemplate = template
	c.invalidMessageParameters = parameters
	return c
}

// WithProhibitedMessage sets the violation message template for prohibited IP case.
// You can set custom template parameters for injecting its values into the final message.
// Also, you can use default parameters:
//
//	{{ value }} - the current (invalid) value.
func (c IPConstraint) WithProhibitedMessage(template string, parameters ...validation.TemplateParameter) IPConstraint {
	c.prohibitedMessageTemplate = template
	c.prohibitedMessageParameters = parameters
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c IPConstraint) When(condition bool) IPConstraint {
	c.isIgnored = !condition
	return c
}

// WhenGroups enables conditional validation of the constraint by using the validation groups.
func (c IPConstraint) WhenGroups(groups ...string) IPConstraint {
	c.groups = groups
	return c
}

func (c IPConstraint) ValidateString(ctx context.Context, validator *validation.Validator, value *string) error {
	if c.isIgnored || validator.IsIgnoredForGroups(c.groups...) || value == nil || *value == "" {
		return nil
	}

	return c.validateIP(ctx, validator, *value)
}

func (c IPConstraint) validateIP(ctx context.Context, validator *validation.Validator, value string) error {
	err := c.validate(value, c.restrictions...)
	if err == nil {
		return nil
	}

	var builder *validation.ViolationBuilder
	var parameters validation.TemplateParameterList

	if errors.Is(err, validate.ErrProhibited) {
		builder = validator.BuildViolation(ctx, c.prohibitedErr, c.prohibitedMessageTemplate)
		parameters = c.prohibitedMessageParameters
	} else {
		builder = validator.BuildViolation(ctx, c.invalidErr, c.invalidMessageTemplate)
		parameters = c.invalidMessageParameters
	}

	return builder.
		WithParameters(
			parameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: value},
			)...,
		).
		Create()
}
