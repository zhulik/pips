package apply

import (
	"context"

	"github.com/zhulik/pips"
)

// EachC creates a concurrent each stage.
func EachC[I any](concurrency int, eacher eacher[I]) pips.Stage {
	return MapC(concurrency, func(_ context.Context, item I) (I, error) {
		return item, eacher(context.Background(), item)
	})
}
