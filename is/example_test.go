package is_test

import (
	"fmt"
	"regexp"

	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/validate"
)

func ExampleEAN8() {
	fmt.Println(is.EAN8("42345671"))
	fmt.Println(is.EAN8("47195127"))
	fmt.Println(is.EAN8("00000000")) // zeros only are prohibited
	fmt.Println(is.EAN8("42345670")) // invalid checksum
	fmt.Println(is.EAN8("A4234671")) // contains non-digit
	// Output:
	// true
	// true
	// false
	// false
	// false
}

func ExampleEAN13() {
	fmt.Println(is.EAN13("4719512002889"))
	fmt.Println(is.EAN13("9782868890061"))
	fmt.Println(is.EAN13("0000000000000")) // zeros only are prohibited
	fmt.Println(is.EAN13("4006381333932")) // invalid checksum
	fmt.Println(is.EAN13("A782868890061")) // contains non-digit
	// Output:
	// true
	// true
	// false
	// false
	// false
}

func ExampleUPCA() {
	fmt.Println(is.UPCA("614141000036"))
	fmt.Println(is.UPCA("123456789999"))
	fmt.Println(is.UPCA("000000000000")) // zeros only are prohibited
	fmt.Println(is.UPCA("614141000037")) // invalid checksum
	fmt.Println(is.UPCA("A14141000036")) // contains non-digit
	// Output:
	// true
	// true
	// false
	// false
	// false
}

func ExampleUPCE() {
	fmt.Println(is.UPCE("123456"))   // 6-digit is always valid
	fmt.Println(is.UPCE("1234505"))  // 7-digit with last check digit
	fmt.Println(is.UPCE("01234505")) // 8-digit with first zero and last check digit
	fmt.Println(is.UPCE("00000000")) // zeros only are prohibited
	fmt.Println(is.UPCE("11234505")) // non-zero number system is prohibited
	fmt.Println(is.UPCE("01234501")) // invalid checksum
	fmt.Println(is.UPCE("A2345673")) // contains non-digit
	fmt.Println(is.UPCE("12345"))    // invalid length
	// Output:
	// true
	// true
	// true
	// false
	// false
	// false
	// false
	// false
}

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

func ExampleInList() {
	fmt.Println(is.InList("foo", nil))
	fmt.Println(is.InList("foo", []string{"bar", "baz"}))
	fmt.Println(is.InList("foo", []string{"bar", "baz", "foo"}))
	fmt.Println(is.InList(2, []int{1, 2, 3}))
	// Output:
	// false
	// false
	// true
	// true
}

func ExampleUnique() {
	fmt.Println(is.Unique([]string{}))
	fmt.Println(is.Unique([]string{"one", "two", "three"}))
	fmt.Println(is.Unique([]string{"one", "two", "one"}))
	fmt.Println(is.Unique([]int{1, 2, 1}))
	// Output:
	// true
	// true
	// false
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
	fmt.Println(is.URL("https://example.com"))                                                    // valid absolute URL
	fmt.Println(is.URL("ftp://example.com", validate.RestrictURLSchemas("http", "https", "ftp"))) // valid URL with custom schema
	fmt.Println(is.URL("example.com"))                                                            // invalid URL
	fmt.Println(is.URL("//example.com", validate.RestrictURLSchemas("")))                         // valid relative URL
	fmt.Println(is.URL("http://example.com", validate.RestrictURLHosts("sample.com")))            // not matching host
	fmt.Println(                                                                                  // matching by regexp
		is.URL(
			"http://sub.example.com",
			validate.RestrictURLHostByPattern(regexp.MustCompile(`^.*\.example\.com$`)),
		),
	)
	// Output:
	// true
	// true
	// false
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
