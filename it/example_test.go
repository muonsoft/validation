package it_test

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

func ExampleIsEAN8() {
	err := validator.Validate(context.Background(), validation.String("42345670", it.IsEAN8()))
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid EAN-8.
}

func ExampleIsEAN13() {
	err := validator.Validate(context.Background(), validation.String("4006381333932", it.IsEAN13()))
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid EAN-13.
}

func ExampleIsUPCA() {
	err := validator.Validate(context.Background(), validation.String("614141000037", it.IsUPCA()))
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid UPC-A.
}

func ExampleIsUPCE() {
	err := validator.Validate(context.Background(), validation.String("01234501", it.IsUPCE()))
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid UPC-E.
}

func ExampleIsNotBlank() {
	fmt.Println(validator.Validate(context.Background(), validation.String("", it.IsNotBlank())))
	fmt.Println(validator.Validate(context.Background(), validation.Countable(len([]string{}), it.IsNotBlank())))
	fmt.Println(validator.Validate(context.Background(), validation.Comparable[string]("", it.IsNotBlank())))
	// Output:
	// violation: This value should not be blank.
	// violation: This value should not be blank.
	// violation: This value should not be blank.
}

func ExampleIsNotBlankNumber() {
	fmt.Println(validator.Validate(context.Background(), validation.Number[int](0, it.IsNotBlankNumber[int]())))
	fmt.Println(validator.Validate(context.Background(), validation.Number[float64](0.0, it.IsNotBlankNumber[float64]())))
	// Output:
	// violation: This value should not be blank.
	// violation: This value should not be blank.
}

func ExampleIsNotBlankComparable() {
	fmt.Println(validator.Validate(context.Background(), validation.Comparable[int](0, it.IsNotBlankComparable[int]())))
	fmt.Println(validator.Validate(context.Background(), validation.Comparable[string]("", it.IsNotBlankComparable[string]())))
	// Output:
	// violation: This value should not be blank.
	// violation: This value should not be blank.
}

func ExampleIsBlank() {
	fmt.Println(validator.Validate(context.Background(), validation.String("foo", it.IsBlank())))
	fmt.Println(validator.Validate(context.Background(), validation.Countable(len([]string{"foo"}), it.IsBlank())))
	fmt.Println(validator.Validate(context.Background(), validation.Comparable[string]("foo", it.IsBlank())))
	// Output:
	// violation: This value should be blank.
	// violation: This value should be blank.
	// violation: This value should be blank.
}

func ExampleIsBlankNumber() {
	fmt.Println(validator.Validate(context.Background(), validation.Number[int](1, it.IsBlankNumber[int]())))
	fmt.Println(validator.Validate(context.Background(), validation.Number[float64](1.1, it.IsBlankNumber[float64]())))
	// Output:
	// violation: This value should be blank.
	// violation: This value should be blank.
}

func ExampleIsBlankComparable() {
	fmt.Println(validator.Validate(context.Background(), validation.Comparable[int](1, it.IsBlankComparable[int]())))
	fmt.Println(validator.Validate(context.Background(), validation.Comparable[string]("foo", it.IsBlankComparable[string]())))
	// Output:
	// violation: This value should be blank.
	// violation: This value should be blank.
}

func ExampleIsNotNil() {
	var s *string
	fmt.Println(validator.Validate(context.Background(), validation.Nil(s == nil, it.IsNotNil())))
	fmt.Println(validator.Validate(context.Background(), validation.NilString(s, it.IsNotNil())))
	fmt.Println(validator.Validate(context.Background(), validation.NilComparable[string](s, it.IsNotNil())))
	// Output:
	// violation: This value should not be nil.
	// violation: This value should not be nil.
	// violation: This value should not be nil.
}

func ExampleIsNotNilNumber() {
	var n *int
	var f *float64
	fmt.Println(validator.Validate(context.Background(), validation.NilNumber[int](n, it.IsNotNilNumber[int]())))
	fmt.Println(validator.Validate(context.Background(), validation.NilNumber[float64](f, it.IsNotNilNumber[float64]())))
	// Output:
	// violation: This value should not be nil.
	// violation: This value should not be nil.
}

