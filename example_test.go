package validation_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	mslanguage "github.com/muonsoft/language"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message/translations/russian"
	"github.com/muonsoft/validation/validator"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

func ExampleValue() {
	v := ""
	err := validator.Validate(validation.Value(v, it.IsNotBlank()))
	fmt.Println(err)
	// Output:
	// violation: This value should not be blank.
}

func ExamplePropertyValue() {
	v := Book{Title: ""}
	err := validator.Validate(
		validation.PropertyValue("title", v.Title, it.IsNotBlank()),
	)
	fmt.Println(err)
	// Output:
	// violation at 'title': This value should not be blank.
}

func ExampleBool() {
	v := false
	err := validator.Validate(validation.Bool(&v, it.IsTrue()))
	fmt.Println(err)
	// Output:
	// violation: This value should be true.
}

func ExampleBoolProperty() {
	v := struct {
		IsPublished bool
	}{
		IsPublished: false,
	}
	err := validator.Validate(
		validation.BoolProperty("isPublished", &v.IsPublished, it.IsTrue()),
	)
	fmt.Println(err)
	// Output:
	// violation at 'isPublished': This value should be true.
}

func ExampleNumber() {
	v := 5
	err := validator.Validate(validation.Number(&v, it.IsGreaterThanInteger(5)))
	fmt.Println(err)
	// Output:
	// violation: This value should be greater than 5.
}

func ExampleNumberProperty() {
	v := struct {
		Count int
	}{
		Count: 5,
	}
	err := validator.Validate(
		validation.NumberProperty("count", &v.Count, it.IsGreaterThanInteger(5)),
	)
	fmt.Println(err)
	// Output:
	// violation at 'count': This value should be greater than 5.
}

func ExampleString() {
	v := ""
	err := validator.Validate(validation.String(&v, it.IsNotBlank()))
	fmt.Println(err)
	// Output:
	// violation: This value should not be blank.
}

func ExampleStringProperty() {
	v := Book{Title: ""}
	err := validator.Validate(
		validation.StringProperty("title", &v.Title, it.IsNotBlank()),
	)
	fmt.Println(err)
	// Output:
	// violation at 'title': This value should not be blank.
}

func ExampleStrings() {
	v := []string{"foo", "bar", "baz", "foo"}
	err := validator.Validate(
		validation.Strings(v, it.HasUniqueValues()),
	)
	fmt.Println(err)
	// Output:
	// violation: This collection should contain only unique elements.
}

func ExampleStringsProperty() {
	v := Book{Keywords: []string{"foo", "bar", "baz", "foo"}}
	err := validator.Validate(
		validation.StringsProperty("keywords", v.Keywords, it.HasUniqueValues()),
	)
	fmt.Println(err)
	// Output:
	// violation at 'keywords': This collection should contain only unique elements.
}

func ExampleIterable() {
	v := make([]string, 0)
	err := validator.Validate(validation.Iterable(v, it.IsNotBlank()))
	fmt.Println(err)
	// Output:
	// violation: This value should not be blank.
}

func ExampleIterableProperty() {
	v := Product{Tags: []string{}}
	err := validator.Validate(
		validation.IterableProperty("tags", v.Tags, it.IsNotBlank()),
	)
	fmt.Println(err)
	// Output:
	// violation at 'tags': This value should not be blank.
}

func ExampleCountable() {
	s := []string{"a", "b"}
	err := validator.Validate(validation.Countable(len(s), it.HasMinCount(3)))
	fmt.Println(err)
	// Output:
	// violation: This collection should contain 3 elements or more.
}

func ExampleCountableProperty() {
	v := Product{Tags: []string{"a", "b"}}
	err := validator.Validate(
		validation.CountableProperty("tags", len(v.Tags), it.HasMinCount(3)),
	)
	fmt.Println(err)
	// Output:
	// violation at 'tags': This collection should contain 3 elements or more.
}

