package validation_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/muonsoft/language"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message/translations/russian"
	"github.com/muonsoft/validation/validator"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/message/catalog"
)

func ExampleNewValidator() {
	validator, err := validation.NewValidator(
		validation.DefaultLanguage(language.English), // passing default language of translations
		validation.Translations(russian.Messages),    // setting up custom or built-in translations
		// validation.SetViolationFactory(userViolationFactory), // if you want to override creation of violations
	)
	// don't forget to check for errors
	if err != nil {
		log.Fatal(err)
	}

	err = validator.Validate(context.Background(), validation.String("", it.IsNotBlank()))
	fmt.Println(err)
	// Output:
	// violation: "This value should not be blank."
}

func ExampleNil() {
	var v []string
	err := validator.Validate(context.Background(), validation.Nil(v == nil, it.IsNotNil()))
	fmt.Println(err)
	// Output:
	// violation: "This value should not be nil."
}

func ExampleNilProperty() {
	v := struct {
		Tags []string
	}{}
	err := validator.Validate(context.Background(), validation.NilProperty("tags", v.Tags == nil, it.IsNotNil()))
	fmt.Println(err)
	// Output:
	// violation at "tags": "This value should not be nil."
}

func ExampleBool() {
	v := false
	err := validator.Validate(context.Background(), validation.Bool(v, it.IsTrue()))
	fmt.Println(err)
	// Output:
	// violation: "This value should be true."
}

func ExampleBoolProperty() {
	v := struct {
		IsPublished bool
	}{
		IsPublished: false,
	}
	err := validator.Validate(
		context.Background(),
		validation.BoolProperty("isPublished", v.IsPublished, it.IsTrue()),
	)
	fmt.Println(err)
	// Output:
	// violation at "isPublished": "This value should be true."
}

func ExampleNilBool() {
	v := false
	err := validator.Validate(context.Background(), validation.NilBool(&v, it.IsTrue()))
	fmt.Println(err)
	// Output:
	// violation: "This value should be true."
}

func ExampleNilBoolProperty() {
	v := struct {
		IsPublished bool
	}{
		IsPublished: false,
	}
	err := validator.Validate(
		context.Background(),
		validation.NilBoolProperty("isPublished", &v.IsPublished, it.IsTrue()),
	)
	fmt.Println(err)
	// Output:
	// violation at "isPublished": "This value should be true."
}

func ExampleNumber_int() {
	v := 5
	err := validator.Validate(
		context.Background(),
		validation.Number[int](v, it.IsGreaterThan(5)),
	)
	fmt.Println(err)
	// Output:
	// violation: "This value should be greater than 5."
}

func ExampleNumber_float() {
	v := 5.5
	err := validator.Validate(
		context.Background(),
		validation.Number[float64](v, it.IsGreaterThan(6.5)),
	)
	fmt.Println(err)
	// Output:
	// violation: "This value should be greater than 6.5."
}

func ExampleNumberProperty_int() {
	v := struct {
		Count int
	}{
		Count: 5,
	}
	err := validator.Validate(
		context.Background(),
		validation.NumberProperty[int]("count", v.Count, it.IsGreaterThan(5)),
	)
	fmt.Println(err)
	// Output:
	// violation at "count": "This value should be greater than 5."
}

func ExampleNumberProperty_float() {
	v := struct {
		Amount float64
	}{
		Amount: 5.5,
	}
	err := validator.Validate(
		context.Background(),
		validation.NumberProperty[float64]("amount", v.Amount, it.IsGreaterThan(6.5)),
	)
	fmt.Println(err)
	// Output:
	// violation at "amount": "This value should be greater than 6.5."
}

func ExampleNilNumber_int() {
	v := 5
	err := validator.Validate(
		context.Background(),
		validation.NilNumber[int](&v, it.IsGreaterThan(5)),
	)
	fmt.Println(err)
	// Output:
	// violation: "This value should be greater than 5."
}

