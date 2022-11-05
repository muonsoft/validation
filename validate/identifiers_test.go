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
		{value: "00000000000000000000000000"},
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

func TestUUID(t *testing.T) {
	tests := []struct {
		value         string
		versions      []byte
		isNonStrict   bool
		isNotNil      bool
		expectedError error
	}{
		{value: "a7a626af-5aeb-11ed-b514-04922656b1d2"},                      // v1
		{value: "a7a626af-5aeb-11ed-b514-04922656b1d2", versions: []byte{1}}, // v1
		{value: "c6437ef1-5b86-3a4e-a071-c2d4ad414e65"},                      // v3
		{value: "c6437ef1-5b86-3a4e-a071-c2d4ad414e65", versions: []byte{3}}, // v3
		{value: "83eab6fd-230b-44fe-b52f-463387bd8788"},                      // v4
		{value: "83eab6fd-230b-44fe-b52f-463387bd8788", versions: []byte{4}}, // v4
		{value: "9b8edca0-90f2-5031-8e5d-3f708834696c"},                      // v5
		{value: "9b8edca0-90f2-5031-8e5d-3f708834696c", versions: []byte{5}}, // v5
		{value: "1ed5aeba-7a62-6c2e-b514-6c6673564acd"},                      // v6
		{value: "1ed5aeba-7a62-6c2e-b514-6c6673564acd", versions: []byte{6}}, // v6
		{value: "018439ff-b24d-75f3-92ca-2e0bf3e90d53"},                      // v7
		{value: "018439ff-b24d-75f3-92ca-2e0bf3e90d53", versions: []byte{7}}, // v7
		{value: "c207f212-bf89-43b5-abd2-e2bfa10ad3aa"},
		{value: "216fff40-98d9-11e3-a5e2-0800200c9a66"},
		{value: "216FFF40-98D9-11E3-A5E2-0800200C9A66"},
		{value: "456daefb-5aa6-41b5-8dbc-068b05a8b201"},
		{value: "456daEFb-5AA6-41B5-8DBC-068B05A8B201"},
		{value: "1eb01932-4c0b-6570-aa34-d179cdf481ae"},
		{value: "334f52e4-abba-180c-baed-116134d32a16"},
		{value: "334f52e4-abba-280c-baed-116134d32a16"},
		{value: "334f52e4-abba-380c-baed-116134d32a16"},
		{value: "334f52e4-abba-480c-baed-116134d32a16"},
		{value: "334f52e4-abba-580c-baed-116134d32a16"},
		{value: "334f52e4-abba-680c-baed-116134d32a16"},
		{value: "334f52e4-abba-780c-baed-116134d32a16"},
		{value: "216fff40-98d9-11e3-a5e2_0800200c9a66", expectedError: validate.ErrInvalidHyphenPlacement},
		{value: "216gff40-98d9-11e3-a5e2-0800200c9a66", expectedError: validate.ErrInvalid},
		{value: "216Gff40-98d9-11e3-a5e2-0800200c9a66", expectedError: validate.ErrInvalid},
		{value: "216fff40-98d9-11e3-a5e-20800200c9a66", expectedError: validate.ErrInvalidHyphenPlacement},
		{value: "216f-ff40-98d9-11e3-a5e2-0800200c9a66", expectedError: validate.ErrTooLong},
		{value: "216fff40-98d9-11e3-a5e2-0800-200c9a66", expectedError: validate.ErrTooLong},
		{value: "216fff40-98d9-11e3-a5e2-0800200c-9a66", expectedError: validate.ErrTooLong},
		{value: "216fff40-98d9-11e3-a5e20800200c9a66", expectedError: validate.ErrTooShort},
		{value: "216fff4098d911e3a5e20800200c9a66", expectedError: validate.ErrTooShort},
		{value: "216fff40-98d9-11e3-a5e2-0800200c9a6", expectedError: validate.ErrTooShort},
		{value: "216fff40-98d9-11e3-a5e2-0800200c9a666", expectedError: validate.ErrTooLong},
		{value: "216fff40-98d9-01e3-a5e2-0800200c9a66", expectedError: validate.ErrInvalidVersion},
		{value: "216fff40-98d9-81e3-a5e2-0800200c9a66", expectedError: validate.ErrInvalidVersion},
		{value: "216fff40-98d9-91e3-a5e2-0800200c9a66", expectedError: validate.ErrInvalidVersion},
		{value: "216fff40-98d9-a1e3-a5e2-0800200c9a66", expectedError: validate.ErrInvalidVersion},
		{value: "216fff40-98d9-b1e3-a5e2-0800200c9a66", expectedError: validate.ErrInvalidVersion},
		{value: "216fff40-98d9-c1e3-a5e2-0800200c9a66", expectedError: validate.ErrInvalidVersion},
		{value: "216fff40-98d9-d1e3-a5e2-0800200c9a66", expectedError: validate.ErrInvalidVersion},
		{value: "216fff40-98d9-e1e3-a5e2-0800200c9a66", expectedError: validate.ErrInvalidVersion},
		{value: "216fff40-98d9-f1e3-a5e2-0800200c9a66", expectedError: validate.ErrInvalidVersion},
		{value: "216fff40-98d9-11e3-a5e2-0800200c9a66", versions: []byte{2, 3, 4, 5}, expectedError: validate.ErrInvalidVersion},
		{value: "216fff40-98d9-21e3-a5e2-0800200c9a66", versions: []byte{1, 3, 4, 5}, expectedError: validate.ErrInvalidVersion},
		{value: "{216fff40-98d9-11e3-a5e2-0800200c9a66}", expectedError: validate.ErrTooLong},
		{value: "[216fff40-98d9-11e3-a5e2-0800200c9a66]", expectedError: validate.ErrTooLong},

		{value: "216fff4098d911e3a5e20800200c9a66", isNonStrict: true},       // No dashes at all
		{value: "{216fff4098d911e3a5e20800200c9a66}", isNonStrict: true},     // No dashes, wrapped with curly braces
		{value: "{216fff40-98d9-11e3-a5e2-0800200c9a66}", isNonStrict: true}, // Wrapped with curly braces
		{value: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", isNonStrict: true},
		{value: "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}", isNonStrict: true},
		{value: "6ba7b8109dad11d180b400c04fd430c8", isNonStrict: true},
		{value: "{6ba7b8109dad11d180b400c04fd430c8}", isNonStrict: true},
		{value: "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8", isNonStrict: true},
		{value: "urn:uuid:6ba7b8109dad11d180b400c04fd430c8", isNonStrict: true},

		{value: "f7e6a259-7c30-4b4e-a9e1-746c5b948566", isNotNil: true},
		{value: "00000000-0000-0000-0000-000000000000"},
		{value: "00000000-0000-0000-0000-000000000000", isNotNil: true, expectedError: validate.ErrIsNil},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			var options []func(o *validate.UUIDOptions)
			if test.versions != nil {
				options = append(options, validate.AllowUUIDVersions(test.versions...))
			}
			if test.isNonStrict {
				options = append(options, validate.AllowNonCanonicalUUIDFormats())
			}
			if test.isNotNil {
				options = append(options, validate.DenyNilUUID())
			}

			err := validate.UUID(test.value, options...)

			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, test.expectedError)
			}
		})
	}
}
