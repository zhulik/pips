package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type flattenStage[T any] struct {
}

func (f flattenStage[T]) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	pips.MapToDChan(ctx, input, output, func(_ context.Context, item any, out chan<- pips.D[any]) error {
		for _, item := range item.([]T) {
			out <- pips.AnyD(item)
		}
		return nil
	})
}

func Flatten[T any]() pips.Stage {
	return flattenStage[T]{}
}
