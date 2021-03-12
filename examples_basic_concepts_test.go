package validation_test

import (
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

func ExampleValidator_Validate_basicValidation() {
	s := ""

	validator, _ := validation.NewValidator()
	err := validator.Validate(validation.String(&s, it.IsNotBlank()))

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation: This value should not be blank.
}

func ExampleValidator_Validate_singletonValidator() {
	// import "github.com/muonsoft/validation/validator"
	s := ""

	err := validator.Validate(validation.String(&s, it.IsNotBlank()))

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation: This value should not be blank.
}

func ExampleValidator_ValidateString_shorthandAlias() {
	s := ""

	err := validator.ValidateString(&s, it.IsNotBlank())

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation: This value should not be blank.
}
