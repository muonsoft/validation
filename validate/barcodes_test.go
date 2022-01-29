package validate_test

import (
	"testing"

	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/assert"
)

func TestEAN8(t *testing.T) {
	tests := []struct {
		code          string
		expectedError string
	}{
		{code: "42345671"},
		{code: "47195127"},
		{code: "96385074"},
		{code: "00000000", expectedError: "contains only zeros"},
		{code: "42345670", expectedError: "invalid checksum"},
		{code: "12345671", expectedError: "invalid checksum"},
		{code: "423456712", expectedError: "unexpected length"},
		{code: "4234.671", expectedError: `contains non-digit: '.'`},
	}
	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			err := validate.EAN8(test.code)

			if test.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.expectedError)
			}
		})
	}
}

func TestEAN13(t *testing.T) {
	tests := []struct {
		code          string
		expectedError string
	}{
		{code: "4719512002889"},
		{code: "9782868890061"},
		{code: "4006381333931"},
		{code: "0000000000000", expectedError: "contains only zeros"},
		{code: "2266111566", expectedError: "unexpected length"},
		{code: "A782868890061", expectedError: `contains non-digit: 'A'`},
		{code: "4006381333932", expectedError: "invalid checksum"},
	}
	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			err := validate.EAN13(test.code)

			if test.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.expectedError)
			}
		})
	}
}

func TestUPCA(t *testing.T) {
	tests := []struct {
		code          string
		expectedError string
	}{
		{code: "614141000036"},
		{code: "123456789999"},
		{code: "000000000000", expectedError: "contains only zeros"},
		{code: "61414100003", expectedError: "unexpected length"},
		{code: "A14141000036", expectedError: `contains non-digit: 'A'`},
		{code: "614141000037", expectedError: "invalid checksum"},
	}
	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			err := validate.UPCA(test.code)

			if test.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.expectedError)
			}
		})
	}
}

func TestUPCE(t *testing.T) {
	tests := []struct {
		code          string
		expectedError string
	}{
		{code: "123456"}, // 6-digit is always valid
		// valid 7-digit codes composed at  https://www.morovia.com/education/utility/upc-ean.asp
		{code: "1234505"}, // XXNNN0 pattern
		{code: "1234514"}, // XXNNN1 pattern
		{code: "1234523"}, // XXNNN2 pattern
		{code: "1234531"}, // XXXNN3 pattern
		{code: "1234543"}, // XXXXN4 pattern
		{code: "1234558"}, // XXXXX5 pattern
		// valid 8-digit codes same as previous
		{code: "01234505"}, // XXNNN0 pattern
		{code: "01234514"}, // XXNNN1 pattern
		{code: "01234523"}, // XXNNN2 pattern
		{code: "01234531"}, // XXXNN3 pattern
		{code: "01234543"}, // XXXXN4 pattern
		{code: "01234558"}, // XXXXX5 pattern
		{code: "01234565"},
		{code: "02345673"},
		{code: "000000", expectedError: "contains only zeros"},
		{code: "0000000", expectedError: "contains only zeros"},
		{code: "00000000", expectedError: "contains only zeros"},
		{code: "11234505", expectedError: "invalid"}, // 1 is restricted
		{code: "01234501", expectedError: "invalid checksum"},
		{code: "023456731", expectedError: "unexpected length"},
		{code: "A2345673", expectedError: `contains non-digit: 'A'`},
		{code: "12345", expectedError: "unexpected length"},
	}
	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			err := validate.UPCE(test.code)

			if test.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.expectedError)
			}
		})
	}
}

func BenchmarkEAN8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validate.EAN8("42345671")
	}
}
