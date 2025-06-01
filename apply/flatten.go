package apply

import (
	"context"

	"github.com/zhulik/pips"
)

// FlattenStage represents a pipeline stage that flattens slices or arrays into individual items.
// It takes slices or arrays as input and emits each element as a separate item in the pipeline.
type FlattenStage Stage[any]

// Run runs the stage.
// It processes input data by taking slices or arrays and emitting each element
// as a separate item in the output channel.
func (f FlattenStage) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	pips.MapToDChan(ctx, input, output, func(_ context.Context, item any, out chan<- pips.D[any]) error {
		iterateAnySlice(item, func(item any) {
			out <- pips.AnyD(item)
		})
		return nil
	})
}

// Flatten creates a flattening stage.
// A flattening stage takes slices or arrays as input and emits each element as a separate item in the pipeline.
// This is useful when you have a pipeline that produces collections of items, but you want to process
// each item individually in subsequent stages.
func Flatten() pips.Stage {
	return FlattenStage{}
}