func ExampleNilNumber_float() {
	v := 5.5
	err := validator.Validate(
		context.Background(),
		validation.NilNumber[float64](&v, it.IsGreaterThan(6.5)),
	)
	fmt.Println(err)
	// Output:
	// violation: "This value should be greater than 6.5."
}

func ExampleNilNumberProperty_int() {
	v := struct {
		Count int
	}{
		Count: 5,
	}
	err := validator.Validate(
		context.Background(),
		validation.NilNumberProperty[int]("count", &v.Count, it.IsGreaterThan(5)),
	)
	fmt.Println(err)
	// Output:
	// violation at "count": "This value should be greater than 5."
}

func ExampleNilNumberProperty_float() {
	v := struct {
		Amount float64
	}{
		Amount: 5.5,
	}
	err := validator.Validate(
		context.Background(),
		validation.NilNumberProperty[float64]("amount", &v.Amount, it.IsGreaterThan(6.5)),
	)
	fmt.Println(err)
	// Output:
	// violation at "amount": "This value should be greater than 6.5."
}

func ExampleString() {
	v := ""
	err := validator.Validate(
		context.Background(),
		validation.String(v, it.IsNotBlank()),
	)
	fmt.Println(err)
	// Output:
	// violation: "This value should not be blank."
}

func ExampleStringProperty() {
	v := struct {
		Title string
	}{
		Title: "",
	}
	err := validator.Validate(
		context.Background(),
		validation.StringProperty("title", v.Title, it.IsNotBlank()),
	)
	fmt.Println(err)
	// Output:
	// violation at "title": "This value should not be blank."
}

func ExampleNilString() {
	v := ""
	err := validator.Validate(
		context.Background(),
		validation.NilString(&v, it.IsNotBlank()),
	)
	fmt.Println(err)
	// Output:
	// violation: "This value should not be blank."
}

func ExampleNilStringProperty() {
	v := struct {
		Title string
	}{
		Title: "",
	}
	err := validator.Validate(
		context.Background(),
		validation.NilStringProperty("title", &v.Title, it.IsNotBlank()),
	)
	fmt.Println(err)
	// Output:
	// violation at "title": "This value should not be blank."
}

func ExampleCountable() {
	s := []string{"a", "b"}
	err := validator.Validate(
		context.Background(),
		validation.Countable(len(s), it.HasMinCount(3)),
	)
	fmt.Println(err)
	// Output:
	// violation: "This collection should contain 3 elements or more."
}

func ExampleCountableProperty() {
	v := struct {
		Tags []string
	}{
		Tags: []string{"a", "b"},
	}
	err := validator.Validate(
		context.Background(),
		validation.CountableProperty("tags", len(v.Tags), it.HasMinCount(3)),
	)
	fmt.Println(err)
	// Output:
	// violation at "tags": "This collection should contain 3 elements or more."
}

func ExampleTime() {
	t := time.Now()
	compared, _ := time.Parse(time.RFC3339, "2006-01-02T15:00:00Z")
	err := validator.Validate(
		context.Background(),
		validation.Time(t, it.IsEarlierThan(compared)),
	)
	fmt.Println(err)
	// Output:
	// violation: "This value should be earlier than 2006-01-02T15:00:00Z."
}

func ExampleTimeProperty() {
	v := struct {
		CreatedAt time.Time
	}{
		CreatedAt: time.Now(),
	}
	compared, _ := time.Parse(time.RFC3339, "2006-01-02T15:00:00Z")
	err := validator.Validate(
		context.Background(),
		validation.TimeProperty("createdAt", v.CreatedAt, it.IsEarlierThan(compared)),
	)
	fmt.Println(err)
	// Output:
	// violation at "createdAt": "This value should be earlier than 2006-01-02T15:00:00Z."
}

func ExampleNilTime() {
	t := time.Now()
	compared, _ := time.Parse(time.RFC3339, "2006-01-02T15:00:00Z")
	err := validator.Validate(
		context.Background(),
		validation.NilTime(&t, it.IsEarlierThan(compared)),
	)
	fmt.Println(err)
	// Output:
	// violation: "This value should be earlier than 2006-01-02T15:00:00Z."
}

