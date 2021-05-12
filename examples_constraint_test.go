package validation_test

import (
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validator"
)

func ExampleNewCustomStringConstraint() {
	validate := func(s string) bool {
		return s == "valid"
	}
	constraint := validation.NewCustomStringConstraint(
		validate,
		"ExampleConstraint", // constraint name
		"exampleCode",       // violation code
		"Unexpected value.", // violation message template
	)

	s := "foo"
	err := validator.ValidateString(&s, constraint)

	fmt.Println(err)
	// Output:
	// violation: Unexpected value.
}
