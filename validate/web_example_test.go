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

func ExampleIP() {
	fmt.Println(validate.IP("123.123.123.123"))                         // valid IPv4
	fmt.Println(validate.IP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")) // valid IPv6
	fmt.Println(validate.IP("123.123.123.345"))                         // invalid
	fmt.Println(validate.IP("192.168.1.0"))                             // non-restricted private IP
	fmt.Println(validate.IP("192.168.1.0", validate.DenyPrivateIP()))   // restricted private IP
	// Output:
	// <nil>
	// <nil>
	// invalid
	// <nil>
	// prohibited
}

func ExampleIPv4() {
	fmt.Println(validate.IPv4("123.123.123.123"))                         // valid IPv4
	fmt.Println(validate.IPv4("2001:0db8:85a3:0000:0000:8a2e:0370:7334")) // invalid (IPv6)
	fmt.Println(validate.IPv4("123.123.123.345"))                         // invalid
	fmt.Println(validate.IPv4("192.168.1.0"))                             // non-restricted private IP
	fmt.Println(validate.IPv4("192.168.1.0", validate.DenyPrivateIP()))   // restricted private IP
	// Output:
	// <nil>
	// invalid
	// invalid
	// <nil>
	// prohibited
}

func ExampleIPv6() {
	fmt.Println(validate.IPv6("2001:0db8:85a3:0000:0000:8a2e:0370:7334"))                           // valid (IPv6)
	fmt.Println(validate.IPv6("123.123.123.123"))                                                   // invalid IPv4
	fmt.Println(validate.IPv6("z001:0db8:85a3:0000:0000:8a2e:0370:7334"))                           // invalid
	fmt.Println(validate.IPv6("fdfe:dcba:9876:ffff:fdc6:c46b:bb8f:7d4c"))                           // non-restricted private IP
	fmt.Println(validate.IPv6("fdfe:dcba:9876:ffff:fdc6:c46b:bb8f:7d4c", validate.DenyPrivateIP())) // restricted private IP
	// Output:
	// <nil>
	// invalid
	// invalid
	// <nil>
	// prohibited
}
