package validation_test

import (
	"context"
	"fmt"
	"regexp"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

type NumericConstraint struct {
	matcher *regexp.Regexp
}

// it is recommended to use semantic constructors for constraints.
func IsNumeric() NumericConstraint {
	return NumericConstraint{matcher: regexp.MustCompile("^[0-9]+$")}
}

func (c NumericConstraint) SetUp() error {
	// you may return errors here on the constraint initialization process
	return nil
}

func (c NumericConstraint) Name() string {
	return "NumericConstraint"
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
	return scope.CreateViolation("notNumeric", "This value should be numeric.")
}

func ExampleValidator_Validate_customConstraint() {
	s := "alpha"

	err := validator.Validate(
		context.Background(),
		validation.String(s, it.IsNotBlank(), IsNumeric()),
	)

	fmt.Println(err)
	// Output:
	// violation: This value should be numeric.
}
