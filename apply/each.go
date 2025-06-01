package apply

import (
	"context"

	"github.com/zhulik/pips"
)

// eacher is a function type that processes an item of type T and returns an error if the processing fails.
// It's used by the Each stage to perform operations on each item in the pipeline without changing the item itself.
type eacher[T any] func(context.Context, T) error

// Each creates an each stage.
// An each stage applies the provided eacher function to each item in the pipeline,
// but passes the original item unchanged to the next stage.
// This is useful for side effects like logging, metrics collection, or other operations
// where you want to process each item without transforming it.
func Each[I any](eacher eacher[I]) pips.Stage {
	return Map(func(ctx context.Context, item I) (I, error) {
		return item, eacher(ctx, item)
	})
}
