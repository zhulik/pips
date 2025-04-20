package pips

import (
	"context"
)

func CastDChan[I any, O any](ctx context.Context, ch <-chan D[I]) <-chan D[O] {
	return MapChan(ctx, ch, CastD[I, O])
}

func MapChan[I any, O any](ctx context.Context, input <-chan I, f func(I) O) <-chan O {
	out := make(chan O)

	go func() {
		MapToChan(ctx, input, out, f)
		close(out)
	}()

	return out
}

// MapToChan maps the input channel to the output channel using the given function. Does not close any channels.
// Blocks.
func MapToChan[I any, O any](ctx context.Context, input <-chan I, output chan<- O, f func(I) O) {
	for {
		select {
		case <-ctx.Done():
			return

		case res, ok := <-input:
			if !ok {
				return
			}
			output <- f(res)
		}
	}
}

type SubscriptionHandler[T any, O any] func(ctx context.Context, item T, out chan<- D[O]) error

func MapDChan[I any, O any](ctx context.Context, input <-chan D[I], h SubscriptionHandler[I, O]) <-chan D[O] {
	out := make(chan D[O])

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
				if res.Error() != nil {
					out <- ErrD[O](res.Error())
					return
				}

				err := h(ctx, res.Value(), out)
				if err != nil {
					out <- ErrD[O](err)
					return
				}
			}
		}
	}()

	return out
}
