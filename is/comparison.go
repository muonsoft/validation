package is

import "math"

// InList returns true if one of the elements of the list is equal to the value.
func InList[T comparable](value T, list []T) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}

	return false
}

// Unique checks that the slice of comparable values are unique.
func Unique[T comparable](values []T) bool {
	if len(values) == 0 {
		return true
	}

	uniques := make(map[T]struct{}, len(values))

	for _, value := range values {
		if _, exists := uniques[value]; exists {
			return false
		}
		uniques[value] = struct{}{}
	}

	return true
}

// DivisibleBy checks that a value is divisible by another value (divisor).
// It checks that the remainder is zero or almost zero with an error of 1e-12.
func DivisibleBy(divisible, divisor float64) bool {
	const epsilon = 1e-12

	remainder := math.Mod(divisible, divisor)
	if remainder < epsilon {
		return true
	}

	return math.Abs(remainder-divisor) < epsilon
}
