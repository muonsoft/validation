package validation_test

import (
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"

	"regexp"
)

func ExampleConditionalConstraint_Then() {
	v := "foo"
	err := validator.ValidateString(
		&v,
		validation.When(true).
			Then(
				it.Matches(regexp.MustCompile(`^\w+$`)),
			),
	)
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleConditionalConstraint_Else() {
	v := "123"
	err := validator.ValidateString(
		&v,
		validation.When(false).
			Then(
				it.Matches(regexp.MustCompile(`^\w+$`)),
			).
			Else(
				it.Matches(regexp.MustCompile(`^\d+$`)),
			),
	)
	fmt.Println(err)
	// Output:
	// <nil>
}
