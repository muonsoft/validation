package test

//
// import (
// 	"context"
// 	"testing"
//
// 	"github.com/muonsoft/validation"
// 	"github.com/muonsoft/validation/code"
// 	"github.com/muonsoft/validation/it"
// 	"github.com/muonsoft/validation/validationtest"
// 	"github.com/muonsoft/validation/validator"
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestValidateEach_WhenSliceOfStrings_ExpectViolationOnEachElement(t *testing.T) {
// 	strings := []string{"", ""}
//
// 	err := validator.Validate(context.Background(), validation.Each(strings, it.IsNotBlank()))
//
// 	validationtest.Assert(t, err).IsViolationList().WithAttributes(
// 		validationtest.ViolationAttributes{Code: code.NotBlank, PropertyPath: "[0]"},
// 		validationtest.ViolationAttributes{Code: code.NotBlank, PropertyPath: "[1]"},
// 	)
// }
//
// func TestValidateEach_WhenMapOfStrings_ExpectViolationOnEachElement(t *testing.T) {
// 	strings := map[string]string{"key1": "", "key2": ""}
//
// 	err := validator.Validate(context.Background(), validation.Each(strings, it.IsNotBlank()))
//
// 	validationtest.Assert(t, err).IsViolationList().Assert(func(tb testing.TB, violations []validation.Violation) {
// 		tb.Helper()
// 		if assert.Len(tb, violations, 2) {
// 			assert.Equal(tb, code.NotBlank, violations[0].Code())
// 			assert.Contains(tb, violations[0].PropertyPath().String(), "key")
// 			assert.Equal(tb, code.NotBlank, violations[1].Code())
// 			assert.Contains(tb, violations[1].PropertyPath().String(), "key")
// 		}
// 	})
// }
//
// func TestValidateEachString_WhenSliceOfStrings_ExpectViolationOnEachElement(t *testing.T) {
// 	strings := []string{"", ""}
//
// 	err := validator.Validate(context.Background(), validation.EachString(strings, it.IsNotBlank()))
//
// 	validationtest.Assert(t, err).IsViolationList().WithAttributes(
// 		validationtest.ViolationAttributes{Code: code.NotBlank, PropertyPath: "[0]"},
// 		validationtest.ViolationAttributes{Code: code.NotBlank, PropertyPath: "[1]"},
// 	)
// }
