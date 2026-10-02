package validate

import (
	"errors"
	"mime"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	// ErrInvalidFileName indicates an invalid file name.
	ErrInvalidFileName = errors.New("invalid file name")
	// ErrInvalidFileExtension indicates a disallowed extension or invalid name.
	ErrInvalidFileExtension = errors.New("invalid file extension")
	// ErrInvalidMIMEType indicates malformed or disallowed MIME metadata.
	ErrInvalidMIMEType = errors.New("invalid MIME type")
	// ErrInvalidContentType indicates disallowed detected content.
	ErrInvalidContentType = errors.New("invalid content type")
	// ErrInvalidFileExtensions indicates an invalid extension allowlist.
	ErrInvalidFileExtensions = errors.New("invalid file extensions")
	// ErrInvalidMIMETypes indicates an invalid MIME allowlist.
	ErrInvalidMIMETypes = errors.New("invalid MIME types")
	// ErrInvalidDetectedContentType indicates a detector contract violation.
	ErrInvalidDetectedContentType = errors.New("invalid detected content type")
)

// FileNameOptions configures file name validation.
type FileNameOptions struct{ windows bool }

// WithWindowsFileNameRestrictions enables conservative Windows name restrictions
// on every operating system. It does not impose filesystem length limits.
func WithWindowsFileNameRestrictions() func(*FileNameOptions) {
	return func(o *FileNameOptions) { o.windows = true }
}

// FileName validates a name, not a path. Empty names are skipped. Invalid UTF-8,
// control characters, separators, dot entries and drive prefixes are rejected.
// Input is never decoded, normalized or trimmed. Length and storage safety are
// separate concerns. WithWindowsFileNameRestrictions adds portable restrictions.
func FileName(name string, options ...func(*FileNameOptions)) error {
	o := FileNameOptions{}
	for _, option := range options {
		option(&o)
	}
	if name == "" {
		return nil
	}
	if invalidBasicFileName(name) {
		return ErrInvalidFileName
	}
	if o.windows && invalidWindowsFileName(name) {
		return ErrInvalidFileName
	}
	return nil
}

func invalidBasicFileName(name string) bool {
	return !utf8.ValidString(name) || name == "." || name == ".." ||
		strings.ContainsAny(name, `/\`) || hasDrivePrefix(name) || strings.IndexFunc(name, unicode.IsControl) >= 0
}

func invalidWindowsFileName(name string) bool {
	return strings.ContainsAny(name, `<>:"|?*`) || strings.HasSuffix(name, " ") ||
		strings.HasSuffix(name, ".") || reservedWindowsName(name)
}

func hasDrivePrefix(name string) bool {
	return len(name) >= 2 && name[1] == ':' && ((name[0] >= 'A' && name[0] <= 'Z') || (name[0] >= 'a' && name[0] <= 'z'))
}

func reservedWindowsName(name string) bool {
	base, _, _ := strings.Cut(name, ".")
	base = strings.ToUpper(strings.TrimRight(base, " "))
	switch base {
	case "CON", "PRN", "AUX", "NUL", "CONIN$", "CONOUT$":
		return true
	}
	if strings.HasPrefix(base, "COM") || strings.HasPrefix(base, "LPT") {
		tail := base[3:]
		return len(tail) == 1 && tail[0] >= '1' && tail[0] <= '9' || tail == "¹" || tail == "²" || tail == "³"
	}
	return false
}

// FileExtension checks an ASCII case-insensitive extension allowlist, accepting
// entries with or without a leading dot and compound extensions such as tar.gz.
// Entries contain nonempty dot-separated ASCII letters, digits, _, + or -.
// A matching suffix needs a nonempty stem. Paths are rejected; no basename is
// extracted. Empty input is skipped, but invalid configuration is still reported.
func FileExtension(name string, extensions ...string) error {
	if len(extensions) == 0 {
		return ErrInvalidFileExtensions
	}
	for _, extension := range extensions {
		if !validExtension(extension) {
			return ErrInvalidFileExtensions
		}
	}
	if name == "" {
		return nil
	}
	if FileName(name) != nil {
		return ErrInvalidFileExtension
	}
	name = asciiLower(name)
	for _, extension := range extensions {
		suffix := "." + asciiLower(strings.TrimPrefix(extension, "."))
		if len(name) > len(suffix) && strings.HasSuffix(name, suffix) {
			return nil
		}
	}
	return ErrInvalidFileExtension
}

func validExtension(extension string) bool {
	extension = strings.TrimPrefix(extension, ".")
	for _, part := range strings.Split(extension, ".") {
		if part == "" {
			return false
		}
		for _, c := range part {
			if !strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+-", c) {
				return false
			}
		}
	}
	return true
}

