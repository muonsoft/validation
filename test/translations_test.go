package test

import (
	"context"
	"strconv"
	"testing"

	"github.com/muonsoft/language"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message/translations/russian"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/assert"
	textlanguage "golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

func TestValidator_Validate_WhenRussianIsDefaultLanguage_ExpectViolationTranslated(t *testing.T) {
	v := newValidator(
		t,
		validation.DefaultLanguage(language.Russian),
		validation.Translations(russian.Messages),
	)

	tests := []struct {
		maxCount        int
		expectedMessage string
	}{
		{0, "Эта коллекция должна содержать 0 элементов или меньше."},
		{1, "Эта коллекция должна содержать 1 элемент или меньше."},
		{2, "Эта коллекция должна содержать 2 элемента или меньше."},
		{5, "Эта коллекция должна содержать 5 элементов или меньше."},
	}
	for _, test := range tests {
		t.Run("plural form for "+strconv.Itoa(test.maxCount), func(t *testing.T) {
			err := v.ValidateCountable(context.Background(), 10, it.HasMaxCount(test.maxCount))

			validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
				t.Helper()
				return assert.Len(t, violations, 1) &&
					assert.Equal(t, test.expectedMessage, violations[0].Message())
			})
		})
	}
}

func TestValidator_Validate_WhenRussianIsPassedViaArgument_ExpectViolationTranslated(t *testing.T) {
	v := newValidator(t, validation.Translations(russian.Messages))

	tests := []struct {
		maxCount        int
		expectedMessage string
	}{
		{0, "Эта коллекция должна содержать 0 элементов или меньше."},
		{1, "Эта коллекция должна содержать 1 элемент или меньше."},
		{2, "Эта коллекция должна содержать 2 элемента или меньше."},
		{5, "Эта коллекция должна содержать 5 элементов или меньше."},
	}
	for _, test := range tests {
		t.Run("plural form for "+strconv.Itoa(test.maxCount), func(t *testing.T) {
			err := v.Validate(
				context.Background(),
				validation.Language(language.Russian),
				validation.Countable(10, it.HasMaxCount(test.maxCount)),
			)

			validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
				t.Helper()
				return assert.Len(t, violations, 1) &&
					assert.Equal(t, test.expectedMessage, violations[0].Message())
			})
		})
	}
}

func TestValidator_Validate_WhenCustomDefaultLanguageAndUndefinedTranslationLanguage_ExpectDefaultLanguageUsed(t *testing.T) {
	v := newValidator(
		t,
		validation.DefaultLanguage(language.Russian),
		validation.Translations(russian.Messages),
	)

	err := v.Validate(
		context.Background(),
		validation.Language(language.Afrikaans),
		validation.String(stringValue(""), it.IsNotBlank()),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, "Значение не должно быть пустым.", violations[0].Message())
	})
}

func TestValidator_Validate_WhenDefaultLanguageIsNotLoaded_ExpectError(t *testing.T) {
	v, err := validation.NewValidator(validation.DefaultLanguage(language.Russian))

	assert.Nil(t, v)
	assert.EqualError(t, err, "default language is not loaded: missing messages for language 'ru'")
}

func TestValidator_Validate_WhenTranslationLanguageInContextArgument_ExpectTranslationLanguageUsed(t *testing.T) {
	v := newValidator(t, validation.Translations(russian.Messages))

	ctx := language.WithContext(context.Background(), language.Russian)
	err := v.Validate(
		ctx,
		validation.String(stringValue(""), it.IsNotBlank()),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, "Значение не должно быть пустым.", violations[0].Message())
	})
}

func TestValidator_Validate_WhenTranslationLanguageInScopedValidator_ExpectTranslationLanguageUsed(t *testing.T) {
	v := newValidator(t, validation.Translations(russian.Messages)).WithLanguage(language.Russian)

	err := v.ValidateString(context.Background(), stringValue(""), it.IsNotBlank())

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, "Значение не должно быть пустым.", violations[0].Message())
	})
}

func TestValidator_Validate_WhenTranslationLanguageInContextOfScopedValidator_ExpectTranslationLanguageUsed(t *testing.T) {
	ctx := language.WithContext(context.Background(), language.Russian)
	v := newValidator(t, validation.Translations(russian.Messages))

	err := v.ValidateString(ctx, stringValue(""), it.IsNotBlank())

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, "Значение не должно быть пустым.", violations[0].Message())
	})
}

func TestValidator_Validate_WhenTranslationLanguageParsedFromAcceptLanguageHeader_ExpectTranslationLanguageUsed(t *testing.T) {
	v := newValidator(t, validation.Translations(russian.Messages))

	matcher := textlanguage.NewMatcher([]language.Tag{language.Russian})
	tag, _ := textlanguage.MatchStrings(matcher, "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	ctx := language.WithContext(context.Background(), tag)
	err := v.Validate(
		ctx,
		validation.String(stringValue(""), it.IsNotBlank()),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, "Значение не должно быть пустым.", violations[0].Message())
	})
}

func TestValidator_Validate_WhenRecursiveValidation_ExpectViolationTranslated(t *testing.T) {
	v := newValidator(
		t,
		validation.DefaultLanguage(language.Russian),
		validation.Translations(russian.Messages),
	)
	values := []mockValidatableString{{value: ""}}

	err := v.ValidateIterable(context.Background(), values, it.IsNotBlank())

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, "Значение не должно быть пустым.", violations[0].Message())
	})
}

func TestValidator_Validate_WhenTranslatableParameter_ExpectParameterTranslated(t *testing.T) {
	validator := newValidator(
		t,
		validation.DefaultLanguage(language.Russian),
		validation.Translations(map[language.Tag]map[string]catalog.Message{
			language.Russian: {
				"The operation is only possible for the {{ role }}.": catalog.String("Операция возможна только для {{ role }}."),
				"administrator role": catalog.String("роли администратора"),
			},
		}),
	)

	v := ""
	err := validator.Validate(
		context.Background(),
		validation.String(
			&v,
			it.IsNotBlank().
				Message(
					"The operation is only possible for the {{ role }}.",
					validation.TemplateParameter{
						Key:              "{{ role }}",
						Value:            "administrator role",
						NeedsTranslation: true,
					},
				),
		),
	)

	assertHasOneViolation(code.NotBlank, "Операция возможна только для роли администратора.")(t, err)
}

func TestValidate_WhenTranslationsLoadedAfterInit_ExpectTranslationsWorking(t *testing.T) {
	err := validator.SetOptions(
		validation.DefaultLanguage(language.Russian),
		validation.Translations(russian.Messages),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer validator.Reset()

	err = validator.ValidateString(context.Background(), stringValue(""), it.IsNotBlank())

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations []validation.Violation) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, "Значение не должно быть пустым.", violations[0].Message())
	})
}
