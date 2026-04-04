package validate

import (
	"errors"
	"strings"
	"testing"
)

func TestCIDR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		opts    []func(*CIDROptions)
		wantErr error
	}{
		{name: "empty ok", value: "", wantErr: nil},
		{name: "ipv4 valid typical", value: "192.168.0.0/24", wantErr: nil},
		{name: "ipv6 valid typical", value: "2001:db8::/32", wantErr: nil},
		{name: "ipv4 prefix zero", value: "0.0.0.0/0", wantErr: nil},
		{name: "ipv4 prefix thirty two", value: "255.255.255.255/32", wantErr: nil},
		{name: "ipv6 prefix zero", value: "::/0", wantErr: nil},
		{name: "ipv6 prefix one hundred twenty eight", value: "2001:db8::1/128", wantErr: nil},
		{name: "ipv6 compressed valid", value: "fe80::1/64", wantErr: nil},
		{name: "no slash", value: "192.168.0.0", wantErr: ErrInvalidCIDR},
		{name: "empty ip part", value: "/24", wantErr: ErrInvalidCIDR},
		{name: "empty prefix", value: "192.168.0.0/", wantErr: ErrInvalidCIDR},
		{name: "non digit prefix", value: "192.168.0.0/xx", wantErr: ErrInvalidCIDR},
		{name: "prefix with slash in rest", value: "192.168.0.0/24/extra", wantErr: ErrInvalidCIDR},
		{name: "invalid ip", value: "999.999.999.999/24", wantErr: ErrInvalidCIDR},
		{name: "prefix too large ipv4 default max", value: "10.0.0.0/33", wantErr: ErrCIDRNetmaskOutOfRange},
		{name: "prefix too large ipv6 default max", value: "2001:db8::/129", wantErr: ErrCIDRNetmaskOutOfRange},
		{name: "prefix below min", value: "10.0.0.0/8", opts: []func(*CIDROptions){CIDRNetmaskRange(16, 32)}, wantErr: ErrCIDRNetmaskOutOfRange},
		{name: "prefix at min inclusive", value: "10.0.0.0/16", opts: []func(*CIDROptions){CIDRNetmaskRange(16, 24)}, wantErr: nil},
		{name: "prefix at max inclusive", value: "10.0.0.0/24", opts: []func(*CIDROptions){CIDRNetmaskRange(16, 24)}, wantErr: nil},
		{name: "prefix one below custom min", value: "10.0.0.0/15", opts: []func(*CIDROptions){CIDRNetmaskRange(16, 24)}, wantErr: ErrCIDRNetmaskOutOfRange},
		{name: "prefix one above custom max", value: "10.0.0.0/25", opts: []func(*CIDROptions){CIDRNetmaskRange(16, 24)}, wantErr: ErrCIDRNetmaskOutOfRange},
		{name: "ipv6 within custom range", value: "2001:db8::/48", opts: []func(*CIDROptions){CIDRNetmaskRange(32, 64)}, wantErr: nil},
		{name: "ipv6 above custom max", value: "2001:db8::/96", opts: []func(*CIDROptions){CIDRNetmaskRange(32, 64)}, wantErr: ErrCIDRNetmaskOutOfRange},
		{name: "ipv6 on ipv4 only", value: "2001:db8::/32", opts: []func(*CIDROptions){CIDRVersion("4")}, wantErr: ErrInvalidCIDR},
		{name: "ipv4 on ipv6 only", value: "10.0.0.0/8", opts: []func(*CIDROptions){CIDRVersion("6")}, wantErr: ErrInvalidCIDR},
		{name: "ipv4 only allows slash thirty two", value: "10.0.0.1/32", opts: []func(*CIDROptions){CIDRVersion("4")}, wantErr: nil},
		{name: "ipv6 only allows slash one hundred twenty eight", value: "fe80::1/128", opts: []func(*CIDROptions){CIDRVersion("6")}, wantErr: nil},
		{name: "unknown version option ignored ipv4 still valid", value: "10.0.0.0/8", opts: []func(*CIDROptions){CIDRVersion("bogus")}, wantErr: nil},
		{name: "unknown version option ignored ipv6 still valid", value: "2001:db8::/32", opts: []func(*CIDROptions){CIDRVersion("bogus")}, wantErr: nil},
		{name: "numeric prefix overflows int", value: "10.0.0.0/" + strings.Repeat("9", 30), wantErr: ErrInvalidCIDR},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := CIDR(tt.value, tt.opts...)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("CIDR(%q): got %v, want %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestCIDR_optionComposition(t *testing.T) {
	t.Parallel()
	// IPv4-only plus tight mask range: valid when both match.
	err := CIDR("192.168.0.0/24", CIDRVersion("4"), CIDRNetmaskRange(8, 24))
	if err != nil {
		t.Fatalf("CIDR: got %v, want nil", err)
	}
	// Same range rejects IPv6 on version filter before mask check.
	err = CIDR("2001:db8::/32", CIDRVersion("4"), CIDRNetmaskRange(8, 64))
	if !errors.Is(err, ErrInvalidCIDR) {
		t.Fatalf("CIDR: got %v, want %v", err, ErrInvalidCIDR)
	}
}

func TestCIDRViolationNetmaskBounds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		value  string
		opts   []func(*CIDROptions)
		wantLo int
		wantHi int
	}{
		{
			name:   "ipv4 caps hi to thirty two when configured max above thirty two",
			value:  "10.0.0.0/99",
			opts:   []func(*CIDROptions){CIDRNetmaskRange(0, 128)},
			wantLo: 0,
			wantHi: 32,
		},
		{
			name:   "ipv6 does not cap hi when max above thirty two",
			value:  "2001:db8::/64",
			opts:   []func(*CIDROptions){CIDRNetmaskRange(0, 128)},
			wantLo: 0,
			wantHi: 128,
		},
		{
			name:   "ipv6 preserves large configured max",
			value:  "2001:db8::/64",
			opts:   []func(*CIDROptions){CIDRNetmaskRange(0, 200)},
			wantLo: 0,
			wantHi: 200,
		},
		{
			name:   "no slash returns configured bounds without parsing ip",
			value:  "not-a-cidr",
			opts:   []func(*CIDROptions){CIDRNetmaskRange(10, 48)},
			wantLo: 10,
			wantHi: 48,
		},
		{
			name:   "unparseable ip left part does not trigger ipv4 cap",
			value:  "not-an-ip/24",
			opts:   []func(*CIDROptions){CIDRNetmaskRange(0, 128)},
			wantLo: 0,
			wantHi: 128,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			lo, hi := CIDRViolationNetmaskBounds(tt.value, tt.opts...)
			if lo != tt.wantLo || hi != tt.wantHi {
				t.Fatalf("CIDRViolationNetmaskBounds(%q): got (%d,%d), want (%d,%d)", tt.value, lo, hi, tt.wantLo, tt.wantHi)
			}
		})
	}
}
