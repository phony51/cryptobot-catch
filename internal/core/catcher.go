package core

import (
	"context"
	"cryptobot-catch/internal/core/cheques"
	"cryptobot-catch/internal/core/cheques/detecting"
	"cryptobot-catch/pkg/cryptobot"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
)

type CatcherOptions struct {
	DetectStrategies []detecting.DetectStrategy
}

type Catcher struct {
	filter    *cheques.Filter
	activator *cheques.Activator
	gaps      *updates.Manager
}

func (c *Catcher) Run(ctx context.Context, client *tg.Client) error {
	go func() { _ = c.activator.Run(ctx) }()
	go func() { _ = c.filter.Run(ctx) }()
	return c.gaps.Run(ctx, client, 0, updates.AuthOptions{})
	//state, err := client.UpdatesGetState(ctx)
	//
	//logger := zap.L()
	//ticker := time.NewTicker(c.pollingInterval)
	//if err != nil {
	//	return err
	//}
	//for {
	//	select {
	//	case <-ctx.Done():
	//	case <-ticker.C:
	//		diff, err := client.UpdatesGetDifference(ctx, &tg.UpdatesGetDifferenceRequest{
	//			Pts:  state.Pts,
	//			Date: state.Date,
	//			Qts:  state.Qts,
	//		})
	//		if err != nil {
	//			continue
	//		}
	//		switch d := diff.(type) {
	//		case *tg.UpdatesDifference:
	//			logger.Debug(d.String())
	//			go func() {
	//				for _, msg := range d.NewMessages {
	//					if msg, ok := msg.(*tg.Message); ok {
	//						c.messages <- msg
	//					}
	//				}
	//			}()
	//			go func() {
	//				for _, u := range d.OtherUpdates {
	//					switch upd := u.(type) {
	//					case *tg.UpdateNewMessage:
	//						if msg, ok := upd.Message.(*tg.Message); ok {
	//							c.messages <- msg
	//						}
	//					case *tg.UpdateNewChannelMessage:
	//						if msg, ok := upd.Message.(*tg.Message); ok {
	//							c.messages <- msg
	//						}
	//					case *tg.UpdateEditMessage:
	//						if msg, ok := upd.Message.(*tg.Message); ok {
	//							c.messages <- msg
	//						}
	//					case *tg.UpdateEditChannelMessage:
	//						if msg, ok := upd.Message.(*tg.Message); ok {
	//							c.messages <- msg
	//						}
	//					}
	//				}
	//			}()
	//			state = &d.State
	//		case *tg.UpdatesDifferenceTooLong:
	//			state, err = client.UpdatesGetState(ctx)
	//			if err != nil {
	//				logger.Debug(err.Error())
	//			}
	//		}
	//	}
	//}
}

func NewCatcher(cryptoBot *cryptobot.CryptoBot, gaps *updates.Manager, messagesCh chan *tg.Message, options *CatcherOptions) *Catcher {
	chequeIDs := make(chan string)

	return &Catcher{
		activator: cheques.NewActivator(cryptoBot, chequeIDs),
		gaps:      gaps,
		filter: cheques.NewFilter(
			options.DetectStrategies,
			messagesCh,
			chequeIDs,
		),
	}
}

//type UpdateHandler struct {
//	Messages chan<- *tg.Message
//}
//
//func (h *UpdateHandler) Handle(ctx context.Context, u tg.UpdatesClass) error {
//	switch u := u.(type) {
//	case *tg.Updates:
//		for _, upd := range u.Updates {
//			fmt.Println(upd)
//			switch upd := upd.(type) {
//			case *tg.UpdateNewMessage:
//				if msg, ok := upd.Message.(*tg.Message); ok {
//					h.Messages <- msg
//				}
//			case *tg.UpdateNewChannelMessage:
//				if msg, ok := upd.Message.(*tg.Message); ok {
//					h.Messages <- msg
//				}
//			case *tg.UpdateEditMessage:
//				if msg, ok := upd.Message.(*tg.Message); ok {
//					h.Messages <- msg
//				}
//			case *tg.UpdateEditChannelMessage:
//
//			}
//		}
//	}
//	return nil
//}
