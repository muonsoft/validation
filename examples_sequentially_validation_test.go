package validation_test

import (
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

func ExampleSequentially() {
	title := "bar"

	err := validator.ValidateString(
		&title,
		validation.Sequentially(
			it.IsBlank(),
			it.HasMinLength(5),
		),
	)

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation: This value should be blank.
}
