package validate_test

import (
	"fmt"

	"github.com/muonsoft/validation/validate"
)

func ExampleURL() {
	fmt.Println(validate.URL("https://example.com"))                       // valid absolute URL
	fmt.Println(validate.URL("ftp://example.com", "http", "https", "ftp")) // valid URL with custom schema
	fmt.Println(validate.URL("example.com"))                               // url without schema
	fmt.Println(validate.URL("http:// example.com/"))                      // invalid URL
	fmt.Println(validate.URL("//example.com", ""))                         // valid relative URL
	// Output:
	// <nil>
	// <nil>
	// unexpected schema
	// parse "http:// example.com/": invalid character " " in host name
	// <nil>
}

func ExampleIP_validIP() {
	err := validate.IP("123.123.123.123")

	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIP_invalidIP() {
	err := validate.IP("123.123.123.345")

	fmt.Println(err)
	// Output:
	// invalid
}

func ExampleIP_restrictedPrivateIP() {
	err := validate.IP("192.168.1.0", validate.DenyPrivateIP())

	fmt.Println(err)
	// Output:
	// prohibited
}

func ExampleIPv4_validIP() {
	err := validate.IPv4("123.123.123.123")

	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIPv4_invalidIP() {
	err := validate.IPv4("123.123.123.345")

	fmt.Println(err)
	// Output:
	// invalid
}

func ExampleIPv4_restrictedPrivateIP() {
	err := validate.IPv4("192.168.1.0", validate.DenyPrivateIP())

	fmt.Println(err)
	// Output:
	// prohibited
}

func ExampleIPv6_validIP() {
	err := validate.IPv6("2001:0db8:85a3:0000:0000:8a2e:0370:7334")

	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleIPv6_invalidIP() {
	err := validate.IPv6("z001:0db8:85a3:0000:0000:8a2e:0370:7334")

	fmt.Println(err)
	// Output:
	// invalid
}

func ExampleIPv6_restrictedPrivateIP() {
	err := validate.IPv6("fdfe:dcba:9876:ffff:fdc6:c46b:bb8f:7d4c", validate.DenyPrivateIP())

	fmt.Println(err)
	// Output:
	// prohibited
}
