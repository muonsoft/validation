package is_test

import (
	"fmt"

	"github.com/muonsoft/validation/is"
)

func ExampleURL_validAbsoluteURL() {
	valid := is.URL("https://example.com")

	fmt.Println(valid)
	// Output:
	// true
}

func ExampleURL_validURLWithCustomProtocol() {
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

func ExampleRelativeURL_validRelativeURL() {
	valid := is.RelativeURL("//example.com")

	fmt.Println(valid)
	// Output:
	// true
}

func ExampleRelativeURL_invalidURL() {
	valid := is.RelativeURL("example.com")

	fmt.Println(valid)
	// Output:
	// false
}
