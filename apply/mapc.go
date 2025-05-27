package apply

import (
	"context"
	"fmt"
	"reflect"

	"github.com/zhulik/pips"
)

// MapC creates a concurrent map stage.
func MapC[I any, O any](concurrency int, mapper mapper[I, O]) pips.Stage {
	return func(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
		semaphore := make(chan any, concurrency)
		midChan := make(chan pips.D[chan pips.D[any]], concurrency)

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

					if anys, ok := item.([]any); ok {
						var x I
						ch <- pips.AnyD(mapper(ctx, convertSlice[I](anys, reflect.TypeOf(x).Elem())))
					} else {
						i, err := pips.TryCast[I](item)
						if err != nil {
							err = fmt.Errorf("failed to cast mapc stage input: %w", err)
							ch <- pips.ErrD[any](err)
							return
						}

						ch <- pips.AnyD(mapper(ctx, i))
					}
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
				v, err := d.Unpack()
				if err != nil {
					return err
				}

				out <- pips.NewD(v)
				return nil
			}
		})
	}
}
