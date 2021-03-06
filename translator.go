package validation

import (
	"fmt"

	"github.com/muonsoft/validation/message/translations/english"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type Translator struct {
	defaultLanguage language.Tag
	messages        *catalog.Builder
	printers        map[language.Tag]*message.Printer
}

func newTranslator() *Translator {
	return &Translator{
		defaultLanguage: language.English,
		messages:        catalog.NewBuilder(),
		printers:        map[language.Tag]*message.Printer{},
	}
}

func (translator *Translator) init() error {
	err := translator.loadMessages(english.Messages)
	if err != nil {
		return err
	}

	return translator.checkDefaultLanguageIsLoaded()
}

func (translator *Translator) loadMessages(messages map[language.Tag]map[string]catalog.Message) error {
	err := translator.setMessagesToCatalog(messages)
	if err != nil {
		return fmt.Errorf("failed to load translations: %w", err)
	}

	for tag := range messages {
		translator.printers[tag] = message.NewPrinter(tag, message.Catalog(translator.messages))
	}

	return nil
}

func (translator *Translator) setMessagesToCatalog(messages map[language.Tag]map[string]catalog.Message) error {
	for tag, tagMessages := range messages {
		for key, msg := range tagMessages {
			err := translator.messages.Set(tag, key, msg)
			if err != nil {
				return fmt.Errorf("failed to set message '%s' for language %s: %w", key, tag, err)
			}
		}
	}

	return nil
}

func (translator *Translator) translate(tag language.Tag, msg string, pluralCount int) string {
	if tag == language.Und {
		tag = translator.defaultLanguage
	}
	printer := translator.printers[tag]
	if printer == nil {
		printer = translator.printers[tag.Parent()]
	}
	if printer == nil {
		printer = translator.printers[translator.defaultLanguage]
	}
	if printer == nil {
		return msg
	}

	return printer.Sprintf(msg, pluralCount)
}

func (translator *Translator) checkDefaultLanguageIsLoaded() error {
	languages := translator.messages.Languages()

	for _, tag := range languages {
		if tag == translator.defaultLanguage {
			return nil
		}
	}

	return fmt.Errorf("%w: missing messages for language '%s'", errDefaultLanguageNotLoaded, translator.defaultLanguage)
}
