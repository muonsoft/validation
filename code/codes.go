// Copyright 2021 Igor Lazarev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package code contains a list of unique, short, and semantic violation codes.
// They can be used to programmatically test for specific violation.
// All code values are protected by backward compatibility rules.
package code

const (
	Blank             = "blank"
	CountExact        = "countExact"
	CountTooFew       = "countTooFew"
	CountTooMany      = "countTooMany"
	Equal             = "equal"
	False             = "false"
	InvalidEAN13      = "invalidEAN13"
	InvalidEAN8       = "invalidEAN8"
	InvalidEmail      = "invalidEmail"
	InvalidHostname   = "invalidHostname"
	InvalidIP         = "invalidIP"
	InvalidJSON       = "invalidJSON"
	InvalidUPCA       = "invalidUPCA"
	InvalidUPCE       = "invalidUPCE"
	InvalidURL        = "invalidURL"
	LengthExact       = "lengthExact"
	LengthTooFew      = "lengthTooFew"
	LengthTooMany     = "lengthTooMany"
	MatchingFailed    = "matchingFailed"
	Nil               = "nil"
	NoSuchChoice      = "noSuchChoice"
	NotBlank          = "notBlank"
	NotEqual          = "notEqual"
	NotInRange        = "notInRange"
	NotInteger        = "notInteger"
	NotNegative       = "notNegative"
	NotNegativeOrZero = "notNegativeOrZero"
	NotNil            = "notNil"
	NotNumeric        = "notNumeric"
	NotPositive       = "notPositive"
	NotPositiveOrZero = "notPositiveOrZero"
	NotUnique         = "notUnique"
	NotValid          = "notValid"
	ProhibitedIP      = "prohibitedIP"
	TooEarly          = "tooEarly"
	TooEarlyOrEqual   = "tooEarlyOrEqual"
	TooHigh           = "tooHigh"
	TooHighOrEqual    = "tooHighOrEqual"
	TooLate           = "tooLate"
	TooLateOrEqual    = "tooLateOrEqual"
	TooLow            = "tooLow"
	TooLowOrEqual     = "tooLowOrEqual"
	True              = "true"
)
