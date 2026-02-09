package it_test

import (
	"context"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validationtest"
	"github.com/muonsoft/validation/validator"
)

// Test that it constraints implement validation.Constraint and work with validation.Each and validation.This.
func TestItConstraints_ImplementConstraint_WhenUsedWithEach(t *testing.T) {
	tests := []struct {
		name       string
		items      []string
		constraint validation.Constraint[string]
		wantErr    bool
		wantPath   string
	}{
		{
			name:       "IsNotBlank with blank",
			items:      []string{""},
			constraint: it.IsNotBlank(),
			wantErr:    true,
			wantPath:   "[0]",
		},
		{
			name:       "HasMinLength too short",
			items:      []string{"ab"},
			constraint: it.HasMinLength(3),
			wantErr:    true,
			wantPath:   "[0]",
		},
		{
			name:       "IsOneOf invalid choice",
			items:      []string{"x"},
			constraint: it.IsOneOf("a", "b"),
			wantErr:    true,
			wantPath:   "[0]",
		},
		{
			name:       "IsNotBlank valid",
			items:      []string{"ok"},
			constraint: it.IsNotBlank(),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(context.Background(), validation.Each(tt.items, tt.constraint))
			if tt.wantErr {
				validationtest.Assert(t, err).IsViolationList().WithOneViolation().WithPropertyPath(tt.wantPath)
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestItConstraints_ImplementConstraint_WhenUsedWithThis(t *testing.T) {
	// NotBlankConstraint[string].Validate via This
	err := validator.Validate(context.Background(), validation.This("", it.IsNotBlank()))
	validationtest.Assert(t, err).IsViolationList().WithOneViolation()

	err = validator.Validate(context.Background(), validation.This("x", it.IsNotBlank()))
	if err != nil {
		t.Errorf("expected no error for valid value, got %v", err)
	}
}

func TestNotBlankConstraint_Validate_MatchesValidateString(t *testing.T) {
	// Validate(ctx, v, x) should behave like ValidateString(ctx, v, &x)
	err1 := validator.Validate(context.Background(), validation.String("", it.IsNotBlank()))
	err2 := validator.Validate(context.Background(), validation.This("", it.IsNotBlank()))
	validationtest.Assert(t, err1).IsViolationList().WithOneViolation()
	validationtest.Assert(t, err2).IsViolationList().WithOneViolation()
}

func TestLengthConstraint_Validate_MatchesValidateString(t *testing.T) {
	err1 := validator.Validate(context.Background(), validation.String("ab", it.HasMinLength(3)))
	err2 := validator.Validate(context.Background(), validation.This("ab", it.HasMinLength(3)))
	validationtest.Assert(t, err1).IsViolationList().WithOneViolation()
	validationtest.Assert(t, err2).IsViolationList().WithOneViolation()
}

func TestChoiceConstraint_Validate_MatchesValidateComparable(t *testing.T) {
	err1 := validator.Validate(context.Background(), validation.Comparable[string]("x", it.IsOneOf("a", "b")))
	err2 := validator.Validate(context.Background(), validation.This("x", it.IsOneOf("a", "b")))
	validationtest.Assert(t, err1).IsViolationList().WithOneViolation()
	validationtest.Assert(t, err2).IsViolationList().WithOneViolation()
}