func ExampleNilTimeProperty() {
	v := struct {
		CreatedAt time.Time
	}{
		CreatedAt: time.Now(),
	}
	compared, _ := time.Parse(time.RFC3339, "2006-01-02T15:00:00Z")
	err := validator.Validate(
		context.Background(),
		validation.NilTimeProperty("createdAt", &v.CreatedAt, it.IsEarlierThan(compared)),
	)
	fmt.Println(err)
	// Output:
	// violation at "createdAt": "This value should be earlier than 2006-01-02T15:00:00Z."
}

func ExampleEachString() {
	v := []string{""}
	err := validator.Validate(
		context.Background(),
		validation.EachString(v, it.IsNotBlank()),
	)
	fmt.Println(err)
	// Output:
	// violation at "[0]": "This value should not be blank."
}

func ExampleEachStringProperty() {
	v := struct {
		Tags []string
	}{
		Tags: []string{""},
	}
	err := validator.Validate(
		context.Background(),
		validation.EachStringProperty("tags", v.Tags, it.IsNotBlank()),
	)
	fmt.Println(err)
	// Output:
	// violation at "tags[0]": "This value should not be blank."
}

func ExampleEachNumber() {
	v := []int{-1, 0, 1}
	err := validator.Validate(
		context.Background(),
		validation.EachNumber[int](v, it.IsPositiveOrZero[int]()),
	)
	fmt.Println(err)
	// Output:
	// violation at "[0]": "This value should be either positive or zero."
}

func ExampleEachNumberProperty() {
	v := struct {
		Metrics []int
	}{
		Metrics: []int{-1, 0, 1},
	}
	err := validator.Validate(
		context.Background(),
		validation.EachNumberProperty[int]("metrics", v.Metrics, it.IsPositiveOrZero[int]()),
	)
	fmt.Println(err)
	// Output:
	// violation at "metrics[0]": "This value should be either positive or zero."
}

func ExampleEachComparable() {
	v := []string{"foo", "bar", "baz"}
	err := validator.Validate(
		context.Background(),
		validation.EachComparable[string](v, it.IsOneOf("foo", "bar", "buz")),
	)
	fmt.Println(err)
	// Output:
	// violation at "[2]": "The value you selected is not a valid choice."
}

func ExampleEachComparableProperty() {
	v := struct {
		Labels []string
	}{
		Labels: []string{"foo", "bar", "baz"},
	}
	err := validator.Validate(
		context.Background(),
		validation.EachComparableProperty[string]("labels", v.Labels, it.IsOneOf("foo", "bar", "buz")),
	)
	fmt.Println(err)
	// Output:
	// violation at "labels[2]": "The value you selected is not a valid choice."
}

func ExampleComparable_string() {
	v := "unknown"
	err := validator.Validate(
		context.Background(),
		validation.Comparable[string](v, it.IsOneOf("foo", "bar", "baz")),
	)
	fmt.Println(err)
	// Output:
	// violation: "The value you selected is not a valid choice."
}

func ExampleComparable_int() {
	v := 4
	err := validator.Validate(
		context.Background(),
		validation.Comparable[int](v, it.IsOneOf(1, 2, 3, 5)),
	)
	fmt.Println(err)
	// Output:
	// violation: "The value you selected is not a valid choice."
}

func ExampleComparableProperty_string() {
	s := struct {
		Enum string
	}{
		Enum: "unknown",
	}
	err := validator.Validate(
		context.Background(),
		validation.ComparableProperty[string]("enum", s.Enum, it.IsOneOf("foo", "bar", "baz")),
	)
	fmt.Println(err)
	// Output:
	// violation at "enum": "The value you selected is not a valid choice."
}

func ExampleComparableProperty_int() {
	s := struct {
		Metric int
	}{
		Metric: 4,
	}
	err := validator.Validate(
		context.Background(),
		validation.ComparableProperty[int]("metric", s.Metric, it.IsOneOf(1, 2, 3, 5)),
	)
	fmt.Println(err)
	// Output:
	// violation at "metric": "The value you selected is not a valid choice."
}

