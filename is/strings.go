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
