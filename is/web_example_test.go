package is_test

import (
	"fmt"

	"github.com/muonsoft/validation/is"
)

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