func ExampleIsNotNilComparable() {
	var n *int
	var s *string
	fmt.Println(validator.Validate(context.Background(), validation.NilComparable[int](n, it.IsNotNilComparable[int]())))
	fmt.Println(validator.Validate(context.Background(), validation.NilComparable[string](s, it.IsNotNilComparable[string]())))
	// Output:
	// violation: This value should not be nil.
	// violation: This value should not be nil.
}

func ExampleIsNil() {
	s := ""
	sp := &s
	fmt.Println(validator.Validate(context.Background(), validation.Nil(sp == nil, it.IsNil())))
	fmt.Println(validator.Validate(context.Background(), validation.NilString(&s, it.IsNil())))
	fmt.Println(validator.Validate(context.Background(), validation.NilComparable[string](&s, it.IsNil())))
	// Output:
	// violation: This value should be nil.
	// violation: This value should be nil.
	// violation: This value should be nil.
}

func ExampleIsNilNumber() {
	n := 0
	f := 0.0
	fmt.Println(validator.Validate(context.Background(), validation.NilNumber[int](&n, it.IsNilNumber[int]())))
	fmt.Println(validator.Validate(context.Background(), validation.NilNumber[float64](&f, it.IsNilNumber[float64]())))
	// Output:
	// violation: This value should be nil.
	// violation: This value should be nil.
}

func ExampleIsNilComparable() {
	n := 0
	s := ""
	fmt.Println(validator.Validate(context.Background(), validation.NilComparable[int](&n, it.IsNilComparable[int]())))
	fmt.Println(validator.Validate(context.Background(), validation.NilComparable[string](&s, it.IsNilComparable[string]())))
	// Output:
	// violation: This value should be nil.
	// violation: This value should be nil.
}

func ExampleIsTrue() {
	err := validator.Validate(context.Background(), validation.Bool(false, it.IsTrue()))
	fmt.Println(err)
	// Output:
	// violation: This value should be true.
}

func ExampleIsFalse() {
	err := validator.Validate(context.Background(), validation.Bool(true, it.IsFalse()))
	fmt.Println(err)
	// Output:
	// violation: This value should be false.
}

func ExampleIsOneOf() {
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Comparable[string]("foo", it.IsOneOf("one", "two", "three"))),
	)
	fmt.Println(validator.Validate(
		context.Background(),
		validation.String("foo", it.IsOneOf("one", "two", "three"))),
	)
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Comparable[int](1, it.IsOneOf(2, 3, 4))),
	)
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[int](1, it.IsOneOf(2, 3, 4))),
	)
	// Output:
	// violation: The value you selected is not a valid choice.
	// violation: The value you selected is not a valid choice.
	// violation: The value you selected is not a valid choice.
	// violation: The value you selected is not a valid choice.
}

func ExampleIsEqualTo() {
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[int](1, it.IsEqualTo(2)),
	))
	fmt.Println(validator.Validate(
		context.Background(),
		validation.String("foo", it.IsEqualTo("bar")),
	))
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Comparable[string]("foo", it.IsEqualTo("bar")),
	))
	// Output:
	// violation: This value should be equal to 2.
	// violation: This value should be equal to "bar".
	// violation: This value should be equal to "bar".
}

func ExampleIsNotEqualTo() {
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[int](1, it.IsNotEqualTo(1)),
	))
	fmt.Println(validator.Validate(
		context.Background(),
		validation.String("foo", it.IsNotEqualTo("foo")),
	))
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Comparable[string]("foo", it.IsNotEqualTo("foo")),
	))
	// Output:
	// violation: This value should not be equal to 1.
	// violation: This value should not be equal to "foo".
	// violation: This value should not be equal to "foo".
}

func ExampleIsLessThan() {
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[int](1, it.IsLessThan(1))),
	)
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[float64](1.1, it.IsLessThan(1.1))),
	)
	// Output:
	// violation: This value should be less than 1.
	// violation: This value should be less than 1.1.
}

func ExampleIsLessThanOrEqual() {
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[int](1, it.IsLessThanOrEqual(0))),
	)
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[float64](1.1, it.IsLessThanOrEqual(0.1))),
	)
	// Output:
	// violation: This value should be less than or equal to 0.
	// violation: This value should be less than or equal to 0.1.
}

func ExampleIsGreaterThan() {
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[int](1, it.IsGreaterThan(1))),
	)
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[float64](1.1, it.IsGreaterThan(1.1))),
	)
	// Output:
	// violation: This value should be greater than 1.
	// violation: This value should be greater than 1.1.
}

