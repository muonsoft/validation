package test

import (
	"context"
	"strconv"
	"testing"

	languagepkg "github.com/muonsoft/language"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message/translations/russian"
	"github.com/muonsoft/validation/validationtest"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestValidator_Validate_WhenRussianIsDefaultLanguage_ExpectViolationTranslated(t *testing.T) {
	validator, err := validation.NewValidator(
		validation.DefaultLanguage(language.Russian),
		validation.Translations(russian.Messages),
	)
	if err != nil {
		t.Fatal(err)
	}

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
			err = validator.ValidateCountable(10, it.HasMaxCount(test.maxCount))

			validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
				t.Helper()
				return assert.Len(t, violations, 1) &&
					assert.Equal(t, test.expectedMessage, violations[0].GetMessage())
			})
		})
	}
}

func TestValidator_Validate_WhenRussianIsPassedViaOption_ExpectViolationTranslated(t *testing.T) {
	validator, err := validation.NewValidator(
		validation.Translations(russian.Messages),
	)
	if err != nil {
		t.Fatal(err)
	}

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
			err = validator.ValidateCountable(
				10,
				it.HasMaxCount(test.maxCount),
				validation.Language(language.Russian),
			)

			validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
				t.Helper()
				return assert.Len(t, violations, 1) &&
					assert.Equal(t, test.expectedMessage, violations[0].GetMessage())
			})
		})
	}
}

func TestValidator_Validate_WhenCustomDefaultLanguageAndUndefinedTranslationLanguage_ExpectDefaultLanguageUsed(t *testing.T) {
	validator, err := validation.NewValidator(
		validation.DefaultLanguage(language.Russian),
		validation.Translations(russian.Messages),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = validator.ValidateString(
		stringValue(""),
		it.IsNotBlank(),
		validation.Language(language.Afrikaans),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, "Значение не должно быть пустым.", violations[0].GetMessage())
	})
}

func TestValidator_Validate_WhenDefaultLanguageIsNotLoaded_ExpectError(t *testing.T) {
	validator, err := validation.NewValidator(validation.DefaultLanguage(language.Russian))

	assert.Nil(t, validator)
	assert.EqualError(t, err, "default language is not loaded: missing messages for language 'ru'")
}

func TestValidate_WhenTranslationsLoadedAfterInit_ExpectTranslationsWorking(t *testing.T) {
	err := validation.OverrideDefaults(
		validation.DefaultLanguage(language.Russian),
		validation.Translations(russian.Messages),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer validation.ResetDefaults()

	err = validation.ValidateString(stringValue(""), it.IsNotBlank())

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, "Значение не должно быть пустым.", violations[0].GetMessage())
	})
}

func TestValidator_Validate_WhenTranslationLanguageInContext_ExpectTranslationLanguageUsed(t *testing.T) {
	validator, err := validation.NewValidator(validation.Translations(russian.Messages))
	if err != nil {
		t.Fatal(err)
	}

	ctx := languagepkg.WithContext(context.Background(), language.Russian)
	err = validator.ValidateString(
		stringValue(""),
		it.IsNotBlank(),
		validation.Context(ctx),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, "Значение не должно быть пустым.", violations[0].GetMessage())
	})
}

func TestValidator_Validate_WhenTranslationLanguageParsedFromAcceptLanguageHeader_ExpectTranslationLanguageUsed(t *testing.T) {
	validator, err := validation.NewValidator(validation.Translations(russian.Messages))
	if err != nil {
		t.Fatal(err)
	}

	matcher := language.NewMatcher([]language.Tag{language.Russian})
	tag, _ := language.MatchStrings(matcher, "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	ctx := languagepkg.WithContext(context.Background(), tag)
	err = validator.ValidateString(
		stringValue(""),
		it.IsNotBlank(),
		validation.Context(ctx),
	)

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, "Значение не должно быть пустым.", violations[0].GetMessage())
	})
}

func TestValidator_Validate_WhenRecursiveValidation_ExpectViolationTranslated(t *testing.T) {
	validator, err := validation.NewValidator(
		validation.DefaultLanguage(language.Russian),
		validation.Translations(russian.Messages),
	)
	if err != nil {
		t.Fatal(err)
	}
	values := []mockValidatableString{{value: ""}}

	err = validator.ValidateIterable(values, it.IsNotBlank())

	validationtest.AssertIsViolationList(t, err, func(t *testing.T, violations validation.ViolationList) bool {
		t.Helper()
		return assert.Len(t, violations, 1) &&
			assert.Equal(t, "Значение не должно быть пустым.", violations[0].GetMessage())
	})
}