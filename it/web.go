package it

import (
	"errors"
	"net"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/message"
	"github.com/muonsoft/validation/validate"
)

// IsEmail is used for simplified validation of an email address. It allows all values
// with an "@" symbol in, and a "." in the second host part of the email address.
func IsEmail() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(
		is.Email,
		code.InvalidEmail,
		message.Templates[code.InvalidEmail],
	)
}

// IsHTML5Email is used for validation of an email address based on pattern for HTML5
// (see https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address).
func IsHTML5Email() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(
		is.HTML5Email,
		code.InvalidEmail,
		message.Templates[code.InvalidEmail],
	)
}

// IsHostname validates that a value is a valid hostname. It checks that:
//	• each label within a valid hostname may be no more than 63 octets long;
//	• the total length of the hostname must not exceed 255 characters;
//	• hostname is fully qualified and include its top-level domain name
//	  (for instance, example.com is valid but example is not);
//	• checks for reserved top-level domains according to RFC 2606
//	  (hostnames containing them are not considered valid:
//	  .example, .invalid, .localhost, and .test).
//
// If you do not want to check for top-level domains use IsLooseHostname version of constraint.
func IsHostname() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(
		is.StrictHostname,
		code.InvalidHostname,
		message.Templates[code.InvalidHostname],
	)
}

// IsLooseHostname validates that a value is a valid hostname. It checks that:
//	• each label within a valid hostname may be no more than 63 octets long;
//	• the total length of the hostname must not exceed 255 characters.
func IsLooseHostname() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(
		is.Hostname,
		code.InvalidHostname,
		message.Templates[code.InvalidHostname],
	)
}

// URLConstraint is used to validate URL string. This constraint doesn’t check that the host of the
// given URL really exists, because the information of the DNS records is not reliable.
//
// This constraint doesn't check the length of the URL. Use LengthConstraint to check the length of the given value.
type URLConstraint struct {
	isIgnored              bool
	supportsRelativeSchema bool
	schemas                []string
	groups                 []string
	code                   string
	messageTemplate        string
	messageParameters      validation.TemplateParameterList
}

// IsURL creates a URLConstraint to validate an URL. By default, constraint checks
// only for the http:// and https:// schemas. Use the WithSchemas method to configure
// the list of expected schemas. Also, you can use WithRelativeSchema to enable support
// of the relative schema (without schema, e.g. "//example.com").
func IsURL() URLConstraint {
	return URLConstraint{
		schemas:         []string{"http", "https"},
		code:            code.InvalidURL,
		messageTemplate: message.Templates[code.InvalidURL],
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

// Code overrides default code for produced violation.
func (c URLConstraint) Code(code string) URLConstraint {
	c.code = code
	return c
}

// Message sets the violation message template. You can set custom template parameters
// for injecting its values into the final message. Also, you can use default parameters:
//
//	{{ value }} - the current (invalid) value.
func (c URLConstraint) Message(template string, parameters ...validation.TemplateParameter) URLConstraint {
	c.messageTemplate = template
	c.messageParameters = parameters
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

func (c URLConstraint) ValidateString(value *string, scope validation.Scope) error {
	if len(c.schemas) == 0 {
		return scope.NewConstraintError("URLConstraint", "empty list of schemas")
	}
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || *value == "" {
		return nil
	}

	schemas := c.schemas
	if c.supportsRelativeSchema {
		schemas = append(schemas, "")
	}
	if is.URL(*value, schemas...) {
		return nil
	}

	return scope.BuildViolation(c.code, c.messageTemplate).
		SetParameters(
			c.messageParameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: *value},
			)...,
		).
		CreateViolation()
}

// IPConstraint is used to validate IP address. You can check for different versions
// and restrict some ranges by additional options.
type IPConstraint struct {
	isIgnored    bool
	validate     func(value string, restrictions ...validate.IPRestriction) error
	restrictions []validate.IPRestriction

	groups []string

	invalidCode    string
	prohibitedCode string

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
		invalidCode:               code.InvalidIP,
		prohibitedCode:            code.ProhibitedIP,
		invalidMessageTemplate:    message.Templates[code.InvalidIP],
		prohibitedMessageTemplate: message.Templates[code.ProhibitedIP],
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

// InvalidCode overrides default code for violation produced on invalid IP case.
func (c IPConstraint) InvalidCode(code string) IPConstraint {
	c.invalidCode = code
	return c
}

// ProhibitedCode overrides default code for violation produced on prohibited IP case.
func (c IPConstraint) ProhibitedCode(code string) IPConstraint {
	c.prohibitedCode = code
	return c
}

// InvalidMessage sets the violation message template for invalid IP case.
// You can set custom template parameters for injecting its values into the final message.
// Also, you can use default parameters:
//
//	{{ value }} - the current (invalid) value.
func (c IPConstraint) InvalidMessage(template string, parameters ...validation.TemplateParameter) IPConstraint {
	c.invalidMessageTemplate = template
	c.invalidMessageParameters = parameters
	return c
}

// ProhibitedMessage sets the violation message template for prohibited IP case.
// You can set custom template parameters for injecting its values into the final message.
// Also, you can use default parameters:
//
//	{{ value }} - the current (invalid) value.
func (c IPConstraint) ProhibitedMessage(template string, parameters ...validation.TemplateParameter) IPConstraint {
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

func (c IPConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || scope.IsIgnored(c.groups...) || value == nil || *value == "" {
		return nil
	}

	return c.validateIP(*value, scope)
}

func (c IPConstraint) validateIP(value string, scope validation.Scope) error {
	err := c.validate(value, c.restrictions...)
	if err == nil {
		return nil
	}

	var builder *validation.ViolationBuilder
	var parameters validation.TemplateParameterList

	if errors.Is(err, validate.ErrProhibited) {
		builder = scope.BuildViolation(c.prohibitedCode, c.prohibitedMessageTemplate)
		parameters = c.prohibitedMessageParameters
	} else {
		builder = scope.BuildViolation(c.invalidCode, c.invalidMessageTemplate)
		parameters = c.invalidMessageParameters
	}

	return builder.
		SetParameters(
			parameters.Prepend(
				validation.TemplateParameter{Key: "{{ value }}", Value: value},
			)...,
		).
		CreateViolation()
}
