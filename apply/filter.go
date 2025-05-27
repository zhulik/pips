package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type filter[T any] func(context.Context, T) (bool, error)

// Filter creates a filter stage.
func Filter[T any](filter filter[T]) pips.Stage {
	return func(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
		pips.MapToDChan(ctx, input, output, func(ctx context.Context, item any, out chan<- pips.D[any]) error {
			i, err := pips.TryCast[T](item)
			if err != nil {
				return err
			}

			keep, err := filter(ctx, i)
			if err != nil {
				return err
			}
			if keep {
				out <- pips.NewD(item)
			}

			return nil
		})
	}
}
