package is_test

import (
	"fmt"
	"testing"

	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/require"
)

func TestCardScheme(t *testing.T) {
	require.True(t, is.CardScheme("", validate.CardSchemeVisa))
	require.True(t, is.CardScheme("4111111111111111", validate.CardSchemeVisa))
	require.True(t, is.CardScheme("5100000000000000", validate.CardSchemeVisa, validate.CardSchemeMastercard))
	require.False(t, is.CardScheme("4111111111111111", validate.CardSchemeMastercard))
	require.False(t, is.CardScheme("bad", validate.CardSchemeVisa))
	require.False(t, is.CardScheme(""))
	require.False(t, is.CardScheme("", "unknown"))
}

func ExampleCardScheme() {
	fmt.Println(is.CardScheme("4111111111111111", validate.CardSchemeVisa))
	// Output:
	// true
}
