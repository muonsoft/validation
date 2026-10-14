package scenario_test

import (
	"fmt"
	"strings"
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validate"
)

func TestPasswordPresence(t *testing.T) {
	calls := 0
	estimate := func(string) validate.PasswordStrengthScore { calls++; return validate.PasswordStrengthVeryStrong }
	want(t, (s.PasswordUpdate{Estimate: estimate}).Validate(t.Context(), newValidator(t)))
	if calls != 0 {
		t.Fatal("called estimator for nil")
	}
	empty := ""
	want(t, (s.PasswordUpdate{Password: &empty}).Validate(t.Context(), newValidator(t)), expected{"password", validation.ErrPasswordTooWeak})
	want(t, (s.PasswordUpdate{Password: &empty, Estimate: estimate}).Validate(t.Context(), newValidator(t)))
	if calls != 1 {
		t.Fatal("empty input not scored exactly once")
	}
}
func TestPasswordScores(t *testing.T) {
	secret := " secret-with-spaces "
	for score := 0; score <= 4; score++ {
		calls := 0
		estimate := func(value string) validate.PasswordStrengthScore {
			calls++
			if value != secret {
				t.Fatal("changed password")
			}
			return validate.PasswordStrengthScore(score)
		}
		err := (s.PasswordUpdate{Password: &secret, Estimate: estimate}).Validate(t.Context(), newValidator(t).AtProperty("account"))
		if score < 3 {
			want(t, err, expected{"account.password", validation.ErrPasswordTooWeak})
			if strings.Contains(err.Error(), secret) {
				t.Fatal("password leaked")
			}
			list, _ := validation.UnwrapViolations(err)
			if strings.Contains(fmt.Sprint(list.First().Violation().Parameters()), secret) {
				t.Fatal("password leaked in parameters")
			}
			if strings.Contains(list.First().Violation().Message(), secret) {
				t.Fatal("password leaked in message")
			}
		} else {
			want(t, err)
		}
		if calls != 1 {
			t.Fatal("estimator call count")
		}
	}
}
func TestPasswordConfiguration(t *testing.T) {
	secret := "x"
	for _, score := range []validate.PasswordStrengthScore{-1, 5} {
		err := (s.PasswordUpdate{Password: &secret, Estimate: func(string) validate.PasswordStrengthScore { return score }}).Validate(t.Context(), newValidator(t))
		if err == nil {
			t.Fatal("invalid score accepted")
		}
		if _, ok := validation.UnwrapViolations(err); ok {
			t.Fatal("configuration reported as input violation")
		}
	}
}
func TestPasswordDefault(t *testing.T) {
	for _, secret := range []string{"abcdefghijklmnopqrstuvwxyz", strings.Repeat("ab", 50)} {
		want(t, (s.PasswordUpdate{Password: &secret}).Validate(t.Context(), newValidator(t)))
	}
	secret := "aaaaaa"
	want(t, (s.PasswordUpdate{Password: &secret}).Validate(t.Context(), newValidator(t)), expected{"password", validation.ErrPasswordTooWeak})
}
