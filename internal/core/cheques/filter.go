package cheques

import (
	"context"
	"github.com/gotd/td/tg"
)

type Filter struct {
	detectStrategies []DetectStrategy
	messages         <-chan *tg.Message
	chequeIDs        chan<- string
}

func NewFilter(detectStrategies []DetectStrategy, messages <-chan *tg.Message, chequeIDs chan<- string) *Filter {
	return &Filter{
		detectStrategies,
		messages,
		chequeIDs,
	}
}

func (cf *Filter) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-cf.messages:
			for _, strategy := range cf.detectStrategies {
				go func() {
					if chequeID, ok := strategy.ChequeID(msg); ok {
						cf.chequeIDs <- chequeID
					}
				}()
			}
		}
	}
}
