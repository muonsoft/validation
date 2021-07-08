package is_test

import (
	"fmt"

	"github.com/muonsoft/validation/is"
)

func ExampleStringInList() {
	fmt.Println(is.StringInList("foo", nil))
	fmt.Println(is.StringInList("foo", []string{"bar", "baz"}))
	fmt.Println(is.StringInList("foo", []string{"bar", "baz", "foo"}))
	// Output:
	// false
	// false
	// true
}
