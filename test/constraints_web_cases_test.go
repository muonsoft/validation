package test

import (
	"net"
	"net/url"
	"regexp"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/message"
	"github.com/muonsoft/validation/validate"
)

var urlConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsURL passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL(),
		assert:          assertNoError,
	},
	{
		name:            "IsURL passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL(),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "IsURL passes on valid URL",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL(),
		stringValue:     stringValue("http://example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsURL violation on invalid URL",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL(),
		stringValue:     stringValue("example.com"),
		assert:          assertHasOneViolation(validation.ErrInvalidURL, message.InvalidURL),
	},
	{
		name:            "IsURL error on empty schemas",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL().WithSchemas(),
		stringValue:     stringValue(""),
		assert:          assertError(`validate by URLConstraint: empty list of schemas`),
	},
	{
		name:            "IsURL passes on valid URL with custom schema",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL().WithSchemas("ftp"),
		stringValue:     stringValue("ftp://example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsURL with relative schema passes on valid relative URL",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL().WithRelativeSchema(),
		stringValue:     stringValue("//example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsURL with relative schema passes on valid absolute URL",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL().WithRelativeSchema(),
		stringValue:     stringValue("https://example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsURL passes on allowed host",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL().WithHosts("example.com"),
		stringValue:     stringValue("https://example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsURL violation on disallowed host",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL().WithHosts("example.com"),
		stringValue:     stringValue("https://sample.com"),
		assert:          assertHasOneViolation(validation.ErrProhibitedURL, message.ProhibitedURL),
	},
	{
		name:            "IsURL passes on allowed host pattern",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL().WithHostMatches(regexp.MustCompile(`^.*\.example.com$`)),
		stringValue:     stringValue("https://sub.example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsURL violation on disallowed host pattern",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL().WithHostMatches(regexp.MustCompile(`^.*\.example.com$`)),
		stringValue:     stringValue("https://sample.com"),
		assert:          assertHasOneViolation(validation.ErrProhibitedURL, message.ProhibitedURL),
	},
	{
		name:            "IsURL passes on not restricted value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL().WithRestriction(func(u *url.URL) bool { return true }),
		stringValue:     stringValue("https://example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsURL violation on restricted url",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL().WithRestriction(func(u *url.URL) bool { return false }),
		stringValue:     stringValue("https://example.com"),
		assert:          assertHasOneViolation(validation.ErrProhibitedURL, message.ProhibitedURL),
	},
	{
		name:            "IsURL violation on invalid URL with custom message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsURL().
			WithError(ErrCustom).
			WithMessage(
				`Unexpected URL "{{ value }}" at {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("example.com"),
		assert:      assertHasOneViolation(ErrCustom, `Unexpected URL "example.com" at parameter.`),
	},
	{
		name:            "IsURL violation on disallowed host with custom message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsURL().
			WithHosts("example.com").
			WithProhibitedError(ErrCustom).
			WithProhibitedMessage(
				`Prohibited URL "{{ value }}" at {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("https://sample.com"),
		assert:      assertHasOneViolation(ErrCustom, `Prohibited URL "https://sample.com" at parameter.`),
	},
	{
		name:            "IsURL passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL().When(false),
		stringValue:     stringValue("example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsURL passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL().WhenGroups(testGroup),
		stringValue:     stringValue("example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsURL violation when condition is true",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsURL().When(true),
		stringValue:     stringValue("example.com"),
		assert:          assertHasOneViolation(validation.ErrInvalidURL, message.InvalidURL),
	},
}

var emailConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsEmail passes on valid email",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsEmail(),
		stringValue:     stringValue("user@example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsEmail violation on invalid email",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsEmail(),
		stringValue:     stringValue("invalid"),
		assert:          assertHasOneViolation(validation.ErrInvalidEmail, message.InvalidEmail),
	},
	{
		name:            "IsHTML5Email passes on valid email",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsHTML5Email(),
		stringValue:     stringValue("user@example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsHTML5Email violation on invalid email",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsHTML5Email(),
		stringValue:     stringValue("invalid"),
		assert:          assertHasOneViolation(validation.ErrInvalidEmail, message.InvalidEmail),
	},
}

var ipConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsIP passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIP(),
		assert:          assertNoError,
	},
	{
		name:            "IsIP passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIP(),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "IsIP passes on valid IP v4",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIP(),
		stringValue:     stringValue("123.123.123.123"),
		assert:          assertNoError,
	},
	{
		name:            "IsIP violation on invalid IP v4",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIP(),
		stringValue:     stringValue("123.123.123.321"),
		assert:          assertHasOneViolation(validation.ErrInvalidIP, message.InvalidIP),
	},
	{
		name:            "IsIPv4 passes on valid IP v4",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIPv4(),
		stringValue:     stringValue("123.123.123.123"),
		assert:          assertNoError,
	},
	{
		name:            "IsIPv4 violation on IP v6",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIPv4(),
		stringValue:     stringValue("2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
		assert:          assertHasOneViolation(validation.ErrInvalidIP, message.InvalidIP),
	},
	{
		name:            "IsIPv6 passes on valid IP v6",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIPv6(),
		stringValue:     stringValue("2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
		assert:          assertNoError,
	},
	{
		name:            "IsIPv6 violation on IP v4",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIPv6(),
		stringValue:     stringValue("123.123.123.123"),
		assert:          assertHasOneViolation(validation.ErrInvalidIP, message.InvalidIP),
	},
	{
		name:            "IsIP violation on private IP",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIP().DenyPrivateIP(),
		stringValue:     stringValue("192.168.1.0"),
		assert:          assertHasOneViolation(validation.ErrProhibitedIP, message.ProhibitedIP),
	},
	{
		name:            "IsIP violation on custom IP",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsIP().DenyIP(func(ip net.IP) bool {
			return ip.IsLoopback()
		}),
		stringValue: stringValue("127.0.0.1"),
		assert:      assertHasOneViolation(validation.ErrProhibitedIP, message.ProhibitedIP),
	},
	{
		name:            "IsIP violation with custom message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsIP().
			WithInvalidError(ErrCustom).
			WithInvalidMessage(
				`Unexpected IP "{{ value }}" at {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("123.123.123.321"),
		assert: assertHasOneViolation(
			ErrCustom,
			`Unexpected IP "123.123.123.321" at parameter.`,
		),
	},
	{
		name:            "IsIP violation with custom restricted message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsIP().
			DenyPrivateIP().
			WithProhibitedError(ErrCustom).
			WithProhibitedMessage(
				`Unexpected IP "{{ value }}" at {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "parameter"},
			),
		stringValue: stringValue("192.168.1.0"),
		assert: assertHasOneViolation(
			ErrCustom,
			`Unexpected IP "192.168.1.0" at parameter.`,
		),
	},
	{
		name:            "IsIP passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIP().When(false),
		stringValue:     stringValue("123.123.123.321"),
		assert:          assertNoError,
	},
	{
		name:            "IsIP passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIP().WhenGroups(testGroup),
		stringValue:     stringValue("123.123.123.321"),
		assert:          assertNoError,
	},
	{
		name:            "IsIP passes when condition is false",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsIP().When(true),
		stringValue:     stringValue("123.123.123.321"),
		assert:          assertHasOneViolation(validation.ErrInvalidIP, message.InvalidIP),
	},
}

var cidrConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsCIDR passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCIDR(),
		assert:          assertNoError,
	},
	{
		name:            "IsCIDR passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCIDR(),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "IsCIDR passes on valid IPv4 CIDR",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCIDR(),
		stringValue:     stringValue("192.168.0.0/24"),
		assert:          assertNoError,
	},
	{
		name:            "IsCIDR passes on valid IPv6 CIDR",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCIDR(),
		stringValue:     stringValue("2001:db8::/32"),
		assert:          assertNoError,
	},
	{
		name:            "IsCIDR violation on missing slash",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCIDR(),
		stringValue:     stringValue("192.168.0.0"),
		assert:          assertHasOneViolation(validation.ErrInvalidCIDR, message.InvalidCIDR),
	},
	{
		name:            "IsCIDR violation on invalid prefix for IPv4",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCIDR(),
		stringValue:     stringValue("10.0.0.0/33"),
		assert: assertHasOneViolation(
			validation.ErrCIDRNetmaskOutOfRange,
			`The value of the netmask should be between 0 and 32.`,
		),
	},
	{
		name:            "IsCIDR violation on IPv6 when IPv4 only",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCIDR().IPv4Only(),
		stringValue:     stringValue("2001:db8::/32"),
		assert:          assertHasOneViolation(validation.ErrInvalidCIDR, message.InvalidCIDR),
	},
	{
		name:            "IsCIDR violation on IPv4 when IPv6 only",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCIDR().IPv6Only(),
		stringValue:     stringValue("10.0.0.0/8"),
		assert:          assertHasOneViolation(validation.ErrInvalidCIDR, message.InvalidCIDR),
	},
	{
		name:            "IsCIDR violation on prefix below WithNetmaskRange",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCIDR().WithNetmaskRange(16, 32),
		stringValue:     stringValue("10.0.0.0/8"),
		assert: assertHasOneViolation(
			validation.ErrCIDRNetmaskOutOfRange,
			`The value of the netmask should be between 16 and 32.`,
		),
	},
	{
		name:            "IsCIDR passes when When(false)",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCIDR().When(false),
		stringValue:     stringValue("not-cidr"),
		assert:          assertNoError,
	},
	{
		name:            "IsCIDR passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsCIDR().WhenGroups(testGroup),
		stringValue:     stringValue("not-cidr"),
		assert:          assertNoError,
	},
	{
		name:            "IsCIDR violation with custom invalid error and message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsCIDR().
			WithInvalidError(ErrCustom).
			WithInvalidMessage(
				`Bad CIDR "{{ value }}" at {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "field"},
			),
		stringValue: stringValue("192.168.0.0"),
		assert: assertHasOneViolation(
			ErrCustom,
			`Bad CIDR "192.168.0.0" at field.`,
		),
	},
}

var macAddressConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsMacAddress passes on nil",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsMacAddress(),
		assert:          assertNoError,
	},
	{
		name:            "IsMacAddress passes on empty value",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsMacAddress(),
		stringValue:     stringValue(""),
		assert:          assertNoError,
	},
	{
		name:            "IsMacAddress passes on colon form",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsMacAddress(),
		stringValue:     stringValue("00:1a:2b:3c:4d:5e"),
		assert:          assertNoError,
	},
	{
		name:            "IsMacAddress passes on hyphen form",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsMacAddress(),
		stringValue:     stringValue("00-1a-2b-3c-4d-5e"),
		assert:          assertNoError,
	},
	{
		name:            "IsMacAddress passes on dot form",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsMacAddress(),
		stringValue:     stringValue("001a.2b3c.4d5e"),
		assert:          assertNoError,
	},
	{
		name:            "IsMacAddress violation on invalid",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsMacAddress(),
		stringValue:     stringValue("not-mac"),
		assert:          assertHasOneViolation(validation.ErrInvalidMAC, message.InvalidMAC),
	},
	{
		name:            "IsMacAddress violation on EUI-64 length",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsMacAddress(),
		stringValue:     stringValue("02:00:5e:10:00:00:00:01"),
		assert:          assertHasOneViolation(validation.ErrInvalidMAC, message.InvalidMAC),
	},
	{
		name:            "IsMacAddress WithType broadcast only passes broadcast",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsMacAddress().WithType(validate.MacAddressTypeBroadcast),
		stringValue:     stringValue("ff:ff:ff:ff:ff:ff"),
		assert:          assertNoError,
	},
	{
		name:            "IsMacAddress WithType broadcast rejects unicast",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsMacAddress().WithType(validate.MacAddressTypeBroadcast),
		stringValue:     stringValue("00:1a:2b:3c:4d:5e"),
		assert:          assertHasOneViolation(validation.ErrInvalidMAC, message.InvalidMAC),
	},
	{
		name:            "IsMacAddress WithType all_no_broadcast rejects broadcast",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsMacAddress().WithType(validate.MacAddressTypeAllNoBroadcast),
		stringValue:     stringValue("ff:ff:ff:ff:ff:ff"),
		assert:          assertHasOneViolation(validation.ErrInvalidMAC, message.InvalidMAC),
	},
	{
		name:            "IsMacAddress passes when When(false)",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsMacAddress().When(false),
		stringValue:     stringValue("bad"),
		assert:          assertNoError,
	},
	{
		name:            "IsMacAddress passes when groups not match",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsMacAddress().WhenGroups(testGroup),
		stringValue:     stringValue("bad"),
		assert:          assertNoError,
	},
	{
		name:            "IsMacAddress violation with custom error and message",
		isApplicableFor: specificValueTypes(stringType),
		constraint: it.IsMacAddress().
			WithError(ErrCustom).
			WithMessage(
				`Bad MAC "{{ value }}" at {{ custom }}.`,
				validation.TemplateParameter{Key: "{{ custom }}", Value: "field"},
			),
		stringValue: stringValue("xx"),
		assert: assertHasOneViolation(
			ErrCustom,
			`Bad MAC "xx" at field.`,
		),
	},
}

var hostnameConstraintTestCases = []ConstraintValidationTestCase{
	{
		name:            "IsHostname passes on valid hostname",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsHostname(),
		stringValue:     stringValue("example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsHostname violation on invalid hostname",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsHostname(),
		stringValue:     stringValue("example-.com"),
		assert:          assertHasOneViolation(validation.ErrInvalidHostname, message.InvalidHostname),
	},
	{
		name:            "IsHostname violation on reserved hostname",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsHostname(),
		stringValue:     stringValue("example.localhost"),
		assert:          assertHasOneViolation(validation.ErrInvalidHostname, message.InvalidHostname),
	},
	{
		name:            "IsLooseHostname passes on valid hostname",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsLooseHostname(),
		stringValue:     stringValue("example.com"),
		assert:          assertNoError,
	},
	{
		name:            "IsLooseHostname violation on invalid hostname",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsLooseHostname(),
		stringValue:     stringValue("example-.com"),
		assert:          assertHasOneViolation(validation.ErrInvalidHostname, message.InvalidHostname),
	},
	{
		name:            "IsLooseHostname passes on reserved hostname",
		isApplicableFor: specificValueTypes(stringType),
		constraint:      it.IsLooseHostname(),
		stringValue:     stringValue("example.localhost"),
		assert:          assertNoError,
	},
}
