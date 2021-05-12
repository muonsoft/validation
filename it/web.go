package it

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/message"
)

// IsEmail is used for simplified validation of an email address. It allows all values
// with an "@" symbol in, and a "." in the second host part of the email address.
func IsEmail() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(
		is.Email,
		"EmailConstraint",
		code.InvalidEmail,
		message.InvalidEmail,
	)
}

// IsHTML5Email is used for validation of an email address based on pattern for HTML5
// (see https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address).
func IsHTML5Email() validation.CustomStringConstraint {
	return validation.NewCustomStringConstraint(
		is.HTML5Email,
		"HTML5EmailConstraint",
		code.InvalidEmail,
		message.InvalidEmail,
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
	messageTemplate        string
}

// IsURL creates a URLConstraint to validate an URL. By default, constraint checks
// only for the http:// and https:// schemas. Use the WithSchemas method to configure
// the list of expected schemas. Also, you can use WithRelativeSchema to enable support
// of the relative schema (without schema, e.g. "//example.com").
func IsURL() URLConstraint {
	return URLConstraint{
		schemas:         []string{"http", "https"},
		messageTemplate: message.InvalidURL,
	}
}

// SetUp will return an error if the list of schemas is empty.
func (c URLConstraint) SetUp() error {
	if len(c.schemas) == 0 {
		return errEmptySchemas
	}

	return nil
}

// Name is the constraint name.
func (c URLConstraint) Name() string {
	return "URLConstraint"
}

// WithRelativeSchema enables support of relative URL schema, which means that URL value
// may be treated as relative (without schema, e.g. "//example.com").
func (c URLConstraint) WithRelativeSchema() URLConstraint {
	c.supportsRelativeSchema = true
	return c
}

// WithSchemas is used to set up a list of accepted schemas. For example, if you also consider the ftp:// type URLs
// to be valid, redefine the schemas list, listing http, https, and also ftp.
// If the list is empty, then an error will be returned by the SetUp method.
func (c URLConstraint) WithSchemas(schemas ...string) URLConstraint {
	c.schemas = schemas
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
	if c.isIgnored || value == nil || *value == "" {
		return nil
	}

	schemas := c.schemas
	if c.supportsRelativeSchema {
		schemas = append(schemas, "")
	}
	if is.URL(*value, schemas...) {
		return nil
	}

	return scope.BuildViolation(code.InvalidURL, c.messageTemplate).
		AddParameter("{{ value }}", *value).
		CreateViolation()
}
