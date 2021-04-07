package it

import "errors"

var (
	errEmptyChoices = errors.New("empty list of choices")
	errEmptyRegex   = errors.New("nil regex")
	errInvalidRange = errors.New("invalid range")
)
