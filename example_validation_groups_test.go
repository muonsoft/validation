package validation_test

import (
	"context"
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

type User struct {
	Email    string
	Password string
	City     string
}

func (u User) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(
		ctx,
		validation.StringProperty(
			"email",
			u.Email,
			it.IsNotBlank().WhenGroups("registration"),
			it.IsEmail().WhenGroups("registration"),
		),
		validation.StringProperty(
			"password",
			u.Password,
			it.IsNotBlank().WhenGroups("registration"),
			it.HasMinLength(7).WhenGroups("registration"),
		),
		validation.StringProperty(
			"city",
			u.City,
			it.HasMinLength(2), // this constraint belongs to the default group
		),
	)
}

func ExampleValidator_WithGroups() {
	user := User{
		Email:    "invalid email",
		Password: "1234",
		City:     "Z",
	}

	err1 := validator.WithGroups("registration").Validate(context.Background(), validation.Valid(user))
	err2 := validator.Validate(context.Background(), validation.Valid(user))

	if violations, ok := validation.UnwrapViolationList(err1); ok {
		fmt.Println("violations for registration group:")
		for violation := violations.First(); violation != nil; violation = violation.Next() {
			fmt.Println(violation)
		}
	}
	if violations, ok := validation.UnwrapViolationList(err2); ok {
		fmt.Println("violations for default group:")
		for violation := violations.First(); violation != nil; violation = violation.Next() {
			fmt.Println(violation)
		}
	}

	// Output:
	// violations for registration group:
	// violation at 'email': This value is not a valid email address.
	// violation at 'password': This value is too short. It should have 7 characters or more.
	// violations for default group:
	// violation at 'city': This value is too short. It should have 2 characters or more.
}
