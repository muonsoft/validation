package examples

import (
	"context"
	"fmt"

	languagepkg "github.com/muonsoft/language"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message/translations/russian"
	"golang.org/x/text/language"
)

func ExampleValidator_Validate_translationsByDefaultLanguage() {
	validator, _ := validation.NewValidator(
		validation.Translations(russian.Messages),
		validation.DefaultLanguage(language.Russian),
	)

	s := ""
	err := validator.ValidateString(&s, it.IsNotBlank())

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation: Значение не должно быть пустым.
}

func ExampleValidator_Validate_translationsByOption() {
	validator, _ := validation.NewValidator(
		validation.Translations(russian.Messages),
	)

	s := ""
	err := validator.ValidateString(&s, validation.Language(language.Russian), it.IsNotBlank())

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation: Значение не должно быть пустым.
}

func ExampleValidator_Validate_translationsByContextOption() {
	validator, _ := validation.NewValidator(
		validation.Translations(russian.Messages),
	)

	s := ""
	ctx := languagepkg.WithContext(context.Background(), language.Russian)
	err := validator.ValidateString(&s, validation.Context(ctx), it.IsNotBlank())

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation: Значение не должно быть пустым.
}

func ExampleValidator_Validate_translationsByContextValidator() {
	validator, _ := validation.NewValidator(
		validation.Translations(russian.Messages),
	)
	ctx := languagepkg.WithContext(context.Background(), language.Russian)
	validator = validator.WithContext(ctx)

	s := ""
	err := validator.ValidateString(&s, it.IsNotBlank())

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation: Значение не должно быть пустым.
}
