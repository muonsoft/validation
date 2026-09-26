package validate_test

import (
	"errors"
	"testing"

	"github.com/muonsoft/validation/validate"
)

func TestTimezone(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		options []func(*validate.TimezoneOptions)
		wantErr error
	}{
		{name: "empty", value: ""},
		{name: "UTC", value: "UTC"},
		{name: "Europe/Berlin", value: "Europe/Berlin"},
		{name: "America/New_York", value: "America/New_York"},
		{name: "Etc/GMT+5", value: "Etc/GMT+5"},
		{name: "nested identifier", value: "America/Argentina/Buenos_Aires"},
		{name: "legacy alias", value: "US/Eastern"},
		{name: "legacy regional alias", value: "Europe/Kiev"},
		{name: "invalid empty segment", value: "Europe//Berlin", wantErr: validate.ErrInvalidTimezone},
		{name: "invalid dot segment", value: "Europe/./Berlin", wantErr: validate.ErrInvalidTimezone},
		{name: "invalid relative UTC", value: "./UTC", wantErr: validate.ErrInvalidTimezone},
		{name: "invalid posix prefix", value: "posix/Europe/Berlin", wantErr: validate.ErrInvalidTimezone},
		{name: "invalid right prefix", value: "right/Europe/Berlin", wantErr: validate.ErrInvalidTimezone},
		{name: "invalid absolute path", value: "/Europe/Berlin", wantErr: validate.ErrInvalidTimezone},
		{name: "invalid parent segment", value: "Europe/../Europe/Berlin", wantErr: validate.ErrInvalidTimezone},
		{name: "invalid trailing slash", value: "Europe/Berlin/", wantErr: validate.ErrInvalidTimezone},
		{name: "invalid case", value: "europe/berlin", wantErr: validate.ErrInvalidTimezone},
		{name: "invalid unknown", value: "Invalid/Zone", wantErr: validate.ErrInvalidTimezone},
		{name: "invalid Local", value: "Local", wantErr: validate.ErrInvalidTimezone},
		{name: "invalid Factory", value: "Factory", wantErr: validate.ErrInvalidTimezone},
		{name: "invalid abbreviation EST", value: "EST", wantErr: validate.ErrInvalidTimezone},
		{name: "invalid abbreviation GMT", value: "GMT", wantErr: validate.ErrInvalidTimezone},
		{name: "invalid no slash", value: "Berlin", wantErr: validate.ErrInvalidTimezone},
		{
			name:    "zone Europe valid",
			value:   "Europe/Paris",
			options: []func(*validate.TimezoneOptions){validate.WithTimezoneZone(validate.TimezoneZoneEurope)},
		},
		{
			name:    "zone Europe invalid",
			value:   "America/New_York",
			options: []func(*validate.TimezoneOptions){validate.WithTimezoneZone(validate.TimezoneZoneEurope)},
			wantErr: validate.ErrInvalidTimezone,
		},
		{
			name:    "zone America valid",
			value:   "America/Chicago",
			options: []func(*validate.TimezoneOptions){validate.WithTimezoneZone(validate.TimezoneZoneAmerica)},
		},
		{
			name:    "zone Asia valid",
			value:   "Asia/Tokyo",
			options: []func(*validate.TimezoneOptions){validate.WithTimezoneZone(validate.TimezoneZoneAsia)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Timezone(tt.value, tt.options...)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("Timezone(%q): %v", tt.value, err)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Timezone(%q): got %v, want %v", tt.value, err, tt.wantErr)
			}
		})
	}
}
