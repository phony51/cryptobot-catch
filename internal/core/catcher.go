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
