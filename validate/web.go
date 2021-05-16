package validate

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

var (
	ErrUnexpectedSchema = errors.New("unexpected schema")
)

// URL is used to validate that value is a valid URL string. By default (if no schemas are passed),
// the function checks only for the http:// and https:// schemas. Use the schemas argument
// to configure the list of expected schemas. If an empty string is passed as a schema, then
// URL value may be treated as relative (without schema, e.g. "//example.com").
//
// If value is not a valid URL the function will return one of the errors:
//	• parsing error from url.Parse method if value cannot be parsed as an URL;
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
	return isPrivateIP
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

// isPrivateIP reports whether ip is a private address, according to
// RFC 1918 (IPv4 addresses) and RFC 4193 (IPv6 addresses).
//
// This function is ported from golang v1.17.
// See https://github.com/golang/go/blob/4c8f48ed4f3db0e3ba376e6b7a261d26b41d8dd0/src/net/ip.go#L133.
func isPrivateIP(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		// Following RFC 1918, Section 3. Private Address Space which says:
		//   The Internet Assigned Numbers Authority (IANA) has reserved the
		//   following three blocks of the IP address space for private internets:
		//     10.0.0.0        -   10.255.255.255  (10/8 prefix)
		//     172.16.0.0      -   172.31.255.255  (172.16/12 prefix)
		//     192.168.0.0     -   192.168.255.255 (192.168/16 prefix)
		return ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1]&0xf0 == 16) ||
			(ip4[0] == 192 && ip4[1] == 168)
	}
	// Following RFC 4193, Section 8. IANA Considerations which says:
	//   The IANA has assigned the FC00::/7 prefix to "Unique Local Unicast".
	return len(ip) == net.IPv6len && ip[0]&0xfe == 0xfc
}
