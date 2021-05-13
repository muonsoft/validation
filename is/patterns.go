package is

import "regexp"

const (
	looseEmailPattern = `^.+\@\S+\.\S+$`
	html5EmailPattern = "^[a-zA-Z0-9.!#$%&\\'*+\\\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$"
)

var (
	looseEmailRegex = regexp.MustCompile(looseEmailPattern)
	html5EmailRegex = regexp.MustCompile(html5EmailPattern)
)
