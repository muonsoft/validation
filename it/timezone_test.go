package it_test

import (
	"context"
	"errors"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validate"
)

func TestTimezoneConstraint_WithZoneCopiesAreIndependent(t *testing.T) {
	v, err := validation.NewValidator()
	if err != nil {
		t.Fatal(err)
	}
	base := it.IsTimezone().
		WithZone(validate.TimezoneZoneAsia).
		WithZone(validate.TimezoneZoneAfrica).
		WithZone(validate.TimezoneZoneAll)
	europe := base.WithZone(validate.TimezoneZoneEurope)
	america := base.WithZone(validate.TimezoneZoneAmerica)

	for _, tc := range []struct {
		name       string
		constraint it.TimezoneConstraint
		value      string
		wantErr    bool
	}{
		{"base accepts Europe", base, "Europe/Berlin", false},
		{"base accepts America", base, "America/Chicago", false},
		{"Europe accepts Europe", europe, "Europe/Berlin", false},
		{"Europe rejects America", europe, "America/Chicago", true},
		{"America accepts America", america, "America/Chicago", false},
		{"America rejects Europe", america, "Europe/Berlin", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.constraint.Validate(context.Background(), v, tc.value)
			if tc.wantErr {
				if !errors.Is(err, validation.ErrInvalidTimezone) {
					t.Fatalf("got %v, want ErrInvalidTimezone", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}
