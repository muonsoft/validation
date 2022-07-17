package validate

import (
	"errors"
	"net"
	"net/url"
	"regexp"
	"strings"
)

var ErrUnexpectedSchema = errors.New("unexpected schema")

// URL is used to validate that value is a valid URL string. By default, (if no schemas are passed),
// the function checks only for the http:// and https:// schemas. Use the schemas' argument
// to configure the list of expected schemas. If an empty string is passed as a schema, then
// URL value may be treated as relative (without schema, e.g. "//example.com").
//
// If value is not a valid URL the function will return one of the errors:
//	• parsing error from url.Parse method if value cannot be parsed as a URL;
//	• ErrUnexpectedSchema if schema is not matching one of the listed schemas;
//	• ErrInvalid if value is not matching the regular expression.
func URL(value string, schemas ...string) error {
	if len(schemas) == 0 {
		schemas = []string{"http", "https"}
	}
	u, err := url.Parse(value)
	if err != nil {
		return err
	}

	err = validateSchema(u, schemas)
	if err != nil {
		return err
	}

	if !urlRegex.MatchString(value) {
		return ErrInvalid
	}

	return nil
}

func validateSchema(u *url.URL, schemas []string) error {
	for _, schema := range schemas {
		if schema == u.Scheme {
			return nil
		}
	}

	return ErrUnexpectedSchema
}

// IPRestriction can be used to limit valid IP address values.
type IPRestriction func(ip net.IP) bool

// DenyPrivateIP denies using of private IPs according to RFC 1918 (IPv4 addresses)
// and RFC 4193 (IPv6 addresses).
func DenyPrivateIP() IPRestriction {
	return func(ip net.IP) bool {
		return ip.IsPrivate()
	}
}

// IP validates that a value is a valid IP address (IPv4 or IPv6). You can use a list
// of restrictions to additionally check for a restricted range of IPs. For example,
// you can deny using private IP addresses using DenyPrivateIP function.
//
// If value is not valid the function will return one of the errors:
//	• ErrInvalid on invalid IP address;
//	• ErrProhibited on restricted IP address.
func IP(value string, restrictions ...IPRestriction) error {
	return validateIP(value, restrictions...)
}

// IPv4 validates that a value is a valid IPv4 address. You can use a list
// of restrictions to additionally check for a restricted range of IPs. For example,
// you can deny using private IP addresses using DenyPrivateIP function.
//
// If value is not valid the function will return one of the errors:
//	• ErrInvalid on invalid IP address or when using IPv6;
//	• ErrProhibited on restricted IP address.
func IPv4(value string, restrictions ...IPRestriction) error {
	err := validateIP(value, restrictions...)
	if err != nil {
		return err
	}
	if !strings.Contains(value, ".") || strings.Contains(value, ":") {
		return ErrInvalid
	}

	return nil
}

// IPv6 validates that a value is a valid IPv6 address. You can use a list
// of restrictions to additionally check for a restricted range of IPs. For example,
// you can deny using private IP addresses using DenyPrivateIP function.
//
// If value is not valid the function will return one of the errors:
//	• ErrInvalid on invalid IP address or when using IPv4;
//	• ErrProhibited on restricted IP address.
func IPv6(value string, restrictions ...IPRestriction) error {
	err := validateIP(value, restrictions...)
	if err != nil {
		return err
	}
	if !strings.Contains(value, ":") {
		return ErrInvalid
	}

	return nil
}

func validateIP(value string, restrictions ...IPRestriction) error {
	ip := net.ParseIP(value)
	if ip == nil {
		return ErrInvalid
	}
	for _, isProhibited := range restrictions {
		if isProhibited(ip) {
			return ErrProhibited
		}
	}

	return nil
}

