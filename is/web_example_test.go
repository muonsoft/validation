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
