package validate

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

var (
	// ErrInvalidWeekFormat indicates a value outside the YYYY-Www format (weeks 01–53).
	ErrInvalidWeekFormat = errors.New("invalid week format")
	// ErrInvalidWeekNumber indicates week 53 in an ISO year with only 52 weeks.
	ErrInvalidWeekNumber = errors.New("invalid week number")
	// ErrWeekTooEarly indicates a week before the inclusive minimum.
	ErrWeekTooEarly = errors.New("week too early")
	// ErrWeekTooLate indicates a week after the inclusive maximum.
	ErrWeekTooLate = errors.New("week too late")
	// ErrInvalidWeekBounds indicates malformed, nonexistent, or reversed bounds.
	ErrInvalidWeekBounds = errors.New("invalid week bounds")
)

// WeekOptions configures inclusive ISO week bounds for [Week].
type WeekOptions struct{ min, max string }

// WithMinWeek sets the inclusive minimum in YYYY-Www format. Empty removes the bound.
func WithMinWeek(value string) func(*WeekOptions) {
	return func(o *WeekOptions) { o.min = value }
}

// WithMaxWeek sets the inclusive maximum in YYYY-Www format. Empty removes the bound.
func WithMaxWeek(value string) func(*WeekOptions) {
	return func(o *WeekOptions) { o.max = value }
}

var weekPattern = regexp.MustCompile(`^[0-9]{4}-W(0[1-9]|[1-4][0-9]|5[0-3])$`)

// Week validates an ISO week in YYYY-Www format, using literal years 0000–9999
// in the proleptic Gregorian calendar. Week 53 must exist in the specified year.
// Empty values are valid; use [NotBlank] to require a value.
// Invalid options return [ErrInvalidWeekBounds], even for empty input.
// Other failures return [ErrInvalidWeekFormat], [ErrInvalidWeekNumber],
// [ErrWeekTooEarly], or [ErrWeekTooLate].
func Week(value string, options ...func(*WeekOptions)) error {
	var o WeekOptions
	for _, option := range options {
		option(&o)
	}
	if !o.valid() {
		return ErrInvalidWeekBounds
	}
	if value == "" {
		return nil
	}
	if err := checkWeek(value); err != nil {
		return err
	}
	if o.min != "" && value < o.min {
		return ErrWeekTooEarly
	}
	if o.max != "" && value > o.max {
		return ErrWeekTooLate
	}
	return nil
}

func (o WeekOptions) valid() bool {
	if o.min != "" && checkWeek(o.min) != nil {
		return false
	}
	if o.max != "" && checkWeek(o.max) != nil {
		return false
	}
	return o.min == "" || o.max == "" || o.min <= o.max
}

func checkWeek(value string) error {
	if !weekPattern.MatchString(value) {
		return ErrInvalidWeekFormat
	}
	year, _ := strconv.Atoi(value[:4])
	week, _ := strconv.Atoi(value[6:])
	_, lastWeek := time.Date(year, time.December, 28, 0, 0, 0, 0, time.UTC).ISOWeek()
	if week > lastWeek {
		return ErrInvalidWeekNumber
	}
	return nil
}
