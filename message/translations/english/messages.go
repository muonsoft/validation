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
		message.Blank: catalog.String(message.Blank),
		message.CountExact: plural.Selectf(1, "",
			plural.One, "This collection should contain exactly {{ limit }} element.",
			plural.Other, "This collection should contain exactly {{ limit }} elements."),
		message.CountTooFew: plural.Selectf(1, "",
			plural.One, "This collection should contain {{ limit }} element or more.",
			plural.Other, "This collection should contain {{ limit }} elements or more."),
		message.CountTooMany: plural.Selectf(1, "",
			plural.One, "This collection should contain {{ limit }} element or less.",
			plural.Other, "This collection should contain {{ limit }} elements or less."),
		message.Equal:           catalog.String(message.Equal),
		message.False:           catalog.String(message.False),
		message.InvalidEmail:    catalog.String(message.InvalidEmail),
		message.InvalidHostname: catalog.String(message.InvalidHostname),
		message.InvalidIP:       catalog.String(message.InvalidIP),
		message.InvalidURL:      catalog.String(message.InvalidURL),
		message.LengthExact: plural.Selectf(1, "",
			plural.One, "This value should have exactly {{ limit }} character.",
			plural.Other, "This value should have exactly {{ limit }} characters."),
		message.LengthTooFew: plural.Selectf(1, "",
			plural.One, "This value is too short. It should have {{ limit }} character or more.",
			plural.Other, "This value is too short. It should have {{ limit }} characters or more."),
		message.LengthTooMany: plural.Selectf(1, "",
			plural.One, "This value is too long. It should have {{ limit }} character or less.",
			plural.Other, "This value is too long. It should have {{ limit }} characters or less."),
		message.Nil:               catalog.String(message.Nil),
		message.NoSuchChoice:      catalog.String(message.NoSuchChoice),
		message.NotBlank:          catalog.String(message.NotBlank),
		message.NotEqual:          catalog.String(message.NotEqual),
		message.NotInRange:        catalog.String(message.NotInRange),
		message.NotNegative:       catalog.String(message.NotNegative),
		message.NotNegativeOrZero: catalog.String(message.NotNegativeOrZero),
		message.NotNil:            catalog.String(message.NotNil),
		message.NotPositive:       catalog.String(message.NotPositive),
		message.NotPositiveOrZero: catalog.String(message.NotPositiveOrZero),
		message.NotValid:          catalog.String(message.NotValid),
		message.ProhibitedIP:      catalog.String(message.ProhibitedIP),
		message.TooEarly:          catalog.String(message.TooEarly),
		message.TooEarlyOrEqual:   catalog.String(message.TooEarlyOrEqual),
		message.TooHigh:           catalog.String(message.TooHigh),
		message.TooHighOrEqual:    catalog.String(message.TooHighOrEqual),
		message.TooLate:           catalog.String(message.TooLate),
		message.TooLateOrEqual:    catalog.String(message.TooLateOrEqual),
		message.TooLow:            catalog.String(message.TooLow),
		message.TooLowOrEqual:     catalog.String(message.TooLowOrEqual),
		message.True:              catalog.String(message.True),
	},
}
