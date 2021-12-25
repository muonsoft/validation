package is

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