func ExampleIsGreaterThanOrEqual() {
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[int](1, it.IsGreaterThanOrEqual(2))),
	)
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[float64](1.1, it.IsGreaterThanOrEqual(1.2))),
	)
	// Output:
	// violation: This value should be greater than or equal to 2.
	// violation: This value should be greater than or equal to 1.2.
}

func ExampleIsPositive() {
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[int](-1, it.IsPositive[int]())),
	)
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[float64](-1.1, it.IsPositive[float64]())),
	)
	// Output:
	// violation: This value should be positive.
	// violation: This value should be positive.
}

func ExampleIsPositiveOrZero() {
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[int](-1, it.IsPositiveOrZero[int]())),
	)
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[float64](-1.1, it.IsPositiveOrZero[float64]())),
	)
	// Output:
	// violation: This value should be either positive or zero.
	// violation: This value should be either positive or zero.
}

func ExampleIsNegative() {
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[int](1, it.IsNegative[int]())),
	)
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[float64](1.1, it.IsNegative[float64]())),
	)
	// Output:
	// violation: This value should be negative.
	// violation: This value should be negative.
}

func ExampleIsNegativeOrZero() {
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[int](1, it.IsNegativeOrZero[int]())),
	)
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[float64](1.1, it.IsNegativeOrZero[float64]())),
	)
	// Output:
	// violation: This value should be either negative or zero.
	// violation: This value should be either negative or zero.
}

func ExampleIsBetween() {
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[int](1, it.IsBetween(10, 20))),
	)
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Number[float64](1.1, it.IsBetween(10.111, 20.222))),
	)
	// Output:
	// violation: This value should be between 10 and 20.
	// violation: This value should be between 10.111 and 20.222.
}

func ExampleIsEarlierThan() {
	t, _ := time.Parse(time.RFC3339, "2009-02-04T21:00:57-08:00")
	t2, _ := time.Parse(time.RFC3339, "2009-02-03T21:00:57-08:00")
	err := validator.ValidateTime(context.Background(), t, it.IsEarlierThan(t2))
	fmt.Println(err)
	// Output:
	// violation: This value should be earlier than 2009-02-03T21:00:57-08:00.
}

func ExampleIsEarlierThanOrEqual() {
	t, _ := time.Parse(time.RFC3339, "2009-02-04T21:00:57-08:00")
	t2, _ := time.Parse(time.RFC3339, "2009-02-03T21:00:57-08:00")
	err := validator.ValidateTime(context.Background(), t, it.IsEarlierThanOrEqual(t2))
	fmt.Println(err)
	// Output:
	// violation: This value should be earlier than or equal to 2009-02-03T21:00:57-08:00.
}

func ExampleIsLaterThan() {
	t, _ := time.Parse(time.RFC3339, "2009-02-04T21:00:57-08:00")
	t2, _ := time.Parse(time.RFC3339, "2009-02-05T21:00:57-08:00")
	err := validator.ValidateTime(context.Background(), t, it.IsLaterThan(t2))
	fmt.Println(err)
	// Output:
	// violation: This value should be later than 2009-02-05T21:00:57-08:00.
}

func ExampleIsLaterThanOrEqual() {
	t, _ := time.Parse(time.RFC3339, "2009-02-04T21:00:57-08:00")
	t2, _ := time.Parse(time.RFC3339, "2009-02-05T21:00:57-08:00")
	err := validator.ValidateTime(context.Background(), t, it.IsLaterThanOrEqual(t2))
	fmt.Println(err)
	// Output:
	// violation: This value should be later than or equal to 2009-02-05T21:00:57-08:00.
}

func ExampleIsBetweenTime() {
	t, _ := time.Parse(time.RFC3339, "2009-02-04T21:00:57-08:00")
	after, _ := time.Parse(time.RFC3339, "2009-02-05T21:00:57-08:00")
	before, _ := time.Parse(time.RFC3339, "2009-02-06T21:00:57-08:00")
	err := validator.ValidateTime(context.Background(), t, it.IsBetweenTime(after, before))
	fmt.Println(err)
	// Output:
	// violation: This value should be between 2009-02-05T21:00:57-08:00 and 2009-02-06T21:00:57-08:00.
}

