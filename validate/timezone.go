package validate

import (
	"errors"
	"strings"
	"time"
)

// ErrInvalidTimezone is returned by [Timezone] when the value is not a valid IANA timezone identifier.
var ErrInvalidTimezone = errors.New("invalid timezone")

// TimezoneZone restricts valid timezone identifiers to a geographical region,
// aligned with Symfony\Component\Validator\Constraints\Timezone zone option
// (PHP \DateTimeZone::AFRICA, ::AMERICA, and similar constants).
type TimezoneZone string

const (
	// TimezoneZoneAll accepts any IANA timezone identifier known to [time.LoadLocation],
	// excluding implementation-specific names such as "Local".
	TimezoneZoneAll TimezoneZone = ""
	// TimezoneZoneAfrica restricts identifiers to the Africa region (e.g. "Africa/Nairobi").
	TimezoneZoneAfrica TimezoneZone = "Africa"
	// TimezoneZoneAmerica restricts identifiers to the America region (e.g. "America/New_York").
	TimezoneZoneAmerica TimezoneZone = "America"
	// TimezoneZoneAntarctica restricts identifiers to the Antarctica region.
	TimezoneZoneAntarctica TimezoneZone = "Antarctica"
	// TimezoneZoneArctic restricts identifiers to the Arctic region.
	TimezoneZoneArctic TimezoneZone = "Arctic"
	// TimezoneZoneAsia restricts identifiers to the Asia region (e.g. "Asia/Tokyo").
	TimezoneZoneAsia TimezoneZone = "Asia"
	// TimezoneZoneAtlantic restricts identifiers to the Atlantic region.
	TimezoneZoneAtlantic TimezoneZone = "Atlantic"
	// TimezoneZoneAustralia restricts identifiers to the Australia region.
	TimezoneZoneAustralia TimezoneZone = "Australia"
	// TimezoneZoneEurope restricts identifiers to the Europe region (e.g. "Europe/Berlin").
	TimezoneZoneEurope TimezoneZone = "Europe"
	// TimezoneZoneIndian restricts identifiers to the Indian Ocean region.
	TimezoneZoneIndian TimezoneZone = "Indian"
	// TimezoneZonePacific restricts identifiers to the Pacific region.
	TimezoneZonePacific TimezoneZone = "Pacific"
)

// TimezoneOptions configures [Timezone] validation.
type TimezoneOptions struct {
	zone TimezoneZone
}

func newTimezoneOptions() TimezoneOptions {
	return TimezoneOptions{zone: TimezoneZoneAll}
}

// WithTimezoneZone sets the geographical region filter (default [TimezoneZoneAll]).
func WithTimezoneZone(zone TimezoneZone) func(*TimezoneOptions) {
	return func(o *TimezoneOptions) {
		o.zone = zone
	}
}

// Timezone validates whether the value is a known IANA timezone identifier,
// as in Symfony\Component\Validator\Constraints\Timezone.
//
// Validation uses [time.LoadLocation] and requires canonical IANA-style identifiers:
// either exactly "UTC" or a name containing "/" (e.g. "Europe/Berlin", "Etc/GMT+5").
// Implementation-specific names such as "Local", bare abbreviations (e.g. "EST"), and
// unknown identifiers are rejected.
//
// Empty string is considered valid (use [NotBlank] or similar to reject empty values).
//
// Possible errors:
//   - [ErrInvalidTimezone] when the string is unknown, not IANA-shaped, or outside the configured zone.
func Timezone(value string, options ...func(*TimezoneOptions)) error {
	if value == "" {
		return nil
	}
	if !isIANATimezoneIdentifier(value) {
		return ErrInvalidTimezone
	}

	opts := newTimezoneOptions()
	for _, opt := range options {
		opt(&opts)
	}
	if opts.zone != TimezoneZoneAll && !strings.HasPrefix(value, string(opts.zone)+"/") {
		return ErrInvalidTimezone
	}
	return nil
}

func isIANATimezoneIdentifier(value string) bool {
	switch value {
	case "Local", "Factory":
		return false
	}
	if _, err := time.LoadLocation(value); err != nil {
		return false
	}
	if value == "UTC" {
		return true
	}
	return strings.Contains(value, "/")
}
