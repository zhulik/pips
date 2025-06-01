package apply

import (
	"context"
	"time"

	"github.com/zhulik/pips"
)

// BatchConfig represents configuration for Batch and Rebatch stages.
// It contains settings that control the behavior of batching operations.
type BatchConfig struct {
	// FlushInterval defines the maximum time to wait before flushing a batch,
	// even if the batch size hasn't been reached.
	FlushInterval time.Duration
}

// WithFlushInterval sets the flush interval for Batch and Rebatch stages using the provided duration.
// This configurer allows batches to be emitted after a specified time period,
// even if they haven't reached the maximum batch size.
func WithFlushInterval(interval time.Duration) BatchConfigurer {
	return func(config *BatchConfig) {
		config.FlushInterval = interval
	}
}

// BatchConfigurer defines a function type used to configure a BatchConfig instance.
// It's used to provide a functional options pattern for configuring batch stages.
type BatchConfigurer func(*BatchConfig)

// BatchStage represents a pipeline stage that buffers data into batches based on size or flush interval.
// It collects individual items into batches and emits them when either the batch size is reached
// or the flush interval has elapsed.
type BatchStage struct {
	batchSize int
	config    *BatchConfig
}

// Run runs the stage.
// It processes input data by collecting items into batches and sending them to the output channel
// when either the batch size is reached or the flush interval has elapsed.
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
// It collects individual items into batches of the specified size and emits them as a single item.
// The batchSize parameter determines the maximum number of items in each batch.
// Optional configurers can be provided to customize the batching behavior, such as setting a flush interval.
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
