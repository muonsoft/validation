package validate_test

import (
	"testing"

	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/assert"
)

func TestULID(t *testing.T) {
	tests := []struct {
		value         string
		expectedError error
	}{
		{value: "01ARZ3NDEKTSV4RRFFQ69G5FAV"},
		{value: "", expectedError: validate.ErrTooShort},
		{value: "01ARZ3NDEKTSV4RRFFQ69G5FA", expectedError: validate.ErrTooShort},
		{value: "01ARZ3NDEKTSV4RRFFQ69G5FAVA", expectedError: validate.ErrTooLong},
		{value: "01ARZ3NDEKTSV4RRFFQ69G5FAO", expectedError: validate.ErrInvalidCharacters},
		{value: "71ARZ3NDEKTSV4RRFFQ69G5FAV"},
		{value: "81ARZ3NDEKTSV4RRFFQ69G5FAV", expectedError: validate.ErrTooLarge},
		{value: "Z1ARZ3NDEKTSV4RRFFQ69G5FAV", expectedError: validate.ErrTooLarge},
		{value: "not-even-ulid-like", expectedError: validate.ErrTooShort},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			err := validate.ULID(test.value)

			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, test.expectedError)
			}
		})
	}
}
