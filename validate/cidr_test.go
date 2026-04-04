package validate

import (
	"errors"
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
		{name: "ipv4 valid", value: "192.168.0.0/24", wantErr: nil},
		{name: "ipv6 valid", value: "2001:db8::/32", wantErr: nil},
		{name: "no slash", value: "192.168.0.0", wantErr: ErrInvalidCIDR},
		{name: "empty ip part", value: "/24", wantErr: ErrInvalidCIDR},
		{name: "empty prefix", value: "192.168.0.0/", wantErr: ErrInvalidCIDR},
		{name: "non digit prefix", value: "192.168.0.0/xx", wantErr: ErrInvalidCIDR},
		{name: "invalid ip", value: "999.999.999.999/24", wantErr: ErrInvalidCIDR},
		{name: "prefix too large ipv4", value: "10.0.0.0/33", wantErr: ErrCIDRNetmaskOutOfRange},
		{name: "prefix below min", value: "10.0.0.0/8", opts: []func(*CIDROptions){CIDRNetmaskRange(16, 32)}, wantErr: ErrCIDRNetmaskOutOfRange},
		{name: "ipv6 on ipv4 only", value: "2001:db8::/32", opts: []func(*CIDROptions){CIDRVersion("4")}, wantErr: ErrInvalidCIDR},
		{name: "ipv4 on ipv6 only", value: "10.0.0.0/8", opts: []func(*CIDROptions){CIDRVersion("6")}, wantErr: ErrInvalidCIDR},
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

func TestCIDRViolationNetmaskBounds_ipv4CapsMax(t *testing.T) {
	t.Parallel()
	lo, hi := CIDRViolationNetmaskBounds("10.0.0.0/99", CIDRNetmaskRange(0, 128))
	if lo != 0 || hi != 32 {
		t.Fatalf("got lo=%d hi=%d, want 0, 32", lo, hi)
	}
}
