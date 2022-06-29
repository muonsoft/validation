package validation_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

var ErrNotNumeric = errors.New("not numeric")

type NumericConstraint struct {
	matcher *regexp.Regexp
}

// it is recommended to use semantic constructors for constraints.
func IsNumeric() NumericConstraint {
	return NumericConstraint{matcher: regexp.MustCompile("^[0-9]+$")}
}

func (c NumericConstraint) ValidateString(value *string, scope validation.Scope) error {
	// usually, you should ignore empty values
	// to check for an empty value you should use it.NotBlankConstraint
	if value == nil || *value == "" {
		return nil
	}

	if c.matcher.MatchString(*value) {
		return nil
	}

	// use the scope to build violation with translations
	return scope.CreateViolation(ErrNotNumeric, "This value should be numeric.")
}

func ExampleValidator_Validate_customConstraint() {
	s := "alpha"

	err := validator.Validate(
		context.Background(),
		validation.String(s, it.IsNotBlank(), IsNumeric()),
	)

	fmt.Println(err)
	fmt.Println("errors.Is(err, ErrNotNumeric) =", errors.Is(err, ErrNotNumeric))
	// Output:
	// violation: This value should be numeric.
	// errors.Is(err, ErrNotNumeric) = true
}
