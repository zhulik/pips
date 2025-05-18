package apply

import (
	"context"
	"time"

	"github.com/zhulik/pips"
)

// Rebatch creates a rebatching stage.
func Rebatch[T any](batchSize int, configurer ...BatchConfigurer) pips.Stage {
	config := &BatchConfig{}
	for _, c := range configurer {
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

				item := res.Value()

				if anyItems, ok := item.([]any); ok {
					for _, item := range anyItems {
						buffer = append(buffer, item.(T))

						if len(buffer) >= batchSize {
							sendReset()
						}
					}
					continue
				}

				for _, item := range item.([]T) {
					buffer = append(buffer, item)

					if len(buffer) >= batchSize {
						sendReset()
					}
				}
			}
		}
	}
}
