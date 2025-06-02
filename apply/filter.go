package apply

import (
	"context"
	"fmt"

	"github.com/zhulik/pips"
)

type FilterConfig[T any] struct {
	Filter FilterFn[T]
}

// FilterFn is a function type that evaluates an item of type T and returns a boolean indicating
// whether the item should be kept in the pipeline, along with any error that occurred during evaluation.
// It's used by the Filter stage to determine which items should continue through the pipeline.
type FilterFn[T any] func(context.Context, T) (bool, error)

// FilterStage represents a pipeline stage that filters items based on a predicate function.
// It passes only items that satisfy the FilterFn condition to the next stage.
type FilterStage[T any] Stage[FilterConfig[T]]

// Run runs the stage.
// It processes input data by applying the FilterFn function to each item and only passing
// items that satisfy the FilterFn condition to the output channel.
func (f FilterStage[T]) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	pips.MapToDChan(ctx, input, output, func(ctx context.Context, item any, out chan<- pips.D[any]) error {
		i, err := pips.TryCast[T](item)
		if err != nil {
			return fmt.Errorf("failed to cast filter stage input: %w", err)
		}

		keep, err := f.config.Filter(ctx, i)
		if err != nil {
			return err
		}
		if keep {
			out <- pips.NewD(item)
		}

		return nil
	})
}

// Filter creates a FilterFn stage.
// A FilterFn stage evaluates each item in the pipeline using the provided FilterFn function
// and only passes items that satisfy the FilterFn condition (return true) to the next stage.
// Items that don't satisfy the condition are dropped from the pipeline.
// This is useful for removing unwanted items from the data stream based on some criteria.
func Filter[T any](filter FilterFn[T]) pips.Stage {
	return FilterStage[T]{
		config: FilterConfig[T]{
			Filter: filter,
		},
		source: getStageSource(),
	}
}
