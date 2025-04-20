package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type SubscriptionHandler[T any] func(ctx context.Context, item T, out chan<- pips.D[any]) error

func Subscriber[T any](ctx context.Context, input <-chan T, h SubscriptionHandler[T]) <-chan pips.D[any] {
	out := make(chan pips.D[any])

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return

			case res, ok := <-input:
				if !ok {
					return
				}

				err := h(ctx, res, out)
				if err != nil {
					out <- pips.ErrD[any](err)
					return
				}
			}
		}
	}()

	return out
}
