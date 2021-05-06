package it

import "errors"

var (
	errEmptyChoices   = errors.New("empty list of choices")
	errEmptyProtocols = errors.New("empty list of protocols")
	errEmptyRegex     = errors.New("nil regex")
	errInvalidRange   = errors.New("invalid range")
)
