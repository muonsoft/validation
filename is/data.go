package is

import "encoding/json"

// JSON checks that value is a valid JSON string.
func JSON(value string) bool {
	return json.Valid([]byte(value))
}
