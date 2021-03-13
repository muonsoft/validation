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

type Entity struct {
	Name string
}

type EntityRepository struct {
	entities []Entity
}

func (repository *EntityRepository) FindByName(ctx context.Context, name string) ([]Entity, error) {
	contextValue, ok := ctx.Value(exampleKey).(string)
	if !ok {
		return nil, errors.New("context value missing")
	}
	if contextValue != "value" {
		return nil, errors.New("invalid context value")
	}

	found := make([]Entity, 0)

	for _, entity := range repository.entities {
		if entity.Name == name {
			found = append(found, entity)
		}
	}

	return found, nil
}

type UniqueEntityConstraint struct {
	repository *EntityRepository
}

func (c *UniqueEntityConstraint) SetUp(scope *validation.Scope) error {
	return nil
}

func (c *UniqueEntityConstraint) GetName() string {
	return "UniqueEntityConstraint"
}

func (c *UniqueEntityConstraint) ValidateString(value *string, scope validation.Scope) error {
	// usually, you should ignore empty values
	// to check for an empty value you should use it.NotBlankConstraint
	if value == nil || *value == "" {
		return nil
	}

	// you can pass the context value from the scope
	entities, err := c.repository.FindByName(scope.Context(), *value)
	// here you can return a service error so that the validation process
	// is stopped immediately
	if err != nil {
		return err
	}
	if len(entities) == 0 {
		return nil
	}

	// use the scope to build violation with translations
	return scope.
		BuildViolation("notUnique", `Entity with name "{{ name }}" already exists.`).
		// you can inject parameter value to the message here
		SetParameter("{{ name }}", *value).
		GetViolation()
}

func ExampleValidator_Validate_customServiceConstraint() {
	repository := &EntityRepository{entities: []Entity{{"camera"}, {"book"}}}
	isEntityUnique := &UniqueEntityConstraint{repository: repository}

	entity := Entity{Name: "book"}
	ctx := context.WithValue(context.Background(), exampleKey, "value")

	err := validator.Validate(
		// you can pass here the context value to the validation scope
		validation.Context(ctx),
		validation.String(&entity.Name, it.IsNotBlank(), isEntityUnique),
	)

	violations := err.(validation.ViolationList)
	for _, violation := range violations {
		fmt.Println(violation.Error())
	}
	// Output:
	// violation: Entity with name "book" already exists.
}
