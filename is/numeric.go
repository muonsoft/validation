package is

import "strconv"

// Integer checks that string is an integer.
func Integer(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// Number checks that string is a number (a float or an integer).
func Number(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
