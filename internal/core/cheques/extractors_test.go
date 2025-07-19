package cheques

import (
	. "cryptobot-catch/pkg/testing"
	"github.com/gotd/td/tg"
	"testing"
)

type extractorTestCase TestCase[*tg.Message, string]

var (
	textExtractorTestCases = []extractorTestCase{
		{
			&tg.Message{Message: "​ (https://imggen.send.tg/checks/image?asset=USDT&asset_amount=0.01999&fiat=RUB&fiat_amount=1.57&main=asset&v2)Чек\n\nСумма: 🪙 0.019994 USDT (1.57 RUB)\n\nЛюбой может активировать этот чек.\n\nСкопируйте ссылку, чтобы поделиться чеком:\nt.me/send?start=CQ8dBasuJKhG\n\n⚠️ Никогда не делайте скриншот вашего чека и не отправляйте его никому! Ссылку на чек могут использовать мошенники, чтобы получить доступ к вашим средствам."},
			"8dBasuJKhG",
		},

		{
			&tg.Message{Message: "t.me/send?start=CQa123b456YZ"},
			"a123b456YZ",
		},
		{
			&tg.Message{Message: "CQ123456789"},
			"",
		},
	}

	inlineExtractorTestCases = []extractorTestCase{
		{
			&tg.Message{
				ReplyMarkup: &tg.ReplyInlineMarkup{
					Rows: []tg.KeyboardButtonRow{
						{
							Buttons: []tg.KeyboardButtonClass{
								&tg.KeyboardButtonURL{
									URL: "http://t.me/send?start=CQ8dBasuJKhG",
								},
							},
						},
					},
				},
			}, "8dBasuJKhG",
		},
		{
			&tg.Message{
				ReplyMarkup: &tg.ReplyInlineMarkup{
					Rows: []tg.KeyboardButtonRow{
						{
							Buttons: []tg.KeyboardButtonClass{
								&tg.KeyboardButtonURL{
									URL: "http://t.me/send?start=wallet",
								},
							},
						},
					},
				},
			}, "",
		},
	}

	textExtractor   = TextExtractor{}
	inlineExtractor = InlineExtractor{}
)

func TestInlineExtractor_Extract(t *testing.T) {
	var actual string
	for _, c := range inlineExtractorTestCases {
		actual, _ = inlineExtractor.Extract(c.Data)
		if actual != c.Expected {
			t.Errorf("expected %s, got %s for keyboard %s",
				c.Expected, actual, c.Data.ReplyMarkup)
		}
	}
}

func BenchmarkInlineExtractor_Extract(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		inlineExtractor.Extract(inlineExtractorTestCases[0].Data)
	}
}

func TestTextExtractor_Extract(t *testing.T) {
	for _, c := range textExtractorTestCases {
		actual, _ := textExtractor.Extract(c.Data)
		if actual != c.Expected {
			t.Errorf("expected %q, got %q for message %q",
				c.Expected, actual, c.Data.Message)
		}
	}
}

func BenchmarkTextExtractor_Extract(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		textExtractor.Extract(textExtractorTestCases[0].Data)
	}
}
