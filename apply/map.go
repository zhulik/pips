package apply

import (
	"context"
	"reflect"

	"github.com/zhulik/pips"
)

type mapStage[I any, O any] struct {
	mapper func(context.Context, I) (O, error)
}

// Run maps items from the input channel and sends them to the output channel.
func (m mapStage[I, O]) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	pips.MapToDChan(ctx, input, output, func(ctx context.Context, item any, out chan<- pips.D[any]) error {
		var res O
		var err error

		if _, ok := item.([]any); ok {
			var x I
			res, err = m.mapper(ctx, convertSlice[I](item.([]any), reflect.TypeOf(x).Elem()))
		} else {
			res, err = m.mapper(ctx, item.(I))
		}
		if err != nil {
			return err
		}

		out <- pips.AnyD(res)

		return nil
	})
}

// Map creates a map stage.
func Map[I any, O any](mapper func(context.Context, I) (O, error)) pips.Stage {
	return mapStage[I, O]{mapper}
}

// convertSlice converts a slice of any type to a slice of the given type targetElementType.
// T must match targetElementType!
func convertSlice[T any](sourceSlice []any, targetElementType reflect.Type) T {
	targetSliceType := reflect.SliceOf(targetElementType)
	targetSlice := reflect.MakeSlice(targetSliceType, len(sourceSlice), len(sourceSlice))

	for i, item := range sourceSlice {
		sourceValue := reflect.ValueOf(item)

		if sourceValue.Type().AssignableTo(targetElementType) {
			targetSlice.Index(i).Set(sourceValue)
			continue
		}
	}

	return targetSlice.Interface().(T)
}
