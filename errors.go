package validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/muonsoft/validation/message"
)

var (
	ErrInvalidEAN13      = NewError("invalid EAN-13", message.InvalidEAN13)
	ErrInvalidEAN8       = NewError("invalid EAN-8", message.InvalidEAN8)
	ErrInvalidEmail      = NewError("invalid email", message.InvalidEmail)
	ErrInvalidHostname   = NewError("invalid hostname", message.InvalidHostname)
	ErrInvalidIP         = NewError("invalid IP address", message.InvalidIP)
	ErrInvalidJSON       = NewError("invalid JSON", message.InvalidJSON)
	ErrInvalidUPCA       = NewError("invalid UPC-A", message.InvalidUPCA)
	ErrInvalidUPCE       = NewError("invalid UPC-E", message.InvalidUPCE)
	ErrInvalidURL        = NewError("invalid URL", message.InvalidURL)
	ErrIsBlank           = NewError("is blank", message.IsBlank)
	ErrIsEqual           = NewError("is equal", message.IsEqual)
	ErrIsNil             = NewError("is nil", message.IsNil)
	ErrNoSuchChoice      = NewError("no such choice", message.NoSuchChoice)
	ErrNotBlank          = NewError("is not blank", message.NotBlank)
	ErrNotEqual          = NewError("is not equal", message.NotEqual)
	ErrNotExactCount     = NewError("not exact count", message.NotExactCount)
	ErrNotExactLength    = NewError("not exact length", message.NotExactLength)
	ErrNotFalse          = NewError("is not false", message.NotFalse)
	ErrNotInRange        = NewError("is not in range", message.NotInRange)
	ErrNotInteger        = NewError("is not an integer", message.NotInteger)
	ErrNotNegative       = NewError("is not negative", message.NotNegative)
	ErrNotNegativeOrZero = NewError("is not negative or zero", message.NotNegativeOrZero)
	ErrNotNil            = NewError("is not nil", message.NotNil)
	ErrNotNumeric        = NewError("is not numeric", message.NotNumeric)
	ErrNotPositive       = NewError("is not positive", message.NotPositive)
	ErrNotPositiveOrZero = NewError("is not positive or zero", message.NotPositiveOrZero)
	ErrNotTrue           = NewError("is not true", message.NotTrue)
	ErrNotUnique         = NewError("is not unique", message.NotUnique)
	ErrNotValid          = NewError("is not valid", message.NotValid)
	ErrProhibitedIP      = NewError("is prohibited IP", message.ProhibitedIP)
	ErrTooEarly          = NewError("is too early", message.TooEarly)
	ErrTooEarlyOrEqual   = NewError("is too early or equal", message.TooEarlyOrEqual)
	ErrTooFewElements    = NewError("too few elements", message.TooFewElements)
	ErrTooHigh           = NewError("is too high", message.TooHigh)
	ErrTooHighOrEqual    = NewError("is too high or equal", message.TooHighOrEqual)
	ErrTooLate           = NewError("is too late", message.TooLate)
	ErrTooLateOrEqual    = NewError("is too late or equal", message.TooLateOrEqual)
	ErrTooLong           = NewError("is too long", message.TooLong)
	ErrTooLow            = NewError("is too low", message.TooLow)
	ErrTooLowOrEqual     = NewError("is too low or equal", message.TooLowOrEqual)
	ErrTooManyElements   = NewError("too many elements", message.TooManyElements)
	ErrTooShort          = NewError("is too short", message.TooShort)
)

// Error is a base type for static validation error used as an underlying error for Violation.
// It can be used to programmatically test for a specific violation.
// Error code values are protected by backward compatibility rules, template values are not protected.
type Error struct {
	code     string
	template string
}

// NewError creates a static validation error. It should be used to create only package-level errors.
func NewError(code string, template string) *Error {
	return &Error{code: code, template: template}
}

func (err *Error) Error() string    { return err.code }
func (err *Error) Template() string { return err.template }

// ConstraintError is used to return critical error from constraint that immediately
// stops the validation process. It is recommended to use validator.CreateConstraintError() method
// to initiate an error from current validation context.
type ConstraintError struct {
	ConstraintName string
	Path           *PropertyPath
	Description    string
}

func (err *ConstraintError) Error() string {
	var s strings.Builder
	s.WriteString("failed to validate by " + err.ConstraintName)
	if err.Path != nil {
		s.WriteString(` at path "` + err.Path.String() + `"`)
	}
	s.WriteString(": " + err.Description)

	return s.String()
}

// ConstraintAlreadyStoredError is returned when trying to put a constraint
// in the validator store using an existing key.
type ConstraintAlreadyStoredError struct {
	Key string
}

func (err *ConstraintAlreadyStoredError) Error() string {
	return fmt.Sprintf(`constraint with key "%s" already stored`, err.Key)
}

// ConstraintNotFoundError is returned when trying to get a constraint
// from the validator store using a non-existent key.
type ConstraintNotFoundError struct {
	Key  string
	Type string
}

func (err *ConstraintNotFoundError) Error() string {
	return fmt.Sprintf(`constraint by key "%s" of type "%s" is not found`, err.Key, err.Type)
}

var errTranslatorOptionsDenied = errors.New("translation options denied when using custom translator")
