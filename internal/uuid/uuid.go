// Copyright (C) 2013-2018 by Maxim Bublis <b@codemonkey.ru>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package uuid is a port of https://github.com/gofrs/uuid with functions
// only for parsing UUID from string.
package uuid

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	ErrTooShort               = errors.New("too short")
	ErrTooLong                = errors.New("too long")
	ErrInvalidHyphenPlacement = errors.New("invalid hyphen placement")
	ErrInvalid                = errors.New("invalid")
)

// Size of a UUID in bytes.
const Size = 16

// UUID is an array type to represent the value of a UUID, as defined in RFC-4122.
type UUID [Size]byte

// UUID versions.
const (
	_  byte = iota
	V1      // Version 1 (date-time and MAC address)
	_       // Version 2 (date-time and MAC address, DCE security version) [removed]
	V3      // Version 3 (namespace name-based)
	V4      // Version 4 (random)
	V5      // Version 5 (namespace name-based)
	V6      // Version 6 (k-sortable timestamp and random data, field-compatible with v1) [peabody draft]
	V7      // Version 7 (k-sortable timestamp and random data) [peabody draft]
	_       // Version 8 (k-sortable timestamp, meant for custom implementations) [peabody draft] [not implemented]
)

// UUID layout variants.
const (
	VariantNCS byte = iota
	VariantRFC4122
	VariantMicrosoft
	VariantFuture
)

// FromString returns a UUID parsed from the input string.
// Input is expected in a form accepted by UnmarshalText.
func FromString(input string) (UUID, error) {
	u := UUID{}
	err := u.UnmarshalText([]byte(input))
	return u, err
}

func CanonicalFromString(input string) (UUID, error) {
	u := UUID{}
	if len(input) < 36 {
		return u, ErrTooShort
	}
	if len(input) > 36 {
		return u, ErrTooLong
	}

	err := u.decodeCanonical([]byte(input))
	return u, err
}

// String parse helpers.
var (
	urnPrefix  = []byte("urn:uuid:")
	byteGroups = []int{8, 4, 4, 4, 12}
)

// Nil is the nil UUID, as specified in RFC-4122, that has all 128 bits set to
// zero.
var Nil = UUID{}

// IsNil returns if the UUID is equal to the nil UUID.
func (u UUID) IsNil() bool {
	return u == Nil
}

// Version returns the algorithm version used to generate the UUID.
func (u UUID) Version() byte {
	return u[6] >> 4
}

func (u UUID) ValidVersion(versions ...byte) bool {
	if len(versions) == 0 {
		versions = []byte{1, 2, 3, 4, 5, 6, 7}
	}

	for _, version := range versions {
		if u.Version() == version {
			return true
		}
	}

	return false
}

// Variant returns the UUID layout variant.
func (u UUID) Variant() byte {
	switch {
	case (u[8] >> 7) == 0x00:
		return VariantNCS
	case (u[8] >> 6) == 0x02:
		return VariantRFC4122
	case (u[8] >> 5) == 0x06:
		return VariantMicrosoft
	case (u[8] >> 5) == 0x07:
		fallthrough
	default:
		return VariantFuture
	}
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Following formats are supported:
//
//	"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
//	"{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
//	"urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
//	"6ba7b8109dad11d180b400c04fd430c8"
//	"{6ba7b8109dad11d180b400c04fd430c8}",
//	"urn:uuid:6ba7b8109dad11d180b400c04fd430c8"
//
// ABNF for supported UUID text representation follows:
//
//	URN := 'urn'
//	UUID-NID := 'uuid'
//
//	hexdig := '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' |
//	          'a' | 'b' | 'c' | 'd' | 'e' | 'f' |
//	          'A' | 'B' | 'C' | 'D' | 'E' | 'F'
//
//	hexoct := hexdig hexdig
//	2hexoct := hexoct hexoct
//	4hexoct := 2hexoct 2hexoct
//	6hexoct := 4hexoct 2hexoct
//	12hexoct := 6hexoct 6hexoct
//
//	hashlike := 12hexoct
//	canonical := 4hexoct '-' 2hexoct '-' 2hexoct '-' 6hexoct
//
//	plain := canonical | hashlike
//	uuid := canonical | hashlike | braced | urn
//
//	braced := '{' plain '}' | '{' hashlike  '}'
//	urn := URN ':' UUID-NID ':' plain
func (u *UUID) UnmarshalText(text []byte) error {
	if len(text) < 32 {
		return ErrTooShort
	}
	if len(text) > 45 {
		return ErrTooLong
	}

	switch len(text) {
	case 32:
		return u.decodeHashLike(text)
	case 34, 38:
		return u.decodeBraced(text)
	case 36:
		return u.decodeCanonical(text)
	case 41, 45:
		return u.decodeURN(text)
	default:
		return fmt.Errorf("%w: incorrect UUID length %d in string %q", ErrInvalid, len(text), text)
	}
}

// decodeCanonical decodes UUID strings that are formatted as defined in RFC-4122 (section 3):
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8".
func (u *UUID) decodeCanonical(t []byte) error {
	if t[8] != '-' || t[13] != '-' || t[18] != '-' || t[23] != '-' {
		return ErrInvalidHyphenPlacement
	}

	src := t
	dst := u[:]

	for i, byteGroup := range byteGroups {
		if i > 0 {
			src = src[1:] // skip dash
		}
		_, err := hex.Decode(dst[:byteGroup/2], src[:byteGroup])
		if err != nil {
			return err
		}
		src = src[byteGroup:]
		dst = dst[byteGroup/2:]
	}

	return nil
}

// decodeHashLike decodes UUID strings that are using the following format:
//
//	"6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodeHashLike(t []byte) error {
	_, err := hex.Decode(u[:], t)
	return err
}

// decodeBraced decodes UUID strings that are using the following formats:
//
//	"{6ba7b810-9dad-11d1-80b4-00c04fd430c8}"
//	"{6ba7b8109dad11d180b400c04fd430c8}".
func (u *UUID) decodeBraced(t []byte) error {
	l := len(t)

	if t[0] != '{' || t[l-1] != '}' {
		return fmt.Errorf("%w: incorrect UUID format in string %q", ErrInvalid, t)
	}

	return u.decodePlain(t[1 : l-1])
}

// decodeURN decodes UUID strings that are using the following formats:
//
//	"urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
//	"urn:uuid:6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodeURN(t []byte) error {
	total := len(t)

	urnUUIDPrefix := t[:9]

	if !bytes.Equal(urnUUIDPrefix, urnPrefix) {
		return fmt.Errorf("%w: incorrect UUID format in string %q", ErrInvalid, t)
	}

	return u.decodePlain(t[9:total])
}

// decodePlain decodes UUID strings that are using the following formats:
//
//	"6ba7b810-9dad-11d1-80b4-00c04fd430c8" or in hash-like format
//	"6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodePlain(t []byte) error {
	switch len(t) {
	case 32:
		return u.decodeHashLike(t)
	case 36:
		return u.decodeCanonical(t)
	default:
		return fmt.Errorf("%w: incorrect UUID length %d in string %q", ErrInvalid, len(t), t)
	}
}
