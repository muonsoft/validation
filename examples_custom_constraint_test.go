package validation_test

import (
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

func (c NumericConstraint) SetUp(scope *validation.Scope) error {
	// you may return errors here on the constraint initialization process
	return nil
}

func (c NumericConstraint) GetName() string {
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
	return scope.BuildViolation("notNumeric", "This value should be numeric.").GetViolation()
}

func ExampleValidator_Validate_customConstraint() {
	s := "alpha"

	err := validator.Validate(
		validation.String(&s, it.IsNotBlank(), IsNumeric()),
	)

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation: This value should be numeric.
}
