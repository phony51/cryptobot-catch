package detecting

import (
	"fmt"
	"github.com/gotd/td/tg"
	"regexp"
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

type RegexChequeIDDetectStrategy struct{}

var chequeIDPattern = regexp.MustCompile("CQ([A-Za-z0-9]{10})")

func (s *RegexChequeIDDetectStrategy) ChequeID(msg *tg.Message) (string, bool) {
	found := chequeIDPattern.FindStringSubmatch(msg.Message)
	if len(found) != 0 {
		return found[1], true
	}
	return "", false
}
