package scenario_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	s "example.com/validation-eval/scenario"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/message/translations"
	"golang.org/x/text/language"
)

type catalogStub struct {
	calls atomic.Int32
	ctx   context.Context
	fail  error
	t     *testing.T
}

func (c *catalogStub) Lookup(ctx context.Context, code string) (s.Entry, bool, error) {
	c.calls.Add(1)
	if ctx != c.ctx {
		c.t.Error("lost context")
	}
	return s.Entry{Category: "basic"}, code != "missing", c.fail
}
func customValidator(t *testing.T) *validation.Validator {
	t.Helper()
	v, err := validation.NewValidator(s.ValidatorOptions()...)
	if err != nil {
		t.Fatal(err)
	}
	return v
}
func TestConstraintRequiredness(t *testing.T) {
	source := &catalogStub{ctx: t.Context(), t: t}
	c := s.NewCodeConstraint(source, "basic")
	var _ validation.StringConstraint = c
	v := customValidator(t)
	empty := ""
	for _, p := range []*string{nil, &empty} {
		want(t, c.ValidateString(t.Context(), v, p))
		want(t, (s.Signup{Code: p}).Validate(t.Context(), v, c), expected{"code", validation.ErrIsBlank},
			expected{"label", validation.ErrIsBlank})
	}
	if source.calls.Load() != 0 {
		t.Fatal("queried empty")
	}
	code := "known"
	want(t, (s.Signup{Code: &code, Label: "A"}).Validate(t.Context(), v, c))
	want(t, (s.Transfer{Destination: "known", Label: "A"}).Validate(t.Context(), v, c))
	if source.calls.Load() != 2 {
		t.Fatal(source.calls.Load())
	}
}
func TestConstraintErrors(t *testing.T) {
	source := &catalogStub{ctx: t.Context(), t: t}
	c := s.NewCodeConstraint(source, "premium")
	v := customValidator(t).AtProperty("request")
	want(t, (s.Transfer{Destination: "missing"}).Validate(t.Context(), v, c), expected{"request.destinationCode", s.ErrUnknown},
		expected{"request.label", validation.ErrIsBlank})
	want(t, (s.Transfer{Destination: "known", Label: "A"}).Validate(t.Context(), v, c), expected{"request.destinationCode", s.ErrCategory})
	source.fail = errors.New("offline")
	err := (s.Transfer{Destination: "known"}).Validate(t.Context(), v, c)
	if !errors.Is(err, source.fail) {
		t.Fatal(err)
	}
}
func TestConstraintReuseAndTranslations(t *testing.T) {
	source := &catalogStub{ctx: t.Context(), t: t}
	c := s.NewCodeConstraint(source, "premium")
	v := customValidator(t)
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			prefix := fmt.Sprintf("item%d", i)
			local := v.AtProperty(prefix)
			message := "Code must belong to premium."
			if i%2 == 0 {
				local = local.WithLanguage(language.Russian)
				message = "Код должен относиться к категории premium."
			}
			err := (s.Transfer{Destination: "known", Label: "A"}).Validate(t.Context(), local, c)
			want(t, err, expected{prefix + ".destinationCode", s.ErrCategory})
			list, _ := validation.UnwrapViolations(err)
			if list.First().Violation().Message() != message {
				t.Error(list.First().Violation().Message())
			}
		}(i)
	}
	wg.Wait()
	if source.calls.Load() != 12 {
		t.Fatal(source.calls.Load())
	}
	err := (s.Transfer{Destination: "missing", Label: "A"}).Validate(t.Context(), v.WithLanguage(language.Russian), c)
	list, _ := validation.UnwrapViolations(err)
	if list.First().Violation().Message() != "Код не найден." {
		t.Fatal(err)
	}
}

func TestConstraintCallerFactory(t *testing.T) {
	translator, err := translations.NewTranslator()
	if err != nil {
		t.Fatal(err)
	}
	base := validation.NewViolationFactory(translator)
	calls := 0
	factory := validation.NewViolationFunc(func(code error, template string, plural int, parameters []validation.TemplateParameter, path *validation.PropertyPath, lang language.Tag) validation.Violation {
		calls++
		return base.CreateViolation(code, template, plural, parameters, path, lang)
	})
	v, err := validation.NewValidator(validation.SetViolationFactory(factory))
	if err != nil {
		t.Fatal(err)
	}
	source := &catalogStub{ctx: t.Context(), t: t}
	c := s.NewCodeConstraint(source, "premium")
	err = (s.Transfer{Destination: "known", Label: "A"}).Validate(t.Context(), v.AtProperty("outer"), c)
	want(t, err, expected{"outer.destinationCode", s.ErrCategory})
	if calls != 1 {
		t.Fatalf("caller factory used %d times", calls)
	}
}
