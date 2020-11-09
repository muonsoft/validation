package validation

import "strings"

func renderMessage(template string, parameters map[string]string) string {
	message := template

	for key, value := range parameters {
		message = strings.ReplaceAll(message, key, value)
	}

	return message
}
