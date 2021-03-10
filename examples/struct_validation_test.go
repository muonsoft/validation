package examples

import (
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

type Document struct {
	Title    string
	Keywords []string
}

func ExampleValidator_Validate_basicStructValidation() {
	document := Document{
		Title:    "",
		Keywords: []string{""},
	}

	err := validator.Validate(
		validation.StringProperty("title", &document.Title, it.IsNotBlank()),
		validation.CountableProperty("keywords", len(document.Keywords), it.HasCountBetween(2, 10)),
		validation.EachStringProperty("keywords", document.Keywords, it.IsNotBlank()),
	)

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation at 'title': This value should not be blank.
	// violation at 'keywords': This collection should contain 2 elements or more.
	// violation at 'keywords[0]': This value should not be blank.
}