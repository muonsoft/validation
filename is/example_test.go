package is_test

import (
	"fmt"

	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/validate"
)

func ExampleJSON() {
	fmt.Println(is.JSON(`{"valid": true}`)) // valid
	fmt.Println(is.JSON(`"invalid": true`)) // invalid
	// Output:
	// true
	// false
}

func ExampleInteger() {
	fmt.Println(is.Integer("123"))
	fmt.Println(is.Integer("123.123"))
	fmt.Println(is.Integer("-123"))
	fmt.Println(is.Integer("123foo"))
	// Output:
	// true
	// false
	// true
	// false
}

func ExampleNumber() {
	fmt.Println(is.Number("123"))
	fmt.Println(is.Number("123.123"))
	fmt.Println(is.Number("123e123"))
	fmt.Println(is.Number("-123"))
	fmt.Println(is.Number("123foo"))
	// Output:
	// true
	// true
	// true
	// true
	// false
}

func ExampleStringInList() {
	fmt.Println(is.StringInList("foo", nil))
	fmt.Println(is.StringInList("foo", []string{"bar", "baz"}))
	fmt.Println(is.StringInList("foo", []string{"bar", "baz", "foo"}))
	// Output:
	// false
	// false
	// true
}

func ExampleUniqueStrings() {
	fmt.Println(is.UniqueStrings([]string{}))
	fmt.Println(is.UniqueStrings([]string{"one", "two", "three"}))
	fmt.Println(is.UniqueStrings([]string{"one", "two", "one"}))
	// Output:
	// true
	// true
	// false
}

func ExampleEmail() {
	fmt.Println(is.Email("user@example.com"))         // valid
	fmt.Println(is.Email("{}~!@example.com"))         // valid
	fmt.Println(is.Email("пользователь@example.com")) // valid
	fmt.Println(is.Email("user example.com"))         // invalid
	// Output:
	// true
	// true
	// true
	// false
}

func ExampleHTML5Email() {
	fmt.Println(is.HTML5Email("user@example.com"))         // valid
	fmt.Println(is.HTML5Email("{}~!@example.com"))         // valid
	fmt.Println(is.HTML5Email("пользователь@example.com")) // invalid
	fmt.Println(is.HTML5Email("user example.com"))         // invalid
	// Output:
	// true
	// true
	// false
	// false
}

func ExampleURL() {
	fmt.Println(is.URL("https://example.com"))                       // valid absolute URL
	fmt.Println(is.URL("ftp://example.com", "http", "https", "ftp")) // valid URL with custom schema
	fmt.Println(is.URL("example.com"))                               // invalid URL
	fmt.Println(is.URL("//example.com", ""))                         // valid relative URL
	// Output:
	// true
	// true
	// false
	// true
}

func ExampleIP() {
	fmt.Println(is.IP("123.123.123.123"))                         // valid IPv4
	fmt.Println(is.IP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")) // valid IPv6
	fmt.Println(is.IP("123.123.123.345"))                         // invalid
	fmt.Println(is.IP("192.168.1.0"))                             // non-restricted private IP
	fmt.Println(is.IP("192.168.1.0", validate.DenyPrivateIP()))   // restricted private IP
	// Output:
	// true
	// true
	// false
	// true
	// false
}

func ExampleIPv4() {
	fmt.Println(is.IPv4("123.123.123.123"))                         // valid IPv4
	fmt.Println(is.IPv4("2001:0db8:85a3:0000:0000:8a2e:0370:7334")) // invalid (IPv6)
	fmt.Println(is.IPv4("123.123.123.345"))                         // invalid
	fmt.Println(is.IPv4("192.168.1.0"))                             // non-restricted private IP
	fmt.Println(is.IPv4("192.168.1.0", validate.DenyPrivateIP()))   // restricted private IP
	// Output:
	// true
	// false
	// false
	// true
	// false
}

func ExampleIPv6() {
	fmt.Println(is.IPv6("2001:0db8:85a3:0000:0000:8a2e:0370:7334"))                           // valid (IPv6)
	fmt.Println(is.IPv6("123.123.123.123"))                                                   // invalid IPv4
	fmt.Println(is.IPv6("z001:0db8:85a3:0000:0000:8a2e:0370:7334"))                           // invalid
	fmt.Println(is.IPv6("fdfe:dcba:9876:ffff:fdc6:c46b:bb8f:7d4c"))                           // non-restricted private IP
	fmt.Println(is.IPv6("fdfe:dcba:9876:ffff:fdc6:c46b:bb8f:7d4c", validate.DenyPrivateIP())) // restricted private IP
	// Output:
	// true
	// false
	// false
	// true
	// false
}

func ExampleHostname() {
	fmt.Println(is.Hostname("example.com"))       // valid
	fmt.Println(is.Hostname("example.localhost")) // valid
	fmt.Println(is.Hostname("com"))               // valid
	fmt.Println(is.Hostname("example-.com"))      // invalid
	// Output:
	// true
	// true
	// true
	// false
}

func ExampleStrictHostname() {
	fmt.Println(is.StrictHostname("example.com"))       // valid
	fmt.Println(is.StrictHostname("example.localhost")) // reserved
	fmt.Println(is.StrictHostname("com"))               // invalid
	fmt.Println(is.StrictHostname("example-.com"))      // invalid
	// Output:
	// true
	// false
	// false
	// false
}
