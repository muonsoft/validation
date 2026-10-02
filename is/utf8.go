package is

import "github.com/muonsoft/validation/validate"

// UTF8 reports whether value is valid UTF-8, including the empty string.
// Use it to reject malformed or truncated bytes from imports or integrations.
// It checks encoding only, not text meaning or safety, and does not detect,
// convert, or normalize encodings. See [validate.UTF8] for limitations.
func UTF8(value string) bool {
	return validate.UTF8(value) == nil
}
