package scenario

import (
	"strings"
)

func DisplayName(first, last string) string {
	parts := []string{}
	for _, part := range []string{first, last} {
		if part = strings.TrimSpace(part); part != "" {
			parts = append(parts, part)
		}
	}
	if len(parts) == 0 {
		return "Anonymous"
	}
	return strings.Join(parts, " ")
}
