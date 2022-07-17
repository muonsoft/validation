package validate

import (
	"errors"
	"net"
	"net/url"
	"regexp"
	"strings"
)

var (
	ErrRestrictedSchema = errors.New("restricted schema")
	ErrRestrictedHost   = errors.New("restricted host")
)

// URL is used to validate that value is a valid URL string. You can use a list of restrictions
// to additionally check for a restricted set of URLs. By default, if no restrictions are passed,
// the function checks for the http:// and https:// schemas.
//
// Use the RestrictURLSchemas option to configure the list of expected schemas. If an empty string is passed
// as a schema, then URL value may be treated as relative (without schema, e.g. "//example.com").
//
// Use the RestrictURLHosts or RestrictURLHostByPattern option to configure the list of allowed hosts.
//
// If value is not a valid URL the function will return one of the errors:
//	• parsing error from url.Parse method if value cannot be parsed as a URL;
//	• ErrRestrictedSchema if schema is not matching one of the listed schemas;
//	• ErrRestrictedHost if host is not matching one of the listed hosts;
//	• ErrInvalid if value is not matching the regular expression.
func URL(value string, restrictions ...URLRestriction) error {
	if len(restrictions) == 0 {
		restrictions = append(restrictions, RestrictURLSchemas("http", "https"))
	}
	u, err := url.Parse(value)
	if err != nil {
		return err
	}

	for _, check := range restrictions {
		if err = check(u); err != nil {
			return err
		}
	}

	if !urlRegex.MatchString(value) {
		return ErrInvalid
	}

	return nil
}

// URLRestriction can be used to limit valid URL values.
type URLRestriction func(u *url.URL) error

// RestrictURLSchemas make URL validation accepts only the listed schemas.
func RestrictURLSchemas(schemas ...string) URLRestriction {
	return func(u *url.URL) error {
		for _, schema := range schemas {
			if schema == u.Scheme {
				return nil
			}
		}

		return ErrRestrictedSchema
	}
}

// RestrictURLHosts make URL validation accepts only the listed hosts.
func RestrictURLHosts(hosts ...string) URLRestriction {
	return func(u *url.URL) error {
		for _, host := range hosts {
			if host == u.Host {
				return nil
			}
		}

		return ErrRestrictedHost
	}
}

// RestrictURLHostByPattern make URL validation accepts only a value with a host matching by pattern.
func RestrictURLHostByPattern(pattern *regexp.Regexp) URLRestriction {
	return func(u *url.URL) error {
		if pattern.MatchString(u.Host) {
			return nil
		}

		return ErrRestrictedHost
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

// IPRestriction can be used to limit valid IP address values.
type IPRestriction func(ip net.IP) bool

// DenyPrivateIP denies using of private IPs according to RFC 1918 (IPv4 addresses)
// and RFC 4193 (IPv6 addresses).
func DenyPrivateIP() IPRestriction {
	return func(ip net.IP) bool {
		return ip.IsPrivate()
	}
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
