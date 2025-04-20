package pips

import (
	"context"
)

// Stage is a unit of work in the pipeline.
type Stage interface {
	// Run is supposed to start a new goroutine which consumes from the input channel and produces to the
	// output channel. The stage is responsible for closing the output channel.
	// The stage finishes when the context is canceled or the input channel is closed.
	Run(context.Context, <-chan D[any]) <-chan D[any]
}

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
func (p *Pipeline[I, O]) Run(ctx context.Context, input <-chan I) <-chan D[O] {
	inChan := make(chan D[any])

	var prevOut <-chan D[any]

	prevOut = inChan

	for _, stage := range p.stages {
		prevOut = stage.Run(ctx, prevOut)
	}

	outCh := CastDChan[any, O](ctx, prevOut)

	go func() {
		MapToChan(ctx, input, inChan, AnyD)
		defer close(inChan)
	}()

	return outCh
}
