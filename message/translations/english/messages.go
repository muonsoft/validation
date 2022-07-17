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

// Package english contains violation message texts translated into English language.
// Values are not protected by backward compatibility rules and can be changed at any time, even in patch versions.
package english

import (
	"github.com/muonsoft/validation/message"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

var Messages = map[language.Tag]map[string]catalog.Message{
	language.English: {
		message.NotBlank: catalog.String(message.NotBlank),
		message.NotExactCount: plural.Selectf(1, "",
			plural.One, "This collection should contain exactly {{ limit }} element.",
			plural.Other, "This collection should contain exactly {{ limit }} elements."),
		message.TooFewElements: plural.Selectf(1, "",
			plural.One, "This collection should contain {{ limit }} element or more.",
			plural.Other, "This collection should contain {{ limit }} elements or more."),
		message.TooManyElements: plural.Selectf(1, "",
			plural.One, "This collection should contain {{ limit }} element or less.",
			plural.Other, "This collection should contain {{ limit }} elements or less."),
		message.NotEqual:        catalog.String(message.NotEqual),
		message.NotFalse:        catalog.String(message.NotFalse),
		message.InvalidDate:     catalog.String(message.InvalidDate),
		message.InvalidDateTime: catalog.String(message.InvalidDateTime),
		message.InvalidEAN13:    catalog.String(message.InvalidEAN13),
		message.InvalidEAN8:     catalog.String(message.InvalidEAN8),
		message.InvalidEmail:    catalog.String(message.InvalidEmail),
		message.InvalidHostname: catalog.String(message.InvalidHostname),
		message.InvalidIP:       catalog.String(message.InvalidIP),
		message.InvalidJSON:     catalog.String(message.InvalidJSON),
		message.InvalidTime:     catalog.String(message.InvalidTime),
		message.InvalidUPCA:     catalog.String(message.InvalidUPCA),
		message.InvalidUPCE:     catalog.String(message.InvalidUPCE),
		message.InvalidURL:      catalog.String(message.InvalidURL),
		message.NotExactLength: plural.Selectf(1, "",
			plural.One, "This value should have exactly {{ limit }} character.",
			plural.Other, "This value should have exactly {{ limit }} characters."),
		message.TooShort: plural.Selectf(1, "",
			plural.One, "This value is too short. It should have {{ limit }} character or more.",
			plural.Other, "This value is too short. It should have {{ limit }} characters or more."),
		message.TooLong: plural.Selectf(1, "",
			plural.One, "This value is too long. It should have {{ limit }} character or less.",
			plural.Other, "This value is too long. It should have {{ limit }} characters or less."),
		message.NotNil:            catalog.String(message.NotNil),
		message.NoSuchChoice:      catalog.String(message.NoSuchChoice),
		message.IsBlank:           catalog.String(message.IsBlank),
		message.IsEqual:           catalog.String(message.IsEqual),
		message.NotInRange:        catalog.String(message.NotInRange),
		message.NotInteger:        catalog.String(message.NotInteger),
		message.NotNegative:       catalog.String(message.NotNegative),
		message.NotNegativeOrZero: catalog.String(message.NotNegativeOrZero),
		message.IsNil:             catalog.String(message.IsNil),
		message.NotNumeric:        catalog.String(message.NotNumeric),
		message.NotPositive:       catalog.String(message.NotPositive),
		message.NotPositiveOrZero: catalog.String(message.NotPositiveOrZero),
		message.NotUnique:         catalog.String(message.NotUnique),
		message.NotValid:          catalog.String(message.NotValid),
		message.ProhibitedIP:      catalog.String(message.ProhibitedIP),
		message.ProhibitedURL:     catalog.String(message.ProhibitedURL),
		message.TooEarly:          catalog.String(message.TooEarly),
		message.TooEarlyOrEqual:   catalog.String(message.TooEarlyOrEqual),
		message.TooHigh:           catalog.String(message.TooHigh),
		message.TooHighOrEqual:    catalog.String(message.TooHighOrEqual),
		message.TooLate:           catalog.String(message.TooLate),
		message.TooLateOrEqual:    catalog.String(message.TooLateOrEqual),
		message.TooLow:            catalog.String(message.TooLow),
		message.TooLowOrEqual:     catalog.String(message.TooLowOrEqual),
		message.NotTrue:           catalog.String(message.NotTrue),
	},
}
