package apply

import (
	"context"

	"github.com/zhulik/pips"
)

// MapCStage represents a pipeline stage that transforms items using a mapper function with concurrency.
// It takes items of type I as input and emits items of type O as output,
// processing multiple items simultaneously up to the specified concurrency limit.
type MapCStage[I any, O any] struct {
	concurrency int
	mapper      mapper[I, O]
}

// Run runs the stage.
// It processes input data by applying the mapper function to each item concurrently,
// up to the specified concurrency limit, and sending the transformed items to the output channel.
// This implementation uses a semaphore pattern to limit concurrency and goroutines to process items in parallel.
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
// A concurrent map stage transforms each item in the pipeline using the provided mapper function,
// processing multiple items simultaneously up to the specified concurrency limit.
// The concurrency parameter determines the maximum number of items that can be processed in parallel.
// This is useful for performance-intensive transformations or operations that may block,
// such as network requests or file I/O, allowing the pipeline to continue processing other items.
func MapC[I any, O any](concurrency int, mapper mapper[I, O]) pips.Stage {
	return MapCStage[I, O]{
		mapper:      mapper,
		concurrency: concurrency,
	}
}
