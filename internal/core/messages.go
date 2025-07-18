package core

import (
	"context"
	"github.com/gotd/td/tg"
)

type MessagePipe struct {
	messagesCh chan *tg.Message
	d          tg.UpdateDispatcher
}

func NewMessagePipe() MessagePipe {
	d := tg.NewUpdateDispatcher()

	pud := MessagePipe{
		make(chan *tg.Message),
		d,
	}

	d.OnNewMessage(pud.NewMessageHandle)
	d.OnNewChannelMessage(pud.NewChannelMessageHandle)
	d.OnEditMessage(pud.EditMessageHandle)
	d.OnEditChannelMessage(pud.EditChannelMessageHandle)
	return pud
}

func (p *MessagePipe) NewMessageHandle(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
	if msg, ok := update.GetMessage().(*tg.Message); ok {
		p.messagesCh <- msg
	}
	return nil
}

func (p *MessagePipe) NewChannelMessageHandle(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
	if msg, ok := update.GetMessage().(*tg.Message); ok {
		p.messagesCh <- msg
	}
	return nil
}

func (p *MessagePipe) EditMessageHandle(ctx context.Context, e tg.Entities, update *tg.UpdateEditMessage) error {
	if msg, ok := update.GetMessage().(*tg.Message); ok {
		p.messagesCh <- msg
	}
	return nil
}

func (p *MessagePipe) EditChannelMessageHandle(ctx context.Context, e tg.Entities, update *tg.UpdateEditChannelMessage) error {
	if msg, ok := update.GetMessage().(*tg.Message); ok {
		p.messagesCh <- msg
	}
	return nil
}

func (p *MessagePipe) Dispatcher() tg.UpdateDispatcher {
	return p.d
}

func (p *MessagePipe) Start(ctx context.Context) <-chan *tg.Message {
	go func() {
		<-ctx.Done()
		close(p.messagesCh)

	}()
	return p.messagesCh
}
