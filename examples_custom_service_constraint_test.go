package validation_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

type contextKey string

const exampleKey contextKey = "exampleKey"

type TagStorage struct {
	// this might be stored in the database
	tags []string
}

func (storage *TagStorage) FindByName(ctx context.Context, name string) ([]string, error) {
	contextValue, ok := ctx.Value(exampleKey).(string)
	if !ok {
		return nil, errors.New("context value missing")
	}
	if contextValue != "value" {
		return nil, errors.New("invalid context value")
	}

	found := make([]string, 0)

	for _, tag := range storage.tags {
		if tag == name {
			found = append(found, tag)
		}
	}

	return found, nil
}

type ExistingTagConstraint struct {
	storage *TagStorage
}

func (c *ExistingTagConstraint) SetUp() error {
	return nil
}

func (c *ExistingTagConstraint) Name() string {
	return "ExistingTagConstraint"
}

func (c *ExistingTagConstraint) ValidateString(value *string, scope validation.Scope) error {
	// usually, you should ignore empty values
	// to check for an empty value you should use it.NotBlankConstraint
	if value == nil || *value == "" {
		return nil
	}

	// you can pass the context value from the scope
	entities, err := c.storage.FindByName(scope.Context(), *value)
	// here you can return a service error so that the validation process
	// is stopped immediately
	if err != nil {
		return err
	}
	if len(entities) > 0 {
		return nil
	}

	// use the scope to build violation with translations
	return scope.
		BuildViolation("unknownTag", `Tag "{{ value }}" does not exist.`).
		// you can inject parameter value to the message here
		AddParameter("{{ value }}", *value).
		CreateViolation()
}

func ExampleValidator_Validate_customServiceConstraint() {
	storage := &TagStorage{tags: []string{"camera", "book"}}
	isTagExists := &ExistingTagConstraint{storage: storage}

	tag := "movie"
	ctx := context.WithValue(context.Background(), exampleKey, "value")

	err := validator.Validate(
		// you can pass here the context value to the validation scope
		validation.Context(ctx),
		validation.String(&tag, it.IsNotBlank(), isTagExists),
	)

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation: Tag "movie" does not exist.
}
