package is

import "github.com/muonsoft/validation/validate"

// URL is used to check that value is a valid absolute URL string, which means that a protocol (or scheme) is required.
// By default (if no protocols are passed), the function checks only for the http:// and https:// protocols.
// Use the protocols argument to configure the list of expected protocols.
func URL(value string, protocols ...string) bool {
	return validate.URL(value, protocols...) == nil
}

// RelativeURL is used to check that value is a valid absolute or relative URL string. The protocol is considered
// optional when validating the syntax of the given URL. This means that both http:// and https:// are valid
// but also relative URLs that contain no protocol (e.g. //example.com). By default, the function checks
// only for the http:// and https:// protocols. Use the protocols argument to configure
// the list of expected protocols.
func RelativeURL(value string, protocols ...string) bool {
	return validate.RelativeURL(value, protocols...) == nil
}

func Email(value string) bool {
	return looseEmailRegex.MatchString(value)
}

func HTML5Email(value string) bool {
	return html5EmailRegex.MatchString(value)
}
