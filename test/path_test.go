package test

import (
	"testing"

	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/muonsoft/validation/validator"
)

func TestValidate_AtProperty_WhenGivenProperty_ExpectViolationWithProperty(t *testing.T) {
	err := validator.AtProperty("property").ValidateString(stringValue(""), it.IsNotBlank())

	assertHasOneViolation(code.NotBlank, message.NotBlank, "property")(t, err)
}

func TestValidate_AtIndex_WhenGivenIndex_ExpectViolationWithIndex(t *testing.T) {
	err := validator.AtIndex(1).ValidateString(stringValue(""), it.IsNotBlank())

	assertHasOneViolation(code.NotBlank, message.NotBlank, "[1]")(t, err)
}
