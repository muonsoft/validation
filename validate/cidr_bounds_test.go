package validate_test

import (
	"testing"

	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/assert"
)

func TestCIDR_PrefixCannotExceedAddressSize(t *testing.T) {
	for _, tc := range []struct {
		value string
		max   int
		valid bool
	}{
		{"::1/128", 256, true},
		{"::1/129", 256, false},
		{"::1/256", 256, false},
		{"192.168.0.1/32", 256, true},
		{"192.168.0.1/33", 256, false},
	} {
		err := validate.CIDR(tc.value, validate.CIDRNetmaskRange(0, tc.max))
		if tc.valid {
			assert.NoError(t, err)
		} else {
			assert.ErrorIs(t, err, validate.ErrCIDRNetmaskOutOfRange)
		}
	}
	_, hi := validate.CIDRViolationNetmaskBounds("::1/129", validate.CIDRNetmaskRange(0, 256))
	assert.Equal(t, 128, hi)
}
