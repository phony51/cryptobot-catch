package updates

import (
	"context"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

type UpdateXMessage interface {
	GetMessage() tg.MessageClass
}

type UpdateXMessageHandlerFunc[U UpdateXMessage] func(ctx context.Context, e tg.Entities, update U) error

func UpdateXMessageHandler[U UpdateXMessage](messages chan<- *tg.Message) UpdateXMessageHandlerFunc[U] {
	return func(ctx context.Context, e tg.Entities, update U) error {
		if msg, ok := update.GetMessage().(*tg.Message); ok {
			messages <- msg
		}
		return nil
	}
}

func UpdatesManager(messages chan<- *tg.Message) *updates.Manager {
	d := tg.NewUpdateDispatcher()

	d.OnNewMessage(tg.NewMessageHandler(UpdateXMessageHandler[*tg.UpdateNewMessage](messages)))
	d.OnNewChannelMessage(tg.NewChannelMessageHandler(UpdateXMessageHandler[*tg.UpdateNewChannelMessage](messages)))
	d.OnEditMessage(tg.EditMessageHandler(UpdateXMessageHandler[*tg.UpdateEditMessage](messages)))
	d.OnEditChannelMessage(tg.EditChannelMessageHandler(UpdateXMessageHandler[*tg.UpdateEditChannelMessage](messages)))
	return updates.New(updates.Config{
		Handler: d,
		Logger:  zap.L(),
	})
}
