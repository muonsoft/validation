package is_test

import (
	"fmt"
	"testing"

	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/require"
)

func TestWeek(t *testing.T) {
	require.True(t, is.Week(""))
	require.True(t, is.Week("2020-W53"))
	require.False(t, is.Week("2021-W53"))
	require.False(t, is.Week("2020-W1"))
	require.False(t, is.Week("2020-W53", validate.WithMinWeek("2021-W01")))
	require.False(t, is.Week("2021-W01", validate.WithMaxWeek("2020-W53")))
	require.False(t, is.Week("", validate.WithMinWeek("invalid")))
}

func ExampleWeek() {
	fmt.Println(is.Week("2020-W53"))
	fmt.Println(is.Week("2021-W53"))
	// Output:
	// true
	// false
}
