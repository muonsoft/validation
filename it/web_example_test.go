package it_test

import (
	"fmt"

	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

func ExampleIsURL_validURL() {
	v := "http://example.com"
	err := validator.ValidateString(&v, it.IsURL())
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsURL_invalidURL() {
	v := "example.com"
	err := validator.ValidateString(&v, it.IsURL())
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid URL.
}

func ExampleURLConstraint_WithRelativeSchema() {
	v := "//example.com"
	err := validator.ValidateString(&v, it.IsURL().WithRelativeSchema())
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleURLConstraint_WithSchemas() {
	v := "ftp://example.com"
	err := validator.ValidateString(&v, it.IsURL().WithSchemas("http", "https", "ftp"))
	fmt.Println(err)
	// Output:
	// <nil>
}
