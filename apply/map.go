package apply

import (
	"context"
	"reflect"

	"github.com/zhulik/pips"
)

type mapper[I any, O any] func(context.Context, I) (O, error)

type MapperStage[I any, O any] struct {
	mapper mapper[I, O]
}

// Run runs the stage.
func (m MapperStage[I, O]) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	pips.MapToDChan(ctx, input, output, func(ctx context.Context, item any, out chan<- pips.D[any]) error {
		out <- pips.AnyD(mapItemOrSlice(ctx, item, m.mapper))

		return nil
	})
}

// Map creates a map stage.
func Map[I any, O any](mapper mapper[I, O]) pips.Stage {
	return MapperStage[I, O]{
		mapper: mapper,
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

func mapItemOrSlice[I any, O any](ctx context.Context, item any, mapper mapper[I, O]) (O, error) {
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
