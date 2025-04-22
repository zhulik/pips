package apply

import (
	"context"

	"github.com/zhulik/pips"
)

// Pipeline creates a pipeline stage.
func Pipeline[I any, O any](pipeline *pips.Pipeline[I, O]) pips.Stage {
	return func(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
		pips.MapToDChan(
			ctx,
			pipeline.Run(ctx, pips.CastDChan[any, I](ctx, input)),
			output,
			func(_ context.Context, item O, out chan<- pips.D[any]) error {
				out <- pips.AnyD(item)
				return nil
			},
		)
	}
}
