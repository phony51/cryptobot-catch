package extracting

import (
	"context"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
	"sync"
)

type ExtractorWorkerPool struct {
	extractor  Extractor
	wg         sync.WaitGroup
	numWorkers int
}

func (wp *ExtractorWorkerPool) Start(ctx context.Context, messagesCh <-chan *tg.Message) <-chan string {
	chequeIDs := make(chan string)
	logger := zap.L()
	for range wp.numWorkers {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-messagesCh:
					if !ok {
						return
					}
					if id, found := wp.extractor.Extract(msg); found {
						chequeIDs <- id
						logger.Info("cheque caught", zap.String("chequeID", id), zap.String("by", wp.extractor.Name))
					}
				}
			}
		}()
	}

	go func() {
		wp.wg.Wait()
		close(chequeIDs)
	}()
	return chequeIDs
}

func NewExtractorWorkerPool(extractor Extractor, numWorkers int) ExtractorWorkerPool {
	return ExtractorWorkerPool{
		extractor,
		sync.WaitGroup{},
		numWorkers,
	}
}
