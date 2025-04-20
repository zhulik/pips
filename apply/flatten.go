package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type flattenStage struct {
}

func (f flattenStage) Run(ctx context.Context, input <-chan pips.D[any]) <-chan pips.D[any] {
	return pips.MapDChan(ctx, input, func(_ context.Context, item any, out chan<- pips.D[any]) error {
		for _, item := range item.([]any) {
			out <- pips.NewD(item)
		}
		return nil
	})
}

func Flatten() pips.Stage {
	return flattenStage{}
}
