package cheques

import (
	"github.com/gotd/td/tg"
	"strings"
)

type DetectStrategy interface {
	ChequeID(*tg.Message) (string, bool)
}

type InlineDetectStrategy struct{}

func (s *InlineDetectStrategy) ChequeID(msg *tg.Message) (string, bool) {
	if keyboard, ok := msg.ReplyMarkup.(*tg.ReplyInlineMarkup); ok {
		if button, ok := keyboard.Rows[0].Buttons[0].(*tg.KeyboardButtonURL); ok {
			return strings.CutPrefix(button.URL, "http://t.me/send?start=CQ")
		}
	}
	return "", false
}
