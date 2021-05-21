package test

import (
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/it"
)

var numberComparisonTestCases = mergeTestCases(
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

var rangeComparisonTestCases = mergeTestCases(
	isBetweenIntegersTestCases,
	isBetweenFloatsTestCases,
)

var stringComparisonTestCases = mergeTestCases(
	isEqualToStringTestCases,
	isNotEqualToStringTestCases,
)

var timeComparisonTestCases = mergeTestCases(
	isEarlierThanTestCases,
	isEarlierThanOrEqualTestCases,
	isLaterThanTestCases,
	isLaterThanOrEqualTestCases,
)

var isEqualToIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsEqualToInteger passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsEqualToInteger(1),
		assert:          assertNoError,
	},
	{
		name:            "IsEqualToInteger violation on not equal int",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsEqualToInteger(1),
		assert:          assertHasOneViolation(code.Equal, "This value should be equal to 1."),
	},
	{
		name:            "IsEqualToInteger passes on equal int",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsEqualToInteger(1),
		assert:          assertNoError,
	},
	{
		name:            "IsEqualToInteger violation on not equal float",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(0.99),
		constraint:      it.IsEqualToInteger(1),
		assert:          assertHasOneViolation(code.Equal, "This value should be equal to 1."),
	},
	{
		name:            "IsEqualToInteger violation with custom message",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint: it.IsEqualToInteger(1).
			Message(
				`Unexpected value "{{ value }}" at {{ custom }}, expected value is "{{ comparedValue }}".`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(code.Equal, `Unexpected value "0" at parameter, expected value is "1".`),
	},
	{
		name:            "IsEqualToInteger passes when condition is false",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsEqualToInteger(1).When(false),
		assert:          assertNoError,
	},
	{
		name:            "IsEqualToInteger violation when condition is true",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsEqualToInteger(1).When(true),
		assert:          assertHasOneViolation(code.Equal, "This value should be equal to 1."),
	},
}

var isEqualToFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsEqualToFloat passes on nil",
		isApplicableFor: specificValueTypes(floatType),
		constraint:      it.IsEqualToFloat(1.5),
		assert:          assertNoError,
	},
	{
		name:            "IsEqualToFloat violation on not equal float",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(0.5),
		constraint:      it.IsEqualToFloat(1.5),
		assert:          assertHasOneViolation(code.Equal, "This value should be equal to 1.5."),
	},
	{
		name:            "IsEqualToFloat passes on equal float and int",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsEqualToFloat(1),
		assert:          assertNoError,
	},
}

var isNotEqualToIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNotEqualToInteger passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsNotEqualToInteger(1),
		assert:          assertNoError,
	},
	{
		name:            "IsNotEqualToInteger violation on equal int or float",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsNotEqualToInteger(1),
		assert:          assertHasOneViolation(code.NotEqual, "This value should not be equal to 1."),
	},
	{
		name:            "IsNotEqualToInteger passes on not equal float and int",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsNotEqualToInteger(1),
		assert:          assertNoError,
	},
}

var isNotEqualToFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNotEqualToFloat passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsNotEqualToFloat(1),
		assert:          assertNoError,
	},
	{
		name:            "IsNotEqualToFloat violation on equal int or float",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsNotEqualToFloat(1),
		assert:          assertHasOneViolation(code.NotEqual, "This value should not be equal to 1."),
	},
	{
		name:            "IsNotEqualToFloat passes on not equal float and int",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsNotEqualToFloat(1),
		assert:          assertNoError,
	},
}

var isLessThanIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLessThanInteger passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsLessThanInteger(1),
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanInteger violation on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		constraint:      it.IsLessThanInteger(1),
		assert:          assertHasOneViolation(code.TooHigh, "This value should be less than 1."),
	},
	{
		name:            "IsLessThanInteger violation on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsLessThanInteger(1),
		assert:          assertHasOneViolation(code.TooHigh, "This value should be less than 1."),
	},
	{
		name:            "IsLessThanInteger passes on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsLessThanInteger(1),
		assert:          assertNoError,
	},
}

var isLessThanFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLessThanFloat passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsLessThanFloat(1),
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanFloat violation on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		constraint:      it.IsLessThanFloat(1),
		assert:          assertHasOneViolation(code.TooHigh, "This value should be less than 1."),
	},
	{
		name:            "IsLessThanFloat violation on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsLessThanFloat(1),
		assert:          assertHasOneViolation(code.TooHigh, "This value should be less than 1."),
	},
	{
		name:            "IsLessThanFloat passes on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsLessThanFloat(1),
		assert:          assertNoError,
	},
}

var isLessThanOrEqualIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLessThanOrEqualInteger passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsLessThanOrEqualInteger(1),
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanOrEqualInteger violation on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		constraint:      it.IsLessThanOrEqualInteger(1),
		assert:          assertHasOneViolation(code.TooHighOrEqual, "This value should be less than or equal to 1."),
	},
	{
		name:            "IsLessThanOrEqualInteger passes on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsLessThanOrEqualInteger(1),
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanOrEqualInteger passes on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsLessThanOrEqualInteger(1),
		assert:          assertNoError,
	},
}

var isLessThanOrEqualFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLessThanOrEqualFloat passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsLessThanOrEqualFloat(1),
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanOrEqualFloat violation on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		constraint:      it.IsLessThanOrEqualFloat(1),
		assert:          assertHasOneViolation(code.TooHighOrEqual, "This value should be less than or equal to 1."),
	},
	{
		name:            "IsLessThanOrEqualFloat passes on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsLessThanOrEqualFloat(1),
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanOrEqualFloat passes on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsLessThanOrEqualFloat(1),
		assert:          assertNoError,
	},
}

var isGreaterThanIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsGreaterThanInteger passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsGreaterThanInteger(1),
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanInteger violation on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsGreaterThanInteger(2),
		assert:          assertHasOneViolation(code.TooLow, "This value should be greater than 2."),
	},
	{
		name:            "IsGreaterThanInteger violation on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		constraint:      it.IsGreaterThanInteger(2),
		assert:          assertHasOneViolation(code.TooLow, "This value should be greater than 2."),
	},
	{
		name:            "IsGreaterThanInteger passes on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(3),
		floatValue:      floatValue(3),
		constraint:      it.IsGreaterThanInteger(2),
		assert:          assertNoError,
	},
}

var isGreaterThanFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsGreaterThanFloat passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsGreaterThanFloat(1),
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanFloat violation on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsGreaterThanFloat(2),
		assert:          assertHasOneViolation(code.TooLow, "This value should be greater than 2."),
	},
	{
		name:            "IsGreaterThanFloat violation on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		constraint:      it.IsGreaterThanFloat(2),
		assert:          assertHasOneViolation(code.TooLow, "This value should be greater than 2."),
	},
	{
		name:            "IsGreaterThanFloat passes on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(3),
		floatValue:      floatValue(3),
		constraint:      it.IsGreaterThanFloat(2),
		assert:          assertNoError,
	},
}

var isGreaterThanOrEqualIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsGreaterThanOrEqualInteger passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsGreaterThanOrEqualInteger(1),
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanOrEqualInteger violation on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsGreaterThanOrEqualInteger(2),
		assert:          assertHasOneViolation(code.TooLowOrEqual, "This value should be greater than or equal to 2."),
	},
	{
		name:            "IsGreaterThanOrEqualInteger passes on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		constraint:      it.IsGreaterThanOrEqualInteger(2),
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanOrEqualInteger passes on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(3),
		floatValue:      floatValue(3),
		constraint:      it.IsGreaterThanOrEqualInteger(2),
		assert:          assertNoError,
	},
}

var isGreaterThanOrEqualFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsGreaterThanOrEqualFloat passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsGreaterThanOrEqualFloat(1),
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanOrEqualFloat violation on less value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsGreaterThanOrEqualFloat(2),
		assert:          assertHasOneViolation(code.TooLowOrEqual, "This value should be greater than or equal to 2."),
	},
	{
		name:            "IsGreaterThanOrEqualFloat passes on equal value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		constraint:      it.IsGreaterThanOrEqualFloat(2),
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanOrEqualFloat passes on greater value",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(3),
		floatValue:      floatValue(3),
		constraint:      it.IsGreaterThanOrEqualFloat(2),
		assert:          assertNoError,
	},
}

var isPositiveTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsPositive passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsPositive(),
		assert:          assertNoError,
	},
	{
		name:            "IsPositive violation on negative",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		constraint:      it.IsPositive(),
		assert:          assertHasOneViolation(code.NotPositive, "This value should be positive."),
	},
	{
		name:            "IsPositive violation on zero",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsPositive(),
		assert:          assertHasOneViolation(code.NotPositive, "This value should be positive."),
	},
	{
		name:            "IsPositive passes on positive",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsPositive(),
		assert:          assertNoError,
	},
}

var isPositiveOrZeroTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsPositiveOrZero passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsPositiveOrZero(),
		assert:          assertNoError,
	},
	{
		name:            "IsPositiveOrZero violation on negative",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		constraint:      it.IsPositiveOrZero(),
		assert:          assertHasOneViolation(code.NotPositiveOrZero, "This value should be either positive or zero."),
	},
	{
		name:            "IsPositiveOrZero passes on zero",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsPositiveOrZero(),
		assert:          assertNoError,
	},
	{
		name:            "IsPositiveOrZero passes on positive",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsPositiveOrZero(),
		assert:          assertNoError,
	},
}

var isNegativeTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNegative passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsNegative(),
		assert:          assertNoError,
	},
	{
		name:            "IsNegative passes on negative",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		constraint:      it.IsNegative(),
		assert:          assertNoError,
	},
	{
		name:            "IsNegative violation on zero",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsNegative(),
		assert:          assertHasOneViolation(code.NotNegative, "This value should be negative."),
	},
	{
		name:            "IsNegative violation on positive",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsNegative(),
		assert:          assertHasOneViolation(code.NotNegative, "This value should be negative."),
	},
}

var isNegativeOrZeroTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNegativeOrZero passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsNegativeOrZero(),
		assert:          assertNoError,
	},
	{
		name:            "IsNegativeOrZero passes on negative",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		constraint:      it.IsNegativeOrZero(),
		assert:          assertNoError,
	},
	{
		name:            "IsNegativeOrZero passes on zero",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsNegativeOrZero(),
		assert:          assertNoError,
	},
	{
		name:            "IsNegativeOrZero violation on positive",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsNegativeOrZero(),
		assert:          assertHasOneViolation(code.NotNegativeOrZero, "This value should be either negative or zero."),
	},
}

var isBetweenIntegersTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsBetweenIntegers error on equal min and max",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsBetweenIntegers(1, 1),
		assert:          assertError(`failed to set up constraint "RangeConstraint": invalid range`),
	},
	{
		name:            "IsBetweenIntegers error on min greater than max",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsBetweenIntegers(1, 0),
		assert:          assertError(`failed to set up constraint "RangeConstraint": invalid range`),
	},
	{
		name:            "IsBetweenIntegers passes on nil",
		isApplicableFor: specificValueTypes(intType, floatType),
		constraint:      it.IsBetweenIntegers(1, 2),
		assert:          assertNoError,
	},
	{
		name:            "IsBetweenIntegers violation on value less than min",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsBetweenIntegers(1, 2),
		assert:          assertHasOneViolation(code.NotInRange, "This value should be between 1 and 2."),
	},
	{
		name:            "IsBetweenIntegers violation on value greater than max",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(3),
		floatValue:      floatValue(3),
		constraint:      it.IsBetweenIntegers(1, 2),
		assert:          assertHasOneViolation(code.NotInRange, "This value should be between 1 and 2."),
	},
	{
		name:            "IsBetweenIntegers passes on value equal to min",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsBetweenIntegers(1, 2),
		assert:          assertNoError,
	},
	{
		name:            "IsBetweenIntegers passes on value equal to max",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(2),
		floatValue:      floatValue(2),
		constraint:      it.IsBetweenIntegers(1, 2),
		assert:          assertNoError,
	},
	{
		name:            "IsBetweenIntegers violation with custom message",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint: it.IsBetweenIntegers(1, 2).
			Message(
				`Unexpected value "{{ value }}" at {{ custom }}, expected value must be between "{{ min }}" and "{{ max }}".`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(
			code.NotInRange,
			`Unexpected value "0" at parameter, expected value must be between "1" and "2".`,
		),
	},
	{
		name:            "IsBetweenIntegers passes when condition is false",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsBetweenIntegers(1, 2).When(false),
		assert:          assertNoError,
	},
	{
		name:            "IsBetweenIntegers violation when condition is true",
		isApplicableFor: specificValueTypes(intType, floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsBetweenIntegers(1, 2).When(true),
		assert:          assertHasOneViolation(code.NotInRange, "This value should be between 1 and 2."),
	},
}

var isBetweenFloatsTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsBetweenFloats passes on nil",
		isApplicableFor: specificValueTypes(floatType),
		constraint:      it.IsBetweenFloats(1.1, 2.2),
		assert:          assertNoError,
	},
	{
		name:            "IsBetweenFloats violation on value less than min",
		isApplicableFor: specificValueTypes(floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsBetweenFloats(1.1, 2.2),
		assert:          assertHasOneViolation(code.NotInRange, "This value should be between 1.1 and 2.2."),
	},
	{
		name:            "IsBetweenFloats violation on value greater than max",
		isApplicableFor: specificValueTypes(floatType),
		intValue:        intValue(3),
		floatValue:      floatValue(3),
		constraint:      it.IsBetweenFloats(1.1, 2.2),
		assert:          assertHasOneViolation(code.NotInRange, "This value should be between 1.1 and 2.2."),
	},
	{
		name:            "IsBetweenFloats passes on value equal to min",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(1.1),
		constraint:      it.IsBetweenFloats(1.1, 2.2),
		assert:          assertNoError,
	},
	{
		name:            "IsBetweenFloats passes on value equal to max",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(2.2),
		constraint:      it.IsBetweenFloats(1.1, 2.2),
		assert:          assertNoError,
	},
}

var isEqualToStringTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsEqualToString passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsEqualToString("expected"),
		assert:          assertNoError,
	},
	{
		name:            "IsEqualToString violation on not equal value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("actual"),
		constraint:      it.IsEqualToString("expected"),
		assert:          assertHasOneViolation(code.Equal, `This value should be equal to "expected".`),
	},
	{
		name:            "IsEqualToString passes on equal value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("expected"),
		constraint:      it.IsEqualToString("expected"),
		assert:          assertNoError,
	},
	{
		name:            "IsEqualToString violation with custom message",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("actual"),
		constraint: it.IsEqualToString("expected").
			Message(
				`Unexpected value {{ value }} at {{ custom }}, expected value is {{ comparedValue }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(
			code.Equal,
			`Unexpected value "actual" at parameter, expected value is "expected".`,
		),
	},
	{
		name:            "IsEqualToString passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("actual"),
		constraint:      it.IsEqualToString("expected").When(false),
		assert:          assertNoError,
	},
	{
		name:            "IsEqualToString violation when condition is tue",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("actual"),
		constraint:      it.IsEqualToString("expected").When(true),
		assert:          assertHasOneViolation(code.Equal, `This value should be equal to "expected".`),
	},
}

var isNotEqualToStringTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNotEqualToString passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsNotEqualToString("expected"),
		assert:          assertNoError,
	},
	{
		name:            "IsNotEqualToString passes on not equal value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("actual"),
		constraint:      it.IsNotEqualToString("expected"),
		assert:          assertNoError,
	},
	{
		name:            "IsNotEqualToString violation on equal value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("expected"),
		constraint:      it.IsNotEqualToString("expected"),
		assert:          assertHasOneViolation(code.NotEqual, `This value should not be equal to "expected".`),
	},
}

var isEarlierThanTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsEarlierThan passes on nil",
		isApplicableFor: specificValueTypes(timeType),
		constraint:      it.IsEarlierThan(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsEarlierThan violation on greater value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 40, 0, 0, time.UTC)),
		constraint:      it.IsEarlierThan(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertHasOneViolation(code.TooLate, "This value should be earlier than 2021-03-29T12:30:00Z."),
	},
	{
		name:            "IsEarlierThan violation on equal value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		constraint:      it.IsEarlierThan(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertHasOneViolation(code.TooLate, "This value should be earlier than 2021-03-29T12:30:00Z."),
	},
	{
		name:            "IsEarlierThan passes on less value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 29, 29, 0, time.UTC)),
		constraint:      it.IsEarlierThan(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsEarlierThan violation with custom message",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 40, 0, 0, time.UTC)),
		constraint: it.IsEarlierThan(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)).
			Message(
				`Unexpected value "{{ value }}" at {{ custom }}, expected value must be earlier than "{{ comparedValue }}".`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(
			code.TooLate,
			`Unexpected value "2021-03-29T12:40:00Z" at parameter, expected value must be earlier than "2021-03-29T12:30:00Z".`,
		),
	},
	{
		name:            "IsEarlierThan violation with custom message and time layout",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 40, 0, 0, time.UTC)),
		constraint: it.IsEarlierThan(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)).
			Message(`Unexpected value "{{ value }}", expected value must be earlier than "{{ comparedValue }}".`).
			Layout(time.RFC822),
		assert: assertHasOneViolation(
			code.TooLate,
			`Unexpected value "29 Mar 21 12:40 UTC", expected value must be earlier than "29 Mar 21 12:30 UTC".`,
		),
	},
	{
		name:            "IsEarlierThan passes when condition is false",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 40, 0, 0, time.UTC)),
		constraint: it.IsEarlierThan(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)).
			When(false),
		assert: assertNoError,
	},
	{
		name:            "IsEarlierThan violation when condition is true",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 40, 0, 0, time.UTC)),
		constraint: it.IsEarlierThan(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)).
			When(true),
		assert: assertHasOneViolation(code.TooLate, "This value should be earlier than 2021-03-29T12:30:00Z."),
	},
}

var isEarlierThanOrEqualTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsEarlierThanOrEqual passes on nil",
		isApplicableFor: specificValueTypes(timeType),
		constraint:      it.IsEarlierThanOrEqual(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsEarlierThanOrEqual violation on greater value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 40, 0, 0, time.UTC)),
		constraint:      it.IsEarlierThanOrEqual(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertHasOneViolation(code.TooLateOrEqual, "This value should be earlier than or equal to 2021-03-29T12:30:00Z."),
	},
	{
		name:            "IsEarlierThanOrEqual passes on equal value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		constraint:      it.IsEarlierThanOrEqual(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsEarlierThanOrEqual passes on equal value with different time zone",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		constraint:      it.IsEarlierThanOrEqual(time.Date(2021, 03, 29, 8, 30, 0, 0, givenLocation("America/New_York"))),
		assert:          assertNoError,
	},
	{
		name:            "IsEarlierThanOrEqual passes on less value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 29, 29, 0, time.UTC)),
		constraint:      it.IsEarlierThanOrEqual(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
}

var isLaterThanTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLaterThan passes on nil",
		isApplicableFor: specificValueTypes(timeType),
		constraint:      it.IsLaterThan(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsLaterThan passes on greater value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 40, 0, 0, time.UTC)),
		constraint:      it.IsLaterThan(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsLaterThan violation on equal value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		constraint:      it.IsLaterThan(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertHasOneViolation(code.TooEarly, "This value should be later than 2021-03-29T12:30:00Z."),
	},
	{
		name:            "IsLaterThan violation on less value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 29, 29, 0, time.UTC)),
		constraint:      it.IsLaterThan(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertHasOneViolation(code.TooEarly, "This value should be later than 2021-03-29T12:30:00Z."),
	},
}

var isLaterThanOrEqualTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLaterThanOrEqual passes on nil",
		isApplicableFor: specificValueTypes(timeType),
		constraint:      it.IsLaterThanOrEqual(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsLaterThanOrEqual passes on greater value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 40, 0, 0, time.UTC)),
		constraint:      it.IsLaterThanOrEqual(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsLaterThanOrEqual passes on equal value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		constraint:      it.IsLaterThanOrEqual(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsLaterThanOrEqual passes on equal value with different time zone",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		constraint:      it.IsLaterThanOrEqual(time.Date(2021, 03, 29, 8, 30, 0, 0, givenLocation("America/New_York"))),
		assert:          assertNoError,
	},
	{
		name:            "IsLaterThanOrEqual violation on less value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 03, 29, 12, 29, 29, 0, time.UTC)),
		constraint:      it.IsLaterThanOrEqual(time.Date(2021, 03, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertHasOneViolation(code.TooEarlyOrEqual, "This value should be later than or equal to 2021-03-29T12:30:00Z."),
	},
}

var isBetweenTimeTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsBetweenTime error on equal min and max",
		isApplicableFor: specificValueTypes(timeType),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 04, 4, 12, 30, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 04, 4, 12, 30, 0, 0, time.UTC)),
		),
		assert: assertError(`failed to set up constraint "TimeRangeConstraint": invalid range`),
	},
	{
		name:            "IsBetweenTime error on equal min and max in different time zones",
		isApplicableFor: specificValueTypes(timeType),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 04, 4, 12, 30, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 04, 4, 8, 30, 0, 0, givenLocation("America/New_York"))),
		),
		assert: assertError(`failed to set up constraint "TimeRangeConstraint": invalid range`),
	},
	{
		name:            "IsBetweenTime error on min greater than max",
		isApplicableFor: specificValueTypes(timeType),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 04, 4, 12, 40, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 04, 4, 12, 30, 0, 0, time.UTC)),
		),
		assert: assertError(`failed to set up constraint "TimeRangeConstraint": invalid range`),
	},
	{
		name:            "IsBetweenTime passes on nil",
		isApplicableFor: specificValueTypes(timeType),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 04, 4, 12, 30, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 04, 4, 12, 40, 0, 0, time.UTC)),
		),
		assert: assertNoError,
	},
	{
		name:            "IsBetweenTime violation on value less than min",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 04, 4, 12, 20, 0, 0, time.UTC)),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 04, 4, 12, 30, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 04, 4, 12, 40, 0, 0, time.UTC)),
		),
		assert: assertHasOneViolation(code.NotInRange, "This value should be between 2021-04-04T12:30:00Z and 2021-04-04T12:40:00Z."),
	},
	{
		name:            "IsBetweenTime violation on value greater than max",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 04, 4, 12, 50, 0, 0, time.UTC)),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 04, 4, 12, 30, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 04, 4, 12, 40, 0, 0, time.UTC)),
		),
		assert: assertHasOneViolation(code.NotInRange, "This value should be between 2021-04-04T12:30:00Z and 2021-04-04T12:40:00Z."),
	},
	{
		name:            "IsBetweenTime passes on value equal to min",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 04, 4, 12, 30, 0, 0, time.UTC)),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 04, 4, 12, 30, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 04, 4, 12, 40, 0, 0, time.UTC)),
		),
		assert: assertNoError,
	},
	{
		name:            "IsBetweenTime passes on value equal to max",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 04, 4, 12, 40, 0, 0, time.UTC)),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 04, 4, 12, 30, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 04, 4, 12, 40, 0, 0, time.UTC)),
		),
		assert: assertNoError,
	},
	{
		name:            "IsBetweenTime violation with custom message",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 04, 4, 12, 20, 0, 0, time.UTC)),
		constraint: it.
			IsBetweenTime(
				*timeValue(time.Date(2021, 04, 4, 12, 30, 0, 0, time.UTC)),
				*timeValue(time.Date(2021, 04, 4, 12, 40, 0, 0, time.UTC)),
			).
			Message(
				`Unexpected value "{{ value }}" at {{ custom }}, expected value must be between "{{ min }}" and "{{ max }}".`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(
			code.NotInRange,
			`Unexpected value "2021-04-04T12:20:00Z" at parameter, expected value must be between "2021-04-04T12:30:00Z" and "2021-04-04T12:40:00Z".`,
		),
	},
	{
		name:            "IsBetweenTime violation with custom message and time layout",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 04, 4, 12, 20, 0, 0, time.UTC)),
		constraint: it.
			IsBetweenTime(
				*timeValue(time.Date(2021, 04, 4, 12, 30, 0, 0, time.UTC)),
				*timeValue(time.Date(2021, 04, 4, 12, 40, 0, 0, time.UTC)),
			).
			Message(`Unexpected value "{{ value }}", expected value must be between "{{ min }}" and "{{ max }}".`).
			Layout(time.RFC822),
		assert: assertHasOneViolation(
			code.NotInRange,
			`Unexpected value "04 Apr 21 12:20 UTC", expected value must be between "04 Apr 21 12:30 UTC" and "04 Apr 21 12:40 UTC".`,
		),
	},
	{
		name:            "IsBetweenTime passes when condition is false",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 04, 4, 12, 20, 0, 0, time.UTC)),
		constraint: it.
			IsBetweenTime(
				*timeValue(time.Date(2021, 04, 4, 12, 30, 0, 0, time.UTC)),
				*timeValue(time.Date(2021, 04, 4, 12, 40, 0, 0, time.UTC)),
			).
			When(false),
		assert: assertNoError,
	},
	{
		name:            "IsBetweenTime violation when condition is true",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 04, 4, 12, 20, 0, 0, time.UTC)),
		constraint: it.
			IsBetweenTime(
				*timeValue(time.Date(2021, 04, 4, 12, 30, 0, 0, time.UTC)),
				*timeValue(time.Date(2021, 04, 4, 12, 40, 0, 0, time.UTC)),
			).
			When(true),
		assert: assertHasOneViolation(code.NotInRange, "This value should be between 2021-04-04T12:30:00Z and 2021-04-04T12:40:00Z."),
	},
}
