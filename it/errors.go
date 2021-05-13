package it

import "errors"

var (
	errEmptyChoices = errors.New("empty list of choices")
	errEmptySchemas = errors.New("empty list of schemas")
	errEmptyRegex   = errors.New("nil regex")
	errInvalidRange = errors.New("invalid range")
)
