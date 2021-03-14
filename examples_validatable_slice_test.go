package validation_test

import (
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

type Company struct {
	Name    string
	Address string
}

type Companies []Company

func (companies Companies) Validate(validator *validation.Validator) error {
	violations := validation.ViolationList{}

	for i, company := range companies {
		err := validator.AtIndex(i).Validate(
			validation.StringProperty("name", &company.Name, it.IsNotBlank()),
			validation.StringProperty("address", &company.Address, it.IsNotBlank(), it.HasMinLength(3)),
		)
		// appending violations from err
		err = violations.AppendFromError(err)
		// if append returns a non-nil error we should stop validation because an internal error occurred
		if err != nil {
			return err
		}
	}

	// we should always convert ViolationList into error by calling the AsError method
	// otherwise empty violations list will be interpreted as an error
	return violations.AsError()
}

func ExampleValidator_ValidateValidatable_validatableSlice() {
	companies := Companies{
		{"MuonSoft", "London"},
		{"", "x"},
	}

	err := validator.ValidateValidatable(companies)

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation at '[1].name': This value should not be blank.
	// violation at '[1].address': This value is too short. It should have 3 characters or more.
}