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
