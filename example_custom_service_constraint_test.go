package validation_test

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
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

type StockItem struct {
	Name string
	Tags []string
}

func (s StockItem) Validate(ctx context.Context, validator *validation.Validator) error {
	isTagExists, ok := validator.GetConstraint("isTagExists").(validation.StringConstraint)
	if !ok {
		return validation.ConstraintNotFoundError{Key: "isTagExists", Type: "validation.StringConstraint"}
	}

	return validator.Validate(
		ctx,
		validation.StringProperty("name", s.Name, it.IsNotBlank(), it.HasMaxLength(20)),
		validation.EachStringProperty("tags", s.Tags, isTagExists),
	)
}

func ExampleValidator_GetConstraint_customServiceConstraint() {
	storage := &TagStorage{tags: []string{"movie", "book"}}
	isTagExists := &ExistingTagConstraint{storage: storage}

	// custom constraint can be stored in the validator's internal store
	// and can be used later by calling the validator.GetConstraint method
	validator, err := validation.NewValidator(
		validation.StoredConstraint("isTagExists", isTagExists),
	)
	if err != nil {
		log.Fatal(err)
	}

	item := StockItem{
		Name: "War and peace",
		Tags: []string{"book", "camera"},
	}

	err = validator.Validate(
		// you can pass here the context value to the validation scope
		context.WithValue(context.Background(), exampleKey, "value"),
		validation.Valid(item),
	)

	fmt.Println(err)
	// Output:
	// violation at 'tags[1]': Tag "camera" does not exist.
}
