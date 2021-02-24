package russian

import (
	"github.com/muonsoft/validation/message"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

var Messages = map[language.Tag]map[string]catalog.Message{
	language.Russian: {
		message.NotBlank: catalog.String("Значение не должно быть пустым."),
		message.Blank:    catalog.String("Значение должно быть пустым."),
		message.CountTooFew: plural.Selectf(1, "",
			plural.One, "Эта коллекция должна содержать {{ limit }} элемент или больше.",
			plural.Few, "Эта коллекция должна содержать {{ limit }} элемента или больше.",
			plural.Other, "Эта коллекция должна содержать {{ limit }} элементов или больше."),
		message.CountTooMany: plural.Selectf(1, "",
			plural.One, "Эта коллекция должна содержать {{ limit }} элемент или меньше.",
			plural.Few, "Эта коллекция должна содержать {{ limit }} элемента или меньше.",
			plural.Other, "Эта коллекция должна содержать {{ limit }} элементов или меньше."),
		message.CountExact: plural.Selectf(1, "",
			plural.One, "Эта коллекция должна содержать ровно {{ limit }} элемент.",
			plural.Few, "Эта коллекция должна содержать ровно {{ limit }} элемента.",
			plural.Other, "Эта коллекция должна содержать ровно {{ limit }} элементов."),
	},
}
