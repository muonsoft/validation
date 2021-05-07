package validate_test

import (
	"fmt"

	"github.com/muonsoft/validation/validate"
)

func ExampleURL_validAbsoluteURL() {
	err := validate.URL("https://example.com")

	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleURL_validURLWithCustomSchema() {
	err := validate.URL("ftp://example.com", "http", "https", "ftp")

	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleURL_urlWithoutSchema() {
	err := validate.URL("example.com")

	fmt.Println(err)
	// Output:
	// unexpected schema
}

func ExampleURL_invalidURL() {
	err := validate.URL("http:// example.com/")

	fmt.Println(err)
	// Output:
	// parse "http:// example.com/": invalid character " " in host name
}

func ExampleURL_validRelativeURL() {
	err := validate.URL("//example.com", "")

	fmt.Println(err)
	// Output:
	// <nil>
}
