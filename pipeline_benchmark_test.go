package pips_test

import (
	"context"
	"testing"

	"github.com/zhulik/pips"
	"github.com/zhulik/pips/apply"
)

var noOpMap = apply.Map(func(_ context.Context, i int) (int, error) {
	return i, nil
})

var noOpBatchMap = apply.Map(func(_ context.Context, i []int) ([]int, error) {
	return i, nil
})

// BenchmarkPipelineWith3NoOpMappingSteps measures how many items can be processed
// by a pipeline with 3 mapping steps that do nothing.
func BenchmarkPipelineWith3NoOpMappingSteps(b *testing.B) {
	input := make(chan pips.D[int])

	// Create a pipeline with 3 no-op mapping steps
	pipeline := pips.New[int, int]().
		Then(noOpMap).
		Then(apply.Batch(10)).
		Then(apply.Rebatch(30)).
		Then(noOpBatchMap).
		Then(apply.Flatten()).
		Then(noOpMap)

	out := pipeline.Run(context.Background(), input)

	// Reset the timer to exclude setup time
	b.ResetTimer()

	go func() {
		defer close(input)
		for i := 0; i < b.N; i++ {
			input <- pips.NewD(i)
		}
	}()

	// Wait for the result
	err := out.Wait(context.Background())
	if err != nil {
		b.Fatal(err)
	}
}
