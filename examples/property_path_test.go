package examples

import (
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

func ExampleValidator_Validate_passingPropertyPathViaOptions() {
	s := ""

	err := validator.Validate(
		validation.String(
			&s,
			validation.PropertyName("properties"),
			validation.ArrayIndex(1),
			validation.PropertyName("tag"),
			it.IsNotBlank(),
		),
	)

	violation := err.(validation.ViolationList)[0]
	fmt.Println("property path:", violation.GetPropertyPath().Format())
	// Output:
	// property path: properties[1].tag
}

func ExampleValidator_Validate_propertyPathWithScopedValidator() {
	s := ""

	err := validator.
		AtProperty("properties").
		AtIndex(1).
		AtProperty("tag").
		Validate(validation.String(&s, it.IsNotBlank()))

	violation := err.(validation.ViolationList)[0]
	fmt.Println("property path:", violation.GetPropertyPath().Format())
	// Output:
	// property path: properties[1].tag
}

func ExampleValidator_Validate_propertyPathBySpecialArgument() {
	s := ""

	err := validator.Validate(
		// this is an alias for
		// validation.String(&s, validation.PropertyName("property"), it.IsNotBlank()),
		validation.StringProperty("property", &s, it.IsNotBlank()),
	)

	violation := err.(validation.ViolationList)[0]
	fmt.Println("property path:", violation.GetPropertyPath().Format())
	// Output:
	// property path: property
}
