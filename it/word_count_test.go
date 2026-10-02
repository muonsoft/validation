package it_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/muonsoft/language"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message/translations/russian"
	"github.com/stretchr/testify/require"
)

func TestWordCountSegmentation(t *testing.T) {
	v, err := validation.NewValidator()
	require.NoError(t, err)
	tests := []struct {
		value string
		count int
	}{
		{"Hello,world!", 2},
		{"Привет, мир!", 2},
		{"你好世界", 1},
		{"日本語", 1},
		{"ภาษาไทย", 1},
		{"cafe\u0301 déjà vu", 3},
		{"नमस्ते दुनिया", 2},
		{"\u0301\u0308", 0},
		{"don't l’amour O'Reilly", 3},
		{"'hello' ’world’", 2},
		{"a''b a’’b", 4},
		{"well-known foo_bar", 4},
		{"one—two–three", 3},
		{"123 ٤٥٦ abc123", 3},
		{"3.14", 2},
		{"one\u00a0two\u2003three\t四\n五", 5},
		{" \t\r\n\u00a0", 0},
		{"... !!! — _ ' ’", 0},
		{"👩‍💻 😀", 0},
		{"hi😀there", 2},
		{"a\xffb", 2},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			require.NoError(t, v.Validate(context.Background(), validation.String(tt.value, it.HasWordCountBetween(tt.count, tt.count))))
			err := v.Validate(context.Background(), validation.String(tt.value, it.HasMinWordCount(tt.count+1)))
			require.ErrorIs(t, err, validation.ErrTooFewWords)
			if tt.count > 0 {
				err = v.Validate(context.Background(), validation.String(tt.value, it.HasMaxWordCount(tt.count-1)))
				require.ErrorIs(t, err, validation.ErrTooManyWords)
			}
		})
	}
}

func TestWordCountInvalidBounds(t *testing.T) {
	v, err := validation.NewValidator()
	require.NoError(t, err)
	for _, c := range []it.WordCountConstraint{it.HasMinWordCount(-1), it.HasMaxWordCount(-1), it.HasWordCountBetween(2, 1)} {
		err := c.Validate(context.Background(), v, "one")
		require.EqualError(t, err, "validate by WordCountConstraint: bounds must be non-negative and minimum must not exceed maximum")
	}
}

func TestWordCountGroupsAndCopies(t *testing.T) {
	v, err := validation.NewValidator()
	require.NoError(t, err)
	base := it.HasMinWordCount(2).WhenGroups("content")
	customErr := errors.New("custom word count")
	derived := base.WithMinError(customErr).WithMinMessage("custom")
	require.NoError(t, base.Validate(context.Background(), v, "one"))
	require.ErrorIs(t, base.Validate(context.Background(), v.WithGroups("content"), "one"), validation.ErrTooFewWords)
	require.ErrorIs(t, derived.Validate(context.Background(), v.WithGroups("content"), "one"), customErr)
	require.NoError(t, v.Validate(context.Background(), validation.This("one two", it.HasWordCountBetween(2, 2))))
	require.ErrorIs(t, v.Validate(context.Background(), validation.Each([]string{"one", "one two"}, it.HasMaxWordCount(1))), validation.ErrTooManyWords)
}

func TestWordCountTranslations(t *testing.T) {
	v, err := validation.NewValidator(validation.Translations(russian.Messages))
	require.NoError(t, err)
	for _, limit := range []int{0, 1, 2, 5, 11, 21, 22, 25} {
		for _, minimum := range []bool{false, true} {
			t.Run(fmt.Sprintf("%d/min=%t", limit, minimum), func(t *testing.T) {
				c, value, en, ru := it.HasMaxWordCount(limit), strings.Repeat("word ", limit+1), "at most", "не более"
				sentinel := validation.ErrTooManyWords
				if minimum {
					c, value, en, ru = it.HasMinWordCount(limit), "...", "at least", "не менее"
					sentinel = validation.ErrTooFewWords
					if limit == 0 {
						require.NoError(t, c.Validate(context.Background(), v, value))
						return
					}
				}
				englishWord, russianWord := "words", "слов"
				if limit == 1 {
					englishWord = "word"
				}
				if limit%10 == 1 && limit%100 != 11 {
					russianWord = "слова"
				}
				for _, translation := range []struct {
					tag     language.Tag
					message string
				}{
					{language.English, fmt.Sprintf("This value should contain %s %d %s.", en, limit, englishWord)},
					{language.Russian, fmt.Sprintf("Значение должно содержать %s %d %s.", ru, limit, russianWord)},
				} {
					err := c.Validate(context.Background(), v.WithLanguage(translation.tag), value)
					require.ErrorIs(t, err, sentinel)
					require.EqualError(t, err, "violation: \""+translation.message+"\"")
				}
			})
		}
	}
}
