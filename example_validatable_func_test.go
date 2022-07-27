package validation_test

import (
	"context"
	"fmt"
	"net/url"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

type Webstore struct {
	Name  string
	URL   string
	Items WebstoreItems
}

func (w Webstore) Validate(ctx context.Context, validator *validation.Validator) error {
	domain := ""
	u, err := url.Parse(w.URL)
	if err == nil {
		domain = u.Host
	}

	return validator.Validate(
		ctx,
		validation.StringProperty("name", w.Name, it.IsNotBlank(), it.HasMaxLength(100)),
		validation.StringProperty("url", w.URL, it.IsNotBlank(), it.IsURL()),
		validation.CountableProperty("items", len(w.Items), it.IsNotBlank(), it.HasMaxCount(20)),
		validation.ValidProperty(
			"items",
			validation.ValidatableFunc(func(ctx context.Context, validator *validation.Validator) error {
				// passing context value by callback function
				return w.Items.Validate(ctx, validator, domain)
			}),
		),
	)
}

type WebstoreItem struct {
	Name  string
	URL   string
	Price int
}

func (w WebstoreItem) Validate(ctx context.Context, validator *validation.Validator, webstoreDomain string) error {
	return validator.Validate(
		ctx,
		validation.StringProperty("name", w.Name, it.IsNotBlank(), it.HasMaxLength(100)),
		validation.NumberProperty[int]("price", w.Price, it.IsNotBlankNumber[int](), it.IsLessThan[int](10000)),
		validation.StringProperty(
			"url", w.URL,
			it.IsNotBlank(),
			// using webstore domain passed as a function parameter
			it.IsURL().WithHosts(webstoreDomain).WithProhibitedMessage(
				`Webstore item URL domain must match "{{ webstoreDomain }}".`,
				validation.TemplateParameter{Key: "{{ webstoreDomain }}", Value: webstoreDomain},
			),
		),
	)
}

type WebstoreItems []WebstoreItem

func (items WebstoreItems) Validate(ctx context.Context, validator *validation.Validator, webstoreDomain string) error {
	violations := validation.NewViolationList()

	for i, item := range items {
		// passing context value to each item
		err := violations.AppendFromError(item.Validate(ctx, validator.AtIndex(i), webstoreDomain))
		if err != nil {
			return err
		}
	}

	return violations.AsError()
}

func ExampleValidatableFunc_Validate() {
	store := Webstore{
		Name: "Acme store",
		URL:  "https://acme.com/homepage",
		Items: []WebstoreItem{
			{
				Name:  "Book",
				URL:   "https://acme.com/items/the-book",
				Price: 100,
			},
			{
				Name:  "Notepad",
				URL:   "https://store.com/items/notepad",
				Price: 1000,
			},
		},
	}

	err := validator.Validate(context.Background(), validation.Valid(store))

	if violations, ok := validation.UnwrapViolationList(err); ok {
		violations.ForEach(func(i int, violation validation.Violation) error {
			fmt.Println(violation)
			return nil
		})
	}
	// Output:
	// violation at "items[1].url": "Webstore item URL domain must match "acme.com"."
}
