package validation_test

import (
	"context"
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

type Person struct {
	Name    string
	Surname string
	Age     int
}

func (p Person) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(ctx,
		validation.StringProperty("name", p.Name, it.IsNotBlank(), it.HasMaxLength(50)),
		validation.StringProperty("surname", p.Surname, it.IsNotBlank(), it.HasMaxLength(100)),
		validation.NumberProperty[int]("age", p.Age, it.IsBetween(18, 100)),
	)
}

func ExampleValidator_ValidateIt() {
	persons := []Person{
		{
			Name:    "John",
			Surname: "Doe",
			Age:     23,
		},
		{
			Name:    "",
			Surname: "",
			Age:     0,
		},
	}

	for i, person := range persons {
		err := validator.ValidateIt(context.Background(), person)
		if violations, ok := validation.UnwrapViolationList(err); ok {
			fmt.Println("person", i, "is not valid:")
			for violation := violations.First(); violation != nil; violation = violation.Next() {
				fmt.Println(violation)
			}
		}
	}

	// Output:
	// person 1 is not valid:
	// violation at "name": "This value should not be blank."
	// violation at "surname": "This value should not be blank."
	// violation at "age": "This value should be between 18 and 100."
}
