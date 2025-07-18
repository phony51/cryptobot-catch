package core

import (
	"context"
	"cryptobot-catch/internal/core/cheques"
	"cryptobot-catch/internal/core/cheques/extracting"
	"cryptobot-catch/internal/utils"
)

type ExtractorWorkerPoolConfig struct {
	Extractor  extracting.Extractor
	NumWorkers int
}

type Catcher struct {
	messagePipe          MessagePipe
	extractorWorkerPools []extracting.ExtractorWorkerPool
	activator            *cheques.Activator
}

func (c *Catcher) Run(ctx context.Context) error {
	messagesCh := c.messagePipe.Start(ctx)
	numWorkerPools := len(c.extractorWorkerPools)
	chequeIDsChs := make([]<-chan string, numWorkerPools)
	messagesChs := utils.FanOut(messagesCh, numWorkerPools)

	for i := 0; i < numWorkerPools; i++ {
		chequeIDsChs[i] = c.extractorWorkerPools[i].Start(ctx, messagesChs[i])
	}

	chequeIDsCh := utils.FanIn(chequeIDsChs...)
	return c.activator.Run(ctx, chequeIDsCh)
}

func NewCatcher(messagePipe MessagePipe, extractorsWorkerPools []extracting.ExtractorWorkerPool, activator *cheques.Activator) *Catcher {
	return &Catcher{
		messagePipe,
		extractorsWorkerPools,
		activator,
	}
}
