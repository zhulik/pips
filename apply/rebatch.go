package apply

import (
	"context"

	"github.com/zhulik/pips"
)

// Rebatch creates a rebatching stage.
func Rebatch[T any](batchSize int) pips.Stage {
	return func(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
		buffer := make([]any, 0, batchSize)

		sendReset := func() {
			if len(buffer) > 0 {
				output <- pips.AnyD(buffer)
				buffer = make([]any, 0, batchSize)
			}
		}

		defer sendReset()

		pips.MapToDChan(ctx, input, output, func(_ context.Context, item any, _ chan<- pips.D[any]) error {
			for _, item := range item.([]T) {
				buffer = append(buffer, item)

				if len(buffer) >= batchSize {
					sendReset()
				}
			}
			return nil
		})
	}
}
