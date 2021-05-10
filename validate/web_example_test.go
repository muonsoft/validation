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
