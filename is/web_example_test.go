package is_test

import (
	"fmt"

	"github.com/muonsoft/validation/is"
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
