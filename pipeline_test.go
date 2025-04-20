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

	gt6Filter = apply.Filter(func(_ context.Context, i int) (bool, error) {
		return i > 6, nil
	})
)

func TestPipeline(t *testing.T) {
	t.Parallel()

	out := pips.New[any, int]().
		Then(subPipe).
		Then(lenMap).
		Then(apply.Batch(3)).
		Then(apply.Flatten()).
		Then(gt6Filter).
		Run(t.Context(), testhelpers.InputChan())

	testhelpers.RequireSuccessfulPiping(t, out, []int{8, 8, 10})
}
