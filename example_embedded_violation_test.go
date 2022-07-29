package validation_test

import (
	"context"
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

var ErrModificationProhibited = validation.NewError(
	"modification is prohibited",
	"Modification of resource is prohibited.",
)

type AccessViolation struct {
	validation.Violation
	UserID     int
	Permission string
}

func (err *AccessViolation) Error() string {
	return err.Violation.Error()
}

type Blog struct {
	Name    string
	Entries BlogEntries
}

func (b Blog) Validate(ctx context.Context, validator *validation.Validator, userID int) error {
	return validator.Validate(
		ctx,
		validation.StringProperty("name", b.Name, it.IsNotBlank(), it.HasMaxLength(50)),
		validation.ValidProperty(
			"entries",
			validation.ValidatableFunc(func(ctx context.Context, validator *validation.Validator) error {
				// passing user id further
				return b.Entries.Validate(ctx, validator, userID)
			}),
		),
	)
}

type BlogEntry struct {
	AuthorID int
	Title    string
	Text     string
}

func (e BlogEntry) Validate(ctx context.Context, validator *validation.Validator, userID int) error {
	// creating violation with domain payload
	if e.AuthorID != userID {
		return &AccessViolation{
			Violation:  validator.CreateViolation(ctx, ErrModificationProhibited, ErrModificationProhibited.Message()),
			UserID:     userID,
			Permission: "edit",
		}
	}

	return validator.Validate(
		ctx,
		validation.StringProperty("title", e.Title, it.IsNotBlank(), it.HasMaxLength(100)),
		validation.StringProperty("text", e.Text, it.IsNotBlank(), it.HasMaxLength(10000)),
	)
}

type BlogEntries []BlogEntry

func (entries BlogEntries) Validate(ctx context.Context, validator *validation.Validator, userID int) error {
	violations := validation.NewViolationList()

	for i, entry := range entries {
		err := violations.AppendFromError(entry.Validate(ctx, validator.AtIndex(i), userID))
		if err != nil {
			return err
		}
	}

	return violations.AsError()
}

func ExampleValidator_ValidateIt_violationWithPayload() {
	blog := Blog{
		Name: "News blog",
		Entries: []BlogEntry{
			{
				AuthorID: 123,
				Title:    "Good weather",
				Text:     "Good weather is coming!",
			},
			{
				AuthorID: 321,
				Title:    "Secret entry",
				Text:     "This should not be edited!",
			},
		},
	}

	userID := 123 // user id from session
	err := validator.ValidateIt(
		context.Background(),
		validation.ValidatableFunc(func(ctx context.Context, validator *validation.Validator) error {
			return blog.Validate(ctx, validator, userID)
		}),
	)

	if violations, ok := validation.UnwrapViolationList(err); ok {
		violations.ForEach(func(i int, violation validation.Violation) error {
			fmt.Println(violation)
			// unwrap concrete violation from chain
			if accessError, ok := violation.(*AccessViolation); ok {
				fmt.Println("user id:", accessError.UserID)
				fmt.Println("permission:", accessError.Permission)
			}
			return nil
		})
	}
	// Output:
	// violation at "entries[1]": "Modification of resource is prohibited."
	// user id: 123
	// permission: edit
}
