package pips

import (
	"context"
	"fmt"
)

// Stage is a unit of work in the pipeline.
// Stage should consume from the input channel and produce to the
// output channel. The stage must not close channels, must block.
// When using background routines: they must exit when the context is canceled or the input channel is closed.
type Stage func(context.Context, <-chan D[any], chan<- D[any])

// Pipeline is a sequence of stages.
type Pipeline[I any, O any] struct {
	stages []Stage
}

// New creates a new pipeline. The first stage is the input stage, the last stage is the output stage.
func New[I any, O any](stages ...Stage) *Pipeline[I, O] {
	return &Pipeline[I, O]{stages}
}

// Then adds a stage to the pipeline. The stage is added after the last stage.
func (p *Pipeline[I, O]) Then(stages ...Stage) *Pipeline[I, O] {
	p.stages = append(p.stages, stages...)
	return p
}

// Run runs the pipeline in the background.
func (p *Pipeline[I, O]) Run(ctx context.Context, input <-chan D[I]) OutChan[O] {
	inChan := CastDChan[I, any](ctx, input)

	for _, stage := range p.stages {
		newOut := make(chan D[any])

		go p.runStage(ctx, stage, inChan, newOut)

		inChan = newOut
	}

	return CastDChan[any, O](ctx, inChan)
}

func (p *Pipeline[I, O]) runStage(ctx context.Context, stage Stage, in <-chan D[any], out chan<- D[any]) {
	defer close(out)
	defer RecoverPanicAndSendToPipeline(out)

	stage(ctx, in, out)
}

func RecoverPanicAndSendToPipeline[T any](out chan<- D[T]) {
	if r := recover(); r != nil {
		var err error

		if e, ok := r.(error); ok {
			err = fmt.Errorf("%w: %w", ErrPanicInStage, e)
		} else {
			err = fmt.Errorf("%w: %+v", ErrPanicInStage, r)
		}

		out <- ErrD[T](err)
	}
}
