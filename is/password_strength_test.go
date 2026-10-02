package is_test

import (
	"fmt"
	"testing"

	"github.com/muonsoft/validation/is"
	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/require"
)

func TestPasswordStrength(t *testing.T) {
	require.True(t, is.PasswordStrength("Reasonable-pwd"))
	require.False(t, is.PasswordStrength(""))
	require.False(t, is.PasswordStrength("password"))
	require.False(t, is.PasswordStrength("Reasonable-pwd", validate.WithMinPasswordStrength(validate.PasswordStrengthStrong)))
	require.False(t, is.PasswordStrength("secret", validate.WithMinPasswordStrength(0)))
	require.True(t, is.PasswordStrength("custom", validate.WithPasswordStrengthEstimator(func(string) validate.PasswordStrengthScore { return 4 })))
}

func ExamplePasswordStrength() {
	fmt.Println(is.PasswordStrength("Reasonable-pwd"))
	// Output:
	// true
}
