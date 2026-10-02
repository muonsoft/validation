package is_test

import (
	"strings"
	"testing"

	"github.com/muonsoft/validation/is"
	"github.com/stretchr/testify/require"
)

func TestStrictHostname_ReservedDomainsAreCaseInsensitive(t *testing.T) {
	for _, tld := range []string{"test", "example", "invalid", "localhost"} {
		require.False(t, is.StrictHostname("host."+tld))
		require.False(t, is.StrictHostname("host."+strings.ToUpper(tld)))
	}
	require.True(t, is.StrictHostname("host.COM"))
}

func TestHostname_TotalLengthIncludesSeparators(t *testing.T) {
	require.True(t, is.Hostname(strings.Repeat("a.", 127)+"a"))
	require.False(t, is.Hostname(strings.Repeat("a.", 127)+"aa"))
	require.False(t, is.Hostname(strings.Repeat("a.", 254)+"a"))
}
