package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type flattenStage struct {
}

func (f flattenStage) Run(ctx context.Context, input <-chan pips.D[any]) <-chan pips.D[any] {
	return Subscriber(ctx, input, func(_ context.Context, item pips.D[any], out chan<- pips.D[any]) error {
		ary, err := item.Unpack()
		if err != nil {
			return err
		}
		for _, item := range ary.([]any) {
			out <- pips.NewD(item)
		}
		return nil
	})
}

func Flatten() pips.Stage {
	return flattenStage{}
}
