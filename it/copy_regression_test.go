package it_test

import (
	"context"
	"net"
	"net/url"
	"testing"

	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/require"
)

func TestConstraints_DerivedOptionsAreIndependent(t *testing.T) {
	validator, err := validation.NewValidator()
	require.NoError(t, err)
	allowURL := func(*url.URL) bool { return true }
	allowIP := func(net.IP) bool { return false }
	urlBase := it.IsURL().WithRestriction(allowURL).WithRestriction(allowURL).WithRestriction(allowURL)
	deniedURL := urlBase.WithRestriction(func(*url.URL) bool { return false })
	allowedURL := urlBase.WithRestriction(allowURL)
	ipBase := it.IsIP().DenyIP(allowIP).DenyIP(allowIP).DenyIP(allowIP)
	deniedIP := ipBase.DenyPrivateIP()
	allowedIP := ipBase.DenyIP(allowIP)
	cidrBase := it.IsCIDR().WithVersion("all").WithVersion("all").WithVersion("all")
	ipv4 := cidrBase.IPv4Only()
	ipv6 := cidrBase.IPv6Only()
	uuidBase := it.IsUUID().WithVersions(4).WithVersions(4).WithVersions(4)
	uuid4 := uuidBase.WithVersions(4)
	uuid1 := uuidBase.WithVersions(1)
	isbnBase := it.IsISBN().Only10().Only10().Only10()
	isbn10 := isbnBase.Only10()
	isbn13 := isbnBase.Only13()
	macBase := it.IsMacAddress().WithType(validate.MacAddressTypeAll).WithType(validate.MacAddressTypeAll).WithType(validate.MacAddressTypeAll)
	unicast := macBase.WithType(validate.MacAddressTypeUnicastAll)
	multicast := macBase.WithType(validate.MacAddressTypeMulticastAll)
	cases := []struct {
		name, value string
		constraint  validation.StringConstraint
		invalid     bool
	}{
		{"MAC unicast", "00:11:22:33:44:55", unicast, false},
		{"MAC multicast rejects unicast", "00:11:22:33:44:55", multicast, true},
		{"URL deny", "https://example.com", deniedURL, true},
		{"URL allow", "https://example.com", allowedURL, false},
		{"IP deny private", "192.168.1.1", deniedIP, true},
		{"IP allow private", "192.168.1.1", allowedIP, false},
		{"CIDR IPv4", "192.168.1.0/24", ipv4, false},
		{"CIDR IPv6 rejects IPv4", "192.168.1.0/24", ipv6, true},
		{"UUID version 4", "550e8400-e29b-41d4-a716-446655440000", uuid4, false},
		{"UUID version 1 rejects version 4", "550e8400-e29b-41d4-a716-446655440000", uuid1, true},
		{"ISBN 10", "0306406152", isbn10, false},
		{"ISBN 13 rejects ISBN 10", "0306406152", isbn13, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateString(context.Background(), tc.value, tc.constraint)
			if tc.invalid {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDivisibility_ZeroDivisorReturnsConfigurationError(t *testing.T) {
	validator, err := validation.NewValidator()
	require.NoError(t, err)
	cases := []struct {
		name string
		arg  validation.Argument
	}{
		{"integer", validation.Number(1, it.IsDivisibleBy(0))},
		{"unsigned", validation.Number(uint64(1), it.IsDivisibleBy(uint64(0)))},
		{"float", validation.Number(1.0, it.IsDivisibleByFloat(0.0))},
		{"count", validation.Countable(1, it.HasCountDivisibleBy(0))},
		{"empty count", validation.Countable(0, it.HasCountDivisibleBy(0))},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var configErr *validation.ConstraintError
			require.ErrorAs(t, validator.Validate(context.Background(), tc.arg), &configErr)
		})
	}
	require.ErrorIs(t, validator.ValidateFloat(context.Background(), -5, it.IsDivisibleByFloat(2.0)), validation.ErrNotDivisible)
	require.NoError(t, validator.ValidateFloat(context.Background(), -6, it.IsDivisibleByFloat(2.0)))
}
