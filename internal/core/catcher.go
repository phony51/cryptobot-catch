package core

import (
	"context"
	"cryptobot-catch/internal/core/cheques"
	"cryptobot-catch/pkg/cryptobot"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
	"time"
)

type CatchOptions struct {
	PollingInterval  time.Duration
	DetectStrategies []cheques.DetectStrategy
}

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
	ticker := time.NewTicker(1000 * time.Millisecond)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
		case <-ticker.C:
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
				logger.Info(d.String())
				go func() {
					for _, msg := range d.NewMessages {
						if msg, ok := msg.(*tg.Message); ok {
							c.messages <- msg
						}
					}
				}()
				go func() {
					for _, u := range d.OtherUpdates {
						switch upd := u.(type) {
						case *tg.UpdateNewMessage:
							if msg, ok := upd.Message.(*tg.Message); ok {
								c.messages <- msg
							}
						case *tg.UpdateNewChannelMessage:
							if msg, ok := upd.Message.(*tg.Message); ok {
								c.messages <- msg
							}
						case *tg.UpdateEditMessage:
							if msg, ok := upd.Message.(*tg.Message); ok {
								c.messages <- msg
							}
						case *tg.UpdateEditChannelMessage:
							if msg, ok := upd.Message.(*tg.Message); ok {
								c.messages <- msg
							}
						}
					}
				}()
				state = &d.State
			case *tg.UpdatesDifferenceTooLong:
				state, err = client.UpdatesGetState(ctx)
				if err != nil {
					logger.Debug(err.Error())
				}
			}
		}
	}
}

func NewCatcher(cryptoBot *cryptobot.CryptoBot, options *CatchOptions) *Catcher {
	messages := make(chan *tg.Message)
	chequeIDs := make(chan string)
	return &Catcher{
		activator: cheques.NewActivator(cryptoBot, chequeIDs),
		filter: cheques.NewFilter(
			options.DetectStrategies,
			messages,
			chequeIDs,
		),
		messages: messages,
	}
}
