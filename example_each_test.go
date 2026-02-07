package validation_test

import (
	"context"
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

// ExampleFunc demonstrates using a plain function as a [validation.Constraint] via [validation.Func].
// Func implements Constraint[T], so you can use it with [validation.This], [validation.Each], or [validation.EachProperty].
func ExampleFunc() {
	// Wrap a function to validate each string with IsNotBlank
	notBlank := validation.Func[string](func(ctx context.Context, v *validation.Validator, s string) error {
		return v.Validate(ctx, validation.String(s, it.IsNotBlank()))
	})

	err := validator.Validate(context.Background(), validation.This("", notBlank))
	fmt.Println(err)
	// Output:
	// violation: "This value should not be blank."
}

// ExampleEach demonstrates validating each element of a slice with [validation.Constraint] list.
// Violation paths include the element index (e.g. [0], [1]).
func ExampleEach() {
	notBlank := validation.Func[string](func(ctx context.Context, v *validation.Validator, s string) error {
		return v.Validate(ctx, validation.String(s, it.IsNotBlank()))
	})

	items := []string{"ok", "", "valid"}
	err := validator.Validate(context.Background(), validation.Each(items, notBlank))

	if violations, ok := validation.UnwrapViolations(err); ok {
		for el := violations.First(); el != nil; el = el.Next() {
			fmt.Println(el)
		}
	}
	// Output:
	// violation at "[1]": "This value should not be blank."
}

// ExampleEachProperty demonstrates [validation.EachProperty], which adds a property name to the violation path.
// Paths look like "tags[0]", "tags[1]" instead of "[0]", "[1]".
func ExampleEachProperty() {
	notBlank := validation.Func[string](func(ctx context.Context, v *validation.Validator, s string) error {
		return v.Validate(ctx, validation.String(s, it.IsNotBlank()))
	})

	tags := []string{"", ""}
	err := validator.Validate(context.Background(), validation.EachProperty("tags", tags, notBlank))

	if violations, ok := validation.UnwrapViolations(err); ok {
		for el := violations.First(); el != nil; el = el.Next() {
			fmt.Println(el)
		}
	}
	// Output:
	// violation at "tags[0]": "This value should not be blank."
	// violation at "tags[1]": "This value should not be blank."
}

// ExampleEach_withCustomType shows using [validation.Each] with a custom element type and [validation.Func].
func ExampleEach_withCustomType() {
	type Item struct {
		Code string
	}

	// Constraint that each item has a non-blank code
	validCode := validation.Func[Item](func(ctx context.Context, v *validation.Validator, item Item) error {
		return v.Validate(ctx, validation.StringProperty("code", item.Code, it.IsNotBlank()))
	})

	items := []Item{{Code: "A"}, {Code: ""}}
	err := validator.Validate(context.Background(), validation.Each(items, validCode))

	if violations, ok := validation.UnwrapViolations(err); ok {
		for el := violations.First(); el != nil; el = el.Next() {
			fmt.Println(el)
		}
	}
	// Output:
	// violation at "[1].code": "This value should not be blank."
}