func ExampleHasUniqueValues() {
	strings := []string{"foo", "bar", "baz", "foo"}
	ints := []int{1, 2, 3, 1}
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Comparables[string](strings, it.HasUniqueValues[string]()),
	))
	fmt.Println(validator.Validate(
		context.Background(),
		validation.Comparables[int](ints, it.HasUniqueValues[int]()),
	))
	// Output:
	// violation: This collection should contain only unique elements.
	// violation: This collection should contain only unique elements.
}

func ExampleHasMinCount() {
	v := []int{1, 2}
	err := validator.ValidateCountable(context.Background(), len(v), it.HasMinCount(3))
	fmt.Println(err)
	// Output:
	// violation: This collection should contain 3 elements or more.
}

func ExampleHasMaxCount() {
	v := []int{1, 2}
	err := validator.ValidateCountable(context.Background(), len(v), it.HasMaxCount(1))
	fmt.Println(err)
	// Output:
	// violation: This collection should contain 1 element or less.
}

func ExampleHasCountBetween() {
	v := []int{1, 2}
	err := validator.ValidateCountable(context.Background(), len(v), it.HasCountBetween(3, 10))
	fmt.Println(err)
	// Output:
	// violation: This collection should contain 3 elements or more.
}

func ExampleHasExactCount() {
	v := []int{1, 2}
	err := validator.ValidateCountable(context.Background(), len(v), it.HasExactCount(3))
	fmt.Println(err)
	// Output:
	// violation: This collection should contain exactly 3 elements.
}

func ExampleHasMinLength() {
	v := "foo"
	err := validator.ValidateString(context.Background(), v, it.HasMinLength(5))
	fmt.Println(err)
	// Output:
	// violation: This value is too short. It should have 5 characters or more.
}

func ExampleHasMaxLength() {
	v := "foo"
	err := validator.ValidateString(context.Background(), v, it.HasMaxLength(2))
	fmt.Println(err)
	// Output:
	// violation: This value is too long. It should have 2 characters or less.
}

func ExampleHasLengthBetween() {
	v := "foo"
	err := validator.ValidateString(context.Background(), v, it.HasLengthBetween(5, 10))
	fmt.Println(err)
	// Output:
	// violation: This value is too short. It should have 5 characters or more.
}

func ExampleHasExactLength() {
	v := "foo"
	err := validator.ValidateString(context.Background(), v, it.HasExactLength(5))
	fmt.Println(err)
	// Output:
	// violation: This value should have exactly 5 characters.
}

func ExampleMatches() {
	v := "foo123"
	err := validator.ValidateString(context.Background(), v, it.Matches(regexp.MustCompile("^[a-z]+$")))
	fmt.Println(err)
	// Output:
	// violation: This value is not valid.
}

func ExampleDoesNotMatch() {
	v := "foo"
	err := validator.ValidateString(context.Background(), v, it.DoesNotMatch(regexp.MustCompile("^[a-z]+$")))
	fmt.Println(err)
	// Output:
	// violation: This value is not valid.
}

