package validate_test

import (
	"fmt"

	"github.com/muonsoft/validation/validate"
)

func ExampleEAN8() {
	fmt.Println(validate.EAN8("42345671"))
	fmt.Println(validate.EAN8("00000000"))
	fmt.Println(validate.EAN8("42345670"))
	fmt.Println(validate.EAN8("423456712"))
	fmt.Println(validate.EAN8("A4234671"))
	// Output:
	// <nil>
	// contains only zeros
	// invalid checksum
	// unexpected length
	// contains non-digit: 'A'
}

func ExampleEAN13() {
	fmt.Println(validate.EAN13("4719512002889"))
	fmt.Println(validate.EAN13("0000000000000"))
	fmt.Println(validate.EAN13("2266111566"))
	fmt.Println(validate.EAN13("A782868890061"))
	fmt.Println(validate.EAN13("4006381333932"))
	// Output:
	// <nil>
	// contains only zeros
	// unexpected length
	// contains non-digit: 'A'
	// invalid checksum
}

func ExampleUPCA() {
	fmt.Println(validate.UPCA("614141000036"))
	fmt.Println(validate.UPCA("000000000000"))
	fmt.Println(validate.UPCA("61414100003"))
	fmt.Println(validate.UPCA("A14141000036"))
	fmt.Println(validate.UPCA("614141000037"))
	// Output:
	// <nil>
	// contains only zeros
	// unexpected length
	// contains non-digit: 'A'
	// invalid checksum
}

func ExampleUPCE() {
	fmt.Println(validate.UPCE("123456"))   // 6-digit is always valid
	fmt.Println(validate.UPCE("1234505"))  // 7-digit with last check digit
	fmt.Println(validate.UPCE("01234505")) // 8-digit with first zero and last check digit
	fmt.Println(validate.UPCE("00000000"))
	fmt.Println(validate.UPCE("11234505"))
	fmt.Println(validate.UPCE("01234501"))
	fmt.Println(validate.UPCE("023456731"))
	fmt.Println(validate.UPCE("A2345673"))
	fmt.Println(validate.UPCE("12345"))
	// Output:
	// <nil>
	// <nil>
	// <nil>
	// contains only zeros
	// invalid
	// invalid checksum
	// unexpected length
	// contains non-digit: 'A'
	// unexpected length
}

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
