package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type MapCStage[I any, O any] struct {
	concurrency int
	mapper      mapper[I, O]
}

// Run runs the stage.
func (m MapCStage[I, O]) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	semaphore := make(chan any, m.concurrency)
	midChan := make(chan pips.D[chan pips.D[any]], m.concurrency)

	go func() {
		defer close(semaphore)
		defer close(midChan)

		pips.MapToDChan(ctx, input, midChan, func(ctx context.Context, item any, out chan<- pips.D[chan pips.D[any]]) error {
			semaphore <- true

			ch := make(chan pips.D[any])

			go func() {
				defer func() { <-semaphore }()
				defer close(ch)
				defer pips.RecoverPanicAndSendToPipeline(ch)

				ch <- pips.AnyD(mapItemOrSlice(ctx, item, m.mapper))
			}()

			out <- pips.NewD(ch)

			return nil
		})
	}()
	pips.MapToDChan(ctx, midChan, output, func(ctx context.Context, c chan pips.D[any], out chan<- pips.D[any]) error {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case d, ok := <-c:
			if !ok {
				return nil
			}
			out <- d
			return nil
		}
	})
}

// MapC creates a concurrent map stage.
func MapC[I any, O any](concurrency int, mapper mapper[I, O]) pips.Stage {
	return MapCStage[I, O]{
		mapper:      mapper,
		concurrency: concurrency,
	}
}
