package apply

import (
	"context"
	"reflect"

	"github.com/zhulik/pips"
)

type eachStage[I any] struct {
	eacher func(context.Context, I) error
}

// Run maps items from the input channel and sends them to the output channel.
func (m eachStage[I]) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	pips.MapToDChan(ctx, input, output, func(ctx context.Context, item any, out chan<- pips.D[any]) error {
		var err error

		if _, ok := item.([]any); ok {
			var x I
			err = m.eacher(ctx, convertSlice[I](item.([]any), reflect.TypeOf(x).Elem()))
		} else {
			err = m.eacher(ctx, item.(I))
		}
		if err != nil {
			return err
		}

		out <- pips.NewD(item)

		return nil
	})
}

// Each creates an each stage.
func Each[I any](mapper func(context.Context, I) error) pips.Stage {
	return eachStage[I]{mapper}
}
