package is

import "regexp"

const (
	looseEmailPattern = `^.+\@\S+\.\S+$`
	html5EmailPattern = "^[a-zA-Z0-9.!#$%&\\'*+\\\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$"

	// source https://stackoverflow.com/questions/106179/regular-expression-to-match-dns-hostname-or-ip-address
	hostnamePattern = `^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))*$`
)

var (
	looseEmailRegex = regexp.MustCompile(looseEmailPattern)
	html5EmailRegex = regexp.MustCompile(html5EmailPattern)
	hostnameRegex   = regexp.MustCompile(hostnamePattern)
)
