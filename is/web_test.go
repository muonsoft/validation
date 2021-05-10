package is_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/muonsoft/validation/is"
	"github.com/stretchr/testify/assert"
)

func TestEmail_WhenValidEmail_ExpectTrue(t *testing.T) {
	emails := []string{
		"user@example.com",
		"example@example.co.uk",
		"user_name@example.fr",
		"example@example.co..uk",
		"{}~!@!@£$%%^&*().!@£$%^&*()",
		"example@example.co..uk",
		"example@-example.com",
		fmt.Sprintf("example@%s.com", strings.Repeat("x", 64)),
	}
	for _, email := range emails {
		t.Run(email, func(t *testing.T) {
			isValid := is.Email(email)

			assert.True(t, isValid)
		})
	}
}

func TestEmail_WhenInvalidEmail_ExpectFalse(t *testing.T) {
	emails := []string{
		"example",
		"example@",
		"example@localhost",
		"foo@example.com bar",
	}
	for _, email := range emails {
		t.Run(email, func(t *testing.T) {
			isValid := is.Email(email)

			assert.False(t, isValid)
		})
	}
}

func TestHTML5Email_WhenValidEmail_ExpectTrue(t *testing.T) {
	emails := []string{
		"user@example.com",
		"example@example.co.uk",
		"user_name@example.fr",
		"{}~!@example.com",
	}
	for _, email := range emails {
		t.Run(email, func(t *testing.T) {
			isValid := is.HTML5Email(email)

			assert.True(t, isValid)
		})
	}
}

func TestHTML5Email_WhenInvalidEmail_ExpectFalse(t *testing.T) {
	emails := []string{
		"example",
		"example@",
		"example@localhost",
		"example@example.co..uk",
		"foo@example.com bar",
		"example@example.",
		"example@.fr",
		"@example.com",
		"example@example.com;example@example.com",
		"example@.",
		" example@example.com",
		"example@ ",
		" example@example.com ",
		" example @example .com ",
		"example@-example.com",
		fmt.Sprintf("example@%s.com", strings.Repeat("x", 64)),
	}
	for _, email := range emails {
		t.Run(email, func(t *testing.T) {
			isValid := is.HTML5Email(email)

			assert.False(t, isValid)
		})
	}
}

func TestHostname_WhenValidHostname_ExpectTrue(t *testing.T) {
	hostnames := append(append(validMultilevelDomains(), reservedDomains()...), topLevelDomains()...)
	for _, hostname := range hostnames {
		t.Run(hostname, func(t *testing.T) {
			isValid := is.Hostname(hostname)

			assert.True(t, isValid)
		})
	}
}

func TestHostname_WhenInvalidHostname_ExpectFalse(t *testing.T) {
	for _, hostname := range invalidDomains() {
		t.Run(hostname, func(t *testing.T) {
			isValid := is.Hostname(hostname)

			assert.False(t, isValid)
		})
	}
}

func TestStrictHostname_WhenValidHostname_ExpectTrue(t *testing.T) {
	for _, hostname := range validMultilevelDomains() {
		t.Run(hostname, func(t *testing.T) {
			isValid := is.StrictHostname(hostname)

			assert.True(t, isValid)
		})
	}
}

func TestStrictHostname_WhenInvalidHostname_ExpectFalse(t *testing.T) {
	hostnames := append(append(invalidDomains(), reservedDomains()...), topLevelDomains()...)
	for _, hostname := range hostnames {
		t.Run(hostname, func(t *testing.T) {
			isValid := is.StrictHostname(hostname)

			assert.False(t, isValid)
		})
	}
}

func validMultilevelDomains() []string {
	return []string{
		"example.com",
		"example.co.uk",
		"example.fr",
		"xn--diseolatinoamericano-66b.com",
		"xn--ggle-0nda.com",
		"www.xn--simulateur-prt-2kb.fr",
		fmt.Sprintf("%s.com", strings.Repeat("a", 20)),
	}
}

func invalidDomains() []string {
	return []string{
		"acme..com",
		"qq--.com",
		"-example.com",
		"example-.com",
		strings.Repeat("a", 64) + ".com",
		strings.Repeat(strings.Repeat("a", 63)+".", 4) + "host",
	}
}

func reservedDomains() []string {
	return []string{
		"example",
		"foo.example",
		"invalid",
		"bar.invalid",
		"localhost",
		"lol.localhost",
		"test",
		"abc.test",
	}
}

func topLevelDomains() []string {
	return []string{
		"com",
		"net",
		"org",
		"etc",
	}
}
