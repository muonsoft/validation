package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/message"
)

// URLConstraint is used to validate URL string. This constraint doesn’t check that the host of the
// given URL really exists, because the information of the DNS records is not reliable.
//
// This constraint doesn't check the length of the URL. Use LengthConstraint to check the length of the given value.
type URLConstraint struct {
	isIgnored       bool
	isValid         func(value string, protocols ...string) bool
	protocols       []string
	messageTemplate string
}

// IsURL creates a URLConstraint to validate an absolute URL, which means a protocol (or scheme) is required.
// By default, constraint checks only for the http:// and https:// protocols. Use the Protocols method to configure
// the list of expected protocols.
//
// Example
//	v := "http://example.com"
//	err := validator.ValidateString(&v, it.IsURL())
func IsURL() URLConstraint {
	return URLConstraint{
		isValid:         is.URL,
		protocols:       []string{"http", "https"},
		messageTemplate: message.InvalidURL,
	}
}

// IsRelativeURL creates a URLConstraint to validate an absolute or relative URL. The protocol is considered
// optional when validating the syntax of the given URL. This means that both http:// and https:// are valid
// but also relative URLs that contain no protocol (e.g. //example.com).
// By default, constraint checks only for the http:// and https:// protocols. Use the Protocols method to configure
// the list of expected protocols.
//
// Example
//	v := "//example.com"
//	err := validator.ValidateString(&v, it.IsRelativeURL())
func IsRelativeURL() URLConstraint {
	return URLConstraint{
		isValid:         is.RelativeURL,
		protocols:       []string{"http", "https"},
		messageTemplate: message.InvalidURL,
	}
}

// SetUp will return an error if the list of protocols is empty.
func (c URLConstraint) SetUp() error {
	if len(c.protocols) == 0 {
		return errEmptyProtocols
	}

	return nil
}

// Name is the constraint name.
func (c URLConstraint) Name() string {
	return "URLConstraint"
}

// Protocols is used to set up a list of accepted protocols. For example, if you also consider the ftp:// type URLs
// to be valid, redefine the protocols list, listing http, https, and also ftp.
// If the list is empty, then an error will be returned by the SetUp method.
//
// Example
//	v := "ftp://example.com"
//	err := validator.ValidateString(&v, it.IsURL().Protocols("http", "https", "ftp"))
func (c URLConstraint) Protocols(protocols ...string) URLConstraint {
	c.protocols = protocols
	return c
}

// Message sets the violation message template. You can use template parameters
// for injecting its values into the final message:
//
//	{{ value }} - the current (invalid) value.
func (c URLConstraint) Message(message string) URLConstraint {
	c.messageTemplate = message
	return c
}

// When enables conditional validation of this constraint. If the expression evaluates to false,
// then the constraint will be ignored.
func (c URLConstraint) When(condition bool) URLConstraint {
	c.isIgnored = !condition
	return c
}

func (c URLConstraint) ValidateString(value *string, scope validation.Scope) error {
	if c.isIgnored || value == nil || *value == "" || c.isValid(*value, c.protocols...) {
		return nil
	}

	return scope.BuildViolation(code.InvalidURL, c.messageTemplate).
		AddParameter("{{ value }}", *value).
		CreateViolation()
}
