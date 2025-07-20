package core_test

import (
	"context"
	"cryptobot-catch/internal/core"
	"cryptobot-catch/internal/core/cheques"
	"github.com/gotd/td/tg"
	"testing"
)

type MockWallet struct{}

func (w *MockWallet) ActivateCheque(ctx context.Context, chequeID string) error {
	return nil
}

var (
	textExtractor, _   = cheques.NewTextExtractor()
	inlineExtractor    = &cheques.InlineExtractor{}
	inlineChequeUpdate = &tg.UpdateEditChannelMessage{
		Message: &tg.Message{
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
		},
	}
	textChequeUpdate = &tg.UpdateNewChannelMessage{
		Message: &tg.Message{Message: "​ (https://imggen.send.tg/checks/image?asset=USDT&asset_amount=0.01999&fiat=RUB&fiat_amount=1.57&main=asset&v2)Чек\n\nСумма: 🪙 0.019994 USDT (1.57 RUB)\n\nЛюбой может активировать этот чек.\n\nСкопируйте ссылку, чтобы поделиться чеком:\nt.me/send?start=CQ8dBasuJKhG\n\n⚠️ Никогда не делайте скриншот вашего чека и не отправляйте его никому! Ссылку на чек могут использовать мошенники, чтобы получить доступ к вашим средствам."},
	}
)

func BenchmarkCatcher_Inline(b *testing.B) {
	c := core.NewCatcher([]cheques.Extractor{inlineExtractor}, &MockWallet{})
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EditChannelMessageHandle(ctx, inlineChequeUpdate)
	}
	b.ReportAllocs()
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "msg/s")
}

func BenchmarkCatcher_Text(b *testing.B) {
	c := core.NewCatcher([]cheques.Extractor{inlineExtractor}, &MockWallet{})
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.NewChannelMessageHandle(ctx, textChequeUpdate)
	}
	b.ReportAllocs()
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "msg/s")
}

func BenchmarkCatcher_Inline_Text(b *testing.B) {
	c := core.NewCatcher([]cheques.Extractor{inlineExtractor, textExtractor}, &MockWallet{})
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			_ = c.NewChannelMessageHandle(ctx, textChequeUpdate)
		} else {
			_ = c.EditChannelMessageHandle(ctx, inlineChequeUpdate)
		}
	}
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "msg/s")
}