func ExampleNilComparable_string() {
	v := "unknown"
	err := validator.Validate(
		context.Background(),
		validation.NilComparable[string](&v, it.IsOneOf("foo", "bar", "baz")),
	)
	fmt.Println(err)
	// Output:
	// violation: "The value you selected is not a valid choice."
}

func ExampleNilComparable_int() {
	v := 4
	err := validator.Validate(
		context.Background(),
		validation.NilComparable[int](&v, it.IsOneOf(1, 2, 3, 5)),
	)
	fmt.Println(err)
	// Output:
	// violation: "The value you selected is not a valid choice."
}

func ExampleNilComparableProperty_string() {
	s := struct {
		Enum string
	}{
		Enum: "unknown",
	}
	err := validator.Validate(
		context.Background(),
		validation.NilComparableProperty[string]("enum", &s.Enum, it.IsOneOf("foo", "bar", "baz")),
	)
	fmt.Println(err)
	// Output:
	// violation at "enum": "The value you selected is not a valid choice."
}

func ExampleNilComparableProperty_int() {
	s := struct {
		Metric int
	}{
		Metric: 4,
	}
	err := validator.Validate(
		context.Background(),
		validation.NilComparableProperty[int]("metric", &s.Metric, it.IsOneOf(1, 2, 3, 5)),
	)
	fmt.Println(err)
	// Output:
	// violation at "metric": "The value you selected is not a valid choice."
}

func ExampleComparables() {
	v := []string{"foo", "bar", "baz", "foo"}
	err := validator.Validate(
		context.Background(),
		validation.Comparables[string](v, it.HasUniqueValues[string]()),
	)
	fmt.Println(err)
	// Output:
	// violation: "This collection should contain only unique elements."
}

func ExampleComparablesProperty() {
	v := struct {
		Keywords []string
	}{
		Keywords: []string{"foo", "bar", "baz", "foo"},
	}
	err := validator.Validate(
		context.Background(),
		validation.ComparablesProperty[string]("keywords", v.Keywords, it.HasUniqueValues[string]()),
	)
	fmt.Println(err)
	// Output:
	// violation at "keywords": "This collection should contain only unique elements."
}

func ExampleCheck() {
	v := 123
	err := validator.Validate(context.Background(), validation.Check(v > 321))
	fmt.Println(err)
	// Output:
	// violation: "This value is not valid."
}

func ExampleOfStringBy() {
	validate := func(s string) bool {
		return s == "valid"
	}
	errExample := errors.New("exampleCode")

	constraint := validation.OfStringBy(validate).
		WithError(errExample).           // underlying static error
		WithMessage("Unexpected value.") // violation message template

	s := "foo"
	err := validator.Validate(context.Background(), validation.String(s, constraint))

	fmt.Println(err)
	// Output:
	// violation: "Unexpected value."
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
		context.Background(),
		validation.ComparableProperty[string](
			"cardType",
			payment.CardType,
			it.IsOneOf("Visa", "MasterCard"),
		),
		validation.When(payment.CardType == "Visa").
			At(validation.PropertyName("cardNumber")).
			Then(validation.String(payment.CardNumber, it.Matches(visaRegex))).
			Else(validation.String(payment.CardNumber, it.Matches(masterCardRegex))),
	)

	fmt.Println(err)
	// Output:
	// violation at "cardNumber": "This value is not valid."
}

