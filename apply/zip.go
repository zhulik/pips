package apply

import (
	"context"

	"github.com/zhulik/pips"
)

// Zip creates a zip stage.
// A zip stage transforms each item in the pipeline using the provided zipper function,
// but instead of replacing the original item, it pairs the original item with the transformed item.
// It takes items of type I as input and emits pairs (pips.P) of the original item and the transformed item.
// This is useful when you need to keep the original item along with its transformed version
// for later processing or reference.
func Zip[I any, O any](zipper MapperFn[I, O]) pips.Stage {
	return Map(func(ctx context.Context, i I) (pips.P[I, O], error) {
		o, err := zipper(ctx, i)
		if err != nil {
			return nil, err
		}
		return pips.NewP(i, o), nil
	})
}
