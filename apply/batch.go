package apply

import (
	"context"
	"time"

	"github.com/zhulik/pips"
)

// BatchConfig represents configuration for Batch and Rebatch stages
type BatchConfig struct {
	FlushInterval time.Duration
}

// WithFlushInterval sets the flush interval for Batch and Rebatch stages using the provided duration.
func WithFlushInterval(interval time.Duration) BatchConfigurer {
	return func(config *BatchConfig) {
		config.FlushInterval = interval
	}
}

// BatchConfigurer defines a function type used to configure a BatchConfig instance.
type BatchConfigurer func(*BatchConfig)

// BatchStage represents a pipeline stage that buffers data into batches based on size or flush interval.
type BatchStage struct {
	batchSize int
	config    *BatchConfig
}

// Run runs the stage.
func (b BatchStage) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	buffer := make([]any, 0, b.batchSize)

	var flushTicker *time.Ticker
	var tickChan <-chan time.Time

	if b.config.FlushInterval > 0 {
		flushTicker = time.NewTicker(b.config.FlushInterval)
		defer flushTicker.Stop()
		tickChan = flushTicker.C
	} else {
		tickChan = make(chan time.Time)
	}

	sendReset := func() {
		if len(buffer) > 0 {
			output <- pips.AnyD(buffer)
			buffer = make([]any, 0, b.batchSize)
			if flushTicker != nil {
				flushTicker.Reset(b.config.FlushInterval)
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

			buffer = append(buffer, res.Value())

			if len(buffer) >= b.batchSize {
				sendReset()
			}
		}
	}
}

// Batch creates a batching stage.
func Batch(batchSize int, configurers ...BatchConfigurer) pips.Stage {
	config := &BatchConfig{}
	for _, c := range configurers {
		c(config)
	}

	return BatchStage{
		batchSize: batchSize,
		config:    config,
	}
}