func ExampleIsJSON_validJSON() {
	v := `{"valid": true}`
	err := validator.Validate(context.Background(), validation.String(v, it.IsJSON()))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsJSON_invalidJSON() {
	v := `"invalid": true`
	err := validator.Validate(context.Background(), validation.String(v, it.IsJSON()))
	fmt.Println(err)
	// Output:
	// violation: This value should be valid JSON.
}

func ExampleIsInteger_validInteger() {
	v := "123"
	err := validator.Validate(context.Background(), validation.String(v, it.IsInteger()))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsInteger_invalidInteger() {
	v := "foo"
	err := validator.Validate(context.Background(), validation.String(v, it.IsInteger()))
	fmt.Println(err)
	// Output:
	// violation: This value is not an integer.
}

func ExampleIsNumeric_validNumeric() {
	v := "123.123"
	err := validator.Validate(context.Background(), validation.String(v, it.IsNumeric()))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsNumeric_invalidNumeric() {
	v := "foo.bar"
	err := validator.Validate(context.Background(), validation.String(v, it.IsNumeric()))
	fmt.Println(err)
	// Output:
	// violation: This value is not a numeric.
}

func ExampleIsEmail_validEmail() {
	v := "user@example.com"
	err := validator.Validate(context.Background(), validation.String(v, it.IsEmail()))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsEmail_invalidEmail() {
	v := "user example.com"
	err := validator.Validate(context.Background(), validation.String(v, it.IsEmail()))
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid email address.
}

func ExampleIsHTML5Email_validEmail() {
	v := "{}~!@example.com"
	err := validator.Validate(context.Background(), validation.String(v, it.IsEmail()))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsHTML5Email_invalidEmail() {
	v := "@example.com"
	err := validator.Validate(context.Background(), validation.String(v, it.IsEmail()))
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid email address.
}

func ExampleIsHostname_validHostname() {
	v := "example.com"
	err := validator.Validate(context.Background(), validation.String(v, it.IsHostname()))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsHostname_invalidHostname() {
	v := "example-.com"
	err := validator.Validate(context.Background(), validation.String(v, it.IsHostname()))
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid hostname.
}

func ExampleIsHostname_reservedHostname() {
	v := "example.localhost"
	err := validator.Validate(context.Background(), validation.String(v, it.IsHostname()))
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid hostname.
}

func ExampleIsLooseHostname_validHostname() {
	v := "example.com"
	err := validator.Validate(context.Background(), validation.String(v, it.IsLooseHostname()))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsLooseHostname_invalidHostname() {
	v := "example-.com"
	err := validator.Validate(context.Background(), validation.String(v, it.IsLooseHostname()))
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid hostname.
}

func ExampleIsLooseHostname_reservedHostname() {
	v := "example.localhost"
	err := validator.Validate(context.Background(), validation.String(v, it.IsLooseHostname()))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsURL_validURL() {
	v := "http://example.com"
	err := validator.Validate(context.Background(), validation.String(v, it.IsURL()))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsURL_invalidURL() {
	v := "example.com"
	err := validator.Validate(context.Background(), validation.String(v, it.IsURL()))
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid URL.
}

func ExampleURLConstraint_WithRelativeSchema() {
	v := "//example.com"
	err := validator.Validate(context.Background(), validation.String(v, it.IsURL().WithRelativeSchema()))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleURLConstraint_WithSchemas() {
	v := "ftp://example.com"
	err := validator.Validate(context.Background(), validation.String(v, it.IsURL().WithSchemas("http", "https", "ftp")))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsIP_validIP() {
	v := "123.123.123.123"
	err := validator.Validate(context.Background(), validation.String(v, it.IsIP()))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsIP_invalidIP() {
	v := "123.123.123.345"
	err := validator.Validate(context.Background(), validation.String(v, it.IsIP()))
	fmt.Println(err)
	// Output:
	// violation: This is not a valid IP address.
}

func ExampleIsIPv4_validIP() {
	v := "123.123.123.123"
	err := validator.Validate(context.Background(), validation.String(v, it.IsIPv4()))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsIPv4_invalidIP() {
	v := "123.123.123.345"
	err := validator.Validate(context.Background(), validation.String(v, it.IsIPv4()))
	fmt.Println(err)
	// Output:
	// violation: This is not a valid IP address.
}

func ExampleIsIPv6_validIP() {
	v := "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
	err := validator.Validate(context.Background(), validation.String(v, it.IsIPv6()))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsIPv6_invalidIP() {
	v := "z001:0db8:85a3:0000:0000:8a2e:0370:7334"
	err := validator.Validate(context.Background(), validation.String(v, it.IsIPv6()))
	fmt.Println(err)
	// Output:
	// violation: This is not a valid IP address.
}

func ExampleIPConstraint_DenyPrivateIP_restrictedPrivateIPv4() {
	v := "192.168.1.0"
	err := validator.Validate(context.Background(), validation.String(v, it.IsIP().DenyPrivateIP()))
	fmt.Println(err)
	// Output:
	// violation: This IP address is prohibited to use.
}

func ExampleIPConstraint_DenyPrivateIP_restrictedPrivateIPv6() {
	v := "fdfe:dcba:9876:ffff:fdc6:c46b:bb8f:7d4c"
	err := validator.Validate(context.Background(), validation.String(v, it.IsIPv6().DenyPrivateIP()))
	fmt.Println(err)
	// Output:
	// violation: This IP address is prohibited to use.
}

func ExampleIPConstraint_DenyIP() {
	v := "127.0.0.1"
	err := validator.Validate(
		context.Background(),
		validation.String(
			v,
			it.IsIP().DenyIP(func(ip net.IP) bool {
				return ip.IsLoopback()
			}),
		),
	)
	fmt.Println(err)
	// Output:
	// violation: This IP address is prohibited to use.
}
