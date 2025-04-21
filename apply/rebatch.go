package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type rebatchStage[T any] struct {
	size int
}

// Run rebatches items from the input channel into batches of the given size and sends them to the output channel.
func (s rebatchStage[T]) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	buffer := make([]any, 0, s.size)

	sendReset := func() {
		if len(buffer) > 0 {
			output <- pips.AnyD(buffer)
			buffer = make([]any, 0, s.size)
		}
	}

	defer sendReset()

	pips.MapToDChan(ctx, input, output, func(_ context.Context, item any, _ chan<- pips.D[any]) error {
		for _, item := range item.([]T) {
			buffer = append(buffer, item)

			if len(buffer) >= s.size {
				sendReset()
			}
		}
		return nil
	})
}

// Rebatch creates a rebatching stage.
func Rebatch[T any](batchSize int) pips.Stage {
	return rebatchStage[T]{batchSize}
}
