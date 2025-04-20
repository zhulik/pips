package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type flattenStage[T any] struct {
}

// Run flattens items from the input channel and sends them to the output channel.
func (f flattenStage[T]) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	pips.MapToDChan(ctx, input, output, func(_ context.Context, item any, out chan<- pips.D[any]) error {
		for _, item := range item.([]T) {
			out <- pips.AnyD(item)
		}
		return nil
	})
}

// Flatten creates a flattening stage.
func Flatten[T any]() pips.Stage {
	return flattenStage[T]{}
}
