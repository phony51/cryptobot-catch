package cheques

import (
	"context"
	"github.com/gotd/td/tg"
)

type Filter struct {
	Strategies []DetectStrategy
	messages   <-chan *tg.Message
	chequeIDs  chan<- string
}

func NewFilter(strategies []DetectStrategy, messages <-chan *tg.Message, chequeIDs chan<- string) *Filter {
	return &Filter{
		strategies,
		messages,
		chequeIDs,
	}
}

func (cf *Filter) Run(ctx context.Context) error {
	//var once sync.Once
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-cf.messages:
			for _, strategy := range cf.Strategies {
				go func() {
					if chequeID, ok := strategy.ChequeID(msg); ok {
						cf.chequeIDs <- chequeID
					}
				}()
			}
		}
	}
}
