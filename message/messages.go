// The source code of the messages is taken from the Symfony Validator component
// See https://github.com/symfony/validator
//
// Copyright (c) 2004-2021 Fabien Potencier
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package message contains violation message templates. They are not protected by
// backward compatibility rules and can be changed at any time, even in patch versions.
// It is recommended not to use this code directly in your application.
package message

const (
	InvalidDate       = "This value is not a valid date."
	InvalidDateTime   = "This value is not a valid datetime."
	InvalidEAN13      = "This value is not a valid EAN-13."
	InvalidEAN8       = "This value is not a valid EAN-8."
	InvalidEmail      = "This value is not a valid email address."
	InvalidHostname   = "This value is not a valid hostname."
	InvalidIP         = "This is not a valid IP address."
	InvalidJSON       = "This value should be valid JSON."
	InvalidTime       = "This value is not a valid time."
	InvalidULID       = "This is not a valid ULID."
	InvalidUPCA       = "This value is not a valid UPC-A."
	InvalidUPCE       = "This value is not a valid UPC-E."
	InvalidURL        = "This value is not a valid URL."
	InvalidUUID       = "This is not a valid UUID."
	IsBlank           = "This value should not be blank."
	IsEqual           = "This value should not be equal to {{ comparedValue }}."
	IsNil             = "This value should not be nil."
	NoSuchChoice      = "The value you selected is not a valid choice."
	NotBlank          = "This value should be blank."
	NotDivisible      = "This value should be a multiple of {{ comparedValue }}."
	NotDivisibleCount = "The number of elements in this collection should be a multiple of {{ divisibleBy }}."
	NotEqual          = "This value should be equal to {{ comparedValue }}."
	NotExactCount     = "This collection should contain exactly {{ limit }} element(s)."
	NotExactLength    = "This value should have exactly {{ limit }} character(s)."
	NotFalse          = "This value should be false."
	NotInRange        = "This value should be between {{ min }} and {{ max }}."
	NotInteger        = "This value is not an integer."
	NotNegative       = "This value should be negative."
	NotNegativeOrZero = "This value should be either negative or zero."
	NotNil            = "This value should be nil."
	NotNumeric        = "This value is not a numeric."
	NotPositive       = "This value should be positive."
	NotPositiveOrZero = "This value should be either positive or zero."
	NotTrue           = "This value should be true."
	NotUnique         = "This collection should contain only unique elements."
	NotValid          = "This value is not valid."
	ProhibitedIP      = "This IP address is prohibited to use."
	ProhibitedURL     = "This URL is prohibited to use."
	TooEarly          = "This value should be later than {{ comparedValue }}."
	TooEarlyOrEqual   = "This value should be later than or equal to {{ comparedValue }}."
	TooFewElements    = "This collection should contain {{ limit }} element(s) or more."
	TooHigh           = "This value should be less than {{ comparedValue }}."
	TooHighOrEqual    = "This value should be less than or equal to {{ comparedValue }}."
	TooLate           = "This value should be earlier than {{ comparedValue }}."
	TooLateOrEqual    = "This value should be earlier than or equal to {{ comparedValue }}."
	TooLong           = "This value is too long. It should have {{ limit }} character(s) or less."
	TooLow            = "This value should be greater than {{ comparedValue }}."
	TooLowOrEqual     = "This value should be greater than or equal to {{ comparedValue }}."
	TooManyElements   = "This collection should contain {{ limit }} element(s) or less."
	TooShort          = "This value is too short. It should have {{ limit }} character(s) or more."
)
