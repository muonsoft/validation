package validate

import (
	"errors"
	"net/url"
)

var (
	ErrUnexpectedProtocol = errors.New("unexpected protocol")
)

// URL is used to validate that value is a valid absolute URL string, which means that a protocol (or scheme)
// is required. By default (if no protocols are passed), the function checks only for
// the http:// and https:// protocols. Use the protocols argument to configure the list of expected protocols.
//
// If value is not a valid URL the function will return one of the errors:
//	• parsing error from url.Parse method if value cannot be parsed as an URL;
//	• ErrUnexpectedProtocol if protocol is not matching one of the listed protocols;
//	• ErrInvalid if value is not matching the regular expression.
func URL(value string, protocols ...string) error {
	return validateURL(value, false, protocols...)
}

// RelativeURL is used to validate that value is a valid absolute or relative URL string. The protocol is considered
// optional when validating the syntax of the given URL. This means that both http:// and https:// are valid
// but also relative URLs that contain no protocol (e.g. //example.com). By default, the function checks
// only for the http:// and https:// protocols. Use the protocols argument to configure
// the list of expected protocols.
//
// If value is not a valid URL the function will return one of the errors:
//	• parsing error from url.Parse method if value cannot be parsed as an URL;
//	• ErrUnexpectedProtocol if protocol is not matching one of the listed protocols or it is not a relative URL;
//	• ErrInvalid if value is not matching the regular expression.
func RelativeURL(value string, protocols ...string) error {
	return validateURL(value, true, protocols...)
}

func validateURL(value string, isRelative bool, protocols ...string) error {
	if len(protocols) == 0 {
		protocols = []string{"http", "https"}
	}
	u, err := url.Parse(value)
	if err != nil {
		return err
	}

	err = validateProtocol(u, isRelative, protocols)
	if err != nil {
		return err
	}

	if !urlRegex.MatchString(value) {
		return ErrInvalid
	}

	return nil
}

func validateProtocol(u *url.URL, isRelative bool, protocols []string) error {
	if isRelative && u.Scheme == "" {
		return nil
	}

	for _, protocol := range protocols {
		if protocol == u.Scheme {
			return nil
		}
	}

	return ErrUnexpectedProtocol
}
