package apply

import (
	"context"

	"github.com/zhulik/pips"
)

// Flatten creates a flattening stage.
func Flatten[T any]() pips.Stage {
	return func(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
		pips.MapToDChan(ctx, input, output, func(_ context.Context, item any, out chan<- pips.D[any]) error {
			for _, item := range item.([]T) {
				out <- pips.AnyD(item)
			}
			return nil
		})
	}
}
