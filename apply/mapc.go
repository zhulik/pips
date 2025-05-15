package apply

import (
	"context"

	"github.com/zhulik/pips"
)

// MapC creates a concurrent map stage.
func MapC[I any, O any](concurrency int, mapper mapper[I, O]) pips.Stage {
	semaphore := make(chan any, concurrency)

	return pips.New[I, O]().
		Then(
			Map(func(ctx context.Context, i I) (<-chan pips.D[O], error) {
				ch := make(chan pips.D[O], 1)

				go func() {
					semaphore <- nil
					defer func() { <-semaphore }()

					ch <- pips.NewD(mapper(ctx, i))
				}()

				return ch, nil
			}),
		).
		Then(
			Map(func(ctx context.Context, f <-chan pips.D[O]) (O, error) {
				select {
				case <-ctx.Done():
					var o O
					return o, ctx.Err()
				case d := <-f:
					return d.Unpack()
				}
			}),
		).ToStage()
}
