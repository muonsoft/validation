package it

import "errors"

var errEmptyChoices = errors.New("empty list of choices")
var errEmptyPattern = errors.New("empty pattern")
var errInvalidPattern = errors.New("invalid pattern")
