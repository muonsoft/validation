package is

import (
	"strings"

	"github.com/muonsoft/validation/validate"
)

var reservedTopLevelDomains = []string{
	"example",
	"invalid",
	"localhost",
	"test",
}

// Email is used for simplified validation of an email address. It allows all values
// with an "@" symbol in, and a "." in the second host part of the email address.
func Email(value string) bool {
	return looseEmailRegex.MatchString(value)
}

// HTML5Email is used for validation of an email address based on pattern for HTML5
// (see https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address).
func HTML5Email(value string) bool {
	return html5EmailRegex.MatchString(value)
}

// URL is used to check that value is a valid URL string. By default (if no schemas are passed),
// the function checks only for the http:// and https:// schemas. Use the schemas argument
// to configure the list of expected schemas. If an empty string is passed as a schema, then
// URL value may be treated as relative (without schema, e.g. "//example.com").
func URL(value string, schemas ...string) bool {
	return validate.URL(value, schemas...) == nil
}

// IP checks that a value is a valid IP address (IPv4 or IPv6). You can use a list
// of restrictions to additionally check for a restricted range of IPs.
// See validate.IPRestriction for details.
func IP(value string, restrictions ...validate.IPRestriction) bool {
	return validate.IP(value, restrictions...) == nil
}

// IPv4 checks that a value is a valid IPv4 address. You can use a list
// of restrictions to additionally check for a restricted range of IPs.
// See validate.IPRestriction for details.
func IPv4(value string, restrictions ...validate.IPRestriction) bool {
	return validate.IPv4(value, restrictions...) == nil
}

// IPv6 checks that a value is a valid IPv6 address. You can use a list
// of restrictions to additionally check for a restricted range of IPs.
// See validate.IPRestriction for details.
func IPv6(value string, restrictions ...validate.IPRestriction) bool {
	return validate.IPv6(value, restrictions...) == nil
}

// Hostname checks that a value is a valid hostname. It checks that each label
// within a valid hostname may be no more than 63 octets long. Also, it checks that
// the total length of the hostname must not exceed 255 characters.
//
// See StrictHostname for additional checks.
func Hostname(value string) bool {
	return hostnameRegex.MatchString(value) && len(strings.ReplaceAll(value, ".", "")) <= 255
}

// StrictHostname checks that a value is a valid hostname. Beside checks from Hostname function
// it checks that hostname is fully qualified and include its top-level domain name (TLD).
// For instance, example.com is valid but example is not.
//
// Additionally it checks for reserved top-level domains according to RFC 2606 and
// that's why hostnames containing them are not considered valid:
// .example, .invalid, .localhost, and .test.
func StrictHostname(value string) bool {
	if !Hostname(value) {
		return false
	}

	domains := strings.Split(value, ".")
	if len(domains) < 2 {
		return false
	}

	tld := domains[len(domains)-1]
	for _, reservedTLD := range reservedTopLevelDomains {
		if tld == reservedTLD {
			return false
		}
	}

	return true
}
