package cheques

import (
	"context"
	"cryptobot-catch/internal/core/cheques/detecting"
	"fmt"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

type Filter struct {
	detectStrategies []detecting.DetectStrategy
	messages         <-chan *tg.Message
	chequeIDs        chan<- string
}

func NewFilter(detectStrategies []detecting.DetectStrategy, messages <-chan *tg.Message, chequeIDs chan<- string) *Filter {
	return &Filter{
		detectStrategies,
		messages,
		chequeIDs,
	}
}

func (cf *Filter) Run(ctx context.Context) error {
	var mStrategy detecting.MappedDetectStrategy

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-cf.messages:
			for _, strategy := range cf.detectStrategies {
				if chequeID, ok := strategy.ChequeID(msg); ok {
					cf.chequeIDs <- chequeID
					mStrategy = strategy.(detecting.MappedDetectStrategy)
					zap.L().Info("Cheque caught",
						zap.String("chequeID", chequeID),
						zap.String("strategy", fmt.Sprint(mStrategy.Alias())),
					)
				}

			}
		}
	}
}
