package test

import (
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var numberComparisonTestCases = mergeTestCases(
	isLessThanIntegerTestCases,
	isLessThanFloatTestCases,
	isLessThanOrEqualIntegerTestCases,
	isLessThanOrEqualFloatTestCases,
	isGreaterThanIntegerTestCases,
	isGreaterThanFloatTestCases,
	isGreaterThanOrEqualIntegerTestCases,
	isGreaterThanOrEqualFloatTestCases,
	isPositiveTestCases,
	isPositiveFloatTestCases,
	isPositiveOrZeroTestCases,
	isPositiveOrZeroFloatTestCases,
	isNegativeTestCases,
	isNegativeFloatTestCases,
	isNegativeOrZeroTestCases,
	isNegativeOrZeroFloatTestCases,
	isDivisibleTestCases,
	isDivisibleByFloatTestCases,
)

var rangeComparisonTestCases = mergeTestCases(
	isBetweenIntegersTestCases,
	isBetweenFloatsTestCases,
)

var comparableComparisonTestCases = mergeTestCases(
	isEqualToTestCases,
	isNotEqualToTestCases,
)

var timeComparisonTestCases = mergeTestCases(
	isEarlierThanTestCases,
	isEarlierThanOrEqualTestCases,
	isLaterThanTestCases,
	isLaterThanOrEqualTestCases,
)

var isLessThanIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLessThan passes on nil",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsLessThan(1),
		assert:          assertNoError,
	},
	{
		name:            "IsLessThan violation on greater value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(2),
		constraint:      it.IsLessThan(1),
		assert:          assertHasOneViolation(validation.ErrTooHigh, "This value should be less than 1."),
	},
	{
		name:            "IsLessThan violation on equal value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(1),
		constraint:      it.IsLessThan(1),
		assert:          assertHasOneViolation(validation.ErrTooHigh, "This value should be less than 1."),
	},
	{
		name:            "IsLessThan passes on less value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		constraint:      it.IsLessThan(1),
		assert:          assertNoError,
	},
}

var isLessThanFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLessThan (float) passes on nil",
		isApplicableFor: specificValueTypes(floatType),
		constraint:      it.IsLessThan(1.0),
		assert:          assertNoError,
	},
	{
		name:            "IsLessThan (float) violation on greater value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(2),
		constraint:      it.IsLessThan(1.0),
		assert:          assertHasOneViolation(validation.ErrTooHigh, "This value should be less than 1."),
	},
	{
		name:            "IsLessThan (float) violation on equal value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(1),
		constraint:      it.IsLessThan(1.0),
		assert:          assertHasOneViolation(validation.ErrTooHigh, "This value should be less than 1."),
	},
	{
		name:            "IsLessThan (float) passes on less value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(0),
		constraint:      it.IsLessThan(1.0),
		assert:          assertNoError,
	},
}

var isLessThanOrEqualIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLessThanOrEqual passes on nil",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsLessThanOrEqual(1),
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanOrEqual violation on greater value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(2),
		constraint:      it.IsLessThanOrEqual(1),
		assert:          assertHasOneViolation(validation.ErrTooHighOrEqual, "This value should be less than or equal to 1."),
	},
	{
		name:            "IsLessThanOrEqual passes on equal value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(1),
		constraint:      it.IsLessThanOrEqual(1),
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanOrEqual passes on less value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		constraint:      it.IsLessThanOrEqual(1),
		assert:          assertNoError,
	},
}

var isLessThanOrEqualFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLessThanOrEqual (float) passes on nil",
		isApplicableFor: specificValueTypes(floatType),
		constraint:      it.IsLessThanOrEqual(1.0),
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanOrEqual (float) violation on greater value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(2),
		constraint:      it.IsLessThanOrEqual(1.0),
		assert:          assertHasOneViolation(validation.ErrTooHighOrEqual, "This value should be less than or equal to 1."),
	},
	{
		name:            "IsLessThanOrEqual (float) passes on equal value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(1),
		constraint:      it.IsLessThanOrEqual(1.0),
		assert:          assertNoError,
	},
	{
		name:            "IsLessThanOrEqual (float) passes on less value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(0),
		constraint:      it.IsLessThanOrEqual(1.0),
		assert:          assertNoError,
	},
}

var isGreaterThanIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsGreaterThan passes on nil",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsGreaterThan(1),
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThan violation on less value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(1),
		constraint:      it.IsGreaterThan(2),
		assert:          assertHasOneViolation(validation.ErrTooLow, "This value should be greater than 2."),
	},
	{
		name:            "IsGreaterThan violation on equal value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(2),
		constraint:      it.IsGreaterThan(2),
		assert:          assertHasOneViolation(validation.ErrTooLow, "This value should be greater than 2."),
	},
	{
		name:            "IsGreaterThan passes on greater value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(3),
		constraint:      it.IsGreaterThan(2),
		assert:          assertNoError,
	},
}

var isGreaterThanFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsGreaterThan (float) passes on nil",
		isApplicableFor: specificValueTypes(floatType),
		constraint:      it.IsGreaterThan(1.0),
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThan (float) violation on less value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(1),
		constraint:      it.IsGreaterThan(2.0),
		assert:          assertHasOneViolation(validation.ErrTooLow, "This value should be greater than 2."),
	},
	{
		name:            "IsGreaterThan (float) violation on equal value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(2),
		constraint:      it.IsGreaterThan(2.0),
		assert:          assertHasOneViolation(validation.ErrTooLow, "This value should be greater than 2."),
	},
	{
		name:            "IsGreaterThan (float) passes on greater value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(3),
		constraint:      it.IsGreaterThan(2.0),
		assert:          assertNoError,
	},
}

var isGreaterThanOrEqualIntegerTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsGreaterThanOrEqual passes on nil",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsGreaterThanOrEqual(1),
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanOrEqual violation on less value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(1),
		constraint:      it.IsGreaterThanOrEqual(2),
		assert:          assertHasOneViolation(validation.ErrTooLowOrEqual, "This value should be greater than or equal to 2."),
	},
	{
		name:            "IsGreaterThanOrEqual passes on equal value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(2),
		constraint:      it.IsGreaterThanOrEqual(2),
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanOrEqual passes on greater value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(3),
		constraint:      it.IsGreaterThanOrEqual(2),
		assert:          assertNoError,
	},
}

var isGreaterThanOrEqualFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsGreaterThanOrEqual (float) passes on nil",
		isApplicableFor: specificValueTypes(floatType),
		constraint:      it.IsGreaterThanOrEqual(1.0),
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanOrEqual (float) violation on less value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(1),
		constraint:      it.IsGreaterThanOrEqual(2.0),
		assert:          assertHasOneViolation(validation.ErrTooLowOrEqual, "This value should be greater than or equal to 2."),
	},
	{
		name:            "IsGreaterThanOrEqual (float) passes on equal value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(2),
		constraint:      it.IsGreaterThanOrEqual(2.0),
		assert:          assertNoError,
	},
	{
		name:            "IsGreaterThanOrEqual (float) passes on greater value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(3),
		constraint:      it.IsGreaterThanOrEqual(2.0),
		assert:          assertNoError,
	},
}

var isPositiveTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsPositive passes on nil",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsPositive[int](),
		assert:          assertNoError,
	},
	{
		name:            "IsPositive violation on negative",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(-1),
		constraint:      it.IsPositive[int](),
		assert:          assertHasOneViolation(validation.ErrNotPositive, "This value should be positive."),
	},
	{
		name:            "IsPositive violation on zero",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		constraint:      it.IsPositive[int](),
		assert:          assertHasOneViolation(validation.ErrNotPositive, "This value should be positive."),
	},
	{
		name:            "IsPositive passes on positive",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(1),
		constraint:      it.IsPositive[int](),
		assert:          assertNoError,
	},
}

var isPositiveFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsPositive (float) passes on nil",
		isApplicableFor: specificValueTypes(floatType),
		constraint:      it.IsPositive[float64](),
		assert:          assertNoError,
	},
	{
		name:            "IsPositive (float) violation on negative",
		isApplicableFor: specificValueTypes(floatType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		constraint:      it.IsPositive[float64](),
		assert:          assertHasOneViolation(validation.ErrNotPositive, "This value should be positive."),
	},
	{
		name:            "IsPositive (float) violation on zero",
		isApplicableFor: specificValueTypes(floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsPositive[float64](),
		assert:          assertHasOneViolation(validation.ErrNotPositive, "This value should be positive."),
	},
	{
		name:            "IsPositive (float) passes on positive",
		isApplicableFor: specificValueTypes(floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsPositive[float64](),
		assert:          assertNoError,
	},
}

var isPositiveOrZeroTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsPositiveOrZero passes on nil",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsPositiveOrZero[int](),
		assert:          assertNoError,
	},
	{
		name:            "IsPositiveOrZero violation on negative",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		constraint:      it.IsPositiveOrZero[int](),
		assert:          assertHasOneViolation(validation.ErrNotPositiveOrZero, "This value should be either positive or zero."),
	},
	{
		name:            "IsPositiveOrZero passes on zero",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsPositiveOrZero[int](),
		assert:          assertNoError,
	},
	{
		name:            "IsPositiveOrZero passes on positive",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsPositiveOrZero[int](),
		assert:          assertNoError,
	},
}

var isPositiveOrZeroFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsPositiveOrZero (float) passes on nil",
		isApplicableFor: specificValueTypes(floatType),
		constraint:      it.IsPositiveOrZero[float64](),
		assert:          assertNoError,
	},
	{
		name:            "IsPositiveOrZero (float) violation on negative",
		isApplicableFor: specificValueTypes(floatType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		constraint:      it.IsPositiveOrZero[float64](),
		assert:          assertHasOneViolation(validation.ErrNotPositiveOrZero, "This value should be either positive or zero."),
	},
	{
		name:            "IsPositiveOrZero (float) passes on zero",
		isApplicableFor: specificValueTypes(floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsPositiveOrZero[float64](),
		assert:          assertNoError,
	},
	{
		name:            "IsPositiveOrZero (float) passes on positive",
		isApplicableFor: specificValueTypes(floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsPositiveOrZero[float64](),
		assert:          assertNoError,
	},
}

var isNegativeTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNegative passes on nil",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsNegative[int](),
		assert:          assertNoError,
	},
	{
		name:            "IsNegative passes on negative",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		constraint:      it.IsNegative[int](),
		assert:          assertNoError,
	},
	{
		name:            "IsNegative violation on zero",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsNegative[int](),
		assert:          assertHasOneViolation(validation.ErrNotNegative, "This value should be negative."),
	},
	{
		name:            "IsNegative violation on positive",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsNegative[int](),
		assert:          assertHasOneViolation(validation.ErrNotNegative, "This value should be negative."),
	},
}

var isNegativeFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNegative (float) passes on nil",
		isApplicableFor: specificValueTypes(floatType),
		constraint:      it.IsNegative[float64](),
		assert:          assertNoError,
	},
	{
		name:            "IsNegative (float) passes on negative",
		isApplicableFor: specificValueTypes(floatType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		constraint:      it.IsNegative[float64](),
		assert:          assertNoError,
	},
	{
		name:            "IsNegative (float) violation on zero",
		isApplicableFor: specificValueTypes(floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsNegative[float64](),
		assert:          assertHasOneViolation(validation.ErrNotNegative, "This value should be negative."),
	},
	{
		name:            "IsNegative (float) violation on positive",
		isApplicableFor: specificValueTypes(floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsNegative[float64](),
		assert:          assertHasOneViolation(validation.ErrNotNegative, "This value should be negative."),
	},
}

var isNegativeOrZeroTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNegativeOrZero passes on nil",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsNegativeOrZero[int](),
		assert:          assertNoError,
	},
	{
		name:            "IsNegativeOrZero passes on negative",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		constraint:      it.IsNegativeOrZero[int](),
		assert:          assertNoError,
	},
	{
		name:            "IsNegativeOrZero passes on zero",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsNegativeOrZero[int](),
		assert:          assertNoError,
	},
	{
		name:            "IsNegativeOrZero violation on positive",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsNegativeOrZero[int](),
		assert:          assertHasOneViolation(validation.ErrNotNegativeOrZero, "This value should be either negative or zero."),
	},
}

var isNegativeOrZeroFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNegativeOrZero (float) passes on nil",
		isApplicableFor: specificValueTypes(floatType),
		constraint:      it.IsNegativeOrZero[float64](),
		assert:          assertNoError,
	},
	{
		name:            "IsNegativeOrZero (float) passes on negative",
		isApplicableFor: specificValueTypes(floatType),
		intValue:        intValue(-1),
		floatValue:      floatValue(-1),
		constraint:      it.IsNegativeOrZero[float64](),
		assert:          assertNoError,
	},
	{
		name:            "IsNegativeOrZero (float) passes on zero",
		isApplicableFor: specificValueTypes(floatType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsNegativeOrZero[float64](),
		assert:          assertNoError,
	},
	{
		name:            "IsNegativeOrZero (float) violation on positive",
		isApplicableFor: specificValueTypes(floatType),
		intValue:        intValue(1),
		floatValue:      floatValue(1),
		constraint:      it.IsNegativeOrZero[float64](),
		assert:          assertHasOneViolation(validation.ErrNotNegativeOrZero, "This value should be either negative or zero."),
	},
}

var isDivisibleTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsDivisibleBy passes on nil",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsDivisibleBy[int](10),
		assert:          assertNoError,
	},
	{
		name:            "IsDivisibleBy passes on zero",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		constraint:      it.IsDivisibleBy[int](10),
		assert:          assertNoError,
	},
	{
		name:            "IsDivisibleBy passes on divisible value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(20),
		constraint:      it.IsDivisibleBy[int](10),
		assert:          assertNoError,
	},
	{
		name:            "IsDivisibleBy violation on not divisible value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(5),
		constraint:      it.IsDivisibleBy[int](10),
		assert:          assertHasOneViolation(validation.ErrNotDivisible, "This value should be a multiple of 10."),
	},
}

var isDivisibleByFloatTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsDivisibleByFloat passes on nil",
		isApplicableFor: specificValueTypes(floatType),
		constraint:      it.IsDivisibleByFloat[float64](0.01),
		assert:          assertNoError,
	},
	{
		name:            "IsDivisibleByFloat passes on zero",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(0),
		constraint:      it.IsDivisibleByFloat[float64](0.01),
		assert:          assertNoError,
	},
	{
		name:            "IsDivisibleByFloat passes on divisible value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(5.55),
		constraint:      it.IsDivisibleByFloat[float64](0.01),
		assert:          assertNoError,
	},
	{
		name:            "IsDivisibleByFloat violation on not divisible value",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(5.555),
		constraint:      it.IsDivisibleByFloat[float64](0.01),
		assert:          assertHasOneViolation(validation.ErrNotDivisible, "This value should be a multiple of 0.01."),
	},
}

var isBetweenIntegersTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsBetween error on equal min and max",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsBetween(1, 1),
		assert:          assertError(`validate by RangeConstraint[int]: invalid range`),
	},
	{
		name:            "IsBetween error on min greater than max",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsBetween(1, 0),
		assert:          assertError(`validate by RangeConstraint[int]: invalid range`),
	},
	{
		name:            "IsBetween passes on nil",
		isApplicableFor: specificValueTypes(intType),
		constraint:      it.IsBetween(1, 2),
		assert:          assertNoError,
	},
	{
		name:            "IsBetween violation on value less than min",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		constraint:      it.IsBetween(1, 2),
		assert:          assertHasOneViolation(validation.ErrNotInRange, "This value should be between 1 and 2."),
	},
	{
		name:            "IsBetween violation on value greater than max",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(3),
		constraint:      it.IsBetween(1, 2),
		assert:          assertHasOneViolation(validation.ErrNotInRange, "This value should be between 1 and 2."),
	},
	{
		name:            "IsBetween passes on value equal to min",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(1),
		constraint:      it.IsBetween(1, 2),
		assert:          assertNoError,
	},
	{
		name:            "IsBetween passes on value equal to max",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(2),
		constraint:      it.IsBetween(1, 2),
		assert:          assertNoError,
	},
	{
		name:            "IsBetween violation with custom message",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		constraint: it.IsBetween(1, 2).
			WithError(ErrCustom).
			WithMessage(
				`Unexpected value "{{ value }}" at {{ custom }}, expected value must be between "{{ min }}" and "{{ max }}".`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(
			ErrCustom,
			`Unexpected value "0" at parameter, expected value must be between "1" and "2".`,
		),
	},
	{
		name:            "IsBetween passes when condition is false",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		constraint:      it.IsBetween(1, 2).When(false),
		assert:          assertNoError,
	},
	{
		name:            "IsBetween passes when groups not match",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		floatValue:      floatValue(0),
		constraint:      it.IsBetween(1, 2).WhenGroups(testGroup),
		assert:          assertNoError,
	},
	{
		name:            "IsBetween violation when condition is true",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(0),
		constraint:      it.IsBetween(1, 2).When(true),
		assert:          assertHasOneViolation(validation.ErrNotInRange, "This value should be between 1 and 2."),
	},
}

var isBetweenFloatsTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsBetween (float) passes on nil",
		isApplicableFor: specificValueTypes(floatType),
		constraint:      it.IsBetween(1.1, 2.2),
		assert:          assertNoError,
	},
	{
		name:            "IsBetween (float) violation on value less than min",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(0),
		constraint:      it.IsBetween(1.1, 2.2),
		assert:          assertHasOneViolation(validation.ErrNotInRange, "This value should be between 1.1 and 2.2."),
	},
	{
		name:            "IsBetween (float) violation on value greater than max",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(3),
		constraint:      it.IsBetween(1.1, 2.2),
		assert:          assertHasOneViolation(validation.ErrNotInRange, "This value should be between 1.1 and 2.2."),
	},
	{
		name:            "IsBetween (float) passes on value equal to min",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(1.1),
		constraint:      it.IsBetween(1.1, 2.2),
		assert:          assertNoError,
	},
	{
		name:            "IsBetween (float) passes on value equal to max",
		isApplicableFor: specificValueTypes(floatType),
		floatValue:      floatValue(2.2),
		constraint:      it.IsBetween(1.1, 2.2),
		assert:          assertNoError,
	},
}

var isEqualToTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsEqualTo passes on nil",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		constraint:      it.IsEqualTo("expected"),
		assert:          assertNoError,
	},
	{
		name:            "IsEqualTo violation on not equal value",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue("actual"),
		constraint:      it.IsEqualTo("expected"),
		assert:          assertHasOneViolation(validation.ErrNotEqual, `This value should be equal to "expected".`),
	},
	{
		name:            "IsEqualTo violation on not equal integer value",
		isApplicableFor: specificValueTypes(intType),
		intValue:        intValue(1),
		constraint:      it.IsEqualTo(2),
		assert:          assertHasOneViolation(validation.ErrNotEqual, `This value should be equal to 2.`),
	},
	{
		name:            "IsEqualTo passes on equal value",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue("expected"),
		constraint:      it.IsEqualTo("expected"),
		assert:          assertNoError,
	},
	{
		name:            "IsEqualTo violation with custom message",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue("actual"),
		constraint: it.IsEqualTo("expected").
			WithError(ErrCustom).
			WithMessage(
				`Unexpected value {{ value }} at {{ custom }}, expected value is {{ comparedValue }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(
			ErrCustom,
			`Unexpected value "actual" at parameter, expected value is "expected".`,
		),
	},
	{
		name:            "IsEqualTo passes when condition is false",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue("actual"),
		constraint:      it.IsEqualTo("expected").When(false),
		assert:          assertNoError,
	},
	{
		name:            "IsEqualTo passes when groups not match",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue("actual"),
		constraint:      it.IsEqualTo("expected").WhenGroups(testGroup),
		assert:          assertNoError,
	},
	{
		name:            "IsEqualTo violation when condition is tue",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue("actual"),
		constraint:      it.IsEqualTo("expected").When(true),
		assert:          assertHasOneViolation(validation.ErrNotEqual, `This value should be equal to "expected".`),
	},
}

var isNotEqualToTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsNotEqualTo passes on nil",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		constraint:      it.IsNotEqualTo("expected"),
		assert:          assertNoError,
	},
	{
		name:            "IsNotEqualTo passes on not equal value",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue("actual"),
		constraint:      it.IsNotEqualTo("expected"),
		assert:          assertNoError,
	},
	{
		name:            "IsNotEqualTo violation on equal value",
		isApplicableFor: specificValueTypes(stringType, comparableType),
		stringValue:     stringValue("expected"),
		constraint:      it.IsNotEqualTo("expected"),
		assert:          assertHasOneViolation(validation.ErrIsEqual, `This value should not be equal to "expected".`),
	},
}

var isEarlierThanTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsEarlierThan passes on nil",
		isApplicableFor: specificValueTypes(timeType),
		constraint:      it.IsEarlierThan(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsEarlierThan violation on greater value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 40, 0, 0, time.UTC)),
		constraint:      it.IsEarlierThan(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertHasOneViolation(validation.ErrTooLate, "This value should be earlier than 2021-03-29T12:30:00Z."),
	},
	{
		name:            "IsEarlierThan violation on equal value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		constraint:      it.IsEarlierThan(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertHasOneViolation(validation.ErrTooLate, "This value should be earlier than 2021-03-29T12:30:00Z."),
	},
	{
		name:            "IsEarlierThan passes on less value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 29, 29, 0, time.UTC)),
		constraint:      it.IsEarlierThan(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsEarlierThan violation with custom message",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 40, 0, 0, time.UTC)),
		constraint: it.IsEarlierThan(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)).
			WithError(ErrCustom).
			WithMessage(
				`Unexpected value "{{ value }}" at {{ custom }}, expected value must be earlier than "{{ comparedValue }}".`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(
			ErrCustom,
			`Unexpected value "2021-03-29T12:40:00Z" at parameter, expected value must be earlier than "2021-03-29T12:30:00Z".`,
		),
	},
	{
		name:            "IsEarlierThan violation with custom message and time layout",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 40, 0, 0, time.UTC)),
		constraint: it.IsEarlierThan(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)).
			WithMessage(`Unexpected value "{{ value }}", expected value must be earlier than "{{ comparedValue }}".`).
			WithLayout(time.RFC822),
		assert: assertHasOneViolation(
			validation.ErrTooLate,
			`Unexpected value "29 Mar 21 12:40 UTC", expected value must be earlier than "29 Mar 21 12:30 UTC".`,
		),
	},
	{
		name:            "IsEarlierThan passes when condition is false",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 40, 0, 0, time.UTC)),
		constraint:      it.IsEarlierThan(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)).When(false),
		assert:          assertNoError,
	},
	{
		name:            "IsEarlierThan passes when groups not match",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 40, 0, 0, time.UTC)),
		constraint:      it.IsEarlierThan(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)).WhenGroups(testGroup),
		assert:          assertNoError,
	},
	{
		name:            "IsEarlierThan violation when condition is true",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 40, 0, 0, time.UTC)),
		constraint:      it.IsEarlierThan(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)).When(true),
		assert:          assertHasOneViolation(validation.ErrTooLate, "This value should be earlier than 2021-03-29T12:30:00Z."),
	},
}

var isEarlierThanOrEqualTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsEarlierThanOrEqual passes on nil",
		isApplicableFor: specificValueTypes(timeType),
		constraint:      it.IsEarlierThanOrEqual(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsEarlierThanOrEqual violation on greater value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 40, 0, 0, time.UTC)),
		constraint:      it.IsEarlierThanOrEqual(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertHasOneViolation(validation.ErrTooLateOrEqual, "This value should be earlier than or equal to 2021-03-29T12:30:00Z."),
	},
	{
		name:            "IsEarlierThanOrEqual passes on equal value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		constraint:      it.IsEarlierThanOrEqual(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsEarlierThanOrEqual passes on equal value with different time zone",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		constraint:      it.IsEarlierThanOrEqual(time.Date(2021, 0o3, 29, 8, 30, 0, 0, givenLocation("America/New_York"))),
		assert:          assertNoError,
	},
	{
		name:            "IsEarlierThanOrEqual passes on less value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 29, 29, 0, time.UTC)),
		constraint:      it.IsEarlierThanOrEqual(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
}

var isLaterThanTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLaterThan passes on nil",
		isApplicableFor: specificValueTypes(timeType),
		constraint:      it.IsLaterThan(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsLaterThan passes on greater value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 40, 0, 0, time.UTC)),
		constraint:      it.IsLaterThan(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsLaterThan violation on equal value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		constraint:      it.IsLaterThan(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertHasOneViolation(validation.ErrTooEarly, "This value should be later than 2021-03-29T12:30:00Z."),
	},
	{
		name:            "IsLaterThan violation on less value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 29, 29, 0, time.UTC)),
		constraint:      it.IsLaterThan(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertHasOneViolation(validation.ErrTooEarly, "This value should be later than 2021-03-29T12:30:00Z."),
	},
}

var isLaterThanOrEqualTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLaterThanOrEqual passes on nil",
		isApplicableFor: specificValueTypes(timeType),
		constraint:      it.IsLaterThanOrEqual(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsLaterThanOrEqual passes on greater value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 40, 0, 0, time.UTC)),
		constraint:      it.IsLaterThanOrEqual(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsLaterThanOrEqual passes on equal value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		constraint:      it.IsLaterThanOrEqual(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertNoError,
	},
	{
		name:            "IsLaterThanOrEqual passes on equal value with different time zone",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		constraint:      it.IsLaterThanOrEqual(time.Date(2021, 0o3, 29, 8, 30, 0, 0, givenLocation("America/New_York"))),
		assert:          assertNoError,
	},
	{
		name:            "IsLaterThanOrEqual violation on less value",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o3, 29, 12, 29, 29, 0, time.UTC)),
		constraint:      it.IsLaterThanOrEqual(time.Date(2021, 0o3, 29, 12, 30, 0, 0, time.UTC)),
		assert:          assertHasOneViolation(validation.ErrTooEarlyOrEqual, "This value should be later than or equal to 2021-03-29T12:30:00Z."),
	},
}

var isBetweenTimeTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsBetweenTime error on equal min and max",
		isApplicableFor: specificValueTypes(timeType),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
		),
		assert: assertError(`validate by TimeRangeConstraint: invalid range`),
	},
	{
		name:            "IsBetweenTime error on equal min and max in different time zones",
		isApplicableFor: specificValueTypes(timeType),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 0o4, 4, 8, 30, 0, 0, givenLocation("America/New_York"))),
		),
		assert: assertError(`validate by TimeRangeConstraint: invalid range`),
	},
	{
		name:            "IsBetweenTime error on min greater than max",
		isApplicableFor: specificValueTypes(timeType),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 0o4, 4, 12, 40, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
		),
		assert: assertError(`validate by TimeRangeConstraint: invalid range`),
	},
	{
		name:            "IsBetweenTime passes on nil",
		isApplicableFor: specificValueTypes(timeType),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 0o4, 4, 12, 40, 0, 0, time.UTC)),
		),
		assert: assertNoError,
	},
	{
		name:            "IsBetweenTime violation on value less than min",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o4, 4, 12, 20, 0, 0, time.UTC)),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 0o4, 4, 12, 40, 0, 0, time.UTC)),
		),
		assert: assertHasOneViolation(validation.ErrNotInRange, "This value should be between 2021-04-04T12:30:00Z and 2021-04-04T12:40:00Z."),
	},
	{
		name:            "IsBetweenTime violation on value greater than max",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o4, 4, 12, 50, 0, 0, time.UTC)),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 0o4, 4, 12, 40, 0, 0, time.UTC)),
		),
		assert: assertHasOneViolation(validation.ErrNotInRange, "This value should be between 2021-04-04T12:30:00Z and 2021-04-04T12:40:00Z."),
	},
	{
		name:            "IsBetweenTime passes on value equal to min",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 0o4, 4, 12, 40, 0, 0, time.UTC)),
		),
		assert: assertNoError,
	},
	{
		name:            "IsBetweenTime passes on value equal to max",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o4, 4, 12, 40, 0, 0, time.UTC)),
		constraint: it.IsBetweenTime(
			*timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
			*timeValue(time.Date(2021, 0o4, 4, 12, 40, 0, 0, time.UTC)),
		),
		assert: assertNoError,
	},
	{
		name:            "IsBetweenTime violation with custom message",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o4, 4, 12, 20, 0, 0, time.UTC)),
		constraint: it.
			IsBetweenTime(
				*timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
				*timeValue(time.Date(2021, 0o4, 4, 12, 40, 0, 0, time.UTC)),
			).
			WithError(ErrCustom).
			WithMessage(
				`Unexpected value "{{ value }}" at {{ custom }}, expected value must be between "{{ min }}" and "{{ max }}".`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		assert: assertHasOneViolation(
			ErrCustom,
			`Unexpected value "2021-04-04T12:20:00Z" at parameter, expected value must be between "2021-04-04T12:30:00Z" and "2021-04-04T12:40:00Z".`,
		),
	},
	{
		name:            "IsBetweenTime violation with custom message and time layout",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o4, 4, 12, 20, 0, 0, time.UTC)),
		constraint: it.
			IsBetweenTime(
				*timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
				*timeValue(time.Date(2021, 0o4, 4, 12, 40, 0, 0, time.UTC)),
			).
			WithMessage(`Unexpected value "{{ value }}", expected value must be between "{{ min }}" and "{{ max }}".`).
			WithLayout(time.RFC822),
		assert: assertHasOneViolation(
			validation.ErrNotInRange,
			`Unexpected value "04 Apr 21 12:20 UTC", expected value must be between "04 Apr 21 12:30 UTC" and "04 Apr 21 12:40 UTC".`,
		),
	},
	{
		name:            "IsBetweenTime passes when condition is false",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o4, 4, 12, 20, 0, 0, time.UTC)),
		constraint: it.
			IsBetweenTime(
				*timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
				*timeValue(time.Date(2021, 0o4, 4, 12, 40, 0, 0, time.UTC)),
			).
			When(false),
		assert: assertNoError,
	},
	{
		name:            "IsBetweenTime passes when groups not match",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o4, 4, 12, 20, 0, 0, time.UTC)),
		constraint: it.
			IsBetweenTime(
				*timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
				*timeValue(time.Date(2021, 0o4, 4, 12, 40, 0, 0, time.UTC)),
			).
			WhenGroups(testGroup),
		assert: assertNoError,
	},
	{
		name:            "IsBetweenTime violation when condition is true",
		isApplicableFor: specificValueTypes(timeType),
		timeValue:       timeValue(time.Date(2021, 0o4, 4, 12, 20, 0, 0, time.UTC)),
		constraint: it.
			IsBetweenTime(
				*timeValue(time.Date(2021, 0o4, 4, 12, 30, 0, 0, time.UTC)),
				*timeValue(time.Date(2021, 0o4, 4, 12, 40, 0, 0, time.UTC)),
			).
			When(true),
		assert: assertHasOneViolation(validation.ErrNotInRange, "This value should be between 2021-04-04T12:30:00Z and 2021-04-04T12:40:00Z."),
	},
}

var hasUniqueValuesTestCases = []ConstraintValidationTestCase{
	{
		name:            "HasUniqueValues passes on nil",
		isApplicableFor: specificValueTypes(stringsType),
		constraint:      it.HasUniqueValues[string](),
		assert:          assertNoError,
	},
	{
		name:            "HasUniqueValues passes on empty value",
		isApplicableFor: specificValueTypes(stringsType),
		constraint:      it.HasUniqueValues[string](),
		stringsValue:    []string{},
		assert:          assertNoError,
	},
	{
		name:            "HasUniqueValues passes on unique values",
		isApplicableFor: specificValueTypes(stringsType),
		constraint:      it.HasUniqueValues[string](),
		stringsValue:    []string{"one", "two", "three"},
		assert:          assertNoError,
	},
	{
		name:            "HasUniqueValues violation on duplicated values",
		isApplicableFor: specificValueTypes(stringsType),
		constraint:      it.HasUniqueValues[string](),
		stringsValue:    []string{"one", "two", "one"},
		assert:          assertHasOneViolation(validation.ErrNotUnique, message.NotUnique),
	},
	{
		name:            "HasUniqueValues violation with custom message",
		isApplicableFor: specificValueTypes(stringsType),
		constraint: it.HasUniqueValues[string]().
			WithError(ErrCustom).
			WithMessage(
				`Not unique values at {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringsValue: []string{"one", "two", "one"},
		assert:       assertHasOneViolation(ErrCustom, `Not unique values at parameter.`),
	},
	{
		name:            "HasUniqueValues passes when condition is false",
		isApplicableFor: specificValueTypes(stringsType),
		constraint:      it.HasUniqueValues[string]().When(false),
		stringsValue:    []string{"one", "two", "one"},
		assert:          assertNoError,
	},
	{
		name:            "HasUniqueValues passes when groups not match",
		isApplicableFor: specificValueTypes(stringsType),
		constraint:      it.HasUniqueValues[string]().WhenGroups(testGroup),
		stringsValue:    []string{"one", "two", "one"},
		assert:          assertNoError,
	},
	{
		name:            "HasUniqueValues violation when condition is true",
		isApplicableFor: specificValueTypes(stringsType),
		constraint:      it.HasUniqueValues[string]().When(true),
		stringsValue:    []string{"one", "two", "one"},
		assert:          assertHasOneViolation(validation.ErrNotUnique, message.NotUnique),
	},
}
