package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type filterStage[T any] struct {
	filter func(context.Context, T) (bool, error)
}

// Run filters items from the input channel and sends them to the output channel.
func (s filterStage[I]) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	pips.MapToDChan(ctx, input, output, func(ctx context.Context, item any, out chan<- pips.D[any]) error {
		keep, err := s.filter(ctx, item.(I))
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
func Filter[I any](filter func(context.Context, I) (bool, error)) pips.Stage {
	return filterStage[I]{filter}
}
