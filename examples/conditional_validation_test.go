package examples

import (
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

func ExampleValidator_Validate_conditionalValidationOnConstraint() {
	notes := []struct {
		Title    string
		IsPublic bool
		Text     string
	}{
		{Title: "published note", IsPublic: true, Text: "text of published note"},
		{Title: "draft note", IsPublic: true, Text: ""},
	}

	for i, note := range notes {
		err := validator.Validate(
			validation.StringProperty("name", &note.Title, it.IsNotBlank()),
			validation.StringProperty("text", &note.Text, it.IsNotBlank().When(note.IsPublic)),
		)
		if err != nil {
			violations := err.(validation.ViolationList)
			for _, violation := range violations {
				fmt.Printf("error on note %d: %s", i, violation.Error())
			}
		}
	}

	// Output:
	// error on note 1: violation at 'text': This value should not be blank.
}
