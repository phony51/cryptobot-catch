package extracting

import (
	"github.com/gotd/td/tg"
	"regexp"
	"strings"
)

type ExtractFunc func(message *tg.Message) (string, bool)
type Extractor struct {
	Name    string
	Extract ExtractFunc
}

var (
	chequeIDPattern = regexp.MustCompile("CQ([A-Za-z0-9]{10})")
	InlineExtractor = Extractor{
		"inline",
		func(message *tg.Message) (string, bool) {
			if keyboard, ok := message.ReplyMarkup.(*tg.ReplyInlineMarkup); ok {
				if button, ok := keyboard.Rows[0].Buttons[0].(*tg.KeyboardButtonURL); ok {
					return strings.CutPrefix(button.URL, "http://t.me/send?start=CQ")
				}
			}
			return "", false
		},
	}
	TextExtractor = Extractor{
		"text",
		func(message *tg.Message) (string, bool) {
			found := chequeIDPattern.FindStringSubmatch(message.Message)
			if len(found) != 0 {
				return found[1], true
			}
			return "", false
		},
	}
)
