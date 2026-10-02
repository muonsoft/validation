package validate_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/require"
)

func TestEstimatePasswordStrength(t *testing.T) {
	tests := []struct {
		password string
		score    validate.PasswordStrengthScore
	}{
		{"", 0},
		{"password", 0},
		{"How-is-this", 1},
		{"Reasonable-pwd", 2},
		{"Good password?", 1},
		{"This 1s a very g00d Pa55word! ;-)", 4},
		{"pudding-smack-👌🏼-fox-😎", 4},
		{strings.Repeat("a", 1000), 0},
		{strings.Repeat("ab", 50), 3},
		{"abcdefghijklmnopqrstuvwxyz", 4},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) { require.Equal(t, tt.score, validate.EstimatePasswordStrength(tt.password)) })
	}
}

func TestPasswordStrengthCategoryThresholds(t *testing.T) {
	// Four unique bytes give exactly two estimated bits per repeated byte.
	// These lengths are immediately below and above each threshold for each pool.
	tests := []struct {
		chars   string
		lengths []int
	}{
		{"0123", []int{28, 38, 48, 58}},
		{"abcd", []int{25, 35, 45, 55}},
		{"ABCD", []int{25, 35, 45, 55}},
		{" !@#", []int{24, 34, 44, 54}},
		{"\x00\x01\x1f\x7f", []int{24, 34, 44, 54}},
		{"\x80\x81\xfe\xff", []int{20, 30, 40, 50}},
	}
	for i, tt := range tests {
		for boundary, length := range tt.lengths {
			t.Run(fmt.Sprintf("%d/%d", i, boundary), func(t *testing.T) {
				input := strings.Repeat(tt.chars, 20)
				require.Equal(t, validate.PasswordStrengthScore(boundary), validate.EstimatePasswordStrength(input[:length-1]))
				require.Equal(t, validate.PasswordStrengthScore(boundary+1), validate.EstimatePasswordStrength(input[:length]))
			})
		}
	}
}

func TestPasswordStrengthOptions(t *testing.T) {
	require.ErrorIs(t, validate.PasswordStrength(""), validate.ErrPasswordTooWeak)
	require.NoError(t, validate.PasswordStrength("Reasonable-pwd"))
	for score := validate.PasswordStrengthVeryWeak; score <= validate.PasswordStrengthVeryStrong; score++ {
		for minimum := validate.PasswordStrengthWeak; minimum <= validate.PasswordStrengthVeryStrong; minimum++ {
			calls := 0
			err := validate.PasswordStrength("secret", validate.WithMinPasswordStrength(minimum), validate.WithPasswordStrengthEstimator(func(value string) validate.PasswordStrengthScore {
				calls++
				require.Equal(t, "secret", value)
				return score
			}))
			require.Equal(t, 1, calls)
			if score < minimum {
				require.ErrorIs(t, err, validate.ErrPasswordTooWeak)
			} else {
				require.NoError(t, err)
			}
		}
	}
	for _, minimum := range []validate.PasswordStrengthScore{-1, 0, 5} {
		require.ErrorIs(t, validate.PasswordStrength("", validate.WithMinPasswordStrength(minimum), validate.WithPasswordStrengthEstimator(func(string) validate.PasswordStrengthScore { t.Fatal("must not run"); return 0 })), validate.ErrInvalidPasswordStrengthMinimum)
	}
	for _, score := range []validate.PasswordStrengthScore{-1, 5} {
		require.ErrorIs(t, validate.PasswordStrength("secret", validate.WithPasswordStrengthEstimator(func(string) validate.PasswordStrengthScore { return score })), validate.ErrInvalidPasswordStrengthScore)
	}
	require.NoError(t, validate.PasswordStrength("Reasonable-pwd", validate.WithPasswordStrengthEstimator(nil)))
}

func ExampleEstimatePasswordStrength() {
	fmt.Println(validate.EstimatePasswordStrength("Reasonable-pwd"))
	// Output:
	// 2
}

func ExamplePasswordStrength() {
	fmt.Println(validate.PasswordStrength("Reasonable-pwd"))
	// Output:
	// <nil>
}
