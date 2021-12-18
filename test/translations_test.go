package test

import (
	"context"
	"strconv"
	"testing"

	"github.com/muonsoft/language"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
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
			err := v.Validate(context.Background(), validation.Countable(10, it.HasMaxCount(test.maxCount)))

			validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithMessage(test.expectedMessage)
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

			validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithMessage(test.expectedMessage)
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
		validation.String("", it.IsNotBlank()),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithMessage("Значение не должно быть пустым.")
}

func TestValidator_Validate_WhenDefaultLanguageIsNotLoaded_ExpectError(t *testing.T) {
	v, err := validation.NewValidator(validation.DefaultLanguage(language.Russian))

	assert.Nil(t, v)
	assert.EqualError(t, err, `failed to set up default translator: default language is not loaded: missing messages for language "ru"`)
}

func TestValidator_Validate_WhenTranslationLanguageInContextArgument_ExpectTranslationLanguageUsed(t *testing.T) {
	v := newValidator(t, validation.Translations(russian.Messages))

	ctx := language.WithContext(context.Background(), language.Russian)
	err := v.Validate(
		ctx,
		validation.String("", it.IsNotBlank()),
	)

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithMessage("Значение не должно быть пустым.")
}

func TestValidator_Validate_WhenTranslationLanguageInScopedValidator_ExpectTranslationLanguageUsed(t *testing.T) {
	v := newValidator(t, validation.Translations(russian.Messages)).WithLanguage(language.Russian)

	err := v.Validate(context.Background(), validation.String("", it.IsNotBlank()))

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithMessage("Значение не должно быть пустым.")
}

func TestValidator_Validate_WhenTranslationLanguageInContextOfScopedValidator_ExpectTranslationLanguageUsed(t *testing.T) {
	ctx := language.WithContext(context.Background(), language.Russian)
	v := newValidator(t, validation.Translations(russian.Messages))

	err := v.Validate(ctx, validation.String("", it.IsNotBlank()))

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithMessage("Значение не должно быть пустым.")
}

func TestValidator_Validate_WhenTranslationLanguageParsedFromAcceptLanguageHeader_ExpectTranslationLanguageUsed(t *testing.T) {
	v := newValidator(t, validation.Translations(russian.Messages))

	matcher := textlanguage.NewMatcher([]language.Tag{language.Russian})
	tag, _ := textlanguage.MatchStrings(matcher, "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	ctx := language.WithContext(context.Background(), tag)
	err := v.Validate(ctx, validation.String("", it.IsNotBlank()))

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithMessage("Значение не должно быть пустым.")
}

func TestValidator_Validate_WhenRecursiveValidation_ExpectViolationTranslated(t *testing.T) {
	v := newValidator(
		t,
		validation.DefaultLanguage(language.Russian),
		validation.Translations(russian.Messages),
	)
	values := []mockValidatableString{{value: ""}}

	err := v.Validate(context.Background(), validation.Iterable(values, it.IsNotBlank()))

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithMessage("Значение не должно быть пустым.")
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
			v,
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
	err := validator.SetUp(
		validation.DefaultLanguage(language.Russian),
		validation.Translations(russian.Messages),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer validator.SetUp()

	err = validator.Validate(context.Background(), validation.String("", it.IsNotBlank()))

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithMessage("Значение не должно быть пустым.")
}

func TestValidate_WhenTranslatorIsOverridden_ExpectTranslationsByOverriddenTranslator(t *testing.T) {
	translator := mockTranslator{translate: func(tag textlanguage.Tag, msg string, pluralCount int) string {
		if msg == message.Templates[code.NotBlank] {
			return "expected message"
		}
		return "unexpected message"
	}}
	validator := newValidator(t, validation.SetTranslator(translator))

	err := validator.ValidateString(context.Background(), "", it.IsNotBlank())

	validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithMessage("expected message")
}

func TestValidate_WhenTranslatorIsOverriddenAndTranslationsPasses_ExpectError(t *testing.T) {
	translator := mockTranslator{}

	validator, err := validation.NewValidator(
		validation.SetTranslator(translator),
		validation.Translations(russian.Messages),
	)

	assert.Nil(t, validator)
	assert.EqualError(t, err, "translation options denied when using custom translator")
}
