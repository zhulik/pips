package apply

import (
	"context"
	"reflect"

	"github.com/zhulik/pips"
)

type mapStage[I any, O any] struct {
	mapper func(context.Context, I) (O, error)
}

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

func Map[I any, O any](mapper func(context.Context, I) (O, error)) pips.Stage {
	return mapStage[I, O]{mapper}
}

func convertSlice[T any](sourceSlice []any, targetElementType reflect.Type) T {
	// Create a new slice of the target element type
	targetSliceType := reflect.SliceOf(targetElementType)
	targetSlice := reflect.MakeSlice(targetSliceType, len(sourceSlice), len(sourceSlice))

	// Convert each element and set it in the target slice
	for i, item := range sourceSlice {
		sourceValue := reflect.ValueOf(item)

		// If source value is directly assignable to target type
		if sourceValue.Type().AssignableTo(targetElementType) {
			targetSlice.Index(i).Set(sourceValue)
			continue
		}

		// If source value is convertible to target type
		if sourceValue.Type().ConvertibleTo(targetElementType) {
			convertedValue := sourceValue.Convert(targetElementType)
			targetSlice.Index(i).Set(convertedValue)
			continue
		}

		panic("error")
	}

	// Return the target slice as interface{}
	return targetSlice.Interface().(T)
}
