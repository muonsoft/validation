package translations

import (
	"fmt"

	"github.com/muonsoft/validation/message/translations/english"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

// Translator is used as a mechanism for message translations. It is based on the "golang.org/x/text" package.
type Translator struct {
	defaultLanguage language.Tag
	messages        *catalog.Builder
	printers        map[language.Tag]*message.Printer
}

type TranslatorOption func(translator *Translator) error

// DefaultLanguage option is used to set up the default language for translation of violation messages.
func DefaultLanguage(tag language.Tag) TranslatorOption {
	return func(translator *Translator) error {
		translator.defaultLanguage = tag
		return nil
	}
}

// SetTranslations option is used to load translation messages into the translator.
//
// By default, all violation messages are generated in the English language with pluralization capabilities.
// To use a custom language you have to load translations on translator initialization.
// Built-in translations are available in the sub-packages of the package "github.com/muonsoft/message/translations".
// The translation mechanism is provided by the "golang.org/x/text" package (be aware, it has no stable version yet).
func SetTranslations(messages map[language.Tag]map[string]catalog.Message) TranslatorOption {
	return func(translator *Translator) error {
		return translator.setMessages(messages)
	}
}

// NewTranslator creates an instance of the Translator for violation messages.
func NewTranslator(options ...TranslatorOption) (*Translator, error) {
	translator := &Translator{
		defaultLanguage: language.English,
		messages:        catalog.NewBuilder(),
		printers:        map[language.Tag]*message.Printer{},
	}
	err := translator.setMessages(english.Messages)
	if err != nil {
		return nil, err
	}

	for _, setOption := range options {
		err = setOption(translator)
		if err != nil {
			return nil, err
		}
	}

	err = translator.checkDefaultLanguageIsLoaded()
	if err != nil {
		return nil, err
	}

	for _, tag := range translator.messages.Languages() {
		translator.printers[tag] = message.NewPrinter(tag, message.Catalog(translator.messages))
	}

	return translator, nil
}

func (translator *Translator) Translate(tag language.Tag, message string, pluralCount int) string {
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
		return message
	}

	return printer.Sprintf(message, pluralCount)
}

func (translator *Translator) setMessages(messages map[language.Tag]map[string]catalog.Message) error {
	err := translator.setMessagesToCatalog(messages)
	if err != nil {
		return fmt.Errorf("failed to load translations: %w", err)
	}

	return nil
}

func (translator *Translator) setMessagesToCatalog(messages map[language.Tag]map[string]catalog.Message) error {
	for tag, tagMessages := range messages {
		for key, msg := range tagMessages {
			err := translator.messages.Set(tag, key, msg)
			if err != nil {
				return fmt.Errorf(`failed to set message "%s" for language %s: %w`, key, tag, err)
			}
		}
	}

	return nil
}

func (translator *Translator) checkDefaultLanguageIsLoaded() error {
	languages := translator.messages.Languages()

	for _, tag := range languages {
		if tag == translator.defaultLanguage {
			return nil
		}
	}

	return fmt.Errorf(`%w: missing messages for language "%s"`, errDefaultLanguageNotLoaded, translator.defaultLanguage)
}