func ExampleWhenArgument_Then() {
	v := "foo"
	err := validator.Validate(
		context.Background(),
		validation.When(true).Then(
			validation.String(v, it.Matches(regexp.MustCompile(`^\w+$`))),
		),
	)
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleWhenArgument_Else() {
	v := "123"
	err := validator.Validate(
		context.Background(),
		validation.When(true).
			Then(validation.String(v, it.Matches(regexp.MustCompile(`^\w+$`)))).
			Else(validation.String(v, it.Matches(regexp.MustCompile(`^\d+$`)))),
	)
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleSequentially() {
	title := "bar"

	err := validator.Validate(
		context.Background(),
		validation.Sequentially(
			validation.String(title, it.IsBlank()),       // validation will fail on first argument
			validation.String(title, it.HasMinLength(5)), // this argument will be ignored
		),
	)

	fmt.Println(err)
	// Output:
	// violation: "This value should be blank."
}

func ExampleAtLeastOneOf() {
	banners := []struct {
		Name      string
		Keywords  []string
		Companies []string
		Brands    []string
	}{
		{Name: "Acme banner", Companies: []string{"Acme"}},
		{Name: "Empty banner"},
	}

	for _, banner := range banners {
		err := validator.Validate(
			context.Background(),
			validation.AtLeastOneOf(
				validation.CountableProperty("keywords", len(banner.Keywords), it.IsNotBlank()),
				validation.CountableProperty("companies", len(banner.Companies), it.IsNotBlank()),
				validation.CountableProperty("brands", len(banner.Brands), it.IsNotBlank()),
			),
		)
		if violations, ok := validation.UnwrapViolationList(err); ok {
			fmt.Println("banner", banner.Name, "is not valid:")
			for violation := violations.First(); violation != nil; violation = violation.Next() {
				fmt.Println(violation)
			}
		}
	}

	// Output:
	// banner Empty banner is not valid:
	// violation at "keywords": "This value should not be blank."
	// violation at "companies": "This value should not be blank."
	// violation at "brands": "This value should not be blank."
}

func ExampleAll() {
	book := struct {
		Name string
		Tags []string
	}{
		Name: "Very long book name",
		Tags: []string{"Fiction", "Thriller", "Science", "Fantasy"},
	}

	err := validator.Validate(
		context.Background(),
		validation.Sequentially(
			// this block passes
			validation.All(
				validation.StringProperty("name", book.Name, it.IsNotBlank()),
				validation.CountableProperty("tags", len(book.Tags), it.IsNotBlank()),
			),
			// this block fails
			validation.All(
				validation.StringProperty("name", book.Name, it.HasMaxLength(10)),
				validation.CountableProperty("tags", len(book.Tags), it.HasMaxCount(3)),
			),
		),
	)

	if violations, ok := validation.UnwrapViolationList(err); ok {
		for violation := violations.First(); violation != nil; violation = violation.Next() {
			fmt.Println(violation)
		}
	}
	// Output:
	// violation at "name": "This value is too long. It should have 10 characters or less."
	// violation at "tags": "This collection should contain 3 elements or less."
}

func ExampleValidator_Validate_basicValidation() {
	validator, err := validation.NewValidator()
	if err != nil {
		log.Fatal(err)
	}

	s := ""
	err = validator.Validate(context.Background(), validation.String(s, it.IsNotBlank()))

	fmt.Println(err)
	fmt.Println("errors.Is(err, validation.ErrIsBlank) =", errors.Is(err, validation.ErrIsBlank))
	// Output:
	// violation: "This value should not be blank."
	// errors.Is(err, validation.ErrIsBlank) = true
}

func ExampleValidator_Validate_singletonValidator() {
	s := ""

	err := validator.Validate(context.Background(), validation.String(s, it.IsNotBlank()))

	fmt.Println(err)
	// Output:
	// violation: "This value should not be blank."
}

func ExampleValidator_ValidateBool() {
	v := false
	err := validator.ValidateBool(context.Background(), v, it.IsTrue())
	fmt.Println(err)
	// Output:
	// violation: "This value should be true."
}

func ExampleValidator_ValidateInt() {
	v := 5
	err := validator.ValidateInt(context.Background(), v, it.IsGreaterThan(5))
	fmt.Println(err)
	// Output:
	// violation: "This value should be greater than 5."
}

func ExampleValidator_ValidateFloat() {
	v := 5.5
	err := validator.ValidateFloat(context.Background(), v, it.IsGreaterThan(6.5))
	fmt.Println(err)
	// Output:
	// violation: "This value should be greater than 6.5."
}

func ExampleValidator_ValidateString() {
	err := validator.ValidateString(context.Background(), "", it.IsNotBlank())

	fmt.Println(err)
	// Output:
	// violation: "This value should not be blank."
}

func ExampleValidator_ValidateStrings() {
	v := []string{"foo", "bar", "baz", "foo"}
	err := validator.ValidateStrings(context.Background(), v, it.HasUniqueValues[string]())
	fmt.Println(err)
	// Output:
	// violation: "This collection should contain only unique elements."
}

func ExampleValidator_ValidateCountable() {
	s := []string{"a", "b"}
	err := validator.ValidateCountable(context.Background(), len(s), it.HasMinCount(3))
	fmt.Println(err)
	// Output:
	// violation: "This collection should contain 3 elements or more."
}

func ExampleValidator_ValidateTime() {
	t := time.Now()
	compared, _ := time.Parse(time.RFC3339, "2006-01-02T15:00:00Z")
	err := validator.ValidateTime(context.Background(), t, it.IsEarlierThan(compared))
	fmt.Println(err)
	// Output:
	// violation: "This value should be earlier than 2006-01-02T15:00:00Z."
}

func ExampleValidator_ValidateEachString() {
	v := []string{""}
	err := validator.ValidateEachString(context.Background(), v, it.IsNotBlank())
	fmt.Println(err)
	// Output:
	// violation at "[0]": "This value should not be blank."
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
		context.Background(),
		validation.StringProperty("title", document.Title, it.IsNotBlank()),
		validation.CountableProperty("keywords", len(document.Keywords), it.HasCountBetween(5, 10)),
		validation.ComparablesProperty[string]("keywords", document.Keywords, it.HasUniqueValues[string]()),
		validation.EachStringProperty("keywords", document.Keywords, it.IsNotBlank()),
	)

	if violations, ok := validation.UnwrapViolationList(err); ok {
		for violation := violations.First(); violation != nil; violation = violation.Next() {
			fmt.Println(violation)
		}
	}
	// Output:
	// violation at "title": "This value should not be blank."
	// violation at "keywords": "This collection should contain 5 elements or more."
	// violation at "keywords": "This collection should contain only unique elements."
	// violation at "keywords[0]": "This value should not be blank."
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
			context.Background(),
			validation.StringProperty("name", note.Title, it.IsNotBlank()),
			validation.StringProperty("text", note.Text, it.IsNotBlank().When(note.IsPublic)),
		)
		if violations, ok := validation.UnwrapViolationList(err); ok {
			for violation := violations.First(); violation != nil; violation = violation.Next() {
				fmt.Printf("error on note %d: %s", i, violation)
			}
		}
	}

	// Output:
	// error on note 1: violation at "text": "This value should not be blank."
}

