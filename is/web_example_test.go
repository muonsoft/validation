package is_test

import (
	"fmt"

	"github.com/muonsoft/validation/is"
)

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
