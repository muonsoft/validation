package validation_test

import (
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

type File struct {
	IsDocument   bool
	DocumentName string
	Name         string
}

func ExampleWhen_structValidationWithConditionalConstraint() {
	file := File{
		IsDocument: true,
		Name:       "file name",
	}

	err := validator.Validate(
		validation.StringProperty(
			"name",
			&file.Name,
			it.IsNotBlank(),
		),
		validation.StringProperty(
			"documentName",
			&file.DocumentName,
			validation.When(file.IsDocument).
				Then(it.IsNotBlank()),
		),
	)

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation at 'documentName': This value should not be blank.
}
