package validation_test

import (
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

func ExampleValidator_Validate_validationInterruptedAtFirstViolation() {
	title := "aaa"

	err := validator.Validate(
		validation.StringProperty(
			"title",
			&title,
			validation.Sequentially(
				it.IsBlank(),
				it.HasMinLength(5),
			),
		),
	)

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation at 'title': This value should be blank.
}
