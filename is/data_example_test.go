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
