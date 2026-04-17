package test

import (
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
)

var identifierConstraintsTestCases = mergeTestCases(
	ulidConstraintTestCases,
	uuidConstraintTestCases,
	currencyConstraintTestCases,
	ibanConstraintTestCases,
	bicConstraintTestCases,
	isinConstraintTestCases,
	isbnConstraintTestCases,
	issnConstraintTestCases,
	luhnConstraintTestCases,
)

var currencyConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsCurrency passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCurrency(),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "IsCurrency passes on valid ISO code",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("EUR"),
		constraint:      it.IsCurrency(),
		assert:          assertNoError,
	},
	{
		name:            "IsCurrency passes on lowercase code",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("usd"),
		constraint:      it.IsCurrency(),
		assert:          assertNoError,
	},
	{
		name:            "IsCurrency violation on unknown code",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("ZZZ"),
		constraint:      it.IsCurrency(),
		assert:          assertHasOneViolation(validation.ErrInvalidCurrency, message.InvalidCurrency),
	},
	{
		name:            "IsCurrency violation on wrong length",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("EU"),
		constraint:      it.IsCurrency(),
		assert:          assertHasOneViolation(validation.ErrInvalidCurrency, message.InvalidCurrency),
	},
	{
		name:            "IsCurrency violation with custom error and message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsCurrency().
			WithError(ErrCustom).
			WithMessage(
				`Invalid value "{{ value }}" for {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("UUU"),
		assert:      assertHasOneViolation(ErrCustom, `Invalid value "UUU" for parameter.`),
	},
	{
		name:            "IsCurrency passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCurrency().When(false),
		stringValue:     stringValue("ZZZ"),
		assert:          assertNoError,
	},
	{
		name:            "IsCurrency violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCurrency().When(true),
		stringValue:     stringValue("ZZZ"),
		assert:          assertHasOneViolation(validation.ErrInvalidCurrency, message.InvalidCurrency),
	},
	{
		name:            "IsCurrency passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCurrency().WhenGroups(testGroup),
		stringValue:     stringValue("ZZZ"),
		assert:          assertNoError,
	},
}

var ulidConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsULID passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		constraint:      it.IsULID(),
		assert:          assertNoError,
	},
	{
		name:            "IsULID violation on invalid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("81ARZ3NDEKTSV4RRFFQ69G5FAV"),
		constraint:      it.IsULID(),
		assert:          assertHasOneViolation(validation.ErrInvalidULID, message.InvalidULID),
	},
}

var ibanConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsIBAN passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIBAN(),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "IsIBAN passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("DE89370400440532013000"),
		constraint:      it.IsIBAN(),
		assert:          assertNoError,
	},
	{
		name:            "IsIBAN passes on spaced value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("CH93 0076 2011 6238 5295 7"),
		constraint:      it.IsIBAN(),
		assert:          assertNoError,
	},
	{
		name:            "IsIBAN violation on invalid checksum",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("DE89370400440532013001"),
		constraint:      it.IsIBAN(),
		assert:          assertHasOneViolation(validation.ErrInvalidIBAN, message.InvalidIBAN),
	},
	{
		name:            "IsIBAN violation on unsupported country",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("US64SVBX1101057138"),
		constraint:      it.IsIBAN(),
		assert:          assertHasOneViolation(validation.ErrInvalidIBAN, message.InvalidIBAN),
	},
	{
		name:            "IsIBAN violation with given error and message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsIBAN().
			WithError(ErrCustom).
			WithMessage(
				`Invalid value "{{ value }}" for {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("bad-iban"),
		assert:      assertHasOneViolation(ErrCustom, `Invalid value "bad-iban" for parameter.`),
	},
	{
		name:            "IsIBAN passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIBAN().When(false),
		stringValue:     stringValue("bad"),
		assert:          assertNoError,
	},
	{
		name:            "IsIBAN violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIBAN().When(true),
		stringValue:     stringValue("bad"),
		assert:          assertHasOneViolation(validation.ErrInvalidIBAN, message.InvalidIBAN),
	},
	{
		name:            "IsIBAN passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIBAN().WhenGroups(testGroup),
		stringValue:     stringValue("bad"),
		assert:          assertNoError,
	},
}

var bicConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsBIC passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsBIC(),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "IsBIC passes on valid 8-char BIC",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("DEUTDEFF"),
		constraint:      it.IsBIC(),
		assert:          assertNoError,
	},
	{
		name:            "IsBIC passes when spaces stripped",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("DEUT DE FF"),
		constraint:      it.IsBIC(),
		assert:          assertNoError,
	},
	{
		name:            "IsBIC passes with CaseInsensitive on lowercase",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("deutdeff"),
		constraint:      it.IsBIC().CaseInsensitive(),
		assert:          assertNoError,
	},
	{
		name:            "IsBIC violation on strict lowercase",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("deutdeff"),
		constraint:      it.IsBIC(),
		assert:          assertHasOneViolation(validation.ErrInvalidBIC, message.InvalidBIC),
	},
	{
		name:            "IsBIC violation on wrong length",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("DEUTDEF"),
		constraint:      it.IsBIC(),
		assert:          assertHasOneViolation(validation.ErrInvalidBIC, message.InvalidBIC),
	},
	{
		name:            "IsBIC violation on unknown country",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("DEUTZZFF"),
		constraint:      it.IsBIC(),
		assert:          assertHasOneViolation(validation.ErrInvalidBIC, message.InvalidBIC),
	},
	{
		name:            "IsBIC violation on IBAN country mismatch",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("DEUTDEFF"),
		constraint:      it.IsBIC().WithIBAN("GB29NWBK60161331926819"),
		assert: assertHasOneViolation(
			validation.ErrBICIBANCountryMismatch,
			"This Business Identifier Code (BIC) is not associated with IBAN GB29NWBK60161331926819.",
		),
	},
	{
		name:            "IsBIC passes when IBAN country matches",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("DEUTDEFF"),
		constraint:      it.IsBIC().WithIBAN("DE89370400440532013000"),
		assert:          assertNoError,
	},
	{
		name:            "IsBIC violation with custom IBAN mismatch error and message",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("DEUTDEFF"),
		constraint: it.IsBIC().
			WithIBAN("GB29NWBK60161331926819").
			WithIBANError(ErrCustom).
			WithIBANMessage(`BIC "{{ value }}" is not linked to {{ iban }}.`),
		assert: assertHasOneViolation(
			ErrCustom,
			`BIC "DEUTDEFF" is not linked to GB29NWBK60161331926819.`,
		),
	},
	{
		name:            "IsBIC violation with given error and message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsBIC().
			WithError(ErrCustom).
			WithMessage(
				`Invalid value "{{ value }}" for {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("bad-bic"),
		assert:      assertHasOneViolation(ErrCustom, `Invalid value "bad-bic" for parameter.`),
	},
	{
		name:            "IsBIC passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsBIC().When(false),
		stringValue:     stringValue("bad"),
		assert:          assertNoError,
	},
	{
		name:            "IsBIC violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsBIC().When(true),
		stringValue:     stringValue("bad"),
		assert:          assertHasOneViolation(validation.ErrInvalidBIC, message.InvalidBIC),
	},
	{
		name:            "IsBIC passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsBIC().WhenGroups(testGroup),
		stringValue:     stringValue("bad"),
		assert:          assertNoError,
	},
}

var isinConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsISIN passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsISIN(),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "IsISIN passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("US0378331005"),
		constraint:      it.IsISIN(),
		assert:          assertNoError,
	},
	{
		name:            "IsISIN passes on lowercase letters",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("us0378331005"),
		constraint:      it.IsISIN(),
		assert:          assertNoError,
	},
	{
		name:            "IsISIN violation on wrong length",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("US037833100"),
		constraint:      it.IsISIN(),
		assert:          assertHasOneViolation(validation.ErrInvalidISIN, message.InvalidISIN),
	},
	{
		name:            "IsISIN violation on invalid pattern",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("123456789101"),
		constraint:      it.IsISIN(),
		assert:          assertHasOneViolation(validation.ErrInvalidISIN, message.InvalidISIN),
	},
	{
		name:            "IsISIN violation on invalid checksum",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("XS2012239364"),
		constraint:      it.IsISIN(),
		assert:          assertHasOneViolation(validation.ErrInvalidISIN, message.InvalidISIN),
	},
	{
		name:            "IsISIN violation with given error and message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsISIN().
			WithError(ErrCustom).
			WithMessage(
				`Invalid value "{{ value }}" for {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("invalid-isin"),
		assert:      assertHasOneViolation(ErrCustom, `Invalid value "invalid-isin" for parameter.`),
	},
	{
		name:            "IsISIN passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsISIN().When(false),
		stringValue:     stringValue("bad"),
		assert:          assertNoError,
	},
	{
		name:            "IsISIN violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsISIN().When(true),
		stringValue:     stringValue("bad"),
		assert:          assertHasOneViolation(validation.ErrInvalidISIN, message.InvalidISIN),
	},
	{
		name:            "IsISIN passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsISIN().WhenGroups(testGroup),
		stringValue:     stringValue("bad"),
		assert:          assertNoError,
	},
}

var isbnConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsISBN passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsISBN(),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "IsISBN passes on valid ISBN-10",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("0-45122-5244"),
		constraint:      it.IsISBN(),
		assert:          assertNoError,
	},
	{
		name:            "IsISBN passes on valid ISBN-13",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("978-0451225245"),
		constraint:      it.IsISBN(),
		assert:          assertNoError,
	},
	{
		name:            "IsISBN violation on neither format",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("978272344228"),
		constraint:      it.IsISBN(),
		assert:          assertHasOneViolation(validation.ErrInvalidISBN, message.InvalidISBN),
	},
	{
		name:            "IsISBN Only10 passes",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("0321812700"),
		constraint:      it.IsISBN().Only10(),
		assert:          assertNoError,
	},
	{
		name:            "IsISBN Only10 violation uses ISBN-10 message",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("1234567890"),
		constraint:      it.IsISBN().Only10(),
		assert:          assertHasOneViolation(validation.ErrInvalidISBN10, message.InvalidISBN10),
	},
	{
		name:            "IsISBN Only13 violation uses ISBN-13 message",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2723442284"),
		constraint:      it.IsISBN().Only13(),
		assert:          assertHasOneViolation(validation.ErrInvalidISBN13, message.InvalidISBN13),
	},
	{
		name:            "IsISBN violation with given error and message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsISBN().
			WithError(ErrCustom).
			WithMessage(
				`Invalid value "{{ value }}" for {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("bad-isbn"),
		assert:      assertHasOneViolation(ErrCustom, `Invalid value "bad-isbn" for parameter.`),
	},
	{
		name:            "IsISBN passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsISBN().When(false),
		stringValue:     stringValue("bad"),
		assert:          assertNoError,
	},
	{
		name:            "IsISBN violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsISBN().When(true),
		stringValue:     stringValue("bad"),
		assert:          assertHasOneViolation(validation.ErrInvalidISBN, message.InvalidISBN),
	},
	{
		name:            "IsISBN passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsISBN().WhenGroups(testGroup),
		stringValue:     stringValue("bad"),
		assert:          assertNoError,
	},
}

var issnConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsISSN passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsISSN(),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "IsISSN passes on valid value with hyphen",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("0317-8471"),
		constraint:      it.IsISSN(),
		assert:          assertNoError,
	},
	{
		name:            "IsISSN passes on valid value without hyphen",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("03178471"),
		constraint:      it.IsISSN(),
		assert:          assertNoError,
	},
	{
		name:            "IsISSN passes on check digit X",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2434-561X"),
		constraint:      it.IsISSN(),
		assert:          assertNoError,
	},
	{
		name:            "IsISSN passes on lowercase check digit x",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("2434-561x"),
		constraint:      it.IsISSN(),
		assert:          assertNoError,
	},
	{
		name:            "IsISSN violation on wrong length",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("0317-847"),
		constraint:      it.IsISSN(),
		assert:          assertHasOneViolation(validation.ErrInvalidISSN, message.InvalidISSN),
	},
	{
		name:            "IsISSN violation on invalid hyphen placement",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("123-45678"),
		constraint:      it.IsISSN(),
		assert:          assertHasOneViolation(validation.ErrInvalidISSN, message.InvalidISSN),
	},
	{
		name:            "IsISSN violation on invalid checksum",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("0317-8470"),
		constraint:      it.IsISSN(),
		assert:          assertHasOneViolation(validation.ErrInvalidISSN, message.InvalidISSN),
	},
	{
		name:            "IsISSN passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsISSN().When(false),
		stringValue:     stringValue("bad"),
		assert:          assertNoError,
	},
	{
		name:            "IsISSN violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsISSN().When(true),
		stringValue:     stringValue("bad"),
		assert:          assertHasOneViolation(validation.ErrInvalidISSN, message.InvalidISSN),
	},
	{
		name:            "IsISSN passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsISSN().WhenGroups(testGroup),
		stringValue:     stringValue("bad"),
		assert:          assertNoError,
	},
}

var luhnConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsLUHN passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsLUHN(),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "IsLUHN passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("79927398713"),
		constraint:      it.IsLUHN(),
		assert:          assertNoError,
	},
	{
		name:            "IsLUHN violation on invalid checksum",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("79927398710"),
		constraint:      it.IsLUHN(),
		assert:          assertHasOneViolation(validation.ErrInvalidLUHN, message.InvalidLUHN),
	},
	{
		name:            "IsLUHN violation on only zeros",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("0000000000000000"),
		constraint:      it.IsLUHN(),
		assert:          assertHasOneViolation(validation.ErrInvalidLUHN, message.InvalidLUHN),
	},
	{
		name:            "IsLUHN violation on non-digit",
		isApplicableFor: specificValueTypes(stringType),
		stringValue:     stringValue("12345a"),
		constraint:      it.IsLUHN(),
		assert:          assertHasOneViolation(validation.ErrInvalidLUHN, message.InvalidLUHN),
	},
	{
		name:            "IsLUHN violation with given error and message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsLUHN().
			WithError(ErrCustom).
			WithMessage(
				`Invalid value "{{ value }}" for {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("bad"),
		assert:      assertHasOneViolation(ErrCustom, `Invalid value "bad" for parameter.`),
	},
	{
		name:            "IsLUHN passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsLUHN().When(false),
		stringValue:     stringValue("79927398710"),
		assert:          assertNoError,
	},
	{
		name:            "IsLUHN violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsLUHN().When(true),
		stringValue:     stringValue("79927398710"),
		assert:          assertHasOneViolation(validation.ErrInvalidLUHN, message.InvalidLUHN),
	},
	{
		name:            "IsLUHN passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsLUHN().WhenGroups(testGroup),
		stringValue:     stringValue("79927398710"),
		assert:          assertNoError,
	},
}

var uuidConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsUUID passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID(),
		assert:          assertNoError,
	},
	{
		name:            "IsUUID passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID(),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "IsUUID passes on valid value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID(),
		stringValue:     stringValue("661eeca0-bc27-4ecc-8f69-6ffb7b1d5a92"),
		assert:          assertNoError,
	},
	{
		name:            "IsUUID violation on invalid value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID(),
		stringValue:     stringValue("invalid"),
		assert:          assertHasOneViolation(validation.ErrInvalidUUID, message.InvalidUUID),
	},
	{
		name:            "IsUUID passes on nil value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID(),
		stringValue:     stringValue("00000000-0000-0000-0000-000000000000"),
		assert:          assertNoError,
	},
	{
		name:            "IsUUID violation on nil value when not nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID().NotNil(),
		stringValue:     stringValue("00000000-0000-0000-0000-000000000000"),
		assert:          assertHasOneViolation(validation.ErrInvalidUUID, message.InvalidUUID),
	},
	{
		name:            "IsUUID violation on non-canonical value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID(),
		stringValue:     stringValue("661eeca0bc274ecc8f696ffb7b1d5a92"),
		assert:          assertHasOneViolation(validation.ErrInvalidUUID, message.InvalidUUID),
	},
	{
		name:            "IsUUID passes on non-canonical value when allowed",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID().NonCanonical(),
		stringValue:     stringValue("661eeca0bc274ecc8f696ffb7b1d5a92"),
		assert:          assertNoError,
	},
	{
		name:            "IsUUID violation on not allowed version",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID().WithVersions(1),
		stringValue:     stringValue("661eeca0-bc27-4ecc-8f69-6ffb7b1d5a92"),
		assert:          assertHasOneViolation(validation.ErrInvalidUUID, message.InvalidUUID),
	},
	{
		name:            "IsUUID violation with given error and message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsUUID().
			WithError(ErrCustom).
			WithMessage(
				`Invalid value "{{ value }}" for {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("invalid"),
		assert:      assertHasOneViolation(ErrCustom, `Invalid value "invalid" for parameter.`),
	},
	{
		name:            "IsUUID passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID().When(false),
		stringValue:     stringValue("invalid"),
		assert:          assertNoError,
	},
	{
		name:            "IsUUID violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID().When(true),
		stringValue:     stringValue("invalid"),
		assert:          assertHasOneViolation(validation.ErrInvalidUUID, message.InvalidUUID),
	},
	{
		name:            "IsUUID passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsUUID().WhenGroups(testGroup),
		stringValue:     stringValue("invalid"),
		assert:          assertNoError,
	},
}
