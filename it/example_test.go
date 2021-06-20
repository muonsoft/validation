package it_test

import (
	"fmt"
	"net"

	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validator"
)

func ExampleHasUniqueValues() {
	v := []string{"foo", "bar", "baz", "foo"}
	err := validator.ValidateStrings(v, it.HasUniqueValues())
	fmt.Println(err)
	// Output:
	// violation: This collection should contain only unique elements.
}

func ExampleIsJSON_validJSON() {
	v := `{"valid": true}`
	err := validator.ValidateString(&v, it.IsJSON())
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsJSON_invalidJSON() {
	v := `"invalid": true`
	err := validator.ValidateString(&v, it.IsJSON())
	fmt.Println(err)
	// Output:
	// violation: This value should be valid JSON.
}

func ExampleIsEmail_validEmail() {
	v := "user@example.com"
	err := validator.ValidateString(&v, it.IsEmail())
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsEmail_invalidEmail() {
	v := "user example.com"
	err := validator.ValidateString(&v, it.IsEmail())
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid email address.
}

func ExampleIsHTML5Email_validEmail() {
	v := "{}~!@example.com"
	err := validator.ValidateString(&v, it.IsEmail())
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsHTML5Email_invalidEmail() {
	v := "@example.com"
	err := validator.ValidateString(&v, it.IsEmail())
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid email address.
}

func ExampleIsHostname_validHostname() {
	v := "example.com"
	err := validator.ValidateString(&v, it.IsHostname())
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsHostname_invalidHostname() {
	v := "example-.com"
	err := validator.ValidateString(&v, it.IsHostname())
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid hostname.
}

func ExampleIsHostname_reservedHostname() {
	v := "example.localhost"
	err := validator.ValidateString(&v, it.IsHostname())
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid hostname.
}

func ExampleIsLooseHostname_validHostname() {
	v := "example.com"
	err := validator.ValidateString(&v, it.IsLooseHostname())
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsLooseHostname_invalidHostname() {
	v := "example-.com"
	err := validator.ValidateString(&v, it.IsLooseHostname())
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid hostname.
}

func ExampleIsLooseHostname_reservedHostname() {
	v := "example.localhost"
	err := validator.ValidateString(&v, it.IsLooseHostname())
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsURL_validURL() {
	v := "http://example.com"
	err := validator.ValidateString(&v, it.IsURL())
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsURL_invalidURL() {
	v := "example.com"
	err := validator.ValidateString(&v, it.IsURL())
	fmt.Println(err)
	// Output:
	// violation: This value is not a valid URL.
}

func ExampleURLConstraint_WithRelativeSchema() {
	v := "//example.com"
	err := validator.ValidateString(&v, it.IsURL().WithRelativeSchema())
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleURLConstraint_WithSchemas() {
	v := "ftp://example.com"
	err := validator.ValidateString(&v, it.IsURL().WithSchemas("http", "https", "ftp"))
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsIP_validIP() {
	v := "123.123.123.123"
	err := validator.ValidateString(&v, it.IsIP())
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsIP_invalidIP() {
	v := "123.123.123.345"
	err := validator.ValidateString(&v, it.IsIP())
	fmt.Println(err)
	// Output:
	// violation: This is not a valid IP address.
}

func ExampleIsIPv4_validIP() {
	v := "123.123.123.123"
	err := validator.ValidateString(&v, it.IsIPv4())
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsIPv4_invalidIP() {
	v := "123.123.123.345"
	err := validator.ValidateString(&v, it.IsIPv4())
	fmt.Println(err)
	// Output:
	// violation: This is not a valid IP address.
}

func ExampleIsIPv6_validIP() {
	v := "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
	err := validator.ValidateString(&v, it.IsIPv6())
	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIsIPv6_invalidIP() {
	v := "z001:0db8:85a3:0000:0000:8a2e:0370:7334"
	err := validator.ValidateString(&v, it.IsIPv6())
	fmt.Println(err)
	// Output:
	// violation: This is not a valid IP address.
}

func ExampleIPConstraint_DenyPrivateIP_restrictedPrivateIPv4() {
	v := "192.168.1.0"
	err := validator.ValidateString(&v, it.IsIP().DenyPrivateIP())
	fmt.Println(err)
	// Output:
	// violation: This IP address is prohibited to use.
}

func ExampleIPConstraint_DenyPrivateIP_restrictedPrivateIPv6() {
	v := "fdfe:dcba:9876:ffff:fdc6:c46b:bb8f:7d4c"
	err := validator.ValidateString(&v, it.IsIPv6().DenyPrivateIP())
	fmt.Println(err)
	// Output:
	// violation: This IP address is prohibited to use.
}

func ExampleIPConstraint_DenyIP() {
	v := "127.0.0.1"
	err := validator.ValidateString(&v, it.IsIP().DenyIP(func(ip net.IP) bool {
		return ip.IsLoopback()
	}))
	fmt.Println(err)
	// Output:
	// violation: This IP address is prohibited to use.
}