func asciiLower(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 'A' && r <= 'Z' {
			return r + ('a' - 'A')
		}
		return r
	}, s)
}

// MIMEType validates concrete type/subtype metadata, optionally against allowed
// types. Valid parameters are ignored in comparisons. Allowlist entries must
// have no parameters or wildcards. Empty input is skipped, even with an allowlist;
// invalid configuration is always reported. No aliases are resolved.
func MIMEType(value string, types ...string) error {
	if !validMIMETypes(types) {
		return ErrInvalidMIMETypes
	}
	if value == "" {
		return nil
	}
	actual, _, err := parseMIMEType(value)
	if err != nil {
		return ErrInvalidMIMEType
	}
	if len(types) == 0 || matchesMIMEType(actual, types) {
		return nil
	}
	return ErrInvalidMIMEType
}

func parseMIMEType(value string) (string, map[string]string, error) {
	actual, params, err := mime.ParseMediaType(value)
	if err != nil || !strings.Contains(actual, "/") || strings.Contains(actual, "*") {
		return "", nil, ErrInvalidMIMEType
	}
	return actual, params, nil
}

func validMIMETypes(types []string) bool {
	for _, value := range types {
		_, params, err := parseMIMEType(value)
		if err != nil || len(params) != 0 || strings.Contains(value, ";") {
			return false
		}
	}
	return true
}

func matchesMIMEType(actual string, types []string) bool {
	for _, value := range types {
		expected, _, _ := parseMIMEType(value)
		if actual == expected {
			return true
		}
	}
	return false
}

// ContentTypeDetector identifies content without modifying data. Shared detectors
// must support concurrent calls. Results must be concrete MIME types; parameters
// are allowed. Detection does not establish file integrity or safety.
type ContentTypeDetector func([]byte) string

// ContentTypeOptions configures content detection.
type ContentTypeOptions struct{ detector ContentTypeDetector }

// WithContentTypeDetector replaces detection. Nil restores http.DetectContentType.
// Custom detectors receive the entire supplied slice without truncation.
func WithContentTypeDetector(detector ContentTypeDetector) func(*ContentTypeOptions) {
	return func(o *ContentTypeOptions) { o.detector = detector }
}

// ContentType checks detected content against a nonempty allowlist of concrete
// MIME types without parameters. Nil and empty data are skipped. The default
// detector inspects at most 512 bytes; unknown content is application/octet-stream
// and only passes when explicitly allowed. No I/O or complete-file parsing occurs.
// Invalid lists return ErrInvalidMIMETypes; invalid detector results return
// ErrInvalidDetectedContentType; mismatches return ErrInvalidContentType.
func ContentType(data []byte, types []string, options ...func(*ContentTypeOptions)) error {
	o := ContentTypeOptions{}
	for _, option := range options {
		option(&o)
	}
	if len(types) == 0 || !validMIMETypes(types) {
		return ErrInvalidMIMETypes
	}
	if len(data) == 0 {
		return nil
	}
	if o.detector == nil {
		o.detector = http.DetectContentType
	}
	actual, _, err := parseMIMEType(o.detector(data))
	if err != nil {
		return ErrInvalidDetectedContentType
	}
	if !matchesMIMEType(actual, types) {
		return ErrInvalidContentType
	}
	return nil
}
