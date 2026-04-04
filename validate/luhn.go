package validate

// LUHN validates whether the value passes the Luhn (mod 10) checksum.
// The string must contain only ASCII digits (0–9); spaces and other characters are not stripped,
// matching Symfony\Component\Validator\Constraints\Luhn.
//
// Empty string is considered valid (use [NotBlank] or similar to reject empty values).
//
// Possible errors:
//   - [ErrContainsNonDigit] when the value contains a non-digit character;
//   - [ErrInvalidChecksum] when the checksum is invalid or the value is all zeros (checksum 0).
//
// See https://en.wikipedia.org/wiki/Luhn_algorithm.
func LUHN(value string) error {
	if value == "" {
		return nil
	}
	for i := 0; i < len(value); i++ {
		if value[i] < '0' || value[i] > '9' {
			return ErrContainsNonDigit
		}
	}
	if !luhnValidDigits(value) {
		return ErrInvalidChecksum
	}

	return nil
}

// luhnValidDigits implements the same checksum as Symfony's LuhnValidator for a decimal string.
func luhnValidDigits(digits string) bool {
	if len(digits) == 0 {
		return false
	}
	for i := 0; i < len(digits); i++ {
		if digits[i] < '0' || digits[i] > '9' {
			return false
		}
	}

	checkSum := 0
	length := len(digits)
	for i := length - 1; i >= 0; i-- {
		d := int(digits[i] - '0')
		if (i%2)^(length%2) != 0 {
			checkSum += d
		} else {
			doubled := d * 2
			checkSum += doubled/10 + doubled%10
		}
	}

	return checkSum != 0 && checkSum%10 == 0
}
