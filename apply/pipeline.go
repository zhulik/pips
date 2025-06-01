package apply

import (
	"context"

	"github.com/zhulik/pips"
)

type PipelineConfig[I any, O any] struct {
	Pipeline *pips.Pipeline[I, O]
}

// PipelineStage represents a pipeline stage that embeds another pipeline.
// It allows for composition of pipelines, where one pipeline can be used as a stage within another pipeline.
// It takes items of type I as input and emits items of type O as output, using the embedded pipeline for processing.
type PipelineStage[I any, O any] Stage[PipelineConfig[I, O]]

// Run runs the stage.
// It processes input data by passing it to the embedded pipeline and then forwarding
// the results from that pipeline to the output channel of the current stage.
// This enables pipeline composition, where complex data processing can be broken down into smaller, reusable pipelines.
func (p PipelineStage[I, O]) Run(ctx context.Context, input <-chan pips.D[any], output chan<- pips.D[any]) {
	pips.MapToDChan(
		ctx,
		p.config.Pipeline.Run(ctx, pips.CastDChan[any, I](ctx, input)),
		output,
		func(_ context.Context, item O, out chan<- pips.D[any]) error {
			out <- pips.AnyD(item)
			return nil
		},
	)
}

// Pipeline creates a pipeline stage.
// A pipeline stage embeds another pipeline within the current pipeline, allowing for pipeline composition.
// It takes a pipeline that processes items of type I and produces items of type O,
// and returns a stage that can be used in another pipeline.
// This enables complex data processing to be broken down into smaller, reusable pipelines
// that can be composed together to form more complex processing flows.
func Pipeline[I any, O any](pipeline *pips.Pipeline[I, O]) pips.Stage {
	return PipelineStage[I, O]{
		config: PipelineConfig[I, O]{Pipeline: pipeline},
	}
}
