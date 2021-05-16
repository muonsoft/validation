package it_test

import (
	"fmt"

	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

func ExampleIsJSON_validJSON() {
	v := `{"valid": true}`
	err := validator.ValidateString(&v, it.IsJSON())
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsJSON_invalidJSON() {
	v := `"invalid": true`
	err := validator.ValidateString(&v, it.IsJSON())
	fmt.Println(err)
	// Output:
	// violation: This value should be valid JSON.
}
