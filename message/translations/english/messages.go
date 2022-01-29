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
	"github.com/muonsoft/validation/code"
	"github.com/muonsoft/validation/message"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

var Messages = map[language.Tag]map[string]catalog.Message{
	language.English: {
		message.Templates[code.Blank]: catalog.String(message.Templates[code.Blank]),
		message.Templates[code.CountExact]: plural.Selectf(1, "",
			plural.One, "This collection should contain exactly {{ limit }} element.",
			plural.Other, "This collection should contain exactly {{ limit }} elements."),
		message.Templates[code.CountTooFew]: plural.Selectf(1, "",
			plural.One, "This collection should contain {{ limit }} element or more.",
			plural.Other, "This collection should contain {{ limit }} elements or more."),
		message.Templates[code.CountTooMany]: plural.Selectf(1, "",
			plural.One, "This collection should contain {{ limit }} element or less.",
			plural.Other, "This collection should contain {{ limit }} elements or less."),
		message.Templates[code.Equal]:           catalog.String(message.Templates[code.Equal]),
		message.Templates[code.False]:           catalog.String(message.Templates[code.False]),
		message.Templates[code.InvalidEAN13]:    catalog.String(message.Templates[code.InvalidEAN13]),
		message.Templates[code.InvalidEAN8]:     catalog.String(message.Templates[code.InvalidEAN8]),
		message.Templates[code.InvalidEmail]:    catalog.String(message.Templates[code.InvalidEmail]),
		message.Templates[code.InvalidHostname]: catalog.String(message.Templates[code.InvalidHostname]),
		message.Templates[code.InvalidIP]:       catalog.String(message.Templates[code.InvalidIP]),
		message.Templates[code.InvalidJSON]:     catalog.String(message.Templates[code.InvalidJSON]),
		message.Templates[code.InvalidUPCA]:     catalog.String(message.Templates[code.InvalidUPCA]),
		message.Templates[code.InvalidUPCE]:     catalog.String(message.Templates[code.InvalidUPCE]),
		message.Templates[code.InvalidURL]:      catalog.String(message.Templates[code.InvalidURL]),
		message.Templates[code.LengthExact]: plural.Selectf(1, "",
			plural.One, "This value should have exactly {{ limit }} character.",
			plural.Other, "This value should have exactly {{ limit }} characters."),
		message.Templates[code.LengthTooFew]: plural.Selectf(1, "",
			plural.One, "This value is too short. It should have {{ limit }} character or more.",
			plural.Other, "This value is too short. It should have {{ limit }} characters or more."),
		message.Templates[code.LengthTooMany]: plural.Selectf(1, "",
			plural.One, "This value is too long. It should have {{ limit }} character or less.",
			plural.Other, "This value is too long. It should have {{ limit }} characters or less."),
		message.Templates[code.Nil]:               catalog.String(message.Templates[code.Nil]),
		message.Templates[code.NoSuchChoice]:      catalog.String(message.Templates[code.NoSuchChoice]),
		message.Templates[code.NotBlank]:          catalog.String(message.Templates[code.NotBlank]),
		message.Templates[code.NotEqual]:          catalog.String(message.Templates[code.NotEqual]),
		message.Templates[code.NotInRange]:        catalog.String(message.Templates[code.NotInRange]),
		message.Templates[code.NotInteger]:        catalog.String(message.Templates[code.NotInteger]),
		message.Templates[code.NotNegative]:       catalog.String(message.Templates[code.NotNegative]),
		message.Templates[code.NotNegativeOrZero]: catalog.String(message.Templates[code.NotNegativeOrZero]),
		message.Templates[code.NotNil]:            catalog.String(message.Templates[code.NotNil]),
		message.Templates[code.NotNumeric]:        catalog.String(message.Templates[code.NotNumeric]),
		message.Templates[code.NotPositive]:       catalog.String(message.Templates[code.NotPositive]),
		message.Templates[code.NotPositiveOrZero]: catalog.String(message.Templates[code.NotPositiveOrZero]),
		message.Templates[code.NotUnique]:         catalog.String(message.Templates[code.NotUnique]),
		message.Templates[code.NotValid]:          catalog.String(message.Templates[code.NotValid]),
		message.Templates[code.ProhibitedIP]:      catalog.String(message.Templates[code.ProhibitedIP]),
		message.Templates[code.TooEarly]:          catalog.String(message.Templates[code.TooEarly]),
		message.Templates[code.TooEarlyOrEqual]:   catalog.String(message.Templates[code.TooEarlyOrEqual]),
		message.Templates[code.TooHigh]:           catalog.String(message.Templates[code.TooHigh]),
		message.Templates[code.TooHighOrEqual]:    catalog.String(message.Templates[code.TooHighOrEqual]),
		message.Templates[code.TooLate]:           catalog.String(message.Templates[code.TooLate]),
		message.Templates[code.TooLateOrEqual]:    catalog.String(message.Templates[code.TooLateOrEqual]),
		message.Templates[code.TooLow]:            catalog.String(message.Templates[code.TooLow]),
		message.Templates[code.TooLowOrEqual]:     catalog.String(message.Templates[code.TooLowOrEqual]),
		message.Templates[code.True]:              catalog.String(message.Templates[code.True]),
	},
}
