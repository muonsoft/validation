package is_test

import (
	"fmt"

	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/validate"
)

func ExampleEmail_validEmail() {
	valid := is.Email("user@example.com")

	fmt.Println(valid)
	// Output:
	// true
}

func ExampleEmail_invalidEmail() {
	valid := is.Email("user example.com")

	fmt.Println(valid)
	// Output:
	// false
}

func ExampleHTML5Email() {
	valid := is.Email("{}~!@example.com")

	fmt.Println(valid)
	// Output:
	// true
}

func ExampleURL_validAbsoluteURL() {
	valid := is.URL("https://example.com")

	fmt.Println(valid)
	// Output:
	// true
}

func ExampleURL_validURLWithCustomSchema() {
	valid := is.URL("ftp://example.com", "http", "https", "ftp")

	fmt.Println(valid)
	// Output:
	// true
}

func ExampleURL_invalidURL() {
	valid := is.URL("example.com")

	fmt.Println(valid)
	// Output:
	// false
}

func ExampleURL_validRelativeURL() {
	valid := is.URL("//example.com", "")

	fmt.Println(valid)
	// Output:
	// true
}

func ExampleIP_validIP() {
	valid := is.IP("123.123.123.123")

	fmt.Println(valid)
	// Output:
	// true
}

func ExampleIP_invalidIP() {
	valid := is.IP("123.123.123.345")

	fmt.Println(valid)
	// Output:
	// false
}

func ExampleIP_restrictedPrivateIP() {
	valid := is.IP("192.168.1.0", validate.DenyPrivateIP())

	fmt.Println(valid)
	// Output:
	// false
}

func ExampleIPv4_validIP() {
	valid := is.IPv4("123.123.123.123")

	fmt.Println(valid)
	// Output:
	// true
}

func ExampleIPv4_invalidIP() {
	valid := is.IPv4("123.123.123.345")

	fmt.Println(valid)
	// Output:
	// false
}

func ExampleIPv4_restrictedPrivateIP() {
	valid := is.IPv4("192.168.1.0", validate.DenyPrivateIP())

	fmt.Println(valid)
	// Output:
	// false
}

func ExampleIPv6_validIP() {
	valid := is.IPv6("2001:0db8:85a3:0000:0000:8a2e:0370:7334")

	fmt.Println(valid)
	// Output:
	// true
}

func ExampleIPv6_invalidIP() {
	valid := is.IPv6("z001:0db8:85a3:0000:0000:8a2e:0370:7334")

	fmt.Println(valid)
	// Output:
	// false
}

func ExampleIPv6_restrictedPrivateIP() {
	valid := is.IPv6("fdfe:dcba:9876:ffff:fdc6:c46b:bb8f:7d4c", validate.DenyPrivateIP())

	fmt.Println(valid)
	// Output:
	// false
}