func ExampleValidator_Validate_passingPropertyPathViaOptions() {
	s := ""

	err := validator.Validate(
		context.Background(),
		validation.String(s, it.IsNotBlank()).At(
			validation.PropertyName("properties"),
			validation.ArrayIndex(1),
			validation.PropertyName("tag"),
		),
	)

	violation := err.(*validation.ViolationList).First()
	fmt.Println("property path:", violation.PropertyPath().String())
	// Output:
	// property path: properties[1].tag
}

func ExampleValidator_Validate_propertyPathWithContextValidator() {
	s := ""

	err := validator.
		AtProperty("properties").
		AtIndex(1).
		AtProperty("tag").
		Validate(context.Background(), validation.String(s, it.IsNotBlank()))

	violation := err.(*validation.ViolationList).First()
	fmt.Println("property path:", violation.PropertyPath().String())
	// Output:
	// property path: properties[1].tag
}

func ExampleValidator_Validate_propertyPathBySpecialArgument() {
	s := ""

	err := validator.Validate(
		context.Background(),
		// this is an alias for
		// validation.String(s, it.IsNotBlank()).At(validation.PropertyName("property")),
		validation.StringProperty("property", s, it.IsNotBlank()),
	)

	violation := err.(*validation.ViolationList).First()
	fmt.Println("property path:", violation.PropertyPath().String())
	// Output:
	// property path: property
}

