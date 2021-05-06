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

func ExampleURL_validURLWithCustomProtocol() {
	err := validate.URL("ftp://example.com", "http", "https", "ftp")

	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleURL_urlWithoutProtocol() {
	err := validate.URL("example.com")

	fmt.Println(err)
	// Output:
	// unexpected protocol
}

func ExampleURL_invalidURL() {
	err := validate.URL("http:// example.com/")

	fmt.Println(err)
	// Output:
	// parse "http:// example.com/": invalid character " " in host name
}

func ExampleRelativeURL_validRelativeURL() {
	err := validate.RelativeURL("//example.com")

	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleRelativeURL_invalidURL() {
	err := validate.RelativeURL("example.com")

	fmt.Println(err)
	// Output:
	// invalid
}