func ExampleTime() {
	t := time.Now()
	compared, _ := time.Parse(time.RFC3339, "2006-01-02T15:00:00Z")
	err := validator.Validate(
		validation.Time(&t, it.IsEarlierThan(compared)),
	)
	fmt.Println(err)
	// Output:
	// violation: This value should be earlier than 2006-01-02T15:00:00Z.
}

func ExampleTimeProperty() {
	v := struct {
		CreatedAt time.Time
	}{
		CreatedAt: time.Now(),
	}
	compared, _ := time.Parse(time.RFC3339, "2006-01-02T15:00:00Z")
	err := validator.Validate(
		validation.TimeProperty("createdAt", &v.CreatedAt, it.IsEarlierThan(compared)),
	)
	fmt.Println(err)
	// Output:
	// violation at 'createdAt': This value should be earlier than 2006-01-02T15:00:00Z.
}

func ExampleEach() {
	v := []string{""}
	err := validator.Validate(validation.Each(v, it.IsNotBlank()))
	fmt.Println(err)
	// Output:
	// violation at '[0]': This value should not be blank.
}

func ExampleEachProperty() {
	v := Product{Tags: []string{""}}
	err := validator.Validate(
		validation.EachProperty("tags", v.Tags, it.IsNotBlank()),
	)
	fmt.Println(err)
	// Output:
	// violation at 'tags[0]': This value should not be blank.
}

func ExampleEachString() {
	v := []string{""}
	err := validator.Validate(validation.EachString(v, it.IsNotBlank()))
	fmt.Println(err)
	// Output:
	// violation at '[0]': This value should not be blank.
}

func ExampleEachStringProperty() {
	v := Product{Tags: []string{""}}
	err := validator.Validate(
		validation.EachStringProperty("tags", v.Tags, it.IsNotBlank()),
	)
	fmt.Println(err)
	// Output:
	// violation at 'tags[0]': This value should not be blank.
}

func ExampleNewCustomStringConstraint() {
	validate := func(s string) bool {
		return s == "valid"
	}
	constraint := validation.NewCustomStringConstraint(
		validate,
		"ExampleConstraint", // constraint name
		"exampleCode",       // violation code
		"Unexpected value.", // violation message template
	)

	s := "foo"
	err := validator.ValidateString(&s, constraint)

	fmt.Println(err)
	// Output:
	// violation: Unexpected value.
}

func ExampleWhen() {
	visaRegex := regexp.MustCompile("^4[0-9]{12}(?:[0-9]{3})?$")
	masterCardRegex := regexp.MustCompile("^(5[1-5][0-9]{14}|2(22[1-9][0-9]{12}|2[3-9][0-9]{13}|[3-6][0-9]{14}|7[0-1][0-9]{13}|720[0-9]{12}))$")

	payment := struct {
		CardType   string
		CardNumber string
		Amount     int
	}{
		CardType:   "Visa",
		CardNumber: "4111",
		Amount:     1000,
	}

	err := validator.Validate(
		validation.StringProperty(
			"cardType",
			&payment.CardType,
			it.IsOneOfStrings("Visa", "MasterCard"),
		),
		validation.StringProperty(
			"cardNumber",
			&payment.CardNumber,
			validation.
				When(payment.CardType == "Visa").
				Then(it.Matches(visaRegex)).
				Else(it.Matches(masterCardRegex)),
		),
	)

	fmt.Println(err)
	// Output:
	// violation at 'cardNumber': This value is not valid.
}

