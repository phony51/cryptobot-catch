package core

import (
	"context"
	"cryptobot-catch/internal/core/cheques"
	"cryptobot-catch/pkg/cryptobot"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

type Catcher struct {
	filter    *cheques.Filter
	activator *cheques.Activator
	messages  chan<- *tg.Message
}

func (c *Catcher) Run(ctx context.Context, client *tg.Client) error {
	go func() { _ = c.activator.Run(ctx) }()
	go func() { _ = c.filter.Run(ctx) }()
	state, err := client.UpdatesGetState(ctx)
	logger := zap.L()
	logger.Info(state.String())
	if err != nil {
		return err
	}

	for {
		diff, err := client.UpdatesGetDifference(ctx, &tg.UpdatesGetDifferenceRequest{
			Pts:  state.Pts,
			Date: state.Date,
			Qts:  state.Qts,
		})
		if err != nil {
			logger.Info(err.Error())
			continue
		}
		switch d := diff.(type) {
		case *tg.UpdatesDifference:
			for _, msg := range d.NewMessages {
				if v, ok := msg.(*tg.Message); ok {
					c.messages <- v
				}
			}
			for _, u := range d.OtherUpdates {
				switch upd := u.(type) {
				case *tg.UpdateEditMessage:
					if v, ok := upd.Message.(*tg.Message); ok {
						c.messages <- v
					}
				case *tg.UpdateEditChannelMessage:
					if v, ok := upd.Message.(*tg.Message); ok {
						c.messages <- v
					}
				}
			}
			state = &d.State
		case *tg.UpdatesDifferenceTooLong:
			state, err = client.UpdatesGetState(ctx)
			if err != nil {
				logger.Info(err.Error())
				continue
			}
		}
	}
}

func NewCatcher(cryptoBot *cryptobot.CryptoBot, strategies ...cheques.DetectStrategy) *Catcher {
	messages := make(chan *tg.Message)
	chequeIDs := make(chan string)
	return &Catcher{
		activator: cheques.NewActivator(cryptoBot, chequeIDs),
		filter: cheques.NewFilter(
			strategies,
			messages,
			chequeIDs,
		),
		messages: messages,
	}
}
