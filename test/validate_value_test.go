package test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/muonsoft/validation/validator"
)

func TestValidateValue_WhenValueOfType_ExpectValueValidated(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{"bool", false},
		{"int8", int8(0)},
		{"uint8", uint8(0)},
		{"float32", float32(0)},
		{"string", ""},
		{"bool pointer", boolValue(false)},
		{"int64 pointer", intValue(0)},
		{"uint64 pointer", uintValue(0)},
		{"float64 pointer", floatValue(0)},
		{"string pointer", stringValue("")},
		{"bool nil", nilBool},
		{"int64 nil", nilInt},
		{"uint64 nil", nilUint},
		{"float64 nil", nilFloat},
		{"string nil", nilString},
		{"time nil", nilTime},
		{"empty time", emptyTime},
		{"empty array", emptyArray},
		{"empty slice", emptySlice},
		{"empty map", emptyMap},
		{"empty array pointer", &emptyArray},
		{"empty slice pointer", &emptySlice},
		{"empty map pointer", &emptyMap},
		{"empty time pointer", &emptyTime},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validator.Validate(
				context.Background(),
				validation.Value(test.value, validation.PropertyName("property"), it.IsNotBlank()),
			)

			assertHasOneViolationAtPath(code.NotBlank, message.NotBlank, "property")(t, err)
		})
	}
}
