package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type FlattenStage struct{}

// Run runs the stage.
func (f FlattenStage) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	pips.MapToDChan(ctx, input, output, func(_ context.Context, item any, out chan<- pips.D[any]) error {
		iterateAnySlice(item, func(item any) {
			out <- pips.AnyD(item)
		})
		return nil
	})
}

// Flatten creates a flattening stage.
func Flatten() pips.Stage {
	return FlattenStage{}
}
