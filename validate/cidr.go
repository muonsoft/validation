package validate

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

var (
	// ErrInvalidCidr is returned when the value is not valid CIDR notation (missing slash,
	// non-digit prefix, invalid IP, or IP version does not match the configured version).
	ErrInvalidCidr = errors.New("invalid CIDR")
	// ErrCidrNetmaskOutOfRange is returned when the prefix length is outside the allowed
	// [netmaskMin, netmaskMax] range (after capping max to 32 for IPv4 when needed),
	// matching Symfony\Component\Validator\Constraints\CidrValidator.
	ErrCidrNetmaskOutOfRange = errors.New("CIDR netmask out of range")
)

// CidrOptions configures [Cidr] validation.
type CidrOptions struct {
	version    string
	netmaskMin int
	netmaskMax int
}

const (
	cidrVersion4   = "4"
	cidrVersion6   = "6"
	cidrVersionAll = "all"
)

func newCidrOptions() CidrOptions {
	return CidrOptions{
		version:    cidrVersionAll,
		netmaskMin: 0,
		netmaskMax: 128,
	}
}

// CidrVersion sets which IP versions are accepted: "4" (IPv4 only), "6" (IPv6 only), or "all".
// The default is "all". Invalid values are ignored (default remains).
func CidrVersion(version string) func(*CidrOptions) {
	return func(o *CidrOptions) {
		switch version {
		case cidrVersion4, cidrVersion6, cidrVersionAll:
			o.version = version
		}
	}
}

// CidrNetmaskRange sets the inclusive allowed range for the CIDR prefix length.
// Defaults are 0–128 for "all", 0–32 for IPv4-only, 0–128 for IPv6-only (Symfony Cidr defaults).
func CidrNetmaskRange(netmaskMin, netmaskMax int) func(*CidrOptions) {
	return func(o *CidrOptions) {
		o.netmaskMin = netmaskMin
		o.netmaskMax = netmaskMax
	}
}

// CidrViolationNetmaskBounds returns the configured inclusive lower bound and the effective upper bound
// for the prefix length in Symfony-style out-of-range messages (caps upper bound at 32 for IPv4 when netmaskMax > 32).
func CidrViolationNetmaskBounds(value string, options ...func(*CidrOptions)) (lo, hi int) {
	o := newCidrOptions()
	for _, opt := range options {
		opt(&o)
	}
	lo = o.netmaskMin
	hi = o.netmaskMax
	ipStr, _, ok := strings.Cut(value, "/")
	if !ok {
		return lo, hi
	}
	if ip := net.ParseIP(ipStr); ip != nil && cidrIPVersion(ip) == 4 && hi > 32 {
		hi = 32
	}
	return lo, hi
}

// Cidr validates that value is a valid CIDR notation (IP/prefix), aligned with
// Symfony\Component\Validator\Constraints\Cidr and CidrValidator.
//
// Empty string is valid (use [NotBlank] or similar to reject empty values).
//
// Possible errors:
//   - [ErrInvalidCidr] for malformed notation, invalid IP, or version mismatch;
//   - [ErrCidrNetmaskOutOfRange] when the prefix is outside the configured netmask range
//     (for IPv4 addresses, effective max is capped at 32 if netmaskMax > 32, like Symfony).
func Cidr(value string, options ...func(*CidrOptions)) error {
	if value == "" {
		return nil
	}

	opts := newCidrOptions()
	for _, opt := range options {
		opt(&opts)
	}

	ipStr, bitsStr, err := cidrSplitNotation(value)
	if err != nil {
		return err
	}
	prefix, err := cidrParsePrefix(bitsStr)
	if err != nil {
		return err
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ErrInvalidCidr
	}

	ver := cidrIPVersion(ip)
	if cidrVersionMismatch(opts.version, ver) {
		return ErrInvalidCidr
	}

	return cidrCheckPrefixRange(prefix, ver, opts)
}

func cidrCheckPrefixRange(prefix, ver int, opts CidrOptions) error {
	maxMask := opts.netmaskMax
	if ver == 4 && maxMask > 32 {
		maxMask = 32
	}
	if prefix < opts.netmaskMin || prefix > maxMask {
		return ErrCidrNetmaskOutOfRange
	}
	return nil
}

func cidrSplitNotation(value string) (ipStr, bitsStr string, err error) {
	ipStr, bitsStr, ok := strings.Cut(value, "/")
	if !ok || ipStr == "" || bitsStr == "" {
		return "", "", ErrInvalidCidr
	}
	return ipStr, bitsStr, nil
}

func cidrParsePrefix(bitsStr string) (int, error) {
	if !isDecimalString(bitsStr) {
		return 0, ErrInvalidCidr
	}
	prefix, err := strconv.Atoi(bitsStr)
	if err != nil {
		return 0, ErrInvalidCidr
	}
	return prefix, nil
}

func cidrVersionMismatch(version string, ver int) bool {
	switch version {
	case cidrVersion4:
		return ver != 4
	case cidrVersion6:
		return ver != 6
	default:
		return false
	}
}

func cidrIPVersion(ip net.IP) int {
	if ip.To4() != nil {
		return 4
	}
	return 6
}

func isDecimalString(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
