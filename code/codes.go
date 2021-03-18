// Copyright 2021 Igor Lazarev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package code contains a list of unique, short, and semantic violation codes.
// They can be used to programmatically test for specific violation.
// All code values are protected by backward compatibility rules.
package code

const (
	Blank          = "blank"
	CountExact     = "countExact"
	CountTooFew    = "countTooFew"
	CountTooMany   = "countTooMany"
	False          = "false"
	LengthExact    = "lengthExact"
	LengthTooFew   = "lengthTooFew"
	LengthTooMany  = "lengthTooMany"
	MatchingFailed = "matchingFailed"
	Nil            = "nil"
	NoSuchChoice   = "noSuchChoice"
	NotBlank       = "notBlank"
	NotNil         = "notNil"
	True           = "true"
)
