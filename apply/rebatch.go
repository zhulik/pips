package apply

import (
	"context"
	"time"

	"github.com/zhulik/pips"
)

type RebatchStage struct {
	batchSize int
	config    *BatchConfig
}

// Run runs the stage.
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
