package is

// StringInList returns true if one of the elements of the list is equal to the string.
func StringInList(s string, list []string) bool {
	for _, ls := range list {
		if ls == s {
			return true
		}
	}

	return false
}

// UniqueStrings checks that slice of strings has unique values.
func UniqueStrings(values []string) bool {
	if len(values) == 0 {
		return true
	}

	uniques := make(map[string]struct{}, len(values))

	for _, value := range values {
		if _, exists := uniques[value]; exists {
			return false
		}
		uniques[value] = struct{}{}
	}

	return true
}
