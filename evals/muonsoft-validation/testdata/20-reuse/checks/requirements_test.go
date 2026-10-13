package scenario_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
)

var regionError = validation.NewError("region", "Region is disabled.")

type regionRule struct {
	calls *int
	t     *testing.T
}

func (r regionRule) ValidateString(ctx context.Context, v *validation.Validator, value *string) error {
	*r.calls++
	if ctx != r.t.Context() {
		r.t.Fatal("context")
	}
	if value != nil && *value != "" && *value != "north" {
		return v.CreateViolation(ctx, regionError, regionError.Message())
	}
	return nil
}
func validProfile() s.Profile {
	return s.Profile{DisplayName: "Ada", Email: "ada@example.com", Age: 18, Address: s.Address{Street: "A", Region: "north"}}
}
func TestReuseAccumulation(t *testing.T) {
	calls := 0
	quota := -1
	x := s.Profile{Email: "bad", Preferences: &s.Preferences{Theme: "", Quota: &quota, Tags: []string{"", "x", "x", "toolonggg"}}, Address: s.Address{Region: "south"}}
	err := x.Validate(t.Context(), newValidator(t).AtProperty("account"), regionRule{&calls, t})
	want(t, err, expected{"account.displayName", validation.ErrIsBlank},
		expected{"account.email", validation.ErrInvalidEmail},
		expected{"account.age", validation.ErrNotInRange},
		expected{"account.address.street", validation.ErrIsBlank},
		expected{"account.address.region", regionError},
		expected{"account.preferences.theme", validation.ErrNoSuchChoice},
		expected{"account.preferences.quota", validation.ErrNotInRange},
		expected{"account.preferences.tags", validation.ErrTooManyElements},
		expected{"account.preferences.tags", validation.ErrNotUnique},
		expected{"account.preferences.tags[0]", validation.ErrIsBlank},
		expected{"account.preferences.tags[3]", validation.ErrTooLong})
	if calls != 1 {
		t.Fatal(calls)
	}
}
func TestReuseBoundaries(t *testing.T) {
	calls := 0
	rule := regionRule{&calls, t}
	v := newValidator(t)
	x := validProfile()
	want(t, x.Validate(t.Context(), v, rule))
	quota := 100
	x.Preferences = &s.Preferences{Theme: "dark", Quota: &quota, Tags: []string{"界界界界界界界界"}}
	x.Age = 120
	x.DisplayName = strings.Repeat("界", 20)
	want(t, x.Validate(t.Context(), v, rule))
	quota = 101
	x.Age = 121
	x.DisplayName += "界"
	want(t, x.Validate(t.Context(), v, rule), expected{"displayName", validation.ErrTooLong},
		expected{"age", validation.ErrNotInRange},
		expected{"preferences.quota", validation.ErrNotInRange})
	x = validProfile()
	x.Preferences = &s.Preferences{Theme: "light"}
	want(t, x.Validate(t.Context(), v, rule))
}

type stringRule func(context.Context, *validation.Validator, *string) error

func (f stringRule) ValidateString(ctx context.Context, v *validation.Validator, p *string) error {
	return f(ctx, v, p)
}
func TestReuseDependencyErrors(t *testing.T) {
	failure := errors.New("offline")
	x := validProfile()
	err := x.Validate(t.Context(), newValidator(t), stringRule(func(context.Context, *validation.Validator, *string) error { return failure }))
	if !errors.Is(err, failure) {
		t.Fatal(err)
	}
}
