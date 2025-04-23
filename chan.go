package pips

import (
	"context"
)

type OutChan[T any] <-chan D[T]

func (o OutChan[T]) Collect(ctx context.Context) ([]T, error) {
	var res []T

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case r, ok := <-o:
			if !ok {
				break loop
			}
			if r.Error() != nil {
				return nil, r.Error()
			}
			res = append(res, r.Value())
		}
	}

	return res, nil
}

func (o OutChan[T]) Wait(ctx context.Context) error {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case r, ok := <-o:
			if !ok {
				break loop
			}
			if r.Error() != nil {
				return r.Error()
			}
		}
	}
	return nil
}

func MapInputChan[I any, O any](ctx context.Context, ch <-chan I, f func(context.Context, I) (O, error)) <-chan D[O] {
	input := make(chan D[O])

	go func() {
		for r := range ch {
			input <- NewD(f(ctx, r))
		}
		close(input)
	}()

	return input
}

func CastDChan[I any, O any](ctx context.Context, input <-chan D[I]) <-chan D[O] {
	out := make(chan D[O])

	go func() {
		MapToDChan(ctx, input, out, func(_ context.Context, item I, out chan<- D[O]) error {
			out <- CastD[I, O](NewD(item))
			return nil
		})
		close(out)
	}()

	return out
}

type SubscriptionHandler[T any, O any] func(ctx context.Context, item T, out chan<- D[O]) error

func MapToDChan[I any, O any](ctx context.Context, input <-chan D[I], output chan<- D[O], h SubscriptionHandler[I, O]) {
	for {
		select {
		case <-ctx.Done():
			return

		case res, ok := <-input:
			if !ok {
				return
			}
			if res.Error() != nil {
				output <- ErrD[O](res.Error())
			}

			err := h(ctx, res.Value(), output)
			if err != nil {
				output <- ErrD[O](err)
				return
			}
		}
	}
}
