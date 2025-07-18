package cheques

import (
	"context"
	"cryptobot-catch/pkg/cryptobot"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

type Activator struct {
	cryptoBot     *cryptobot.CryptoBot
	maxActivators int
}

func NewActivator(cryptoBot *cryptobot.CryptoBot, maxActivators int) *Activator {
	return &Activator{
		cryptoBot,
		maxActivators,
	}
}

func (a *Activator) Run(ctx context.Context, chequeIDsCh <-chan string) error {
	sem := semaphore.NewWeighted(3)
	logger := zap.L()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case chequeID := <-chequeIDsCh:
			if err := sem.Acquire(ctx, 1); err != nil {
				return err
			}
			go func() {
				defer sem.Release(1)
				_ = a.cryptoBot.ActivateCheque(ctx, chequeID)
				logger.Info("cheque activated", zap.String("chequeID", chequeID))
			}()
		}
	}
}
