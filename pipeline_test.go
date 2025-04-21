package pips_test

import (
	"context"
	"testing"

	"github.com/zhulik/pips"
	"github.com/zhulik/pips/apply"
	"github.com/zhulik/pips/testhelpers"
)

var (
	subPipe = apply.Pipeline(pips.New[string, string](apply.Map(func(_ context.Context, s string) (string, error) {
		return s + s, nil
	})))

	lenMap = apply.Map(func(_ context.Context, s string) (int, error) {
		return len(s), nil
	})

	nothingMap = apply.Map(func(_ context.Context, s []int) ([]int, error) {
		return s, nil
	})

	gt6Filter = apply.Filter(func(_ context.Context, i int) (bool, error) {
		return i > 6, nil
	})
)

func TestPipeline_Result(t *testing.T) {
	t.Parallel()

	out := pips.New[any, int]().
		Then(subPipe).
		Then(lenMap).
		Then(apply.Batch[int](3)).
		Then(nothingMap).
		Then(apply.Flatten[int]()).
		Then(gt6Filter).
		Run(t.Context(), testhelpers.InputChan())

	testhelpers.RequireSuccessfulPiping(t, out, []int{8, 8, 10})
}

func TestPipeline_NoResult(t *testing.T) {
	t.Parallel()

	out := pips.New[any, any]().
		Then(apply.Map[string, any](func(_ context.Context, _ string) (any, error) {
			return nil, nil
		})).
		Run(t.Context(), testhelpers.InputChan())

	testhelpers.RequireSuccessfulPiping(t, out, []any{nil, nil, nil, nil})
}
