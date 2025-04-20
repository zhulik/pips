package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type filterStage[T any] struct {
	filter func(context.Context, T) (bool, error)
}

func (s filterStage[I]) Run(ctx context.Context, input <-chan pips.D[any]) <-chan pips.D[any] {
	return Subscriber(ctx, input, func(ctx context.Context, item pips.D[any], out chan<- pips.D[any]) error {
		keep, err := s.filter(ctx, item.Value().(I))
		if err != nil {
			return err
		}
		if !keep {
			return nil
		}

		out <- item
		return nil
	})
}

func Filter[I any](filter func(context.Context, I) (bool, error)) pips.Stage {
	return filterStage[I]{filter}
}
