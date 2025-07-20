package core

import (
	"context"
	"cryptobot-catch/internal/core/cheques"
	"cryptobot-catch/pkg/wallets"
	"github.com/gotd/td/tg"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type UpdateAnyMessage interface {
	tg.UpdateClass
	GetMessage() tg.MessageClass
}

type updateHandler = func(context.Context, UpdateAnyMessage) error

type Catcher struct {
	handlers   map[uint32]updateHandler
	extractors []cheques.Extractor
	wallet     wallets.Wallet
	logger     *zap.Logger
	api        *tg.Client
}

func (c *Catcher) SetAPI(api *tg.Client) {
	c.api = api
}

//func (c *Catcher) NewMessageHandle(ctx context.Context, update UpdateAnyMessage) error {
//	if msg, ok := update.GetMessage().(*tg.Message); ok {
//		for i := 0; i < len(c.extractors); i++ {
//			if chequeID, found := c.extractors[i].Extract(msg); found {
//				err := c.wallet.ActivateCheque(ctx, chequeID)
//				if err != nil {
//					c.logger.Error("failed to activate cheque", zap.Error(err))
//				}
//				c.logger.Debug("caught cheque",
//					zap.String("id", chequeID),
//					zap.String("extractor", c.extractors[i].Name()),
//					zap.String("updateType", update.TypeName()),
//				)
//				return err
//			}
//		}
//	}
//	return nil
//}
//
//func (c *Catcher) NewChannelMessageHandle(ctx context.Context, update UpdateAnyMessage) error {
//	if msg, ok := update.GetMessage().(*tg.Message); ok {
//		for i := 0; i < len(c.extractors); i++ {
//			if chequeID, found := c.extractors[i].Extract(msg); found {
//				err := c.wallet.ActivateCheque(ctx, chequeID)
//				if err != nil {
//					c.logger.Error("failed to activate cheque", zap.Error(err))
//				}
//				c.logger.Debug("caught cheque",
//					zap.String("id", chequeID),
//					zap.String("extractor", c.extractors[i].Name()),
//					zap.String("updateType", update.TypeName()),
//				)
//				return err
//			}
//		}
//	}
//	return nil
//}
//
//func (c *Catcher) EditMessageHandle(ctx context.Context, update UpdateAnyMessage) error {
//	if msg, ok := update.GetMessage().(*tg.Message); ok {
//		for i := 0; i < len(c.extractors); i++ {
//			if chequeID, found := c.extractors[i].Extract(msg); found {
//				err := c.wallet.ActivateCheque(ctx, chequeID)
//				if err != nil {
//					c.logger.Error("failed to activate cheque", zap.Error(err))
//				}
//				c.logger.Debug("caught cheque",
//					zap.String("id", chequeID),
//					zap.String("extractor", c.extractors[i].Name()),
//					zap.String("updateType", update.TypeName()),
//				)
//				return err
//			}
//		}
//	}
//	return nil
//}
//
//func (c *Catcher) EditChannelMessageHandle(ctx context.Context, update UpdateAnyMessage) error {
//	if msg, ok := update.GetMessage().(*tg.Message); ok {
//		for i := 0; i < len(c.extractors); i++ {
//			if chequeID, found := c.extractors[i].Extract(msg); found {
//				err := c.wallet.ActivateCheque(ctx, chequeID)
//				if err != nil {
//					c.logger.Error("failed to activate cheque", zap.Error(err))
//				}
//				c.logger.Debug("caught cheque",
//					zap.String("id", chequeID),
//					zap.String("extractor", c.extractors[i].Name()),
//					zap.String("updateType", update.TypeName()),
//				)
//				return err
//			}
//		}
//	}
//	return nil
//}

func (c *Catcher) UpdateAnyMessageHandle(ctx context.Context, update UpdateAnyMessage) error {
	if msg, ok := update.GetMessage().(*tg.Message); ok {
		for i := 0; i < len(c.extractors); i++ {
			if chequeID, found := c.extractors[i].Extract(msg); found {
				err := c.wallet.ActivateCheque(ctx, chequeID)
				if err != nil {
					c.logger.Error("failed to activate cheque", zap.Error(err))
				}
				c.logger.Debug("caught cheque",
					zap.String("id", chequeID),
					zap.String("extractor", c.extractors[i].Name()),
					zap.String("updateType", update.TypeName()),
				)
				return err
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
	case *tg.UpdatesTooLong:
		state, err := c.api.UpdatesGetState(ctx)
		if err != nil {
			return err
		}
		_, err = c.api.UpdatesGetDifference(ctx, &tg.UpdatesGetDifferenceRequest{
			Pts:  state.Pts,
			Qts:  state.Qts,
			Date: state.Date,
		})
		if err != nil {
			return err
		}
	default:
		c.logger.Debug("unhandled updates", zap.Any("updates", upd))
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
		zap.L(),
		nil,
	}
	handlers[tg.UpdateNewMessageTypeID] = c.UpdateAnyMessageHandle
	handlers[tg.UpdateEditMessageTypeID] = c.UpdateAnyMessageHandle
	handlers[tg.UpdateNewChannelMessageTypeID] = c.UpdateAnyMessageHandle
	handlers[tg.UpdateEditChannelMessageTypeID] = c.UpdateAnyMessageHandle

	return c
}
