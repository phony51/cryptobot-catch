package core

import (
	"context"
	"cryptobot-catch/internal/core/cheques"
	"cryptobot-catch/pkg/cryptobot"
	"github.com/gotd/td/tg"
)

type Catcher struct {
	d          *tg.UpdateDispatcher
	extractors []cheques.Extractor
	cryptoBot  *cryptobot.CryptoBot
}

func (c *Catcher) NewMessageHandle(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
	if msg, ok := update.GetMessage().(*tg.Message); ok {
		for i := 0; i < len(c.extractors); i++ {
			if chequeID, found := c.extractors[i].Extract(msg); found {
				return c.cryptoBot.ActivateCheque(ctx, chequeID)
			}
		}
	}
	return nil
}

func (c *Catcher) NewChannelMessageHandle(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
	if msg, ok := update.GetMessage().(*tg.Message); ok {
		for i := 0; i < len(c.extractors); i++ {
			if chequeID, found := c.extractors[i].Extract(msg); found {
				return c.cryptoBot.ActivateCheque(ctx, chequeID)
			}
		}
	}
	return nil
}

func (c *Catcher) EditMessageHandle(ctx context.Context, e tg.Entities, update *tg.UpdateEditMessage) error {
	if msg, ok := update.GetMessage().(*tg.Message); ok {
		for i := 0; i < len(c.extractors); i++ {
			if chequeID, found := c.extractors[i].Extract(msg); found {
				return c.cryptoBot.ActivateCheque(ctx, chequeID)
			}
		}
	}
	return nil
}

func (c *Catcher) EditChannelMessageHandle(ctx context.Context, e tg.Entities, update *tg.UpdateEditChannelMessage) error {
	if msg, ok := update.GetMessage().(*tg.Message); ok {
		for i := 0; i < len(c.extractors); i++ {
			if chequeID, found := c.extractors[i].Extract(msg); found {
				return c.cryptoBot.ActivateCheque(ctx, chequeID)
			}
		}
	}
	return nil
}

func (c *Catcher) Dispatcher() *tg.UpdateDispatcher {
	return c.d
}

func NewCatcher(extractors []cheques.Extractor, cryptoBot *cryptobot.CryptoBot) *Catcher {
	d := tg.NewUpdateDispatcher()
	c := &Catcher{
		&d,
		extractors,
		cryptoBot,
	}

	d.OnNewMessage(c.NewMessageHandle)
	d.OnNewChannelMessage(c.NewChannelMessageHandle)
	d.OnEditMessage(c.EditMessageHandle)
	d.OnEditChannelMessage(c.EditChannelMessageHandle)

	return c
}
