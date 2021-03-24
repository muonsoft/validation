package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
)

var numberConstraintTestCases = mergeTestCases(
	isEqualToIntegerTestCases,
	isEqualToFloatTestCases,
	isNotEqualToIntegerTestCases,
	isNotEqualToFloatTestCases,
	isLessThanIntegerTestCases,
	isLessThanFloatTestCases,
	isLessThanOrEqualIntegerTestCases,
	isLessThanOrEqualFloatTestCases,
	isGreaterThanIntegerTestCases,
	isGreaterThanFloatTestCases,
	isGreaterThanOrEqualIntegerTestCases,
	isGreaterThanOrEqualFloatTestCases,
	isPositiveTestCases,
	isPositiveOrZeroTestCases,
	isNegativeTestCases,
	isNegativeOrZeroTestCases,
)

var stringConstraintTestCases = mergeTestCases(
	isEqualToStringTestCases,
	isNotEqualToStringTestCases,
)

var isEqualToIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsEqualToInteger passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsEqualToInteger(1)},
		assert:          assertNoError,
	},
	{
		name:            "IsEqualToInteger violation on not equal int",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		options:         []validation.Option{it.IsEqualToInteger(1)},
		assert:          assertHasOneViolation(code.Equal, "This value should be equal to 1.", ""),
	},
	{
		name:            "IsEqualToInteger passes on equal int",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsEqualToInteger(1)},
		assert:          assertNoError,
	},
	{
		name:            "IsEqualToInteger violation on not equal float",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(0.99),
		options:         []validation.Option{it.IsEqualToInteger(1)},
		assert:          assertHasOneViolation(code.Equal, "This value should be equal to 1.", ""),
	},
	{
		name:            "IsEqualToInteger violation with custom message",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		options: []validation.Option{
			it.IsEqualToInteger(1).Message(`Unexpected value "{{ value }}", expected value is "{{ comparedValue }}".`),
		},
		assert: assertHasOneViolation(code.Equal, `Unexpected value "0", expected value is "1".`, ""),
	},
	{
		name:            "IsEqualToInteger passes when condition is false",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		options:         []validation.Option{it.IsEqualToInteger(1).When(false)},
		assert:          assertNoError,
	},
	{
		name:            "IsEqualToInteger violation when condition is true",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		options:         []validation.Option{it.IsEqualToInteger(1).When(true)},
		assert:          assertHasOneViolation(code.Equal, "This value should be equal to 1.", ""),
	},
}

var isEqualToFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsEqualToFloat passes on nil",
		isApplicableFor: specificValueTypes(floatType),
		options:         []validation.Option{it.IsEqualToFloat(1.5)},
		assert:          assertNoError,
	},
	{
		name:            "IsEqualToFloat violation on not equal float",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(0.5),
		options:         []validation.Option{it.IsEqualToFloat(1.5)},
		assert:          assertHasOneViolation(code.Equal, "This value should be equal to 1.5.", ""),
	},
	{
		name:            "IsEqualToFloat passes on equal float and int",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsEqualToFloat(1)},
		assert:          assertNoError,
	},
}

var isNotEqualToIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNotEqualToInteger passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsNotEqualToInteger(1)},
		assert:          assertNoError,
	},
	{
		name:            "IsNotEqualToInteger violation on equal int or float",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsNotEqualToInteger(1)},
		assert:          assertHasOneViolation(code.NotEqual, "This value should not be equal to 1.", ""),
	},
	{
		name:            "IsNotEqualToInteger passes on not equal float and int",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		options:         []validation.Option{it.IsNotEqualToInteger(1)},
		assert:          assertNoError,
	},
}

var isNotEqualToFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNotEqualToFloat passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsNotEqualToFloat(1)},
		assert:          assertNoError,
	},
	{
		name:            "IsNotEqualToFloat violation on equal int or float",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsNotEqualToFloat(1)},
		assert:          assertHasOneViolation(code.NotEqual, "This value should not be equal to 1.", ""),
	},
	{
		name:            "IsNotEqualToFloat passes on not equal float and int",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		options:         []validation.Option{it.IsNotEqualToFloat(1)},
		assert:          assertNoError,
	},
}

var isLessThanIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLessThanInteger passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsLessThanInteger(1)},
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanInteger violation on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		options:         []validation.Option{it.IsLessThanInteger(1)},
		assert:          assertHasOneViolation(code.TooHigh, "This value should be less than 1.", ""),
	},
	{
		name:            "IsLessThanInteger violation on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsLessThanInteger(1)},
		assert:          assertHasOneViolation(code.TooHigh, "This value should be less than 1.", ""),
	},
	{
		name:            "IsLessThanInteger passes on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		options:         []validation.Option{it.IsLessThanInteger(1)},
		assert:          assertNoError,
	},
}

var isLessThanFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLessThanFloat passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsLessThanFloat(1)},
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanFloat violation on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		options:         []validation.Option{it.IsLessThanFloat(1)},
		assert:          assertHasOneViolation(code.TooHigh, "This value should be less than 1.", ""),
	},
	{
		name:            "IsLessThanFloat violation on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsLessThanFloat(1)},
		assert:          assertHasOneViolation(code.TooHigh, "This value should be less than 1.", ""),
	},
	{
		name:            "IsLessThanFloat passes on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		options:         []validation.Option{it.IsLessThanFloat(1)},
		assert:          assertNoError,
	},
}

var isLessThanOrEqualIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLessThanOrEqualInteger passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsLessThanOrEqualInteger(1)},
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanOrEqualInteger violation on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		options:         []validation.Option{it.IsLessThanOrEqualInteger(1)},
		assert:          assertHasOneViolation(code.TooHighOrEqual, "This value should be less than or equal to 1.", ""),
	},
	{
		name:            "IsLessThanOrEqualInteger passes on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsLessThanOrEqualInteger(1)},
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanOrEqualInteger passes on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		options:         []validation.Option{it.IsLessThanOrEqualInteger(1)},
		assert:          assertNoError,
	},
}

var isLessThanOrEqualFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLessThanOrEqualFloat passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsLessThanOrEqualFloat(1)},
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanOrEqualFloat violation on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		options:         []validation.Option{it.IsLessThanOrEqualFloat(1)},
		assert:          assertHasOneViolation(code.TooHighOrEqual, "This value should be less than or equal to 1.", ""),
	},
	{
		name:            "IsLessThanOrEqualFloat passes on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsLessThanOrEqualFloat(1)},
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanOrEqualFloat passes on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		options:         []validation.Option{it.IsLessThanOrEqualFloat(1)},
		assert:          assertNoError,
	},
}

var isGreaterThanIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsGreaterThanInteger passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsGreaterThanInteger(1)},
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanInteger violation on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsGreaterThanInteger(2)},
		assert:          assertHasOneViolation(code.TooLow, "This value should be greater than 2.", ""),
	},
	{
		name:            "IsGreaterThanInteger violation on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		options:         []validation.Option{it.IsGreaterThanInteger(2)},
		assert:          assertHasOneViolation(code.TooLow, "This value should be greater than 2.", ""),
	},
	{
		name:            "IsGreaterThanInteger passes on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(3),
		floatValue:      floatValue(3),
		options:         []validation.Option{it.IsGreaterThanInteger(2)},
		assert:          assertNoError,
	},
}

var isGreaterThanFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsGreaterThanFloat passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsGreaterThanFloat(1)},
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanFloat violation on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsGreaterThanFloat(2)},
		assert:          assertHasOneViolation(code.TooLow, "This value should be greater than 2.", ""),
	},
	{
		name:            "IsGreaterThanFloat violation on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		options:         []validation.Option{it.IsGreaterThanFloat(2)},
		assert:          assertHasOneViolation(code.TooLow, "This value should be greater than 2.", ""),
	},
	{
		name:            "IsGreaterThanFloat passes on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(3),
		floatValue:      floatValue(3),
		options:         []validation.Option{it.IsGreaterThanFloat(2)},
		assert:          assertNoError,
	},
}

var isGreaterThanOrEqualIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsGreaterThanOrEqualInteger passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsGreaterThanOrEqualInteger(1)},
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanOrEqualInteger violation on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsGreaterThanOrEqualInteger(2)},
		assert:          assertHasOneViolation(code.TooLowOrEqual, "This value should be greater than or equal to 2.", ""),
	},
	{
		name:            "IsGreaterThanOrEqualInteger passes on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		options:         []validation.Option{it.IsGreaterThanOrEqualInteger(2)},
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanOrEqualInteger passes on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(3),
		floatValue:      floatValue(3),
		options:         []validation.Option{it.IsGreaterThanOrEqualInteger(2)},
		assert:          assertNoError,
	},
}

var isGreaterThanOrEqualFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsGreaterThanOrEqualFloat passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsGreaterThanOrEqualFloat(1)},
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanOrEqualFloat violation on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsGreaterThanOrEqualFloat(2)},
		assert:          assertHasOneViolation(code.TooLowOrEqual, "This value should be greater than or equal to 2.", ""),
	},
	{
		name:            "IsGreaterThanOrEqualFloat passes on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		options:         []validation.Option{it.IsGreaterThanOrEqualFloat(2)},
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanOrEqualFloat passes on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(3),
		floatValue:      floatValue(3),
		options:         []validation.Option{it.IsGreaterThanOrEqualFloat(2)},
		assert:          assertNoError,
	},
}

var isPositiveTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsPositive passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsPositive()},
		assert:          assertNoError,
	},
	{
		name:            "IsPositive violation on negative",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		options:         []validation.Option{it.IsPositive()},
		assert:          assertHasOneViolation(code.NotPositive, "This value should be positive.", ""),
	},
	{
		name:            "IsPositive violation on zero",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		options:         []validation.Option{it.IsPositive()},
		assert:          assertHasOneViolation(code.NotPositive, "This value should be positive.", ""),
	},
	{
		name:            "IsPositive passes on positive",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsPositive()},
		assert:          assertNoError,
	},
}

var isPositiveOrZeroTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsPositiveOrZero passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsPositiveOrZero()},
		assert:          assertNoError,
	},
	{
		name:            "IsPositiveOrZero violation on negative",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		options:         []validation.Option{it.IsPositiveOrZero()},
		assert:          assertHasOneViolation(code.NotPositiveOrZero, "This value should be either positive or zero.", ""),
	},
	{
		name:            "IsPositiveOrZero passes on zero",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		options:         []validation.Option{it.IsPositiveOrZero()},
		assert:          assertNoError,
	},
	{
		name:            "IsPositiveOrZero passes on positive",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsPositiveOrZero()},
		assert:          assertNoError,
	},
}

var isNegativeTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNegative passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsNegative()},
		assert:          assertNoError,
	},
	{
		name:            "IsNegative passes on negative",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		options:         []validation.Option{it.IsNegative()},
		assert:          assertNoError,
	},
	{
		name:            "IsNegative violation on zero",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		options:         []validation.Option{it.IsNegative()},
		assert:          assertHasOneViolation(code.NotNegative, "This value should be negative.", ""),
	},
	{
		name:            "IsNegative violation on positive",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsNegative()},
		assert:          assertHasOneViolation(code.NotNegative, "This value should be negative.", ""),
	},
}

var isNegativeOrZeroTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNegativeOrZero passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		options:         []validation.Option{it.IsNegativeOrZero()},
		assert:          assertNoError,
	},
	{
		name:            "IsNegativeOrZero passes on negative",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		options:         []validation.Option{it.IsNegativeOrZero()},
		assert:          assertNoError,
	},
	{
		name:            "IsNegativeOrZero passes on zero",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		options:         []validation.Option{it.IsNegativeOrZero()},
		assert:          assertNoError,
	},
	{
		name:            "IsNegativeOrZero violation on positive",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		options:         []validation.Option{it.IsNegativeOrZero()},
		assert:          assertHasOneViolation(code.NotNegativeOrZero, "This value should be either negative or zero.", ""),
	},
}

var isEqualToStringTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsEqualToString passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsEqualToString("expected")},
		assert:          assertNoError,
	},
	{
		name:            "IsEqualToString violation on not equal value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("actual"),
		options:         []validation.Option{it.IsEqualToString("expected")},
		assert:          assertHasOneViolation(code.Equal, `This value should be equal to "expected".`, ""),
	},
	{
		name:            "IsEqualToString passes on equal value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("expected"),
		options:         []validation.Option{it.IsEqualToString("expected")},
		assert:          assertNoError,
	},
	{
		name:            "IsEqualToString violation with custom message",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("actual"),
		options: []validation.Option{
			it.IsEqualToString("expected").Message(`Unexpected value {{ value }}, expected value is {{ comparedValue }}.`),
		},
		assert: assertHasOneViolation(code.Equal, `Unexpected value "actual", expected value is "expected".`, ""),
	},
	{
		name:            "IsEqualToString passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("actual"),
		options:         []validation.Option{it.IsEqualToString("expected").When(false)},
		assert:          assertNoError,
	},
	{
		name:            "IsEqualToString violation when condition is tue",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("actual"),
		options:         []validation.Option{it.IsEqualToString("expected").When(true)},
		assert:          assertHasOneViolation(code.Equal, `This value should be equal to "expected".`, ""),
	},
}

var isNotEqualToStringTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNotEqualToString passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		options:         []validation.Option{it.IsNotEqualToString("expected")},
		assert:          assertNoError,
	},
	{
		name:            "IsNotEqualToString passes on not equal value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("actual"),
		options:         []validation.Option{it.IsNotEqualToString("expected")},
		assert:          assertNoError,
	},
	{
		name:            "IsNotEqualToString violation on equal value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("expected"),
		options:         []validation.Option{it.IsNotEqualToString("expected")},
		assert:          assertHasOneViolation(code.NotEqual, `This value should not be equal to "expected".`, ""),
	},
}