const (
	urlSchema     = `([a-z]*:)?//`
	urlBasicAuth  = `(((?:[\_\.\pL\pN-]|%[0-9A-Fa-f]{2})+:)?((?:[\_\.\pL\pN-]|%[0-9A-Fa-f]{2})+)@)?`
	urlDomainName = `([\pL\pN\pS\-\_\.])+(\.?([\pL\pN]|xn\-\-[\pL\pN-]+)+\.?)`
	ipv4          = `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`
	ipv6          = `\[(?:(?:(?:(?:(?:(?:(?:[0-9a-f]{1,4})):){6})(?:(?:(?:(?:(?:[0-9a-f]{1,4})):(?:(?:[0-9a-f]{1,4})))|(?:(?:(?:(?:(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9]))\.){3}(?:(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9])))))))|(?:(?:::(?:(?:(?:[0-9a-f]{1,4})):){5})(?:(?:(?:(?:(?:[0-9a-f]{1,4})):(?:(?:[0-9a-f]{1,4})))|(?:(?:(?:(?:(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9]))\.){3}(?:(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9])))))))|(?:(?:(?:(?:(?:[0-9a-f]{1,4})))?::(?:(?:(?:[0-9a-f]{1,4})):){4})(?:(?:(?:(?:(?:[0-9a-f]{1,4})):(?:(?:[0-9a-f]{1,4})))|(?:(?:(?:(?:(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9]))\.){3}(?:(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9])))))))|(?:(?:(?:(?:(?:(?:[0-9a-f]{1,4})):){0,1}(?:(?:[0-9a-f]{1,4})))?::(?:(?:(?:[0-9a-f]{1,4})):){3})(?:(?:(?:(?:(?:[0-9a-f]{1,4})):(?:(?:[0-9a-f]{1,4})))|(?:(?:(?:(?:(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9]))\.){3}(?:(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9])))))))|(?:(?:(?:(?:(?:(?:[0-9a-f]{1,4})):){0,2}(?:(?:[0-9a-f]{1,4})))?::(?:(?:(?:[0-9a-f]{1,4})):){2})(?:(?:(?:(?:(?:[0-9a-f]{1,4})):(?:(?:[0-9a-f]{1,4})))|(?:(?:(?:(?:(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9]))\.){3}(?:(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9])))))))|(?:(?:(?:(?:(?:(?:[0-9a-f]{1,4})):){0,3}(?:(?:[0-9a-f]{1,4})))?::(?:(?:[0-9a-f]{1,4})):)(?:(?:(?:(?:(?:[0-9a-f]{1,4})):(?:(?:[0-9a-f]{1,4})))|(?:(?:(?:(?:(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9]))\.){3}(?:(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9])))))))|(?:(?:(?:(?:(?:(?:[0-9a-f]{1,4})):){0,4}(?:(?:[0-9a-f]{1,4})))?::)(?:(?:(?:(?:(?:[0-9a-f]{1,4})):(?:(?:[0-9a-f]{1,4})))|(?:(?:(?:(?:(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9]))\.){3}(?:(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9])))))))|(?:(?:(?:(?:(?:(?:[0-9a-f]{1,4})):){0,5}(?:(?:[0-9a-f]{1,4})))?::)(?:(?:[0-9a-f]{1,4})))|(?:(?:(?:(?:(?:(?:[0-9a-f]{1,4})):){0,6}(?:(?:[0-9a-f]{1,4})))?::))))\]`
	urlHost       = `(` + urlDomainName + `|` + ipv4 + `|` + ipv6 + `)`
	urlPort       = `(:[0-9]+)?`
	urlPath       = `(?:/(?:[\pL\pN\-._\~!$&\'()*+,;=:@]|%[0-9A-Fa-f]{2})*)*`
	urlQuery      = `(?:\?(?:[\pL\pN\-._\~!$&\'\[\]()*+,;=:@/?]|%[0-9A-Fa-f]{2})*)?`
	urlFragment   = `(?:\#(?:[\pL\pN\-._\~!$&\'()*+,;=:@/?]|%[0-9A-Fa-f]{2})*)?`
	urlPattern    = `(?i)^` + urlSchema + urlBasicAuth + urlHost + urlPort + urlPath + urlQuery + urlFragment + `$`
)

var urlRegex = regexp.MustCompile(urlPattern)
