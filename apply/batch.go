package apply

import (
	"context"
	"time"

	"github.com/zhulik/pips"
)

type BatchConfig struct {
	FlushInterval time.Duration
}

func WithFlushInterval(interval time.Duration) BatchConfigurer {
	return func(config *BatchConfig) {
		config.FlushInterval = interval
	}
}

type BatchConfigurer func(*BatchConfig)

// Batch creates a batching stage.
func Batch(batchSize int, configurers ...BatchConfigurer) pips.Stage {
	config := &BatchConfig{}
	for _, c := range configurers {
		c(config)
	}

	return func(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
		buffer := make([]any, 0, batchSize)

		var flushTicker *time.Ticker
		var tickChan <-chan time.Time

		if config.FlushInterval > 0 {
			flushTicker = time.NewTicker(config.FlushInterval)
			defer flushTicker.Stop()
			tickChan = flushTicker.C
		} else {
			tickChan = make(chan time.Time)
		}

		sendReset := func() {
			if len(buffer) > 0 {
				output <- pips.AnyD(buffer)
				buffer = make([]any, 0, batchSize)
				if flushTicker != nil {
					flushTicker.Reset(config.FlushInterval)
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

				if len(buffer) >= batchSize {
					sendReset()
				}
			}
		}
	}
}
