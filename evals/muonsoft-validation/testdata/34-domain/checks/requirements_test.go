package scenario_test

import (
	"context"
	"errors"
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
	"golang.org/x/text/language"
)

type contextKey struct{}

type tenantTranslator struct{}

func (tenantTranslator) Translate(lang language.Tag, template string, _ int) string {
	if lang == language.Russian && template == s.ErrTenantEmail.Message() {
		return "Адрес запрещён."
	}
	return template
}

type trackingFactory struct {
	validation.ViolationFactory
	calls int
}

func (f *trackingFactory) CreateViolation(err error, template string, count int, parameters []validation.TemplateParameter, path *validation.PropertyPath, lang language.Tag) validation.Violation {
	f.calls++
	return f.ViolationFactory.CreateViolation(err, template, count, parameters, path, lang)
}

func TestDomainPrerequisites(t *testing.T) {
	calls := 0
	rule := s.EmailPolicy{Check: func(context.Context, *validation.Validator, string) error { calls++; return nil }}
	want(t, (s.Signup{}).Validate(t.Context(), newValidator(t), rule), expected{"name", validation.ErrIsBlank}, expected{"email", validation.ErrIsBlank})
	want(t, (s.Signup{Email: "bad"}).Validate(t.Context(), newValidator(t), rule), expected{"name", validation.ErrIsBlank}, expected{"email", validation.ErrInvalidEmail})
	if calls != 0 {
		t.Fatal("tenant policy ran before email prerequisites")
	}
}
func TestDomainContextAndAccumulation(t *testing.T) {
	ctx := context.WithValue(t.Context(), contextKey{}, "tenant-7")
	factory := &trackingFactory{ViolationFactory: validation.NewViolationFactory(tenantTranslator{})}
	validator, err := validation.NewValidator(validation.SetViolationFactory(factory))
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	rule := s.EmailPolicy{Check: func(got context.Context, v *validation.Validator, value string) error {
		calls++
		if got.Value(contextKey{}) != "tenant-7" || value != "person@example.com" {
			t.Fatal("lost original context or value")
		}
		if v.IsIgnoredForGroups("tenant") {
			t.Fatal("lost caller groups")
		}
		return v.CreateViolation(got, s.ErrTenantEmail, s.ErrTenantEmail.Message())
	}}
	err = (s.Signup{Email: "person@example.com"}).Validate(ctx, validator.WithLanguage(language.Russian).WithGroups(validation.DefaultGroup, "tenant").AtProperty("signup"), rule)
	want(t, err, expected{"signup.name", validation.ErrIsBlank}, expected{"signup.email", s.ErrTenantEmail})
	wantPath(t, err, "signup.email", validation.PropertyName("signup"), validation.PropertyName("email"))
	if factory.calls != 2 {
		t.Fatalf("lost violation factory: calls=%d", factory.calls)
	}
	list, _ := validation.UnwrapViolations(err)
	for _, item := range list.AsSlice() {
		if errors.Is(item, s.ErrTenantEmail) && item.Message() != "Адрес запрещён." {
			t.Fatal("lost translation or language")
		}
	}
	if calls != 1 {
		t.Fatalf("calls=%d", calls)
	}
}
func TestDomainFatalAndReplacement(t *testing.T) {
	x := s.Signup{Name: "Name", Email: "person@example.com"}
	failure := errors.New("tenant backend unavailable")
	calls := 0
	rule := s.EmailPolicy{Check: func(context.Context, *validation.Validator, string) error { calls++; return failure }}
	err := x.Validate(t.Context(), newValidator(t), rule)
	if !errors.Is(err, failure) {
		t.Fatalf("lost failure: %v", err)
	}
	if _, ok := validation.UnwrapViolations(err); ok {
		t.Fatal("technical error became violation")
	}
	rule.Check = func(context.Context, *validation.Validator, string) error { calls++; return nil }
	want(t, x.Validate(t.Context(), newValidator(t), rule))
	if calls != 2 {
		t.Fatal("ignored replacement policy")
	}
}
