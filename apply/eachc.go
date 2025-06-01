package apply

import (
	"context"

	"github.com/zhulik/pips"
)

// EachC creates a concurrent each stage.
// A concurrent each stage applies the provided eacher function to each item in the pipeline
// with a specified level of concurrency, but passes the original item unchanged to the next stage.
// The concurrency parameter determines the maximum number of items that can be processed simultaneously.
// This is useful for side effects like logging, metrics collection, or other operations
// where you want to process each item without transforming it, but need to do so concurrently
// for performance reasons, especially for operations that may block or take significant time.
func EachC[I any](concurrency int, eacher eacher[I]) pips.Stage {
	return MapC(concurrency, func(ctx context.Context, item I) (I, error) {
		return item, eacher(ctx, item)
	})
}
