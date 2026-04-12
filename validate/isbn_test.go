package validate_test

import (
	"errors"
	"testing"

	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/assert"
)

func TestISBN(t *testing.T) {
	valid10 := []string{
		"2723442284",
		"2723442276",
		"2723455041",
		"2070546810",
		"2711858839",
		"2756406767",
		"2870971648",
		"226623854X",
		"2851806424",
		"0321812700",
		"0-45122-5244",
		"0-4712-92311",
		"0-9752298-0-X",
	}
	valid13 := []string{
		"978-2723442282",
		"978-2723442275",
		"978-2723455046",
		"978-2070546817",
		"978-2711858835",
		"978-2756406763",
		"978-2870971642",
		"978-2266238540",
		"978-2851806420",
		"978-0321812704",
		"978-0451225245",
		"978-0471292319",
	}

	for _, v := range append(valid10, valid13...) {
		assert.NoError(t, validate.ISBN(v), v)
	}
	for _, v := range valid10 {
		assert.NoError(t, validate.ISBN(v, validate.ISBNOnly10()), v)
	}
	for _, v := range valid13 {
		assert.NoError(t, validate.ISBN(v, validate.ISBNOnly13()), v)
	}

	t.Run("invalid ISBN-10", func(t *testing.T) {
		cases := []struct {
			value string
			want  error
		}{
			{"27234422841", validate.ErrISBNTooLong},
			{"272344228", validate.ErrISBNTooShort},
			{"0-4712-9231", validate.ErrISBNTooShort},
			{"1234567890", validate.ErrISBNChecksumFailed},
			{"0987656789", validate.ErrISBNChecksumFailed},
			{"7-35622-5444", validate.ErrISBNChecksumFailed},
			{"0-4X19-92611", validate.ErrISBNChecksumFailed},
			{"0_45122_5244", validate.ErrISBNInvalidCharacters},
			{"2870#971#648", validate.ErrISBNInvalidCharacters},
			{"0-9752298-0-x", validate.ErrISBNInvalidCharacters},
			{"1A34567890", validate.ErrISBNInvalidCharacters},
		}
		for _, tc := range cases {
			err := validate.ISBN(tc.value, validate.ISBNOnly10())
			assert.ErrorIs(t, err, tc.want, tc.value)
		}
	})

	t.Run("invalid ISBN-13", func(t *testing.T) {
		cases := []struct {
			value string
			want  error
		}{
			{"978-27234422821", validate.ErrISBNTooLong},
			{"978-272344228", validate.ErrISBNTooShort},
			{"978-2723442-82", validate.ErrISBNTooShort},
			{"978-2723442281", validate.ErrISBNChecksumFailed},
			{"978-0321513774", validate.ErrISBNChecksumFailed},
			{"979-0431225385", validate.ErrISBNChecksumFailed},
			{"980-0474292319", validate.ErrISBNChecksumFailed},
			{"0-4X19-92619812", validate.ErrISBNInvalidCharacters},
			{"978_2723442282", validate.ErrISBNInvalidCharacters},
			{"978#2723442282", validate.ErrISBNInvalidCharacters},
			{"978-272C442282", validate.ErrISBNInvalidCharacters},
		}
		for _, tc := range cases {
			err := validate.ISBN(tc.value, validate.ISBNOnly13())
			assert.ErrorIs(t, err, tc.want, tc.value)
		}
	})

	t.Run("explicit type length mismatch", func(t *testing.T) {
		assert.ErrorIs(t, validate.ISBN("978-2723442282", validate.ISBNOnly10()), validate.ErrISBNTooLong)
		assert.ErrorIs(t, validate.ISBN("2723442284", validate.ISBNOnly13()), validate.ErrISBNTooShort)
	})

	t.Run("any mode type not recognized", func(t *testing.T) {
		// 11 digits: too long for ISBN-10, too short for ISBN-13
		assert.ErrorIs(t, validate.ISBN("978272344228"), validate.ErrISBNTypeNotRecognized)
		assert.ErrorIs(t, validate.ISBN("97827234422821"), validate.ErrISBNTooLong)
	})

	t.Run("empty is valid", func(t *testing.T) {
		assert.NoError(t, validate.ISBN(""))
	})
}

func TestISBN_wrapsErrors(t *testing.T) {
	err := validate.ISBN("bad", validate.ISBNOnly10())
	assert.True(t, errors.Is(err, validate.ErrISBNInvalidCharacters))
}