func ExampleValidator_At() {
	books := []struct {
		Title string
	}{
		{Title: ""},
	}

	err := validator.At(validation.PropertyName("books"), validation.ArrayIndex(0)).Validate(
		context.Background(),
		validation.StringProperty("title", books[0].Title, it.IsNotBlank()),
	)

	violation := err.(*validation.ViolationList).First()
	fmt.Println("property path:", violation.PropertyPath().String())
	// Output:
	// property path: books[0].title
}

func ExampleValidator_AtProperty() {
	book := struct {
		Title string
	}{
		Title: "",
	}

	err := validator.AtProperty("book").Validate(
		context.Background(),
		validation.StringProperty("title", book.Title, it.IsNotBlank()),
	)

	violation := err.(*validation.ViolationList).First()
	fmt.Println("property path:", violation.PropertyPath().String())
	// Output:
	// property path: book.title
}

func ExampleValidator_AtIndex() {
	books := []struct {
		Title string
	}{
		{Title: ""},
	}

	err := validator.AtIndex(0).Validate(
		context.Background(),
		validation.StringProperty("title", books[0].Title, it.IsNotBlank()),
	)

	violation := err.(*validation.ViolationList).First()
	fmt.Println("property path:", violation.PropertyPath().String())
	// Output:
	// property path: [0].title
}

func ExampleValidator_WithLanguage() {
	validator, err := validation.NewValidator(validation.Translations(russian.Messages))
	if err != nil {
		log.Fatal(err)
	}

	s := ""
	err = validator.WithLanguage(language.Russian).Validate(
		context.Background(),
		validation.String(s, it.IsNotBlank()),
	)

	fmt.Println(err)
	// Output:
	// violation: "Значение не должно быть пустым."
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
	err = validator.Validate(context.Background(), validation.String(s, it.IsNotBlank()))

	fmt.Println(err)
	// Output:
	// violation: "Значение не должно быть пустым."
}

func ExampleValidator_Validate_translationsByContextualValidator() {
	validator, err := validation.NewValidator(
		validation.Translations(russian.Messages),
	)
	if err != nil {
		log.Fatal(err)
	}

	s := ""
	err = validator.WithLanguage(language.Russian).Validate(
		context.Background(),
		validation.String(s, it.IsNotBlank()),
	)

	fmt.Println(err)
	// Output:
	// violation: "Значение не должно быть пустым."
}

func ExampleValidator_Validate_translationsByContextArgument() {
	validator, err := validation.NewValidator(
		validation.Translations(russian.Messages),
	)
	if err != nil {
		log.Fatal(err)
	}

	s := ""
	ctx := language.WithContext(context.Background(), language.Russian)
	err = validator.Validate(
		ctx,
		validation.String(s, it.IsNotBlank()),
	)

	fmt.Println(err)
	// Output:
	// violation: "Значение не должно быть пустым."
}

func ExampleTranslations() {
	validator, err := validation.NewValidator(
		validation.Translations(russian.Messages),
	)
	if err != nil {
		log.Fatal(err)
	}

	s := ""
	ctx := language.WithContext(context.Background(), language.Russian)
	err = validator.Validate(
		ctx,
		validation.String(s, it.IsNotBlank()),
	)

	fmt.Println(err)
	// Output:
	// violation: "Значение не должно быть пустым."
}

func ExampleValidator_Validate_customizingErrorMessage() {
	s := ""

	err := validator.Validate(
		context.Background(),
		validation.String(s, it.IsNotBlank().WithMessage("this value is required")),
	)

	fmt.Println(err)
	// Output:
	// violation: "this value is required"
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
	err = validator.WithLanguage(language.Russian).Validate(
		context.Background(),
		validation.Countable(len(tags), it.HasMinCount(1).WithMinMessage(customMessage)),
	)

	fmt.Println(err)
	// Output:
	// violation: "теги должны содержать 1 элемент и более"
}

