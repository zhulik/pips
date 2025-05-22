package apply

import (
	"context"
	"reflect"

	"github.com/zhulik/pips"
)

type mapper[I any, O any] func(context.Context, I) (O, error)

// Map creates a map stage.
func Map[I any, O any](mapper mapper[I, O]) pips.Stage {
	return func(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
		pips.MapToDChan(ctx, input, output, func(ctx context.Context, item any, out chan<- pips.D[any]) error {
			var res O
			var err error

			if anys, ok := item.([]any); ok {
				var x I
				res, err = mapper(ctx, convertSlice[I](anys, reflect.TypeOf(x).Elem()))
			} else {
				var i I
				i, err = pips.TryCast[I](item)
				if err == nil {
					res, err = mapper(ctx, i)
				}
			}
			if err != nil {
				return err
			}

			out <- pips.AnyD(res)

			return nil
		})
	}
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
