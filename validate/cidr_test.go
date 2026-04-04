package validate

import (
	"errors"
	"testing"
)

func TestCidr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		opts    []func(*CidrOptions)
		wantErr error
	}{
		{name: "empty ok", value: "", wantErr: nil},
		{name: "ipv4 valid", value: "192.168.0.0/24", wantErr: nil},
		{name: "ipv6 valid", value: "2001:db8::/32", wantErr: nil},
		{name: "no slash", value: "192.168.0.0", wantErr: ErrInvalidCidr},
		{name: "empty ip part", value: "/24", wantErr: ErrInvalidCidr},
		{name: "empty prefix", value: "192.168.0.0/", wantErr: ErrInvalidCidr},
		{name: "non digit prefix", value: "192.168.0.0/xx", wantErr: ErrInvalidCidr},
		{name: "invalid ip", value: "999.999.999.999/24", wantErr: ErrInvalidCidr},
		{name: "prefix too large ipv4", value: "10.0.0.0/33", wantErr: ErrCidrNetmaskOutOfRange},
		{name: "prefix below min", value: "10.0.0.0/8", opts: []func(*CidrOptions){CidrNetmaskRange(16, 32)}, wantErr: ErrCidrNetmaskOutOfRange},
		{name: "ipv6 on ipv4 only", value: "2001:db8::/32", opts: []func(*CidrOptions){CidrVersion("4")}, wantErr: ErrInvalidCidr},
		{name: "ipv4 on ipv6 only", value: "10.0.0.0/8", opts: []func(*CidrOptions){CidrVersion("6")}, wantErr: ErrInvalidCidr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Cidr(tt.value, tt.opts...)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Cidr(%q): got %v, want %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestCidrViolationNetmaskBounds_ipv4CapsMax(t *testing.T) {
	t.Parallel()
	lo, hi := CidrViolationNetmaskBounds("10.0.0.0/99", CidrNetmaskRange(0, 128))
	if lo != 0 || hi != 32 {
		t.Fatalf("got lo=%d hi=%d, want 0, 32", lo, hi)
	}
}
