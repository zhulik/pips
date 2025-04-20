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
