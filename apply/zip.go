package apply

import (
	"context"
	"reflect"

	"github.com/zhulik/pips"
)

// Zip creates a zip stage.
func Zip[I any, O any](zipper mapper[I, O]) pips.Stage {
	return func(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
		pips.MapToDChan(ctx, input, output, func(ctx context.Context, item any, out chan<- pips.D[any]) error {
			var res O
			var err error

			castItem, err := pips.TryCast[I](item)
			if err != nil {
				return err
			}

			if anys, ok := item.([]any); ok {
				var x I
				res, err = zipper(ctx, convertSlice[I](anys, reflect.TypeOf(x).Elem()))
			} else {
				res, err = zipper(ctx, castItem)
			}
			if err != nil {
				return err
			}

			out <- pips.AnyD(pips.NewP(castItem, res))

			return nil
		})
	}
}
