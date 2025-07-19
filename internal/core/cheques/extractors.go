package cheques

import (
	"github.com/gotd/td/tg"
	"regexp"
	"strings"
)

type ExtractFunc func(message *tg.Message) (string, bool)

type Extractor interface {
	Name() string
	Extract(message *tg.Message) (string, bool)
}

type InlineExtractor struct{}

func (ie *InlineExtractor) Name() string {
	return "inline"
}

const inlineChequeURLPrefix = "http://t.me/send?start=CQ"

func (ie *InlineExtractor) Extract(message *tg.Message) (string, bool) {
	if keyboard, ok := message.ReplyMarkup.(*tg.ReplyInlineMarkup); ok {
		if button, ok := keyboard.Rows[0].Buttons[0].(*tg.KeyboardButtonURL); ok {
			if strings.HasPrefix(button.URL, inlineChequeURLPrefix) {
				return button.URL[len(inlineChequeURLPrefix):], true
			}
		}
	}
	return "", false
}

const chequePrefix = "CQ"

var chequeIDPattern = regexp.MustCompile(chequePrefix + "[A-Za-z0-9]{10}")

type TextExtractor struct{}

func (te *TextExtractor) Name() string {
	return "text"
}

func (te *TextExtractor) Extract(message *tg.Message) (string, bool) {
	found := chequeIDPattern.FindString(message.Message)
	if found != "" {
		return found[len(chequePrefix):], true
	}
	return "", false
}
