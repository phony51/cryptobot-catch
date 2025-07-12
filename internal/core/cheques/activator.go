package cheques

import (
	"context"
	"cryptobot-catch/pkg/cryptobot"
	"go.uber.org/zap"
)

type Activator struct {
	cryptoBot  *cryptobot.CryptoBot
	chequesIDs <-chan string
}

func NewActivator(cryptoBot *cryptobot.CryptoBot, chequesIDs <-chan string) *Activator {
	return &Activator{
		cryptoBot:  cryptoBot,
		chequesIDs: chequesIDs,
	}
}

func (a *Activator) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case chequeID := <-a.chequesIDs:
			zap.L().Info(chequeID)
			go func() { _ = a.cryptoBot.ActivateCheque(ctx, chequeID) }()
		}
	}
}
