package validation_test

import (
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

func ExampleValidator_Validate_customizingErrorMessage() {
	s := ""

	err := validator.ValidateString(&s, it.IsNotBlank().Message("this value is required"))

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation: this value is required
}

func ExampleValidator_Validate_translationForCustomMessage() {
	const customMessage = "tags should contain more than {{ limit }} element(s)"
	validator, _ := validation.NewValidator(
		validation.Translations(map[language.Tag]map[string]catalog.Message{
			language.Russian: {
				customMessage: plural.Selectf(1, "",
					plural.One, "теги должны содержать {{ limit }} элемент и более",
					plural.Few, "теги должны содержать более {{ limit }} элемента",
					plural.Other, "теги должны содержать более {{ limit }} элементов"),
			},
		}),
	)

	var tags []string
	err := validator.Validate(
		validation.Language(language.Russian),
		validation.Iterable(tags, it.HasMinCount(1).MinMessage(customMessage)),
	)

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation: теги должны содержать 1 элемент и более
}
