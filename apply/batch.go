package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type batchStage struct {
	size int
}

func (s batchStage) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	buffer := make([]any, 0, s.size)

	sendReset := func() {
		if len(buffer) > 0 {
			output <- pips.AnyD(buffer)
			buffer = make([]any, 0, s.size)
		}
	}

	defer sendReset()

	pips.MapToDChan(ctx, input, output, func(_ context.Context, item any, _ chan<- pips.D[any]) error {
		buffer = append(buffer, item)

		if len(buffer) >= s.size {
			sendReset()
		}
		return nil
	})
}

func Batch(batchSize int) pips.Stage {
	return batchStage{batchSize}
}
