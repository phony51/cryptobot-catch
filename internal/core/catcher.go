package core

import (
	"context"
	"cryptobot-catch/internal/core/cheques"
	"cryptobot-catch/pkg/wallets"
	"github.com/gotd/td/tg"
	"go.uber.org/multierr"
)

type UpdateAnyMessage interface {
	GetMessage() tg.MessageClass
}

type updateHandler = func(context.Context, UpdateAnyMessage) error

type Catcher struct {
	handlers   map[uint32]updateHandler
	extractors []cheques.Extractor
	wallet     wallets.Wallet
}

func (c *Catcher) NewMessageHandle(ctx context.Context, update UpdateAnyMessage) error {
	if msg, ok := update.GetMessage().(*tg.Message); ok {
		for i := 0; i < len(c.extractors); i++ {
			if chequeID, found := c.extractors[i].Extract(msg); found {
				return c.wallet.ActivateCheque(ctx, chequeID)
			}
		}
	}
	return nil
}

func (c *Catcher) NewChannelMessageHandle(ctx context.Context, update UpdateAnyMessage) error {
	if msg, ok := update.GetMessage().(*tg.Message); ok {
		for i := 0; i < len(c.extractors); i++ {
			if chequeID, found := c.extractors[i].Extract(msg); found {
				return c.wallet.ActivateCheque(ctx, chequeID)
			}
		}
	}
	return nil
}

func (c *Catcher) EditMessageHandle(ctx context.Context, update UpdateAnyMessage) error {
	if msg, ok := update.GetMessage().(*tg.Message); ok {
		for i := 0; i < len(c.extractors); i++ {
			if chequeID, found := c.extractors[i].Extract(msg); found {
				return c.wallet.ActivateCheque(ctx, chequeID)
			}
		}
	}
	return nil
}

func (c *Catcher) EditChannelMessageHandle(ctx context.Context, update UpdateAnyMessage) error {
	if msg, ok := update.GetMessage().(*tg.Message); ok {
		for i := 0; i < len(c.extractors); i++ {
			if chequeID, found := c.extractors[i].Extract(msg); found {
				return c.wallet.ActivateCheque(ctx, chequeID)
			}
		}
	}
	return nil
}

func (c *Catcher) Handle(ctx context.Context, u tg.UpdatesClass) error {
	var upds []tg.UpdateClass
	switch upd := u.(type) {
	case *tg.Updates:
		upds = upd.Updates
	case *tg.UpdatesCombined:
		upds = upd.Updates
	case *tg.UpdateShort:
		upds = []tg.UpdateClass{upd.Update}
	default:
		return nil
	}

	var err error
	for i := 0; i < len(upds); i++ {
		if h, ok := c.handlers[upds[i].TypeID()]; ok {
			multierr.AppendInto(&err, h(ctx, upds[i].(UpdateAnyMessage)))
		}
	}
	return err
}

func NewCatcher(extractors []cheques.Extractor, wallet wallets.Wallet) *Catcher {
	handlers := make(map[uint32]updateHandler)

	c := &Catcher{
		handlers,
		extractors,
		wallet,
	}
	handlers[tg.UpdateNewMessageTypeID] = c.NewMessageHandle
	handlers[tg.UpdateEditMessageTypeID] = c.EditMessageHandle
	handlers[tg.UpdateNewChannelMessageTypeID] = c.NewChannelMessageHandle
	handlers[tg.UpdateEditChannelMessageTypeID] = c.EditChannelMessageHandle

	return c
}
