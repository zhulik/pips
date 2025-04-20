package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type pipelineStage[I any, O any] struct {
	pipeline *pips.Pipeline[I, O]
}

func (p pipelineStage[I, O]) Run(ctx context.Context, input <-chan pips.D[any]) <-chan pips.D[any] {
	return pips.CastDChan[O, any](ctx, p.pipeline.Run(ctx, pips.CastDChan[any, I](ctx, input)))
}

func Pipeline[I any, O any](pipeline *pips.Pipeline[I, O]) pips.Stage {
	return pipelineStage[I, O]{pipeline}
}
