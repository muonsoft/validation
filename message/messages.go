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

import "github.com/muonsoft/validation/code"

var Templates = map[string]string{
	code.Blank:             "This value should be blank.",
	code.CountExact:        "This collection should contain exactly {{ limit }} element(s).",
	code.CountTooFew:       "This collection should contain {{ limit }} element(s) or more.",
	code.CountTooMany:      "This collection should contain {{ limit }} element(s) or less.",
	code.Equal:             "This value should be equal to {{ comparedValue }}.",
	code.False:             "This value should be false.",
	code.InvalidEmail:      "This value is not a valid email address.",
	code.InvalidHostname:   "This value is not a valid hostname.",
	code.InvalidIP:         "This is not a valid IP address.",
	code.InvalidJSON:       "This value should be valid JSON.",
	code.InvalidURL:        "This value is not a valid URL.",
	code.LengthExact:       "This value should have exactly {{ limit }} character(s).",
	code.LengthTooFew:      "This value is too short. It should have {{ limit }} character(s) or more.",
	code.LengthTooMany:     "This value is too long. It should have {{ limit }} character(s) or less.",
	code.Nil:               "This value should be nil.",
	code.NoSuchChoice:      "The value you selected is not a valid choice.",
	code.NotBlank:          "This value should not be blank.",
	code.NotEqual:          "This value should not be equal to {{ comparedValue }}.",
	code.NotInRange:        "This value should be between {{ min }} and {{ max }}.",
	code.NotNegative:       "This value should be negative.",
	code.NotNegativeOrZero: "This value should be either negative or zero.",
	code.NotNil:            "This value should not be nil.",
	code.NotPositive:       "This value should be positive.",
	code.NotPositiveOrZero: "This value should be either positive or zero.",
	code.NotUnique:         "This collection should contain only unique elements.",
	code.NotValid:          "This value is not valid.",
	code.ProhibitedIP:      "This IP address is prohibited to use.",
	code.TooEarly:          "This value should be later than {{ comparedValue }}.",
	code.TooEarlyOrEqual:   "This value should be later than or equal to {{ comparedValue }}.",
	code.TooHigh:           "This value should be less than {{ comparedValue }}.",
	code.TooHighOrEqual:    "This value should be less than or equal to {{ comparedValue }}.",
	code.TooLate:           "This value should be earlier than {{ comparedValue }}.",
	code.TooLateOrEqual:    "This value should be earlier than or equal to {{ comparedValue }}.",
	code.TooLow:            "This value should be greater than {{ comparedValue }}.",
	code.TooLowOrEqual:     "This value should be greater than or equal to {{ comparedValue }}.",
	code.True:              "This value should be true.",
}