func ExampleValidator_CreateViolation() {
	validator, err := validation.NewValidator()
	if err != nil {
		log.Fatal(err)
	}
	errClient := errors.New("client error")

	violation := validator.CreateViolation(
		context.Background(),
		errClient,
		"Client message.",
		validation.PropertyName("properties"),
		validation.ArrayIndex(1),
	)

	fmt.Println(violation.Error())
	// Output:
	// violation at "properties[1]": "Client message."
}

func ExampleValidator_BuildViolation_buildingViolation() {
	validator, err := validation.NewValidator()
	if err != nil {
		log.Fatal(err)
	}
	errClient := errors.New("client error")

	violation := validator.BuildViolation(context.Background(), errClient, "Client message with {{ parameter }}.").
		WithParameter("{{ parameter }}", "value").
		Create()

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
	errClient := errors.New("client error")

	violation := validator.WithLanguage(language.Russian).
		BuildViolation(context.Background(), errClient, "The operation is only possible for the {{ role }}.").
		WithParameters(validation.TemplateParameter{
			Key:              "{{ role }}",
			Value:            "administrator role",
			NeedsTranslation: true,
		}).
		Create()

	fmt.Println(violation.Message())
	// Output:
	// Операция возможна только для роли администратора.
}

func ExampleValidator_BuildViolationList() {
	validator, err := validation.NewValidator()
	if err != nil {
		log.Fatal(err)
	}
	errFirst := errors.New("error 1")
	errSecond := errors.New("error 2")

	builder := validator.BuildViolationList(context.Background())
	builder.BuildViolation(errFirst, "Client message with {{ parameter 1 }}.").
		WithParameter("{{ parameter 1 }}", "value 1").
		AtProperty("properties").AtIndex(0).
		Add()
	builder.BuildViolation(errSecond, "Client message with {{ parameter 2 }}.").
		WithParameter("{{ parameter 2 }}", "value 2").
		AtProperty("properties").AtIndex(1).
		Add()
	violations := builder.Create()

	violations.ForEach(func(i int, violation validation.Violation) error {
		fmt.Println(violation.Error())
		return nil
	})
	fmt.Println("errors.Is(violations, errFirst) =", errors.Is(violations, errFirst))
	fmt.Println("errors.Is(violations, errSecond) =", errors.Is(violations, errSecond))
	// Output:
	// violation at "properties[0]": "Client message with value 1."
	// violation at "properties[1]": "Client message with value 2."
	// errors.Is(violations, errFirst) = true
	// errors.Is(violations, errSecond) = true
}

func ExampleViolationList_First() {
	violations := validation.NewViolationList(
		validator.BuildViolation(context.Background(), validation.ErrNotValid, "foo").Create(),
		validator.BuildViolation(context.Background(), validation.ErrNotValid, "bar").Create(),
	)

	for violation := violations.First(); violation != nil; violation = violation.Next() {
		fmt.Println(violation)
	}
	// Output:
	// violation: "foo"
	// violation: "bar"
}

func ExampleViolationList_AppendFromError_addingViolation() {
	violations := validation.NewViolationList()
	err := validator.BuildViolation(context.Background(), validation.ErrNotValid, "foo").Create()

	appendErr := violations.AppendFromError(err)

	fmt.Println("append error:", appendErr)
	fmt.Println("violations:", violations)
	// Output:
	// append error: <nil>
	// violations: violation: "foo"
}

func ExampleViolationList_AppendFromError_addingViolationList() {
	violations := validation.NewViolationList()
	err := validation.NewViolationList(
		validator.BuildViolation(context.Background(), validation.ErrNotValid, "foo").Create(),
		validator.BuildViolation(context.Background(), validation.ErrNotValid, "bar").Create(),
	)

	appendErr := violations.AppendFromError(err)

	fmt.Println("append error:", appendErr)
	fmt.Println("violations:", violations)
	// Output:
	// append error: <nil>
	// violations: violations: #0: "foo"; #1: "bar"
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
