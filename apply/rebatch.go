package apply

import (
	"context"
	"time"

	"github.com/zhulik/pips"
)

// RebatchStage represents a pipeline stage that rebatches data from existing batches into new batches.
// It takes batches (slices) as input, flattens them, and then creates new batches of the specified size.
// This is useful when you need to change the batch size in the middle of a pipeline.
type RebatchStage struct {
	batchSize int
	config    *BatchConfig
}

// Run runs the stage.
// It processes input data by taking batches (slices), flattening them into individual items,
// and then collecting those items into new batches of the specified size.
// Like BatchStage, it also supports flushing batches based on a time interval.
func (r RebatchStage) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	buffer := make([]any, 0, r.batchSize)

	var flushTicker *time.Ticker
	var tickChan <-chan time.Time

	if r.config.FlushInterval > 0 {
		flushTicker = time.NewTicker(r.config.FlushInterval)
		defer flushTicker.Stop()
		tickChan = flushTicker.C
	} else {
		tickChan = make(chan time.Time)
	}

	sendReset := func() {
		if len(buffer) > 0 {
			output <- pips.AnyD(buffer)
			buffer = make([]any, 0, r.batchSize)
			if flushTicker != nil {
				flushTicker.Reset(r.config.FlushInterval)
			}
		}
	}

	defer sendReset()

	for {
		select {
		case <-ctx.Done():
			return

		case <-tickChan:
			sendReset()
			continue

		case res, ok := <-input:
			if !ok {
				return
			}
			if res.Error() != nil {
				output <- pips.ErrD[any](res.Error())
				return
			}

			iterateAnySlice(res.Value(), func(item any) {
				buffer = append(buffer, item)

				if len(buffer) >= r.batchSize {
					sendReset()
				}
			})
		}
	}
}

// Rebatch creates a rebatching stage.
// A rebatching stage takes batches (slices) as input, flattens them into individual items,
// and then collects those items into new batches of the specified size.
// The batchSize parameter determines the maximum number of items in each new batch.
// Optional configurers can be provided to customize the rebatching behavior, such as setting a flush interval.
// This is useful when you need to change the batch size in the middle of a pipeline,
// or when you need to process batches from an upstream source but with a different batch size.
func Rebatch(batchSize int, configurers ...BatchConfigurer) pips.Stage {
	config := &BatchConfig{}
	for _, c := range configurers {
		c(config)
	}

	return RebatchStage{
		batchSize: batchSize,
		config:    config,
	}
}
