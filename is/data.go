package is

import "encoding/json"

// JSON checks that value is a valid JSON string.
func JSON(value string) bool {
	return json.Valid([]byte(value))
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
