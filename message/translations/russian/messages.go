// The source code of the messages is taken from the Symfony Validator component
// See https://github.com/symfony/validator/blob/5.x/Resources/translations/validators.ru.xlf
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

// Package russian contains violation message texts translated into Russian language.
// Values are not protected by backward compatibility rules and can be changed at any time, even in patch versions.
package russian

import (
	"github.com/muonsoft/validation/message"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

var Messages = map[language.Tag]map[string]catalog.Message{
	language.Russian: {
		message.Blank: catalog.String("Значение должно быть пустым."),
		message.CountExact: plural.Selectf(1, "",
			plural.One, "Эта коллекция должна содержать ровно {{ limit }} элемент.",
			plural.Few, "Эта коллекция должна содержать ровно {{ limit }} элемента.",
			plural.Other, "Эта коллекция должна содержать ровно {{ limit }} элементов."),
		message.CountTooFew: plural.Selectf(1, "",
			plural.One, "Эта коллекция должна содержать {{ limit }} элемент или больше.",
			plural.Few, "Эта коллекция должна содержать {{ limit }} элемента или больше.",
			plural.Other, "Эта коллекция должна содержать {{ limit }} элементов или больше."),
		message.CountTooMany: plural.Selectf(1, "",
			plural.One, "Эта коллекция должна содержать {{ limit }} элемент или меньше.",
			plural.Few, "Эта коллекция должна содержать {{ limit }} элемента или меньше.",
			plural.Other, "Эта коллекция должна содержать {{ limit }} элементов или меньше."),
		message.Equal:           catalog.String("Значение должно быть равно {{ comparedValue }}."),
		message.False:           catalog.String("Значение должно быть ложным."),
		message.InvalidEmail:    catalog.String("Значение адреса электронной почты недопустимо."),
		message.InvalidHostname: catalog.String("Значение не является корректным именем хоста."),
		message.InvalidIP:       catalog.String("Значение не является допустимым IP адресом."),
		message.InvalidJSON:     catalog.String("Значение должно быть корректным JSON."),
		message.InvalidURL:      catalog.String("Значение не является допустимым URL."),
		message.LengthExact: plural.Selectf(1, "",
			plural.One, "Значение должно быть равно {{ limit }} символу.",
			plural.Few, "Значение должно быть равно {{ limit }} символам.",
			plural.Other, "Значение должно быть равно {{ limit }} символам."),
		message.LengthTooFew: plural.Selectf(1, "",
			plural.One, "Значение слишком короткое. Должно быть равно {{ limit }} символу или больше.",
			plural.Few, "Значение слишком короткое. Должно быть равно {{ limit }} символам или больше.",
			plural.Other, "Значение слишком короткое. Должно быть равно {{ limit }} символам или больше."),
		message.LengthTooMany: plural.Selectf(1, "",
			plural.One, "Значение слишком длинное. Должно быть равно {{ limit }} символу или меньше.",
			plural.Few, "Значение слишком длинное. Должно быть равно {{ limit }} символам или меньше.",
			plural.Other, "Значение слишком длинное. Должно быть равно {{ limit }} символам или меньше."),
		message.Nil:               catalog.String("Значение должно быть nil."),
		message.NoSuchChoice:      catalog.String("Выбранное Вами значение недопустимо."),
		message.NotBlank:          catalog.String("Значение не должно быть пустым."),
		message.NotEqual:          catalog.String("Значение не должно быть равно {{ comparedValue }}."),
		message.NotInRange:        catalog.String("Значение должно быть между {{ min }} и {{ max }}."),
		message.NotNegative:       catalog.String("Значение должно быть отрицательным."),
		message.NotNegativeOrZero: catalog.String("Значение должно быть отрицательным или равным нулю."),
		message.NotNil:            catalog.String("Значение не должно быть nil."),
		message.NotPositive:       catalog.String("Значение должно быть положительным."),
		message.NotPositiveOrZero: catalog.String("Значение должно быть положительным или равным нулю."),
		message.NotValid:          catalog.String("Значение недопустимо."),
		message.ProhibitedIP:      catalog.String("Этот IP-адрес запрещено использовать."),
		message.TooEarly:          catalog.String("Значение должно быть позже чем {{ comparedValue }}."),
		message.TooEarlyOrEqual:   catalog.String("Значение должно быть позже или равно {{ comparedValue }}."),
		message.TooHigh:           catalog.String("Значение должно быть меньше чем {{ comparedValue }}."),
		message.TooHighOrEqual:    catalog.String("Значение должно быть меньше или равно {{ comparedValue }}."),
		message.TooLate:           catalog.String("Значение должно быть раньше чем {{ comparedValue }}."),
		message.TooLateOrEqual:    catalog.String("Значение должно быть раньше или равно {{ comparedValue }}."),
		message.TooLow:            catalog.String("Значение должно быть больше чем {{ comparedValue }}."),
		message.TooLowOrEqual:     catalog.String("Значение должно быть больше или равно {{ comparedValue }}."),
		message.True:              catalog.String("Значение должно быть истинным."),
	},
}
