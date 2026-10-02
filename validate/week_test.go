package validate_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/require"
)

func TestWeek(t *testing.T) {
	tests := []struct {
		value string
		want  error
	}{
		{"", nil},
		{"2020-W53", nil},
		{"2015-W53", nil},
		{"2026-W53", nil},
		{"2021-W01", nil},
		{"2021-W52", nil},
		{"0000-W01", nil},
		{"0001-W01", nil},
		{"9999-W52", nil},
		{"2021-W53", validate.ErrInvalidWeekNumber},
		{"2022-W53", validate.ErrInvalidWeekNumber},
		{"2020-W00", validate.ErrInvalidWeekFormat},
		{"2020-W54", validate.ErrInvalidWeekFormat},
		{"2020-W1", validate.ErrInvalidWeekFormat},
		{"2020-w01", validate.ErrInvalidWeekFormat},
		{"20-W01", validate.ErrInvalidWeekFormat},
		{"10000-W01", validate.ErrInvalidWeekFormat},
		{"-001-W01", validate.ErrInvalidWeekFormat},
		{"２０２０-W01", validate.ErrInvalidWeekFormat},
		{"2020-W０１", validate.ErrInvalidWeekFormat},
		{"2020-W01\n", validate.ErrInvalidWeekFormat},
		{" 2020-W01", validate.ErrInvalidWeekFormat},
		{"2020-W01 ", validate.ErrInvalidWeekFormat},
		{"2020-W01-1", validate.ErrInvalidWeekFormat},
		{"2020-01-01", validate.ErrInvalidWeekFormat},
		{" ", validate.ErrInvalidWeekFormat},
		{"\xff", validate.ErrInvalidWeekFormat},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) { require.ErrorIs(t, validate.Week(tt.value), tt.want) })
	}
}

func TestWeekBounds(t *testing.T) {
	tests := []struct {
		name, value, min, max string
		want                  error
	}{
		{"min equal", "2020-W53", "2020-W53", "", nil},
		{"max equal", "2021-W01", "", "2021-W01", nil},
		{"cross year below", "2020-W53", "2021-W01", "", validate.ErrWeekTooEarly},
		{"cross year above", "2021-W01", "", "2020-W53", validate.ErrWeekTooLate},
		{"same year below", "2020-W01", "2020-W02", "", validate.ErrWeekTooEarly},
		{"same year above", "2020-W03", "", "2020-W02", validate.ErrWeekTooLate},
		{"exact", "2020-W53", "2020-W53", "2020-W53", nil},
		{"empty value", "", "2020-W01", "2021-W01", nil},
		{"invalid min", "", "2020-W00", "", validate.ErrInvalidWeekBounds},
		{"nonexistent min", "2020-W01", "2021-W53", "", validate.ErrInvalidWeekBounds},
		{"invalid max", "2020-W01", "", "oops", validate.ErrInvalidWeekBounds},
		{"nonexistent max", "2020-W01", "", "2021-W53", validate.ErrInvalidWeekBounds},
		{"reversed", "", "2021-W01", "2020-W53", validate.ErrInvalidWeekBounds},
		{"format before bounds", "oops", "2020-W01", "2021-W01", validate.ErrInvalidWeekFormat},
		{"number before bounds", "2021-W53", "2020-W01", "2021-W01", validate.ErrInvalidWeekNumber},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorIs(t, validate.Week(tt.value, validate.WithMinWeek(tt.min), validate.WithMaxWeek(tt.max)), tt.want)
		})
	}
	require.NoError(t, validate.Week("2020-W01", validate.WithMinWeek("2021-W01"), validate.WithMinWeek("")))
	require.NoError(t, validate.Week("2021-W01", validate.WithMaxWeek("2020-W01"), validate.WithMaxWeek("")))
}

// Check every literal year against the ISO rule: a year has 53 weeks when
// January 1 is Thursday, or Wednesday in a leap year.
func TestWeek53AllYears(t *testing.T) {
	for year := 0; year <= 9999; year++ {
		weekday := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC).Weekday()
		leap := year%4 == 0 && (year%100 != 0 || year%400 == 0)
		want := validate.ErrInvalidWeekNumber
		if weekday == time.Thursday || weekday == time.Wednesday && leap {
			want = nil
		}
		require.ErrorIs(t, validate.Week(fmt.Sprintf("%04d-W53", year)), want, "year %d", year)
	}
}

func ExampleWeek() {
	fmt.Println(validate.Week("2020-W53"))
	fmt.Println(validate.Week("2021-W53"))
	// Output:
	// <nil>
	// invalid week number
}
