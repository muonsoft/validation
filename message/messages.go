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

// Package message contains violation message texts. They are not protected by
// backward compatibility rules and can be changed at any time, even in patch versions.
package message

const (
	Blank             = "This value should be blank."
	CountExact        = "This collection should contain exactly {{ limit }} element(s)."
	CountTooFew       = "This collection should contain {{ limit }} element(s) or more."
	CountTooMany      = "This collection should contain {{ limit }} element(s) or less."
	Equal             = "This value should be equal to {{ comparedValue }}."
	False             = "This value should be false."
	InvalidEmail      = "This value is not a valid email address."
	InvalidHostname   = "This value is not a valid hostname."
	InvalidIP         = "This is not a valid IP address."
	InvalidJSON       = "This value should be valid JSON."
	InvalidURL        = "This value is not a valid URL."
	LengthExact       = "This value should have exactly {{ limit }} character(s)."
	LengthTooFew      = "This value is too short. It should have {{ limit }} character(s) or more."
	LengthTooMany     = "This value is too long. It should have {{ limit }} character(s) or less."
	Nil               = "This value should be nil."
	NoSuchChoice      = "The value you selected is not a valid choice."
	NotBlank          = "This value should not be blank."
	NotEqual          = "This value should not be equal to {{ comparedValue }}."
	NotInRange        = "This value should be between {{ min }} and {{ max }}."
	NotNegative       = "This value should be negative."
	NotNegativeOrZero = "This value should be either negative or zero."
	NotNil            = "This value should not be nil."
	NotPositive       = "This value should be positive."
	NotPositiveOrZero = "This value should be either positive or zero."
	NotValid          = "This value is not valid."
	ProhibitedIP      = "This IP address is prohibited to use."
	TooEarly          = "This value should be later than {{ comparedValue }}."
	TooEarlyOrEqual   = "This value should be later than or equal to {{ comparedValue }}."
	TooHigh           = "This value should be less than {{ comparedValue }}."
	TooHighOrEqual    = "This value should be less than or equal to {{ comparedValue }}."
	TooLate           = "This value should be earlier than {{ comparedValue }}."
	TooLateOrEqual    = "This value should be earlier than or equal to {{ comparedValue }}."
	TooLow            = "This value should be greater than {{ comparedValue }}."
	TooLowOrEqual     = "This value should be greater than or equal to {{ comparedValue }}."
	True              = "This value should be true."
)
