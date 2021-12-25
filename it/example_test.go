package it_test

import (
	"context"
	"fmt"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

func ExampleIsNotBlank() {
	v := ""
	err := validator.Validate(context.Background(), validation.String(v, it.IsNotBlank[string]()))
	fmt.Println(err)
	// Output:
	// violation: This value should not be blank.
}

// func ExampleIsBlank() {
// 	v := "foo"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsBlank()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be blank.
// }
//
// func ExampleIsNotNil() {
// 	var s *string
// 	err := validator.Validate(context.Background(), validation.NilString(s, it.IsNotNil()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should not be nil.
// }
//
// func ExampleIsNil() {
// 	s := ""
// 	err := validator.Validate(context.Background(), validation.NilString(&s, it.IsNil()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be nil.
// }
//
// func ExampleIsTrue() {
// 	b := false
// 	err := validator.Validate(context.Background(), validation.Bool(b, it.IsTrue()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be true.
// }
//
// func ExampleIsFalse() {
// 	b := true
// 	err := validator.Validate(context.Background(), validation.Bool(b, it.IsFalse()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be false.
// }
//
// func ExampleIsOneOfStrings() {
// 	s := "foo"
// 	err := validator.ValidateString(context.Background(), s, it.IsOneOfStrings("one", "two", "three"))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: The value you selected is not a valid choice.
// }
//
// func ExampleIsEqualToInteger() {
// 	v := 1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsEqualToNumber(2))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be equal to 2.
// }
//
// func ExampleIsEqualToFloat() {
// 	v := 1.1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsEqualToFloat(1.2))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be equal to 1.2.
// }
//
// func ExampleIsNotEqualToInteger() {
// 	v := 1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsNotEqualToNumber(1))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should not be equal to 1.
// }
//
// func ExampleIsNotEqualToFloat() {
// 	v := 1.1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsNotEqualToFloat(1.1))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should not be equal to 1.1.
// }
//
// func ExampleIsLessThanInteger() {
// 	v := 1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsLessThan(1))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be less than 1.
// }
//
// func ExampleIsLessThanFloat() {
// 	v := 1.1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsLessThanFloat(1.1))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be less than 1.1.
// }
//
// func ExampleIsLessThanOrEqualInteger() {
// 	v := 1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsLessThanOrEqual(0))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be less than or equal to 0.
// }
//
// func ExampleIsLessThanOrEqualFloat() {
// 	v := 1.1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsLessThanOrEqualFloat(0.1))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be less than or equal to 0.1.
// }
//
// func ExampleIsGreaterThanInteger() {
// 	v := 1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsGreaterThan(1))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be greater than 1.
// }
//
// func ExampleIsGreaterThanFloat() {
// 	v := 1.1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsGreaterThanFloat(1.1))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be greater than 1.1.
// }
//
// func ExampleIsGreaterThanOrEqualInteger() {
// 	v := 1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsGreaterThanOrEqual(2))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be greater than or equal to 2.
// }
//
// func ExampleIsGreaterThanOrEqualFloat() {
// 	v := 1.1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsGreaterThanOrEqualFloat(1.2))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be greater than or equal to 1.2.
// }
//
// func ExampleIsPositive() {
// 	v := -1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsPositive())
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be positive.
// }
//
// func ExampleIsPositiveOrZero() {
// 	v := -1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsPositiveOrZero())
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be either positive or zero.
// }
//
// func ExampleIsNegative() {
// 	v := 1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsNegative())
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be negative.
// }
//
// func ExampleIsNegativeOrZero() {
// 	v := 1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsNegativeOrZero())
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be either negative or zero.
// }
//
// func ExampleIsBetweenIntegers() {
// 	v := 1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsBetween(10, 20))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be between 10 and 20.
// }
//
// func ExampleIsBetweenFloats() {
// 	v := 1.1
// 	err := validator.ValidateNumber(context.Background(), v, it.IsBetweenFloats(10.111, 20.222))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be between 10.111 and 20.222.
// }
//
// func ExampleIsEqualToString() {
// 	v := "foo"
// 	err := validator.ValidateString(context.Background(), v, it.IsEqualToString("bar"))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be equal to "bar".
// }
//
// func ExampleIsNotEqualToString() {
// 	v := "foo"
// 	err := validator.ValidateString(context.Background(), v, it.IsNotEqualToString("foo"))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should not be equal to "foo".
// }
//
// func ExampleIsEarlierThan() {
// 	t, _ := time.Parse(time.RFC3339, "2009-02-04T21:00:57-08:00")
// 	t2, _ := time.Parse(time.RFC3339, "2009-02-03T21:00:57-08:00")
// 	err := validator.ValidateTime(context.Background(), t, it.IsEarlierThan(t2))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be earlier than 2009-02-03T21:00:57-08:00.
// }
//
// func ExampleIsEarlierThanOrEqual() {
// 	t, _ := time.Parse(time.RFC3339, "2009-02-04T21:00:57-08:00")
// 	t2, _ := time.Parse(time.RFC3339, "2009-02-03T21:00:57-08:00")
// 	err := validator.ValidateTime(context.Background(), t, it.IsEarlierThanOrEqual(t2))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be earlier than or equal to 2009-02-03T21:00:57-08:00.
// }
//
// func ExampleIsLaterThan() {
// 	t, _ := time.Parse(time.RFC3339, "2009-02-04T21:00:57-08:00")
// 	t2, _ := time.Parse(time.RFC3339, "2009-02-05T21:00:57-08:00")
// 	err := validator.ValidateTime(context.Background(), t, it.IsLaterThan(t2))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be later than 2009-02-05T21:00:57-08:00.
// }
//
// func ExampleIsLaterThanOrEqual() {
// 	t, _ := time.Parse(time.RFC3339, "2009-02-04T21:00:57-08:00")
// 	t2, _ := time.Parse(time.RFC3339, "2009-02-05T21:00:57-08:00")
// 	err := validator.ValidateTime(context.Background(), t, it.IsLaterThanOrEqual(t2))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be later than or equal to 2009-02-05T21:00:57-08:00.
// }
//
// func ExampleIsBetweenTime() {
// 	t, _ := time.Parse(time.RFC3339, "2009-02-04T21:00:57-08:00")
// 	after, _ := time.Parse(time.RFC3339, "2009-02-05T21:00:57-08:00")
// 	before, _ := time.Parse(time.RFC3339, "2009-02-06T21:00:57-08:00")
// 	err := validator.ValidateTime(context.Background(), t, it.IsBetweenTime(after, before))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be between 2009-02-05T21:00:57-08:00 and 2009-02-06T21:00:57-08:00.
// }
//
// func ExampleHasUniqueValues() {
// 	v := []string{"foo", "bar", "baz", "foo"}
// 	err := validator.Validate(
// 		context.Background(),
// 		validation.Strings(v, it.HasUniqueValues()),
// 	)
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This collection should contain only unique elements.
// }
//
// func ExampleHasMinCount() {
// 	v := []int{1, 2}
// 	err := validator.ValidateIterable(context.Background(), v, it.HasMinCount(3))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This collection should contain 3 elements or more.
// }
//
// func ExampleHasMaxCount() {
// 	v := []int{1, 2}
// 	err := validator.ValidateIterable(context.Background(), v, it.HasMaxCount(1))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This collection should contain 1 element or less.
// }
//
// func ExampleHasCountBetween() {
// 	v := []int{1, 2}
// 	err := validator.ValidateIterable(context.Background(), v, it.HasCountBetween(3, 10))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This collection should contain 3 elements or more.
// }
//
// func ExampleHasExactCount() {
// 	v := []int{1, 2}
// 	err := validator.ValidateIterable(context.Background(), v, it.HasExactCount(3))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This collection should contain exactly 3 elements.
// }
//
// func ExampleHasMinLength() {
// 	v := "foo"
// 	err := validator.ValidateString(context.Background(), v, it.HasMinLength(5))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value is too short. It should have 5 characters or more.
// }
//
// func ExampleHasMaxLength() {
// 	v := "foo"
// 	err := validator.ValidateString(context.Background(), v, it.HasMaxLength(2))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value is too long. It should have 2 characters or less.
// }
//
// func ExampleHasLengthBetween() {
// 	v := "foo"
// 	err := validator.ValidateString(context.Background(), v, it.HasLengthBetween(5, 10))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value is too short. It should have 5 characters or more.
// }
//
// func ExampleHasExactLength() {
// 	v := "foo"
// 	err := validator.ValidateString(context.Background(), v, it.HasExactLength(5))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should have exactly 5 characters.
// }
//
// func ExampleMatches() {
// 	v := "foo123"
// 	err := validator.ValidateString(context.Background(), v, it.Matches(regexp.MustCompile("^[a-z]+$")))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value is not valid.
// }
//
// func ExampleDoesNotMatch() {
// 	v := "foo"
// 	err := validator.ValidateString(context.Background(), v, it.DoesNotMatch(regexp.MustCompile("^[a-z]+$")))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value is not valid.
// }
//
// func ExampleIsJSON_validJSON() {
// 	v := `{"valid": true}`
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsJSON()))
// 	fmt.Println(err)
// 	// Output:
// 	// <nil>
// }
//
// func ExampleIsJSON_invalidJSON() {
// 	v := `"invalid": true`
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsJSON()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value should be valid JSON.
// }
//
// func ExampleIsEmail_validEmail() {
// 	v := "user@example.com"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsEmail()))
// 	fmt.Println(err)
// 	// Output:
// 	// <nil>
// }
//
// func ExampleIsEmail_invalidEmail() {
// 	v := "user example.com"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsEmail()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value is not a valid email address.
// }
//
// func ExampleIsHTML5Email_validEmail() {
// 	v := "{}~!@example.com"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsEmail()))
// 	fmt.Println(err)
// 	// Output:
// 	// <nil>
// }
//
// func ExampleIsHTML5Email_invalidEmail() {
// 	v := "@example.com"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsEmail()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value is not a valid email address.
// }
//
// func ExampleIsHostname_validHostname() {
// 	v := "example.com"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsHostname()))
// 	fmt.Println(err)
// 	// Output:
// 	// <nil>
// }
//
// func ExampleIsHostname_invalidHostname() {
// 	v := "example-.com"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsHostname()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value is not a valid hostname.
// }
//
// func ExampleIsHostname_reservedHostname() {
// 	v := "example.localhost"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsHostname()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value is not a valid hostname.
// }
//
// func ExampleIsLooseHostname_validHostname() {
// 	v := "example.com"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsLooseHostname()))
// 	fmt.Println(err)
// 	// Output:
// 	// <nil>
// }
//
// func ExampleIsLooseHostname_invalidHostname() {
// 	v := "example-.com"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsLooseHostname()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value is not a valid hostname.
// }
//
// func ExampleIsLooseHostname_reservedHostname() {
// 	v := "example.localhost"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsLooseHostname()))
// 	fmt.Println(err)
// 	// Output:
// 	// <nil>
// }
//
// func ExampleIsURL_validURL() {
// 	v := "http://example.com"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsURL()))
// 	fmt.Println(err)
// 	// Output:
// 	// <nil>
// }
//
// func ExampleIsURL_invalidURL() {
// 	v := "example.com"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsURL()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This value is not a valid URL.
// }
//
// func ExampleURLConstraint_WithRelativeSchema() {
// 	v := "//example.com"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsURL().WithRelativeSchema()))
// 	fmt.Println(err)
// 	// Output:
// 	// <nil>
// }
//
// func ExampleURLConstraint_WithSchemas() {
// 	v := "ftp://example.com"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsURL().WithSchemas("http", "https", "ftp")))
// 	fmt.Println(err)
// 	// Output:
// 	// <nil>
// }
//
// func ExampleIsIP_validIP() {
// 	v := "123.123.123.123"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsIP()))
// 	fmt.Println(err)
// 	// Output:
// 	// <nil>
// }
//
// func ExampleIsIP_invalidIP() {
// 	v := "123.123.123.345"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsIP()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This is not a valid IP address.
// }
//
// func ExampleIsIPv4_validIP() {
// 	v := "123.123.123.123"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsIPv4()))
// 	fmt.Println(err)
// 	// Output:
// 	// <nil>
// }
//
// func ExampleIsIPv4_invalidIP() {
// 	v := "123.123.123.345"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsIPv4()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This is not a valid IP address.
// }
//
// func ExampleIsIPv6_validIP() {
// 	v := "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsIPv6()))
// 	fmt.Println(err)
// 	// Output:
// 	// <nil>
// }
//
// func ExampleIsIPv6_invalidIP() {
// 	v := "z001:0db8:85a3:0000:0000:8a2e:0370:7334"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsIPv6()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This is not a valid IP address.
// }
//
// func ExampleIPConstraint_DenyPrivateIP_restrictedPrivateIPv4() {
// 	v := "192.168.1.0"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsIP().DenyPrivateIP()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This IP address is prohibited to use.
// }
//
// func ExampleIPConstraint_DenyPrivateIP_restrictedPrivateIPv6() {
// 	v := "fdfe:dcba:9876:ffff:fdc6:c46b:bb8f:7d4c"
// 	err := validator.Validate(context.Background(), validation.String(v, it.IsIPv6().DenyPrivateIP()))
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This IP address is prohibited to use.
// }
//
// func ExampleIPConstraint_DenyIP() {
// 	v := "127.0.0.1"
// 	err := validator.Validate(
// 		context.Background(),
// 		validation.String(
// 			v,
// 			it.IsIP().DenyIP(func(ip net.IP) bool {
// 				return ip.IsLoopback()
// 			}),
// 		),
// 	)
// 	fmt.Println(err)
// 	// Output:
// 	// violation: This IP address is prohibited to use.
// }
