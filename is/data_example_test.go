package is_test

import (
	"fmt"

	"github.com/muonsoft/validation/is"
)

func ExampleJSON() {
	fmt.Println(is.JSON(`{"valid": true}`)) // valid
	fmt.Println(is.JSON(`"invalid": true`)) // invalid
	// Output:
	// true
	// false
}

func ExampleUniqueStrings() {
	fmt.Println(is.UniqueStrings([]string{}))
	fmt.Println(is.UniqueStrings([]string{"one", "two", "three"}))
	fmt.Println(is.UniqueStrings([]string{"one", "two", "one"}))
	// Output:
	// true
	// true
	// false
}
