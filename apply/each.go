package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type eacher[T any] func(context.Context, T) error

// Each creates an each stage.
func Each[I any](eacher eacher[I]) pips.Stage {
	return Map(func(ctx context.Context, item I) (I, error) {
		return item, eacher(ctx, item)
	})
}
