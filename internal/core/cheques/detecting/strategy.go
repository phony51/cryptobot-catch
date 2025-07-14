package cheques

import (
	"fmt"
	"github.com/gotd/td/tg"
	"regexp"
	"strings"
)

type ChequeDetector interface {
	ChequeID(*tg.Message) (string, bool)
}

type DetectStrategy interface {
	ChequeDetector
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

type RegexFullChequeIDDetectStrategy struct{}

var fullChequeIDPattern = regexp.MustCompile("CQ([A-Za-z0-9]{10})")

func (s *RegexFullChequeIDDetectStrategy) ChequeID(msg *tg.Message) (string, bool) {
	found := fullChequeIDPattern.FindStringSubmatch(msg.Message)
	if len(found) != 0 {
		fmt.Println(found)
		return found[1], true
	}
	return "", false
}
