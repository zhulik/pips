package apply

import (
	"context"
	"fmt"

	"github.com/zhulik/pips"
)

// filter is a function type that evaluates an item of type T and returns a boolean indicating
// whether the item should be kept in the pipeline, along with any error that occurred during evaluation.
// It's used by the Filter stage to determine which items should continue through the pipeline.
type filter[T any] func(context.Context, T) (bool, error)

// FilterStage represents a pipeline stage that filters items based on a predicate function.
// It passes only items that satisfy the filter condition to the next stage.
type FilterStage[T any] struct {
	filter filter[T]
}

// Run runs the stage.
// It processes input data by applying the filter function to each item and only passing
// items that satisfy the filter condition to the output channel.
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
// A filter stage evaluates each item in the pipeline using the provided filter function
// and only passes items that satisfy the filter condition (return true) to the next stage.
// Items that don't satisfy the condition are dropped from the pipeline.
// This is useful for removing unwanted items from the data stream based on some criteria.
func Filter[T any](filter filter[T]) pips.Stage {
	return FilterStage[T]{
		filter: filter,
	}
}
