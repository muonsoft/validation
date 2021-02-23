package translations

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

func LoadMessages(cat *catalog.Builder, messages map[language.Tag]map[string]catalog.Message) error {
	for tag, tagMessages := range messages {
		for key, message := range tagMessages {
			err := cat.Set(tag, key, message)
			if err != nil {
				return fmt.Errorf("failed to set message '%s' for language %s: %w", key, tag, err)
			}
		}
	}

	return nil
}
