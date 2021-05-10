package is

import "encoding/json"

// JSON checks that value is a valid JSON string.
func JSON(value string) bool {
	var data interface{}

	return json.Unmarshal([]byte(value), &data) == nil
}
