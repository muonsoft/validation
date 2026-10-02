package validation_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/muonsoft/language"
	"github.com/muonsoft/validation"
	"github.com/stretchr/testify/require"
)

func TestViolationListJoin_IndependentNodes(t *testing.T) {
	a, b := newViolationWithError(t, errors.New("a")), newViolationWithError(t, errors.New("b"))
	source := validation.NewViolationList(a)
	destination := validation.NewViolationList()
	destination.Join(source)
	destination.Append(b)
	require.Equal(t, []validation.Violation{a}, source.AsSlice())
	source.Append(a)
	require.Equal(t, []validation.Violation{a, b}, destination.AsSlice())
}

func TestViolationListJoin_RepeatedAndSelf(t *testing.T) {
	a := newViolationWithError(t, errors.New("a"))
	for _, self := range []bool{false, true} {
		list := validation.NewViolationList(a)
		var result *validation.ViolationList
		if self {
			list.Join(list)
			result = list
		} else {
			var ok bool
			result, ok = validation.UnwrapViolations(validation.Filter(list, list))
			require.True(t, ok)
		}
		count := 0
		for range result.All() {
			count++
			require.LessOrEqual(t, count, 2, "cyclic list")
		}
		require.Equal(t, 2, count)
		require.Len(t, result.AsSlice(), 2)
	}
}

func TestJoinedErrors_PreserveAllBranches(t *testing.T) {
	a, b := newViolationWithError(t, errors.New("a")), newViolationWithError(t, errors.New("b"))
	fatal := errors.New("backend failure")
	cases := []struct {
		name  string
		err   error
		valid bool
		count int
	}{
		{"violations", errors.Join(a, b), true, 2},
		{"nested lists", fmt.Errorf("outer: %w", errors.Join(a, fmt.Errorf("inner: %w", validation.NewViolationList(b)))), true, 2},
		{"multi wrap", fmt.Errorf("both: %w and %w", a, b), true, 2},
		{"fatal last", errors.Join(a, fatal), false, 0},
		{"fatal first", errors.Join(fatal, a), false, 0},
		{"nested fatal", fmt.Errorf("outer: %w", errors.Join(validation.NewViolationList(a), errors.Join(b, fatal))), false, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			list, ok := validation.UnwrapViolations(tc.err)
			require.Equal(t, tc.valid, ok)
			if tc.valid {
				require.Equal(t, tc.count, list.Len())
				require.ErrorIs(t, list, a.Unwrap())
				require.ErrorIs(t, list, b.Unwrap())
			} else {
				require.Nil(t, list)
				require.ErrorIs(t, validation.Filter(tc.err), fatal)
				target := validation.NewViolationList()
				require.ErrorIs(t, target.AppendFromError(tc.err), fatal)
				require.Zero(t, target.Len())
			}
		})
	}
}

type parameterTranslator struct{}

func (parameterTranslator) Translate(tag language.Tag, message string, _ int) string {
	if message == "field" {
		if tag == language.Russian {
			return "поле"
		}
		return "field_en"
	}
	return message
}

func TestTranslatedParameters_CanBeReusedConcurrently(t *testing.T) {
	validator, err := validation.NewValidator(validation.SetTranslator(parameterTranslator{}))
	require.NoError(t, err)
	parameters := []validation.TemplateParameter{{Key: "{{ label }}", Value: "field", NeedsTranslation: true}}
	arg := validation.Check(false).WithMessage("{{ label }}", parameters...)
	var wg sync.WaitGroup
	for _, tc := range []struct {
		tag  language.Tag
		want string
	}{{language.English, "field_en"}, {language.Russian, "поле"}} {
		localized := validator.WithLanguage(tc.tag)
		// Sequential reuse must also retain the untranslated key.
		list, ok := validation.UnwrapViolations(localized.Validate(context.Background(), arg))
		require.True(t, ok)
		require.Equal(t, tc.want, list.First().Message())
	}
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				_ = validator.Validate(context.Background(), arg)
			}
		}()
	}
	wg.Wait()
	require.Equal(t, "field", parameters[0].Value)
}

type customViolationError struct {
	violation validation.Violation
	cause     error
}

func (e customViolationError) Error() string { return "custom violation" }
func (e customViolationError) Unwrap() error { return e.cause }
func (e customViolationError) As(target any) bool {
	if violation, ok := target.(*validation.Violation); ok {
		*violation = e.violation
		return true
	}
	return false
}

func TestUnwrapViolations_PreservesCustomAs(t *testing.T) {
	violation := newViolationWithError(t, errors.New("code"))
	custom := customViolationError{violation: violation, cause: errors.New("underlying code")}
	list, ok := validation.UnwrapViolations(fmt.Errorf("wrapped: %w", custom))
	require.True(t, ok)
	require.Equal(t, []validation.Violation{violation}, list.AsSlice())
	fatal := errors.New("backend failure")
	list, ok = validation.UnwrapViolations(errors.Join(custom, fatal))
	require.False(t, ok)
	require.Nil(t, list)
}
