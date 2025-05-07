package apply

import (
	"context"
	"reflect"

	"github.com/zhulik/pips"
)

type eacher[T any] func(context.Context, T) error

// Each creates an each stage.
func Each[I any](eacher eacher[I]) pips.Stage {
	return func(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
		pips.MapToDChan(ctx, input, output, func(ctx context.Context, item any, out chan<- pips.D[any]) error {
			var err error

			if anys, ok := item.([]any); ok {
				var x I
				err = eacher(ctx, convertSlice[I](anys, reflect.TypeOf(x).Elem()))
			} else {
				err = eacher(ctx, item.(I))
			}
			if err != nil {
				return err
			}

			out <- pips.NewD(item)

			return nil
		})
	}
}
