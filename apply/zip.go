package apply

import (
	"context"

	"github.com/zhulik/pips"
)

// Zip creates a zip stage.
func Zip[I any, O any](zipper mapper[I, O]) pips.Stage {
	return Map(func(ctx context.Context, i I) (pips.P[I, O], error) {
		o, err := zipper(ctx, i)
		if err != nil {
			return nil, err
		}
		return pips.NewP(i, o), nil
	})
}
