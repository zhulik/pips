package apply

import (
	"context"
	"fmt"

	"github.com/zhulik/pips"
)

type filter[T any] func(context.Context, T) (bool, error)

type FilterStage[T any] struct {
	filter filter[T]
}

// Run runs the stage.
func (f FilterStage[T]) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	pips.MapToDChan(ctx, input, output, func(ctx context.Context, item any, out chan<- pips.D[any]) error {
		i, err := pips.TryCast[T](item)
		if err != nil {
			return fmt.Errorf("failed to cast filter stage input: %w", err)
		}

		keep, err := f.filter(ctx, i)
		if err != nil {
			return err
		}
		if keep {
			out <- pips.NewD(item)
		}

		return nil
	})
}

// Filter creates a filter stage.
func Filter[T any](filter filter[T]) pips.Stage {
	return FilterStage[T]{
		filter: filter,
	}
}
