package cryptobot

import (
	"context"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type CryptoBot struct {
	BotID          int64
	messageBuilder *message.RequestBuilder
}

func NewCryptoBot(sender *message.Sender, peer *tg.InputPeerUser) *CryptoBot {

	return &CryptoBot{
		BotID:          peer.UserID,
		messageBuilder: sender.To(peer),
	}
}

func (c *CryptoBot) ActivateCheque(ctx context.Context, chequeID string) error {
	_, err := c.messageBuilder.Textf(ctx, "/start CQ%s", chequeID)
	if err != nil {
		return err
	}
	return nil
}