func ExampleConditionalConstraint_Then() {
	v := "foo"
	err := validator.ValidateString(
		&v,
		validation.When(true).
			Then(
				it.Matches(regexp.MustCompile(`^\w+$`)),
			),
	)
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleConditionalConstraint_Else() {
	v := "123"
	err := validator.ValidateString(
		&v,
		validation.When(false).
			Then(
				it.Matches(regexp.MustCompile(`^\w+$`)),
			).
			Else(
				it.Matches(regexp.MustCompile(`^\d+$`)),
			),
	)
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleSequentially() {
	title := "bar"

	err := validator.ValidateString(
		&title,
		validation.Sequentially(
			it.IsBlank(),
			it.HasMinLength(5),
		),
	)

	fmt.Println(err)
	// Output:
	// violation: This value should be blank.
}

func ExampleAtLeastOneOf() {
	title := "bar"

	err := validator.ValidateString(
		&title,
		validation.AtLeastOneOf(
			it.IsBlank(),
			it.HasMinLength(5),
		),
	)

	if violations, ok := validation.UnwrapViolationList(err); ok {
		for violation := violations.First(); violation != nil; violation = violation.Next() {
			fmt.Println(violation)
		}
	}
	// Output:
	// violation: This value should be blank.
	// violation: This value is too short. It should have 5 characters or more.
}

func ExampleCompound() {
	title := "bar"
	isEmail := validation.Compound(it.IsEmail(), it.HasLengthBetween(5, 200))

	err := validator.ValidateString(
		&title,
		isEmail,
	)

	if violations, ok := validation.UnwrapViolationList(err); ok {
		for violation := violations.First(); violation != nil; violation = violation.Next() {
			fmt.Println(violation)
		}
	}
	// Output:
	// violation: This value is not a valid email address.
	// violation: This value is too short. It should have 5 characters or more.
}

func ExampleValidator_Validate_basicValidation() {
	s := ""

	validator, err := validation.NewValidator()
	if err != nil {
		log.Fatal(err)
	}
	err = validator.Validate(validation.String(&s, it.IsNotBlank()))

	fmt.Println(err)
	// Output:
	// violation: This value should not be blank.
}

func ExampleValidator_Validate_singletonValidator() {
	s := ""

	err := validator.Validate(validation.String(&s, it.IsNotBlank()))

	fmt.Println(err)
	// Output:
	// violation: This value should not be blank.
}

func ExampleValidator_ValidateString_shorthandAlias() {
	s := ""

	err := validator.ValidateString(&s, it.IsNotBlank())

	fmt.Println(err)
	// Output:
	// violation: This value should not be blank.
}

func ExampleValidator_Validate_basicStructValidation() {
	document := struct {
		Title    string
		Keywords []string
	}{
		Title:    "",
		Keywords: []string{"", "book", "fantasy", "book"},
	}

	err := validator.Validate(
		validation.StringProperty("title", &document.Title, it.IsNotBlank()),
		validation.CountableProperty("keywords", len(document.Keywords), it.HasCountBetween(5, 10)),
		validation.StringsProperty("keywords", document.Keywords, it.HasUniqueValues()),
		validation.EachStringProperty("keywords", document.Keywords, it.IsNotBlank()),
	)

	if violations, ok := validation.UnwrapViolationList(err); ok {
		for violation := violations.First(); violation != nil; violation = violation.Next() {
			fmt.Println(violation)
		}
	}
	// Output:
	// violation at 'title': This value should not be blank.
	// violation at 'keywords': This collection should contain 5 elements or more.
	// violation at 'keywords': This collection should contain only unique elements.
	// violation at 'keywords[0]': This value should not be blank.
}

func ExampleValidator_Validate_conditionalValidationOnConstraint() {
	notes := []struct {
		Title    string
		IsPublic bool
		Text     string
	}{
		{Title: "published note", IsPublic: true, Text: "text of published note"},
		{Title: "draft note", IsPublic: true, Text: ""},
	}

	for i, note := range notes {
		err := validator.Validate(
			validation.StringProperty("name", &note.Title, it.IsNotBlank()),
			validation.StringProperty("text", &note.Text, it.IsNotBlank().When(note.IsPublic)),
		)
		if violations, ok := validation.UnwrapViolationList(err); ok {
			for violation := violations.First(); violation != nil; violation = violation.Next() {
				fmt.Printf("error on note %d: %s", i, violation)
			}
		}
	}

	// Output:
	// error on note 1: violation at 'text': This value should not be blank.
}

func ExampleValidator_Validate_passingPropertyPathViaOptions() {
	s := ""

	err := validator.Validate(
		validation.String(
			&s,
			validation.PropertyName("properties"),
			validation.ArrayIndex(1),
			validation.PropertyName("tag"),
			it.IsNotBlank(),
		),
	)

	violation := err.(*validation.ViolationList).First()
	fmt.Println("property path:", violation.PropertyPath().String())
	// Output:
	// property path: properties[1].tag
}

func ExampleValidator_Validate_propertyPathWithScopedValidator() {
	s := ""

	err := validator.
		AtProperty("properties").
		AtIndex(1).
		AtProperty("tag").
		Validate(validation.String(&s, it.IsNotBlank()))

	violation := err.(*validation.ViolationList).First()
	fmt.Println("property path:", violation.PropertyPath().String())
	// Output:
	// property path: properties[1].tag
}

func ExampleValidator_Validate_propertyPathBySpecialArgument() {
	s := ""

	err := validator.Validate(
		// this is an alias for
		// validation.String(&s, validation.PropertyName("property"), it.IsNotBlank()),
		validation.StringProperty("property", &s, it.IsNotBlank()),
	)

	violation := err.(*validation.ViolationList).First()
	fmt.Println("property path:", violation.PropertyPath().String())
	// Output:
	// property path: property
}

func ExampleValidator_AtProperty() {
	book := &Book{Title: ""}

	err := validator.AtProperty("book").Validate(
		validation.StringProperty("title", &book.Title, it.IsNotBlank()),
	)

	violation := err.(*validation.ViolationList).First()
	fmt.Println("property path:", violation.PropertyPath().String())
	// Output:
	// property path: book.title
}

func ExampleValidator_AtIndex() {
	books := []Book{{Title: ""}}

	err := validator.AtIndex(0).Validate(
		validation.StringProperty("title", &books[0].Title, it.IsNotBlank()),
	)

	violation := err.(*validation.ViolationList).First()
	fmt.Println("property path:", violation.PropertyPath().String())
	// Output:
	// property path: [0].title
}

func ExampleValidator_Validate_translationsByDefaultLanguage() {
	validator, err := validation.NewValidator(
		validation.Translations(russian.Messages),
		validation.DefaultLanguage(language.Russian),
	)
	if err != nil {
		log.Fatal(err)
	}

	s := ""
	err = validator.ValidateString(&s, it.IsNotBlank())

	fmt.Println(err)
	// Output:
	// violation: Значение не должно быть пустым.
}

func ExampleValidator_Validate_translationsByArgument() {
	validator, err := validation.NewValidator(
		validation.Translations(russian.Messages),
	)
	if err != nil {
		log.Fatal(err)
	}

	s := ""
	err = validator.Validate(
		validation.Language(language.Russian),
		validation.String(&s, it.IsNotBlank()),
	)

	fmt.Println(err)
	// Output:
	// violation: Значение не должно быть пустым.
}

func ExampleValidator_Validate_translationsByContextArgument() {
	validator, err := validation.NewValidator(
		validation.Translations(russian.Messages),
	)
	if err != nil {
		log.Fatal(err)
	}

	s := ""
	ctx := mslanguage.WithContext(context.Background(), language.Russian)
	err = validator.Validate(
		validation.Context(ctx),
		validation.String(&s, it.IsNotBlank()),
	)

	fmt.Println(err)
	// Output:
	// violation: Значение не должно быть пустым.
}

func ExampleValidator_Validate_translationsByContextValidator() {
	validator, err := validation.NewValidator(
		validation.Translations(russian.Messages),
	)
	if err != nil {
		log.Fatal(err)
	}
	ctx := mslanguage.WithContext(context.Background(), language.Russian)
	validator = validator.WithContext(ctx)

	s := ""
	err = validator.ValidateString(&s, it.IsNotBlank())

	fmt.Println(err)
	// Output:
	// violation: Значение не должно быть пустым.
}

func ExampleValidator_Validate_customizingErrorMessage() {
	s := ""

	err := validator.ValidateString(&s, it.IsNotBlank().Message("this value is required"))

	fmt.Println(err)
	// Output:
	// violation: this value is required
}

func ExampleValidator_Validate_translationForCustomMessage() {
	const customMessage = "tags should contain more than {{ limit }} element(s)"
	validator, err := validation.NewValidator(
		validation.Translations(map[language.Tag]map[string]catalog.Message{
			language.Russian: {
				customMessage: plural.Selectf(1, "",
					plural.One, "теги должны содержать {{ limit }} элемент и более",
					plural.Few, "теги должны содержать более {{ limit }} элемента",
					plural.Other, "теги должны содержать более {{ limit }} элементов"),
			},
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	var tags []string
	err = validator.Validate(
		validation.Language(language.Russian),
		validation.Iterable(tags, it.HasMinCount(1).MinMessage(customMessage)),
	)

	fmt.Println(err)
	// Output:
	// violation: теги должны содержать 1 элемент и более
}

func ExampleValidator_BuildViolation_buildingViolation() {
	validator, err := validation.NewValidator()
	if err != nil {
		log.Fatal(err)
	}

	violation := validator.BuildViolation("clientCode", "Client message with {{ parameter }}.").
		AddParameter("{{ parameter }}", "value").
		CreateViolation()

	fmt.Println(violation.Message())
	// Output:
	// Client message with value.
}

func ExampleValidator_BuildViolation_translatableParameter() {
	validator, err := validation.NewValidator(
		validation.Translations(map[language.Tag]map[string]catalog.Message{
			language.Russian: {
				"The operation is only possible for the {{ role }}.": catalog.String("Операция возможна только для {{ role }}."),
				"administrator role": catalog.String("роли администратора"),
			},
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	violation := validator.WithLanguage(language.Russian).
		BuildViolation("clientCode", "The operation is only possible for the {{ role }}.").
		SetParameters(validation.TemplateParameter{
			Key:              "{{ role }}",
			Value:            "administrator role",
			NeedsTranslation: true,
		}).
		CreateViolation()

	fmt.Println(violation.Message())
	// Output:
	// Операция возможна только для роли администратора.
}

func ExampleViolationList_First() {
	violations := validation.NewViolationList(
		validator.BuildViolation("", "foo").CreateViolation(),
		validator.BuildViolation("", "bar").CreateViolation(),
	)

	for violation := violations.First(); violation != nil; violation = violation.Next() {
		fmt.Println(violation)
	}
	// Output:
	// violation: foo
	// violation: bar
}

func ExampleViolationList_AppendFromError_addingViolation() {
	violations := validation.NewViolationList()
	err := validator.BuildViolation("", "foo").CreateViolation()

	appendErr := violations.AppendFromError(err)

	fmt.Println("append error:", appendErr)
	fmt.Println("violations:", violations)
	// Output:
	// append error: <nil>
	// violations: violation: foo
}

func ExampleViolationList_AppendFromError_addingViolationList() {
	violations := validation.NewViolationList()
	err := validation.NewViolationList(
		validator.BuildViolation("", "foo").CreateViolation(),
		validator.BuildViolation("", "bar").CreateViolation(),
	)

	appendErr := violations.AppendFromError(err)

	fmt.Println("append error:", appendErr)
	fmt.Println("violations:", violations)
	// Output:
	// append error: <nil>
	// violations: violation: foo; violation: bar
}

func ExampleViolationList_AppendFromError_addingError() {
	violations := validation.NewViolationList()
	err := errors.New("error")

	appendErr := violations.AppendFromError(err)

	fmt.Println("append error:", appendErr)
	fmt.Println("violations length:", violations.Len())
	// Output:
	// append error: error
	// violations length: 0
}
