package apply

import (
	"context"
	"reflect"

	"github.com/zhulik/pips"
)

type MapConfig[I any, O any] struct {
	Mapper MapperFn[I, O]
}

// MapperFn is a function type that transforms an item of type I into an item of type O,
// potentially returning an error if the transformation fails.
// It's used by the Map stage to transform items in the pipeline.
type MapperFn[I any, O any] func(context.Context, I) (O, error)

// MapperStage represents a pipeline stage that transforms items using a MapperFn function.
// It takes items of type I as input and emits items of type O as output.
type MapperStage[I any, O any] struct {
	config MapConfig[I, O]
}

// Run runs the stage.
// It processes input data by applying the MapperFn function to each item
// and sending the transformed items to the output channel.
func (m MapperStage[I, O]) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	pips.MapToDChan(ctx, input, output, func(ctx context.Context, item any, out chan<- pips.D[any]) error {
		out <- pips.AnyD(mapItemOrSlice(ctx, item, m.config.Mapper))

		return nil
	})
}

// Map creates a map stage.
// A map stage transforms each item in the pipeline using the provided MapperFn function.
// It takes items of type I as input and emits items of type O as output.
// This is the most common way to transform data in a pipeline, allowing you to convert
// from one data type to another or modify the content of items.
func Map[I any, O any](mapper MapperFn[I, O]) pips.Stage {
	return MapperStage[I, O]{
		config: MapConfig[I, O]{
			Mapper: mapper,
		},
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

// iterateAnySlice iterates over a slice of any type and calls the provided function for each element.
// It uses reflection to handle slices of any type.
// Panics if the provided item is not a slice.
func iterateAnySlice(item any, fn func(any)) {
	if reflect.TypeOf(item).Kind() == reflect.Slice {
		s := reflect.ValueOf(item)
		for i := 0; i < s.Len(); i++ {
			fn(s.Index(i).Interface())
		}
		return
	}

	panic("not a slice")
}

// mapItemOrSlice applies the MapperFn function to an item or to each element of a slice.
// If the item is a slice, it converts the slice to the expected input type before applying the MapperFn.
// If the item is not a slice, it attempts to cast it to the expected input type before applying the MapperFn.
func mapItemOrSlice[I any, O any](ctx context.Context, item any, mapper MapperFn[I, O]) (O, error) {
	var i I

	if anys, ok := item.([]any); ok {
		return mapper(ctx, convertSlice[I](anys, reflect.TypeOf(i).Elem()))
	}

	i, err := pips.TryCast[I](item)
	if err != nil {
		var o O
		return o, err
	}

	return mapper(ctx, i)
}
