package test

import (
	"fmt"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/message/translations/english"
	"github.com/muonsoft/validation/message/translations/russian"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

func TestAllMessagesTranslated(t *testing.T) {
	allErrors := []*validation.Error{
		validation.ErrInvalidDate,
		validation.ErrInvalidDateTime,
		validation.ErrInvalidEAN13,
		validation.ErrInvalidEAN8,
		validation.ErrInvalidEmail,
		validation.ErrInvalidHostname,
		validation.ErrInvalidIP,
		validation.ErrInvalidJSON,
		validation.ErrInvalidTime,
		validation.ErrInvalidULID,
		validation.ErrInvalidUPCA,
		validation.ErrInvalidUPCE,
		validation.ErrInvalidURL,
		validation.ErrInvalidUUID,
		validation.ErrIsBlank,
		validation.ErrIsEqual,
		validation.ErrIsNil,
		validation.ErrNoSuchChoice,
		validation.ErrNotBlank,
		validation.ErrNotDivisible,
		validation.ErrNotDivisibleCount,
		validation.ErrNotEqual,
		validation.ErrNotExactCount,
		validation.ErrNotExactLength,
		validation.ErrNotFalse,
		validation.ErrNotInRange,
		validation.ErrNotInteger,
		validation.ErrNotNegative,
		validation.ErrNotNegativeOrZero,
		validation.ErrNotNil,
		validation.ErrNotNumeric,
		validation.ErrNotPositive,
		validation.ErrNotPositiveOrZero,
		validation.ErrNotTrue,
		validation.ErrNotUnique,
		validation.ErrNotValid,
		validation.ErrProhibitedIP,
		validation.ErrProhibitedURL,
		validation.ErrTooEarly,
		validation.ErrTooEarlyOrEqual,
		validation.ErrTooFewElements,
		validation.ErrTooHigh,
		validation.ErrTooHighOrEqual,
		validation.ErrTooLate,
		validation.ErrTooLateOrEqual,
		validation.ErrTooLong,
		validation.ErrTooLow,
		validation.ErrTooLowOrEqual,
		validation.ErrTooManyElements,
		validation.ErrTooShort,
	}
	allDictionaries := []map[language.Tag]map[string]catalog.Message{
		english.Messages,
		russian.Messages,
	}

	for _, dictionary := range allDictionaries {
		for languageTag, messages := range dictionary {
			for _, err := range allErrors {
				_, exist := messages[err.Message()]
				if !exist {
					assert.Fail(t, fmt.Sprintf(
						`missing translation for message "%s" and language "%s"`,
						err.Message(),
						languageTag.String(),
					))
				}
			}
		}
	}
}
