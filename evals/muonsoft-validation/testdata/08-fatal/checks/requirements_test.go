package scenario_test

import (
	"context"
	"errors"
	s "example.com/validation-eval/scenario"
	"fmt"
	"github.com/muonsoft/validation"
	"testing"
)

func TestFatalStops(t *testing.T) {
	for _, fatal := range []error{context.Canceled, context.DeadlineExceeded, errors.New("connection lost")} {
		thirdCalls, writes := 0, 0
		ctx := t.Context()
		first := validation.ValidatableFunc(func(ctx context.Context, v *validation.Validator) error {
			return v.CreateViolation(ctx, validation.ErrNotValid, validation.ErrNotValid.Message())
		})
		second := validation.ValidatableFunc(func(got context.Context, _ *validation.Validator) error {
			if got != ctx {
				t.Fatal("context replaced")
			}
			return fmt.Errorf("lookup: %w", fatal)
		})
		third := validation.ValidatableFunc(func(context.Context, *validation.Validator) error { thirdCalls++; return nil })
		err := s.ValidateAndSave(ctx, newValidator(t), first, second, third, func(context.Context) error { writes++; return nil })
		if !errors.Is(err, fatal) {
			t.Fatalf("lost error: %v", err)
		}
		if _, ok := validation.UnwrapViolations(err); ok {
			t.Fatal("masked technical error")
		}
		if thirdCalls != 0 || writes != 0 {
			t.Fatal("work continued after failure")
		}
	}
}
func TestSuccessfulSequence(t *testing.T) {
	calls, writes := 0, 0
	ctx := t.Context()
	check := validation.ValidatableFunc(func(got context.Context, _ *validation.Validator) error {
		if got != ctx {
			t.Fatal("context replaced")
		}
		calls++
		return nil
	})
	err := s.ValidateAndSave(ctx, newValidator(t), check, check, check, func(got context.Context) error {
		if got != ctx {
			t.Fatal("context replaced")
		}
		writes++
		return nil
	})
	wantViolations(t, err, nil)
	if calls != 3 || writes != 1 {
		t.Fatal("wrong calls")
	}
}
func TestViolationsAccumulate(t *testing.T) {
	calls, writes := 0, 0
	check := validation.ValidatableFunc(func(ctx context.Context, v *validation.Validator) error {
		calls++
		return v.AtIndex(calls-1).CreateViolation(ctx, validation.ErrNotValid, validation.ErrNotValid.Message())
	})
	err := s.ValidateAndSave(t.Context(), newValidator(t), check, check, check, func(context.Context) error { writes++; return nil })
	wantViolations(t, err, map[string]error{"[0]": validation.ErrNotValid, "[1]": validation.ErrNotValid, "[2]": validation.ErrNotValid})
	if writes != 0 {
		t.Fatal("saved invalid input")
	}
}
